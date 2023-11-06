package web_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/adapter/web"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetKBWithIDDecoder(t *testing.T) {
	// Given
	var emptyBody []byte
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewGetKBWithIDDecoder(logger)
	givenKBID := "e65d36b3-ca19-4c33-8f59-917ab7399b44"

	getKBWithIDRequest := createHTTPRequest(t, emptyBody, http.MethodGet, "http://anyhost/kbs/"+givenKBID)
	getKBWithIDRequest = mux.SetURLVars(getKBWithIDRequest, map[string]string{
		"id": givenKBID,
	})

	expectedRequest := kbs.KBID("e65d36b3-ca19-4c33-8f59-917ab7399b44")

	// When
	got, err := decoder.Decode(ctx, getKBWithIDRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, got)
}

func TestSearchKBsDecoder(t *testing.T) {
	// Given
	var emptyBody []byte
	ctx := context.TODO()
	logger := newDummyLogger()
	pageSize := "15"
	pageNumber := "1"
	orderBy := "name"
	givenEventID := "drila"
	decoder := web.NewSearchKBsDecoder(logger)

	searchKBsRequest := createHTTPRequest(t, emptyBody, http.MethodGet, "http://anyhost/kbs")
	requestQuery := url.Values{}
	requestQuery.Add("page", pageNumber)
	requestQuery.Add("pagesize", pageSize)
	requestQuery.Add("event-id", givenEventID)
	requestQuery.Add("orderby", orderBy)
	searchKBsRequest.URL.RawQuery = requestQuery.Encode()

	expectedFilter := kbs.QueryFilter{
		EventID:     "drila",
		PageNumber:  1,
		RowsPerPage: 15,
		OrderBy:     "name",
	}

	// When
	got, err := decoder.Decode(ctx, searchKBsRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedFilter, got)
}

func TestCreateKBDecoder(t *testing.T) {
	// Given
	givenCreateBody := []byte(`{"user_id":"drila","username":"alird","content":"drila.alird","event_id":"drila.alird@lemail.com"}`)
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewCreateKBDecoder(logger)
	createKBRequest := createHTTPRequest(t, givenCreateBody, http.MethodPost, "http://anyhost/kbs")
	expectedCreateRequest := &kbs.NewKB{
		UserID:   "drila",
		UserName: "alird",
		Content:  "drila.alird",
		EventID:  "drila.alird@lemail.com",
	}
	// When
	got, err := decoder.Decode(ctx, createKBRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedCreateRequest, got)
}

func TestUpdateKBDecoder(t *testing.T) {
	// Given
	givenUpdateBody := []byte(`{"id":"388df4d7-75a4-4690-af0d-32a73899fdc3","user_id":"drila","username":"alird","content":"drila.alird","event_id":"drila.alird@lemail.com"}`)
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewUpdateKBDecoder(logger)
	updateKBRequest := createHTTPRequest(t, givenUpdateBody, http.MethodPut, "http://anyhost/kbs")
	expectedUpdateRequest := &kbs.UpdateKB{
		ID:       "388df4d7-75a4-4690-af0d-32a73899fdc3",
		UserID:   "drila",
		UserName: "alird",
		Content:  "drila.alird",
		EventID:  "drila.alird@lemail.com",
	}

	// When
	got, err := decoder.Decode(ctx, updateKBRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedUpdateRequest, got)
}

func TestDeleteKBDecoder(t *testing.T) {
	// Given
	var emptyBody []byte
	ctx := context.TODO()
	logger := newDummyLogger()
	decoder := web.NewDeleteKBDecoder(logger)
	givenKBID := "e65d36b3-ca19-4c33-8f59-917ab7399b44"

	deleteKBRequest := createHTTPRequest(t, emptyBody, http.MethodDelete, "http://anyhost/kbs/"+givenKBID)
	deleteKBRequest = mux.SetURLVars(deleteKBRequest, map[string]string{
		"id": givenKBID,
	})

	expectedRequest := kbs.KBID("e65d36b3-ca19-4c33-8f59-917ab7399b44")

	// When
	got, err := decoder.Decode(ctx, deleteKBRequest)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedRequest, got)
}

func createHTTPRequest(t *testing.T, body []byte, httpMethod, url string) *http.Request {
	t.Helper()

	newHTTPRequest, err := http.NewRequest(
		http.MethodGet,
		url,
		bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("unexpected error creating request: %s", err)
	}

	return newHTTPRequest
}

func newDummyLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
