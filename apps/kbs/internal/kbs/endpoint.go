package kbs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type GetKBWithIDEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type CreateKBEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type UpdateKBEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type DeleteKBEndpoint struct {
	service *Service
	logger  *slog.Logger
}

type SearchKBsEndpoint struct {
	service *Service
	logger  *slog.Logger
}

// Endpoints is a wrapper for endpoints
type Endpoints struct {
	GetKBWithIDEndpoint *GetKBWithIDEndpoint
	CreateKBEndpoint    *CreateKBEndpoint
	UpdateKBEndpoint    *UpdateKBEndpoint
	DeleteKBEndpoint    *DeleteKBEndpoint
	SearchKBsEndpoint   *SearchKBsEndpoint
}

// NewEndpoints Create the endpoints for kbs application.
func NewEndpoints(service *Service, logger *slog.Logger) Endpoints {
	return Endpoints{
		CreateKBEndpoint:    MakeCreateKBEndpoint(service, logger),
		UpdateKBEndpoint:    MakeUpdateKBEndpoint(service, logger),
		DeleteKBEndpoint:    MakeDeleteKBEndpoint(service, logger),
		GetKBWithIDEndpoint: MakeGetKBWithIDEndpoint(service, logger),
		SearchKBsEndpoint:   MakeSearchKBsEndpoint(service, logger),
	}
}

// MakeGetKBWithIDEndpoint create endpoint for get a kb with ID service.
func MakeGetKBWithIDEndpoint(srv *Service, logger *slog.Logger) *GetKBWithIDEndpoint {
	newNewEndpoint := GetKBWithIDEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeCreateKBEndpoint create endpoint for create kb service.
func MakeCreateKBEndpoint(srv *Service, logger *slog.Logger) *CreateKBEndpoint {
	newNewEndpoint := CreateKBEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeUpdateKBEndpoint create endpoint for update kb service.
func MakeUpdateKBEndpoint(srv *Service, logger *slog.Logger) *UpdateKBEndpoint {
	newNewEndpoint := UpdateKBEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeDeleteKBEndpoint create endpoint for the delete kb service.
func MakeDeleteKBEndpoint(srv *Service, logger *slog.Logger) *DeleteKBEndpoint {
	newNewEndpoint := DeleteKBEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

// MakeSearchKBsEndpoint kb endpoint to search kbs with filters.
func MakeSearchKBsEndpoint(srv *Service, logger *slog.Logger) *SearchKBsEndpoint {
	newNewEndpoint := SearchKBsEndpoint{
		service: srv,
		logger:  logger,
	}

	return &newNewEndpoint
}

func (g *GetKBWithIDEndpoint) Do(ctx context.Context, request any) (any, error) {
	kbID, ok := request.(KBID)
	if !ok {
		g.logger.Error("invalid kb id", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid kb id")
	}

	kbFound, err := g.service.QueryByID(ctx, kbID)
	if err != nil {
		g.logger.Error(
			"something went wrong trying to get a kb with the given id",
			slog.String("error", err.Error()),
		)
	}

	g.logger.Debug("find kb by id endpoint", slog.String("result", fmt.Sprintf("%+v", kbFound)))

	return newGetKBWithIDResult(kbFound, err), nil
}

func (c *CreateKBEndpoint) Do(ctx context.Context, request any) (any, error) {
	newKB, ok := request.(*NewKB)
	if !ok {
		c.logger.Error("invalid new kb type", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid new kb type")
	}

	newid, err := c.service.Create(ctx, *newKB)
	if err != nil {
		c.logger.Error(
			"something went wrong trying to create a kb with the given id",
			slog.String("error", err.Error()),
		)
	}
	return newCreateKBResult(newid, err), nil
}

func (u *UpdateKBEndpoint) Do(ctx context.Context, request any) (any, error) {
	updateKB, ok := request.(*UpdateKB)
	if !ok {
		u.logger.Error("invalid update kb type", slog.String("request", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid update kb type")
	}

	err := u.service.Update(ctx, *updateKB)
	if err != nil {
		u.logger.Error(
			"something went wrong trying to update a kb with the given id",
			slog.String("error", err.Error()),
		)
	}

	return newUpdateKBResult(err), nil
}

func (d *DeleteKBEndpoint) Do(ctx context.Context, request any) (any, error) {
	kbID, ok := request.(KBID)
	if !ok {
		d.logger.Error("invalid delete kb type", slog.String("received", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid kb id type")
	}

	err := d.service.Delete(ctx, kbID)
	if err != nil {
		d.logger.Error(
			"something went wrong trying to delete a kb with the given id",
			slog.String("error", err.Error()),
		)
	}
	return newDeleteKBResult(err), nil

}

func (s *SearchKBsEndpoint) Do(ctx context.Context, request any) (any, error) {
	kbFilters, ok := request.(QueryFilter)
	if !ok {
		s.logger.Error("invalid kb filters", slog.String("received", fmt.Sprintf("%t", request)))

		return nil, errors.New("invalid kb filters")
	}

	s.logger.Debug("querying kbs", slog.Any("filters", kbFilters))

	searchResult, err := s.service.Query(ctx, kbFilters)
	if err != nil {
		s.logger.Error(
			"something went wrong trying to search kbs with the given filter",
			slog.String("error", err.Error()),
		)
	}

	s.logger.Debug("search kbs endpoint", slog.String("result", fmt.Sprintf("%+v", searchResult)))

	return newSearchKBsDataResult(searchResult, err), nil
}
