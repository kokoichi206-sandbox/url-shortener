package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/kokoichi206-sandbox/url-shortener/handler"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
	"github.com/kokoichi206-sandbox/url-shortener/model/request"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Handler_GetOriginalURL(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		path            string
		makeMockUsecase func(m *MockUsecase)
		wantStatus      int
		wantLocation    string
		wantLog         string
	}{
		"success": {
			path: "/R0D",
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					SearchOriginalURL(gomock.Any(), "R0D").
					Return("https://example.com", nil)
			},
			wantStatus:   http.StatusMovedPermanently,
			wantLocation: "https://example.com",
		},
		"failure: not found": {
			path: "/RXX",
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					SearchOriginalURL(gomock.Any(), "RXX").
					Return("", apperr.ErrShortURLNotFound)
			},
			wantStatus: http.StatusNotFound,
		},
		"failure: server error": {
			path: "/RXX",
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					SearchOriginalURL(gomock.Any(), "RXX").
					Return("", errors.New("usecase error"))
			},
			wantStatus: http.StatusInternalServerError,
			wantLog:    "failed to exec usecase.SearchOriginalURL: usecase error",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := NewMockUsecase(ctrl)
			tc.makeMockUsecase(u)

			b := bytes.NewBuffer([]byte{})
			logger := logger.NewBasicLogger(b, "test", "searchOriginalURL")

			h := handler.New(logger, u)
			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)

			r.GET(
				"/:shortURL",
				handler.HandleWrapper(h.GetOriginalURL, logger),
			)

			req, _ := http.NewRequest(http.MethodGet, tc.path, nil)

			// Act
			r.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.wantStatus, recorder.Code, "status code should be equal")
			assert.Equal(t, tc.wantLocation, recorder.Header().Get("Location"), "location header should be equal")
			assert.True(t, strings.Contains(b.String(), tc.wantLog), "log should contain expected string")
		})
	}
}

func Test_Handler_GenerateURL(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		body            *request.CreateURL
		makeMockUsecase func(m *MockUsecase)
		wantStatus      int
		want            string
		wantLog         string
	}{
		"success": {
			body: &request.CreateURL{
				OriginalURL: "https://example.com",
			},
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					GenerateURL(gomock.Any(), "https://example.com").
					Return("R0D", nil)
			},
			wantStatus: http.StatusOK,
			want:       `{"short_url":"R0D"}`,
		},
		"failure: body empty": {
			makeMockUsecase: func(m *MockUsecase) {},
			wantStatus:      http.StatusBadRequest,
			want:            `{"error":"request body is invalid"}`,
		},
		"failure: usecase error": {
			body: &request.CreateURL{
				OriginalURL: "https://example.com",
			},
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					GenerateURL(gomock.Any(), "https://example.com").
					Return("", errors.New("usecase error"))
			},
			wantStatus: http.StatusInternalServerError,
			want:       `{"error":"internal server error"}`,
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			u := NewMockUsecase(ctrl)
			tc.makeMockUsecase(u)

			b := bytes.NewBuffer([]byte{})
			logger := logger.NewBasicLogger(b, "test", "generateURL")

			h := handler.New(logger, u)
			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)

			r.POST(
				"/api/v1/urls",
				handler.HandleWrapper(h.GenerateURL, logger),
			)

			var buf bytes.Buffer
			if tc.body != nil {
				_ = json.NewEncoder(&buf).Encode(tc.body)
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/urls", &buf)

			// Act
			r.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.wantStatus, recorder.Code, "status code should be equal")
			assert.Equal(t, tc.want, string(recorder.Body.Bytes()), "response body should be equal")
			assert.True(t, strings.Contains(b.String(), tc.wantLog), "log should contain expected string")
		})
	}
}
