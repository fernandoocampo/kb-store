package web

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
	"github.com/gorilla/mux"
)

type GetKBWithIDDecoder struct {
	logger *slog.Logger
}

type SearchKBsDecoder struct {
	logger *slog.Logger
}

type CreateKBDecoder struct {
	logger *slog.Logger
}

type UpdateKBDecoder struct {
	logger *slog.Logger
}

type DeleteKBDecoder struct {
	logger *slog.Logger
}

type KBDecoders struct {
	GetByIDDecoder *GetKBWithIDDecoder
	SearchDecoder  *SearchKBsDecoder
	CreateDecoder  *CreateKBDecoder
	UpdateDecoder  *UpdateKBDecoder
	DeleteDecoder  *DeleteKBDecoder
}

func NewKBDecoders(logger *slog.Logger) KBDecoders {
	newDecoders := KBDecoders{
		GetByIDDecoder: NewGetKBWithIDDecoder(logger),
		SearchDecoder:  NewSearchKBsDecoder(logger),
		CreateDecoder:  NewCreateKBDecoder(logger),
		UpdateDecoder:  NewUpdateKBDecoder(logger),
		DeleteDecoder:  NewDeleteKBDecoder(logger),
	}

	return newDecoders
}

func NewGetKBWithIDDecoder(logger *slog.Logger) *GetKBWithIDDecoder {
	newDecoder := GetKBWithIDDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewSearchKBsDecoder(logger *slog.Logger) *SearchKBsDecoder {
	newDecoder := SearchKBsDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewCreateKBDecoder(logger *slog.Logger) *CreateKBDecoder {
	newDecoder := CreateKBDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewUpdateKBDecoder(logger *slog.Logger) *UpdateKBDecoder {
	newDecoder := UpdateKBDecoder{
		logger: logger,
	}

	return &newDecoder
}

func NewDeleteKBDecoder(logger *slog.Logger) *DeleteKBDecoder {
	newDecoder := DeleteKBDecoder{
		logger: logger,
	}

	return &newDecoder
}

func (g *GetKBWithIDDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	kbIDParam, ok := v["id"]
	if !ok {
		return nil, errors.New("kb ID was not provided")
	}
	return kbs.KBID(kbIDParam), nil
}

func (g *DeleteKBDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	v := mux.Vars(r)
	kbIDParam, ok := v["id"]
	if !ok {
		return nil, errors.New("kb ID was not provided")
	}
	return kbs.KBID(kbIDParam), nil
}

func (s *SearchKBsDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	filterRequest := SearchKBFilter{
		Page:     1,
		PageSize: 10,
	}

	filters := r.URL.Query()

	if v, ok := filters["event-id"]; ok {
		filterRequest.EventID = v[0]
	}

	if v, ok := filters["page"]; ok {
		page, err := strconv.Atoi(v[0])
		if err != nil {
			s.logger.Error("invalid page parameter, it must be an integer", "error", err)
			page = 1
		}
		filterRequest.Page = uint8(page)
	}
	if v, ok := filters["pagesize"]; ok {
		pageSize, err := strconv.Atoi(v[0])
		if err != nil {
			s.logger.Error("level", "ERROR", "invalid page size parameter, it must be an integer", "error", err)
			pageSize = 10
		}
		filterRequest.PageSize = uint8(pageSize)
	}

	if v, ok := filters["orderby"]; ok {
		filterRequest.OrderBy = v[0]
	}

	filter := filterRequest.toSearchKBFilter()

	return filter, nil
}

func (c *CreateKBDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "decoding new kb request")
	var req NewKB
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("level", "ERROR", "new kb request could not be decoded. Request: %q because of: %s", string(body), err.Error())
		return nil, err
	}

	log.Println("level", "DEBUG", "msg", "kb request was decoded", "request", req)

	domainKB := req.toKB()

	return domainKB, nil
}

func (u *UpdateKBDecoder) Decode(ctx context.Context, r *http.Request) (interface{}, error) {
	log.Println("level", "DEBUG", "msg", "decoding update kb request")
	var req UpdateKB
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println("level", "ERROR", "update kb request could not be decoded. Request: %q because of: %s", string(body), err.Error())
		return nil, err
	}

	log.Println("level", "DEBUG", "msg", "kb request was decoded", "request", req)

	domainKB := req.toKB()

	return domainKB, nil
}
