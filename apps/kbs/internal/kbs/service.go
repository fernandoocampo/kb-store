package kbs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

// Storer defines persistence behavior
type Storer interface {
	Save(ctx context.Context, newKB KB) error
	Update(ctx context.Context, kb UpdateKB) error
	Delete(ctx context.Context, kb KB) error
	Query(ctx context.Context, filter QueryFilter) (SearchKBsResult, error)
	// QueryByID find and return a kb with the given id.
	// If kb does not exist it returns a nil kb and nil error.
	QueryByID(ctx context.Context, id KBID) (*KB, error)
}

// ServiceSetup contains service metadata.
type ServiceSetup struct {
	Storer Storer
	Logger *slog.Logger
}

// Service implements kbs business logic.
type Service struct {
	storer Storer
	logger *slog.Logger
}

var (
	errSaveKB    = errors.New("unable to save kb in the repository")
	errQueryKB   = errors.New("unable to query kb")
	errQueryKBs  = errors.New("unable to query kbs")
	errDeleteKB  = errors.New("unable to delete kb")
	errUpdateKB  = errors.New("unable to update kb in the repository")
	errEmptyKBID = errors.New("kb id cannot be empty")
)

// NewService create a new kbs service.
func NewService(settings ServiceSetup) *Service {
	newService := Service{
		logger: settings.Logger,
		storer: settings.Storer,
	}

	return &newService
}

// Create create a kb and store it in a database.
func (s *Service) Create(ctx context.Context, newKB NewKB) (KBID, error) {
	kb := buildNewKB(newKB)

	err := s.storer.Save(ctx, kb)
	if err != nil {
		s.logger.Error("unable to create kb", slog.String("error", err.Error()))

		return EmptyKBID, errSaveKB
	}

	s.logger.Debug(
		"kb was created",
		slog.String("id", kb.ID.String()),
	)

	return kb.ID, nil
}

// Update update a kb in a database.
func (s *Service) Update(ctx context.Context, kb UpdateKB) error {
	err := validKBToUpdate(kb)
	if err != nil {
		return fmt.Errorf("unable to update kb: %w", err)
	}

	kb.fillUpdateTime()

	err = s.storer.Update(ctx, kb)
	if err != nil {
		s.logger.Error("unable to update kb", slog.String("error", err.Error()))

		return errUpdateKB
	}

	return nil
}

func (s *Service) QueryByID(ctx context.Context, id KBID) (*KB, error) {
	if id == EmptyKBID {
		return nil, errEmptyKBID
	}

	kb, err := s.storer.QueryByID(ctx, id)
	if err != nil {
		s.logger.Error(
			"unable to query kb by id",
			slog.String("id", fmt.Sprintf("%+v", id)),
			slog.String("error", err.Error()))

		return nil, errQueryKB
	}

	return kb, nil
}

// Delete detele a kb from database.
func (s *Service) Delete(ctx context.Context, id KBID) error {
	kb, err := s.QueryByID(ctx, id)
	if err != nil {
		return errDeleteKB
	}

	if kb == nil {
		s.logger.Info(
			"unable to delete kb cause it does not exist",
			slog.String("id", fmt.Sprintf("%+v", id)),
		)

		return nil
	}

	err = s.storer.Delete(ctx, *kb)
	if err != nil {
		s.logger.Error("unable to delete kb",
			slog.String("id", fmt.Sprintf("%+v", id)),
			slog.String("error", err.Error()))

		return errUpdateKB
	}

	return nil
}

func (s *Service) Query(ctx context.Context, filter QueryFilter) (SearchKBsResult, error) {
	s.logger.Debug("querying kb on kbs.Service")

	if filter.isInvalid() {
		s.logger.Debug("filter is invalid", slog.String("data", fmt.Sprintf("%+v", filter)))

		return SearchKBsResult{}, nil
	}

	filter.fillDefaultValues()

	result, err := s.storer.Query(ctx, filter)
	if err != nil {
		s.logger.Error(
			"unable to query kbs by filter",
			slog.String("filter", fmt.Sprintf("%+v", filter)),
			slog.String("error", err.Error()))

		return SearchKBsResult{}, errQueryKB
	}

	return result, nil
}
