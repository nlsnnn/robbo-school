package get_test

import (
	"api/internal/http-server/handlers/courses/get"
	"api/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCourseGetter struct {
	mock.Mock
}

func (m *MockCourseGetter) GetCourseByID(ctx context.Context, id int) (storage.Course, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(storage.Course), args.Error(1)
}

func (m *MockCourseGetter) GetCourses(ctx context.Context) ([]storage.Course, error) {
	args := m.Called(ctx)
	return args.Get(0).([]storage.Course), args.Error(1)
}

func TestGetByIDHandler(t *testing.T) {
	cases := []struct {
		name           string
		courseID       string
		respError      string
		mockError      error
		mockResp       storage.Course
		expectedStatus int
	}{
		{
			name:     "Success",
			courseID: "1",
			mockResp: storage.Course{
				Id:          1,
				Title:       "Test Course",
				Description: "Test Description",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID - Not a number",
			courseID:       "abc",
			respError:      "invalid course id",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid ID - Zero",
			courseID:       "0",
			respError:      "invalid course id",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			courseID:       "999",
			respError:      "course not found",
			mockError:      storage.ErrCourseNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Internal Error",
			courseID:       "1",
			respError:      "internal error",
			mockError:      errors.New("unexpected error"),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courseGetterMock := new(MockCourseGetter)

			if tc.respError == "" || tc.mockError != nil {
				if tc.courseID != "abc" && tc.courseID != "0" {
					courseGetterMock.On("GetCourseByID", mock.Anything, mock.AnythingOfType("int")).
						Return(tc.mockResp, tc.mockError).
						Once()
				}
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))

			r := chi.NewRouter()
			r.Get("/courses/{id}", get.GetByID(logger, courseGetterMock))

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/courses/%s", tc.courseID), nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)

			body := rr.Body.String()

			var resp get.CourseResponse

			err := json.Unmarshal([]byte(body), &resp)
			require.NoError(t, err)

			if tc.respError != "" {
				require.Equal(t, tc.respError, resp.Error)
			} else {
				require.Equal(t, "OK", resp.Status)
				require.Equal(t, tc.mockResp.Id, resp.Course.Id)
				require.Equal(t, tc.mockResp.Title, resp.Course.Title)
			}
		})
	}
}
