package dynamodb_test

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"testing"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/adapter/dynamodb"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var integration = flag.Bool("integration", false, "")

func TestSaveKB(t *testing.T) {
	skipNonIntegrationTest(t)

	// Given
	newKB := kbs.KB{
		ID:       newKBID(),
		UserID:   "cb5c9d13-daf8-4720-87eb-80f034b7528f",
		UserName: "mono.mario",
		Content:  "what an amazing show",
		EventID:  "6763fe1b-9391-49f2-acf1-5069e2a9cb21",
	}

	ctx := context.Background()

	store := newStore(ctx, t)

	// When
	err := store.Save(ctx, newKB)

	// Then
	assert.NoError(t, err)
}

func TestFindKBByID(t *testing.T) {
	skipNonIntegrationTest(t)

	// Given
	kbID := newKBID()

	newKB := kbs.KB{
		ID:       kbID,
		UserID:   "cb5c9d13-daf8-4720-87eb-80f034b7528f",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "6763fe1b-9391-49f2-acf1-5069e2a9cb21",
	}

	expectedKB := &kbs.KB{
		ID:       kbID,
		UserID:   "cb5c9d13-daf8-4720-87eb-80f034b7528f",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "6763fe1b-9391-49f2-acf1-5069e2a9cb21",
	}

	ctx := context.Background()

	store := newStore(ctx, t)

	saveKB(t, store, newKB)

	// When
	got, err := store.QueryByID(ctx, kbID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedKB, got)
}

func TestFindKBByEventID(t *testing.T) {
	skipNonIntegrationTest(t)

	// Given
	kbID := newKBID()

	newKB := kbs.KB{
		ID:       kbID,
		UserID:   "c8b21266-ec50-4cef-875d-1b651738729a",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "489cc9d6-28bb-4052-a11f-0ecc8c164613",
	}

	expectedKB := &kbs.KB{
		ID:       kbID,
		UserID:   "c8b21266-ec50-4cef-875d-1b651738729a",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "489cc9d6-28bb-4052-a11f-0ecc8c164613",
	}

	ctx := context.Background()

	store := newStore(ctx, t)

	saveKB(t, store, newKB)

	// When
	got, err := store.QueryByID(ctx, kbID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedKB, got)
}

func TestDeleteKB(t *testing.T) {
	skipNonIntegrationTest(t)

	// Given
	kbID := newKBID()

	kbToDelete := kbs.KB{
		ID:       kbID,
		UserID:   "Mono",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "mono.mario@location.com",
	}

	ctx := context.Background()

	store := newStore(ctx, t)

	saveKB(t, store, kbToDelete)

	// When
	err := store.Delete(ctx, kbToDelete)

	// Then
	assert.NoError(t, err)
	got, err := store.QueryByID(ctx, kbID)
	assert.NoError(t, err)
	assert.Nil(t, got)
}

func TestUpdateKB(t *testing.T) {
	skipNonIntegrationTest(t)

	// Given
	kbID := newKBID()

	existingKB := kbs.KB{
		ID:       kbID,
		UserID:   "Mono",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "mono.mario@location.com",
	}

	expectedKB := kbs.KB{
		ID:       kbID,
		UserID:   "Bear",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "mono.mario@location.com",
	}

	kbToUpdate := kbs.UpdateKB{
		ID:       kbID,
		UserID:   "Bear",
		UserName: "Mario",
		Content:  "mono.mario",
		EventID:  "mono.mario@location.com",
	}

	ctx := context.Background()

	store := newStore(ctx, t)

	saveKB(t, store, existingKB)

	// When
	err := store.Update(ctx, kbToUpdate)

	// Then
	assert.NoError(t, err)
	got, err := store.QueryByID(ctx, kbID)
	assert.NoError(t, err)
	assert.Equal(t, &expectedKB, got)
}

func newStore(ctx context.Context, t *testing.T) *dynamodb.Client {
	t.Helper()

	setup := dynamodb.Setup{
		Logger:   newLogger(),
		Region:   "us-east-1",
		Endpoint: "http://localhost:4566",
	}

	store, err := dynamodb.NewClient(ctx, setup)
	if err != nil {
		t.Fatalf(err.Error())
	}

	return store
}

func saveKB(t *testing.T, store *dynamodb.Client, newKB kbs.KB) {
	t.Helper()

	err := store.Save(context.Background(), newKB)
	if err != nil {
		t.Fatalf("unexpected error saving a new kb: %s", err)
	}
}

func skipNonIntegrationTest(t *testing.T) {
	t.Helper()

	if !*integration {
		t.Skip("this is an integration test")
	}
}

func newLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func newKBID() kbs.KBID {
	return kbs.KBID(uuid.New().String())
}
