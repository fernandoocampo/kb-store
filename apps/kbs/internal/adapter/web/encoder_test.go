package web_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fernandoocampo/kb-store/apps/kbs/internal/adapter/web"
	"github.com/fernandoocampo/kb-store/apps/kbs/internal/kbs"
	"github.com/stretchr/testify/assert"
)

func TestEncodeCreateKB(t *testing.T) {
	// Given
	givenEndpointResult := kbs.CreateKBResult{
		ID:  kbs.KBID("82853922-4481-4a95-8691-30f36c61e45a"),
		Err: "",
	}

	expectedEncodedResult := web.Result{
		Success: true,
		Errors:  nil,
		Data:    "82853922-4481-4a95-8691-30f36c61e45a",
	}

	encoder := web.NewCreateKBEncoder(newDummyLogger())

	ctx := context.TODO()
	recorder := httptest.NewRecorder()

	// When
	err := encoder.Encode(ctx, recorder, givenEndpointResult)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, expectedEncodedResult, createWebResult(t, recorder.Body, kbs.EmptyKBID))
}

func TestEncodeGetKBWithID(t *testing.T) {
	// Given
	givenEndpointResult := kbs.GetKBWithIDResult{
		KB: &kbs.KB{
			ID:       kbs.KBID("82853922-4481-4a95-8691-30f36c61e45a"),
			UserID:   "drila",
			UserName: "alird",
			Content:  "drila.alird",
			EventID:  "drila.alird@lemail.com",
		},
		Err: "",
	}

	expectedEncodedResult := web.Result{
		Success: true,
		Errors:  nil,
		Data: &web.KB{
			ID:       "82853922-4481-4a95-8691-30f36c61e45a",
			UserID:   "drila",
			UserName: "alird",
			Content:  "drila.alird",
			EventID:  "drila.alird@lemail.com",
		},
	}

	encoder := web.NewGetKBWithIDEncoder(newDummyLogger())

	ctx := context.TODO()
	recorder := httptest.NewRecorder()

	// When
	err := encoder.Encode(ctx, recorder, givenEndpointResult)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.Equal(t, expectedEncodedResult, createWebResult(t, recorder.Body, &web.KB{}))
}

func createWebResult(t *testing.T, body io.Reader, data any) web.Result {
	t.Helper()

	var result web.Result
	result.Data = data
	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)

		return result
	}

	return result
}
