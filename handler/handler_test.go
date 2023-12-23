package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
	"github.com/kokoichi206-sandbox/url-shortener/handler"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Handler_Health(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		makeMockUsecase func(m *MockUsecase)
		wantStatus      int
		want            string
		wantLog         string
	}{
		"success": {
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					Health(gomock.Any()).
					Return(nil)
			},
			wantStatus: http.StatusOK,
			want:       `{"health":"ok"}`,
		},
		"failure: ng": {
			makeMockUsecase: func(m *MockUsecase) {
				m.
					EXPECT().
					Health(gomock.Any()).
					Return(errors.New("usecase error"))
			},
			wantStatus: http.StatusInternalServerError,
			want:       `{"error":"internal server error"}`,
			wantLog:    "failed to health check: usecase error",
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
			logger := logger.NewBasicLogger(b, "test", "health")

			h := handler.New(logger, u)
			recorder := httptest.NewRecorder()
			_, r := gin.CreateTestContext(recorder)

			r.GET(
				"/api/v1/health",
				handler.HandleWrapper(h.Health, logger),
			)

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/health", nil)

			// Act
			r.ServeHTTP(recorder, req)

			// Assert
			assert.Equal(t, tc.wantStatus, recorder.Code, "status code should be equal")
			assert.Equal(t, tc.want, string(recorder.Body.Bytes()), "response body should be equal")
			assert.True(t, strings.Contains(b.String(), tc.wantLog), "log should contain expected string")
		})
	}
}
