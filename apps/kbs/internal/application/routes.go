package application

import (
	"net/http"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/adapter/web"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
	"github.com/gorilla/mux"
)

type kbsRouter struct {
	router    *mux.Router
	endpoints kbs.Endpoints
	decoders  web.KBDecoders
	encoders  web.KBEncoders
}

func newKBsRouter(kbsRouter kbsRouter) http.Handler {
	kbsRouter.router.Methods(http.MethodPost).Path("/kbs").Handler(
		web.NewHandler().
			WithEndpoint(kbsRouter.endpoints.CreateKBEndpoint).
			WithDecoder(kbsRouter.decoders.CreateDecoder).
			WithEncoder(kbsRouter.encoders.CreateEncoder),
	)

	kbsRouter.router.Methods(http.MethodPut).Path("/kbs").Handler(
		web.NewHandler().
			WithEndpoint(kbsRouter.endpoints.UpdateKBEndpoint).
			WithDecoder(kbsRouter.decoders.UpdateDecoder).
			WithEncoder(kbsRouter.encoders.UpdateEncoder),
	)

	kbsRouter.router.Methods(http.MethodDelete).Path("/kbs/{id}").Handler(
		web.NewHandler().
			WithEndpoint(kbsRouter.endpoints.DeleteKBEndpoint).
			WithDecoder(kbsRouter.decoders.DeleteDecoder).
			WithEncoder(kbsRouter.encoders.DeleteEncoder),
	)

	kbsRouter.router.Methods(http.MethodGet).Path("/kbs/{id}").Handler(
		web.NewHandler().
			WithEndpoint(kbsRouter.endpoints.GetKBWithIDEndpoint).
			WithDecoder(kbsRouter.decoders.GetByIDDecoder).
			WithEncoder(kbsRouter.encoders.GetByIDEncoder),
	)

	kbsRouter.router.Methods(http.MethodGet).Path("/kbs").Handler(
		web.NewHandler().
			WithEndpoint(kbsRouter.endpoints.SearchKBsEndpoint).
			WithDecoder(kbsRouter.decoders.SearchDecoder).
			WithEncoder(kbsRouter.encoders.SearchEncoder),
	)

	return kbsRouter.router
}
