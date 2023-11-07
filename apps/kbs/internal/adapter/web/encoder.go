package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
)

type GetKBWithIDEncoder struct {
	logger *slog.Logger
}

type SearchKBsEncoder struct {
	logger *slog.Logger
}

type CreateKBEncoder struct {
	logger *slog.Logger
}

type UpdateKBEncoder struct {
	logger *slog.Logger
}

type DeleteKBEncoder struct {
	logger *slog.Logger
}

type KBEncoders struct {
	GetByIDEncoder *GetKBWithIDEncoder
	SearchEncoder  *SearchKBsEncoder
	CreateEncoder  *CreateKBEncoder
	UpdateEncoder  *UpdateKBEncoder
	DeleteEncoder  *DeleteKBEncoder
}

var (
	errUnableToEncodeResult = errors.New("unable to encode the result")
)

func NewKBEncoders(logger *slog.Logger) KBEncoders {
	newEncoders := KBEncoders{
		GetByIDEncoder: NewGetKBWithIDEncoder(logger),
		SearchEncoder:  NewSearchKBsEncoder(logger),
		CreateEncoder:  NewCreateKBEncoder(logger),
		UpdateEncoder:  NewUpdateKBEncoder(logger),
		DeleteEncoder:  NewDeleteKBEncoder(logger),
	}

	return newEncoders
}

func NewGetKBWithIDEncoder(logger *slog.Logger) *GetKBWithIDEncoder {
	newEncoder := GetKBWithIDEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewSearchKBsEncoder(logger *slog.Logger) *SearchKBsEncoder {
	newEncoder := SearchKBsEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewCreateKBEncoder(logger *slog.Logger) *CreateKBEncoder {
	newEncoder := CreateKBEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewUpdateKBEncoder(logger *slog.Logger) *UpdateKBEncoder {
	newEncoder := UpdateKBEncoder{
		logger: logger,
	}

	return &newEncoder
}

func NewDeleteKBEncoder(logger *slog.Logger) *DeleteKBEncoder {
	newEncoder := DeleteKBEncoder{
		logger: logger,
	}

	return &newEncoder
}

func (c *CreateKBEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(kbs.CreateKBResult)
	if !ok {
		c.logger.Error("cannot transform to kbs.CreateKBResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build create kb response")
	}

	err := encodeResultWithJSON(w, toCreateKBResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode create kb result: %w", err)
	}

	return nil
}

func (u *UpdateKBEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(kbs.UpdateKBResult)
	if !ok {
		u.logger.Error("cannot transform to kbs.UpdateKBResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build update kb response")
	}

	err := encodeResultWithJSON(w, toUpdateKBResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode update kb result: %w", err)
	}

	return nil
}

func (u *DeleteKBEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(kbs.DeleteKBResult)
	if !ok {
		u.logger.Error("cannot transform to kbs.DeleteKBResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build delete kb response")
	}

	err := encodeResultWithJSON(w, toDeleteKBResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode delete kb result: %w", err)
	}

	return nil
}

func (g *GetKBWithIDEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(kbs.GetKBWithIDResult)
	if !ok {
		g.logger.Error("cannot transform to kbs.GetKBWithIDResult", "received", fmt.Sprintf("%+v", response))
		return errors.New("cannot build get kb response")
	}

	err := encodeResultWithJSON(w, toGetKBWithIDResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode get kb by id result: %w", err)
	}

	return nil
}

func (s *SearchKBsEncoder) Encode(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	result, ok := response.(kbs.SearchKBsDataResult)
	if !ok {
		s.logger.Error("cannot transform to kbs.SearchKBsDataResult", "received", fmt.Sprintf("%T", response))
		return errors.New("cannot build search kbs response")
	}

	err := encodeResultWithJSON(w, toSearchKBsResponse(result))
	if err != nil {
		return fmt.Errorf("unable to encode search kbs result: %w", err)
	}

	return nil
}

func encodeResultWithJSON(w http.ResponseWriter, kb Result) error {
	w.Header().Set("Content-Type", "application/json")

	if kb.Failed() {
		w.WriteHeader(http.StatusInternalServerError)
	}

	err := json.NewEncoder(w).Encode(kb)
	if err != nil {
		return errUnableToEncodeResult
	}

	return nil
}
