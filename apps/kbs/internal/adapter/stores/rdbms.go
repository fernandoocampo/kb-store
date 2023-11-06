package stores

import (
	"context"
	"log/slog"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
)

type Setup struct {
	Logger *slog.Logger
}

// Store handles logic to persist data from this microservice.
type Store struct {
	logger *slog.Logger
}

func NewStore(setup Setup) *Store {
	newStore := Store{
		logger: setup.Logger,
	}

	return &newStore
}

func (s *Store) Save(ctx context.Context, newKB kbs.KB) error {
	s.logger.Info("Saving new kb in database")
	return nil
}

func (s *Store) Update(ctx context.Context, kb kbs.UpdateKB) error {
	s.logger.Info("Updating new kb in database")
	return nil
}

func (s *Store) Delete(ctx context.Context, kb kbs.KB) error {
	s.logger.Info("Deleting new kb in database")
	return nil
}

func (s *Store) Query(ctx context.Context, filter kbs.QueryFilter) (kbs.SearchKBsResult, error) {
	s.logger.Info("Querying kbs in database")
	result := kbs.SearchKBsResult{
		KBs: []kbs.KB{
			{
				ID:     kbs.KBID("56016eaf-5e15-44db-839c-ef4f7f9df437"),
				UserID: "Drila",
			},
			{
				ID:     kbs.KBID("ec665f5e-da4e-4f51-bc4c-310dd7cc9590"),
				UserID: "Michael",
			},
		},
		Total:       2,
		Page:        filter.PageNumber,
		RowsPerPage: filter.RowsPerPage,
	}
	return result, nil
}

func (s *Store) QueryByID(ctx context.Context, id kbs.KBID) (*kbs.KB, error) {
	s.logger.Info("Querying kb by id in database")
	kb := kbs.KB{
		ID:     kbs.KBID("56016eaf-5e15-44db-839c-ef4f7f9df437"),
		UserID: "Drila",
	}

	return &kb, nil
}
