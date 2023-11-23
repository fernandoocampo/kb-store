package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/adapter/dynamodb"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/adapter/web"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/setups"
)

// Event contains an application event.
type Event struct {
	KB    string
	Error error
}

// Setup contains application metadata
type Setup struct {
	Version    string
	BuildDate  string
	CommitHash string
}

// Server is the server of our application.
type Server struct {
	logger     *slog.Logger
	store      kbs.Storer
	setup      setups.Application
	version    string
	buildDate  string
	commitHash string
}

var (
	errStartingApplication = errors.New("unable to start application")
)

func NewServer() *Server {
	newServer := Server{
		version:    setups.Version,
		buildDate:  setups.BuildDate,
		commitHash: setups.CommitHash,
	}

	return &newServer
}

func (s *Server) Run() error {
	s.notifyStart()

	confError := s.loadConfiguration()
	if confError != nil {
		return errStartingApplication
	}

	loggerError := s.initializeLogger()
	if loggerError != nil {
		return errStartingApplication
	}

	s.logger.Debug("application configuration", "parameters", fmt.Sprintf("%+v", s.setup))

	s.logger.Info("starting database connection")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.createDynamodbStorer(ctx)
	if err != nil {
		return errStartingApplication
	}

	kbServiceSetup := kbs.ServiceSetup{
		Storer: s.store,
		Logger: s.logger,
	}
	kbService := kbs.NewService(kbServiceSetup)
	kbEndpoints := kbs.NewEndpoints(kbService, s.logger)

	eventStream := make(chan Event)
	s.listenToOSSignal(eventStream)
	s.startWebServer(kbEndpoints, eventStream)

	eventKB := <-eventStream
	s.logger.Info("ending server", "event", eventKB.KB)

	if eventKB.Error != nil {
		s.logger.Error("ending server with error", "error", eventKB.Error)

		return errStartingApplication
	}

	return nil
}

func (s *Server) initializeLogger() error {
	logLevel := slog.LevelDebug

	if s.setup.LogLevel == setups.ProductionLog {
		logLevel = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level: logLevel,
	}

	loggerHandler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	logger := slog.New(loggerHandler)

	logger.Info(fmt.Sprintf("using %q log level", handlerOptions.Level.Level().String()))

	slog.SetDefault(logger)

	s.logger = logger

	return nil
}

func (s *Server) notifyStart() {
	log.Println(
		"starting service",
		"version:", s.version,
		"commit:", s.commitHash,
		"build date:", s.buildDate,
	)
}

// Stop stop application, take advantage of this to clean resources
func (s *Server) Stop() {
	s.logger.Info("stopping the application")
}

func (s *Server) listenToOSSignal(eventStream chan<- Event) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		osSignal := (<-c).String()
		event := Event{
			KB: osSignal,
		}
		eventStream <- event
	}()
}

// startWebServer starts the web server.
func (s *Server) startWebServer(kbEndpoints kbs.Endpoints, eventStream chan<- Event) {
	go func() {
		s.logger.Info("starting http server", slog.String("port", s.setup.ApplicationPort))
		router := kbsRouter{
			router:    web.NewRouter(),
			endpoints: kbEndpoints,
			decoders:  web.NewKBDecoders(s.logger),
			encoders:  web.NewKBEncoders(s.logger),
		}
		handler := newKBsRouter(router)
		err := http.ListenAndServe(s.setup.ApplicationPort, handler)
		if err != nil {
			eventStream <- Event{
				KB:    "web server was ended with error",
				Error: err,
			}
			return
		}
		eventStream <- Event{
			KB: "web server was ended",
		}
	}()
}

func (s *Server) loadConfiguration() error {
	applicationSetUp, err := setups.Load()
	if err != nil {
		log.Println("level", "ERROR", "msg", "application setup could not be loaded", "error", err)

		return errors.New("application setup could not be loaded")
	}
	s.setup = applicationSetUp
	return nil
}

func (s *Server) createDynamodbStorer(ctx context.Context) error {
	storeSetup := dynamodb.Setup{
		Logger:   s.logger,
		Region:   s.setup.Repository.Region,
		Endpoint: s.setup.Repository.Endpoint,
	}

	storer, err := dynamodb.NewClient(ctx, storeSetup)
	if err != nil {
		s.logger.Error("unable to create dynamodb client at startup", slog.String("error", err.Error()))

		return errors.New("application dynamodb storage could not be created")
	}

	s.store = storer

	return nil
}
