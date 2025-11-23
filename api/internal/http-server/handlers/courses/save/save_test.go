package save_test

import (
	"api/internal/http-server/handlers/courses/save"
	"api/internal/storage"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCourseSaver struct {
	mock.Mock
}

func (m *MockCourseSaver) SaveCourse(ctx context.Context, course storage.Course) (storage.Course, error) {
	args := m.Called(ctx, course)
	return args.Get(0).(storage.Course), args.Error(1)
}

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string
		title     string
		desc      string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			title: "New Course",
			desc:  "New Description",
		},
		{
			name:      "Empty Title",
			title:     "",
			desc:      "Description",
			respError: "field Title is a required field",
		},
		{
			name:      "Empty Description",
			title:     "Title",
			desc:      "",
			respError: "field Description is a required field",
		},
		{
			name:      "Save Error",
			title:     "Error Course",
			desc:      "Error Description",
			respError: "failed to save course",
			mockError: errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			courseSaverMock := new(MockCourseSaver)

			if tc.respError == "" || tc.mockError != nil {
				courseSaverMock.On("SaveCourse", mock.Anything, mock.MatchedBy(func(c storage.Course) bool {
					return c.Title == tc.title && c.Description == tc.desc
				})).Return(storage.Course{
					Id:          1,
					Title:       tc.title,
					Description: tc.desc,
				}, tc.mockError).Once()
			}

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := save.New(logger, courseSaverMock)

			input := fmt.Sprintf(`{"title": "%s", "description": "%s"}`, tc.title, tc.desc)

			req, err := http.NewRequest(http.MethodPost, "/courses", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)

			body := rr.Body.String()

			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
