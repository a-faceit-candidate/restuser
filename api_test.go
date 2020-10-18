package restuser_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a-faceit-candidate/restuser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_CreateUser(t *testing.T) {
	someUserToCreate := &restuser.User{
		Name:    "Pepe",
		Email:   "pepe@faceit.com",
		Country: "fr",
	}

	someCreatedUser := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-02T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	someErrorResponse := &restuser.ErrorResponse{Message: "everything is wrong"}

	for _, tc := range []struct {
		name                string
		srv                 testServerExpectations
		expectedReturnValue *restuser.User
		expectedError       error
	}{
		{
			name: "happy case",
			srv: testServerExpectations{
				method:          http.MethodPost,
				url:             "/v1/users",
				body:            someUserToCreate,
				responseStatus:  http.StatusCreated,
				responsePayload: someCreatedUser,
			},
			expectedReturnValue: someCreatedUser,
			expectedError:       nil,
		},
		{
			name: "bad request",
			srv: testServerExpectations{
				method:          http.MethodPost,
				url:             "/v1/users",
				body:            someUserToCreate,
				responseStatus:  http.StatusBadRequest,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusBadRequest, Response: someErrorResponse},
		},
		{
			name: "internal error",
			srv: testServerExpectations{
				method:          http.MethodPost,
				url:             "/v1/users",
				body:            someUserToCreate,
				responseStatus:  http.StatusInternalServerError,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusInternalServerError, Response: someErrorResponse},
		},
		{
			name: "unexpected error",
			srv: testServerExpectations{
				method:          http.MethodPost,
				url:             "/v1/users",
				body:            someUserToCreate,
				responseStatus:  http.StatusBadGateway,
				responsePayload: nil,
			},
			expectedReturnValue: nil,
			expectedError:       errors.New("received unexpected status code 502"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := startTestServer(t, tc.srv)
			defer srv.Close()

			api := restuser.New(restuser.Config{URL: srv.URL})
			res, err := api.CreateUser(context.Background(), someUserToCreate)
			assert.Equal(t, tc.expectedReturnValue, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	t.Run("nil user", func(t *testing.T) {
		api := restuser.New(restuser.Config{"http://google.com"})
		_, err := api.CreateUser(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestAPI_UpdateUser(t *testing.T) {
	someUserToUpdate := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	someUpdatedUser := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	someErrorResponse := &restuser.ErrorResponse{Message: "everything is wrong"}

	for _, tc := range []struct {
		name                string
		srv                 testServerExpectations
		expectedReturnValue *restuser.User
		expectedError       error
	}{
		{
			name: "happy case",
			srv: testServerExpectations{
				method:          http.MethodPut,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				body:            someUserToUpdate,
				responseStatus:  http.StatusOK,
				responsePayload: someUpdatedUser,
			},
			expectedReturnValue: someUpdatedUser,
			expectedError:       nil,
		},
		{
			name: "bad request",
			srv: testServerExpectations{
				method:          http.MethodPut,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				body:            someUserToUpdate,
				responseStatus:  http.StatusBadRequest,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusBadRequest, Response: someErrorResponse},
		},
		{
			name: "conflict",
			srv: testServerExpectations{
				method:          http.MethodPut,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				body:            someUserToUpdate,
				responseStatus:  http.StatusConflict,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusConflict, Response: someErrorResponse},
		},
		{
			name: "not found",
			srv: testServerExpectations{
				method:          http.MethodPut,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				body:            someUserToUpdate,
				responseStatus:  http.StatusNotFound,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusNotFound, Response: someErrorResponse},
		},
		{
			name: "internal error",
			srv: testServerExpectations{
				method:          http.MethodPut,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				body:            someUserToUpdate,
				responseStatus:  http.StatusInternalServerError,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusInternalServerError, Response: someErrorResponse},
		},
		{
			name: "unexpected error",
			srv: testServerExpectations{
				method:          http.MethodPut,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				body:            someUserToUpdate,
				responseStatus:  http.StatusBadGateway,
				responsePayload: nil,
			},
			expectedReturnValue: nil,
			expectedError:       errors.New("received unexpected status code 502"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := startTestServer(t, tc.srv)
			defer srv.Close()

			api := restuser.New(restuser.Config{URL: srv.URL})
			res, err := api.UpdateUser(context.Background(), someUserToUpdate)
			assert.Equal(t, tc.expectedReturnValue, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}

	t.Run("nil user", func(t *testing.T) {
		api := restuser.New(restuser.Config{"http://google.com"})
		_, err := api.UpdateUser(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestAPI_GetUser(t *testing.T) {
	someUser := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	someErrorResponse := &restuser.ErrorResponse{Message: "everything is wrong"}

	for _, tc := range []struct {
		name                string
		srv                 testServerExpectations
		expectedReturnValue *restuser.User
		expectedError       error
	}{
		{
			name: "happy case",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusOK,
				responsePayload: someUser,
			},
			expectedReturnValue: someUser,
			expectedError:       nil,
		},
		{
			name: "not found",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusNotFound,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusNotFound, Response: someErrorResponse},
		},
		{
			name: "internal error",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusInternalServerError,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusInternalServerError, Response: someErrorResponse},
		},
		{
			name: "unexpected error",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusBadGateway,
				responsePayload: nil,
			},
			expectedReturnValue: nil,
			expectedError:       errors.New("received unexpected status code 502"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := startTestServer(t, tc.srv)
			defer srv.Close()

			api := restuser.New(restuser.Config{URL: srv.URL})
			res, err := api.GetUser(context.Background(), someUser.ID)
			assert.Equal(t, tc.expectedReturnValue, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAPI_DeleteUser(t *testing.T) {
	const someUserID = "c3e11b46-109c-11eb-adc1-0242ac120002"
	someErrorResponse := &restuser.ErrorResponse{Message: "everything is wrong"}

	for _, tc := range []struct {
		name          string
		srv           testServerExpectations
		expectedError error
	}{
		{
			name: "happy case",
			srv: testServerExpectations{
				method:         http.MethodDelete,
				url:            "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus: http.StatusNoContent,
			},
			expectedError: nil,
		},
		{
			name: "not found",
			srv: testServerExpectations{
				method:          http.MethodDelete,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusNotFound,
				responsePayload: someErrorResponse,
			},
			expectedError: restuser.Error{StatusCode: http.StatusNotFound, Response: someErrorResponse},
		},
		{
			name: "internal error",
			srv: testServerExpectations{
				method:          http.MethodDelete,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusInternalServerError,
				responsePayload: someErrorResponse,
			},
			expectedError: restuser.Error{StatusCode: http.StatusInternalServerError, Response: someErrorResponse},
		},
		{
			name: "unexpected error",
			srv: testServerExpectations{
				method:          http.MethodDelete,
				url:             "/v1/users/c3e11b46-109c-11eb-adc1-0242ac120002",
				responseStatus:  http.StatusBadGateway,
				responsePayload: nil,
			},
			expectedError: errors.New("received unexpected status code 502"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := startTestServer(t, tc.srv)
			defer srv.Close()

			api := restuser.New(restuser.Config{URL: srv.URL})
			err := api.DeleteUser(context.Background(), someUserID)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestAPI_ListUsers(t *testing.T) {
	frenchUser := restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pierre",
		Email:     "pierre@faceit.com",
		Country:   "fr",
	}
	spanishUser := restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120003",
		CreatedAt: "2007-01-02T16:04:05Z",
		UpdatedAt: "2007-01-03T16:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "es",
	}

	someErrorResponse := &restuser.ErrorResponse{Message: "everything is wrong"}

	for _, tc := range []struct {
		name                string
		srv                 testServerExpectations
		params              restuser.ListUsersParams
		expectedReturnValue []restuser.User
		expectedError       error
	}{
		{
			name: "happy case all",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users",
				responseStatus:  http.StatusOK,
				responsePayload: []restuser.User{frenchUser, spanishUser},
			},
			expectedReturnValue: []restuser.User{frenchUser, spanishUser},
			expectedError:       nil,
		},
		{
			name: "happy case filtered",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users?country=es",
				responseStatus:  http.StatusOK,
				responsePayload: []restuser.User{spanishUser},
			},
			params:              restuser.ListUsersParams{Country: "es"},
			expectedReturnValue: []restuser.User{spanishUser},
			expectedError:       nil,
		},
		{
			name: "bad request",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users",
				responseStatus:  http.StatusBadRequest,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusBadRequest, Response: someErrorResponse},
		},
		{
			name: "internal error",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users",
				responseStatus:  http.StatusInternalServerError,
				responsePayload: someErrorResponse,
			},
			expectedReturnValue: nil,
			expectedError:       restuser.Error{StatusCode: http.StatusInternalServerError, Response: someErrorResponse},
		},
		{
			name: "unexpected error",
			srv: testServerExpectations{
				method:          http.MethodGet,
				url:             "/v1/users",
				responseStatus:  http.StatusBadGateway,
				responsePayload: nil,
			},
			expectedReturnValue: nil,
			expectedError:       errors.New("received unexpected status code 502"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := startTestServer(t, tc.srv)
			defer srv.Close()

			api := restuser.New(restuser.Config{URL: srv.URL})
			res, err := api.ListUsers(context.Background(), tc.params)
			assert.Equal(t, tc.expectedReturnValue, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestWithBasePath(t *testing.T) {
	someUser := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	for _, tc := range []struct {
		name string
		do   func(*restuser.API)
	}{
		{
			name: "CreateUser",
			do: func(api *restuser.API) {
				_, _ = api.CreateUser(context.Background(), someUser)
			},
		},
		{
			name: "UpdateUser",
			do: func(api *restuser.API) {
				_, _ = api.UpdateUser(context.Background(), someUser)
			},
		},
		{
			name: "DeleteUser",
			do: func(api *restuser.API) {
				_ = api.DeleteUser(context.Background(), someUser.ID)
			},
		},
		{
			name: "GetUser",
			do: func(api *restuser.API) {
				_, _ = api.GetUser(context.Background(), someUser.ID)
			},
		},
		{
			name: "ListUsers",
			do: func(api *restuser.API) {
				_, _ = api.ListUsers(context.Background(), restuser.ListUsersParams{})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const someDifferentBasePath = "/preproduction/v1"

			var called int64
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.True(t, strings.HasPrefix(req.RequestURI, someDifferentBasePath))
				atomic.AddInt64(&called, 1)
			}))
			defer srv.Close()

			api := restuser.New(restuser.Config{URL: srv.URL}, restuser.WithBasePath(someDifferentBasePath))
			tc.do(api)

			// assert that was called, thead safely
			assert.True(t, atomic.LoadInt64(&called) > 0)
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	someUser := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	for _, tc := range []struct {
		name string
		do   func(*restuser.API)
	}{
		{
			name: "CreateUser",
			do: func(api *restuser.API) {
				_, _ = api.CreateUser(context.Background(), someUser)
			},
		},
		{
			name: "UpdateUser",
			do: func(api *restuser.API) {
				_, _ = api.UpdateUser(context.Background(), someUser)
			},
		},
		{
			name: "DeleteUser",
			do: func(api *restuser.API) {
				_ = api.DeleteUser(context.Background(), someUser.ID)
			},
		},
		{
			name: "GetUser",
			do: func(api *restuser.API) {
				_, _ = api.GetUser(context.Background(), someUser.ID)
			},
		},
		{
			name: "ListUsers",
			do: func(api *restuser.API) {
				_, _ = api.ListUsers(context.Background(), restuser.ListUsersParams{})
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const authorizationHeaderName = "Authorization"

			var called int64
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				assert.Contains(t, req.Header, authorizationHeaderName)
				atomic.AddInt64(&called, 1)
			}))
			defer srv.Close()

			client := &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
				req.Header.Add(authorizationHeaderName, "Bearer foo")
				return http.DefaultTransport.RoundTrip(req)
			})}

			api := restuser.New(restuser.Config{URL: srv.URL}, restuser.WithHTTPClient(client))
			tc.do(api)

			// assert that was called, thead safely
			assert.True(t, atomic.LoadInt64(&called) > 0)
		})
	}
}

func TestAPI_CallsWithContext(t *testing.T) {
	someUser := &restuser.User{
		ID:        "c3e11b46-109c-11eb-adc1-0242ac120002",
		CreatedAt: "2006-01-02T15:04:05Z",
		UpdatedAt: "2006-01-03T15:04:05Z",
		Name:      "Pepe",
		Email:     "pepe@faceit.com",
		Country:   "fr",
	}

	for _, tc := range []struct {
		name string
		do   func(context.Context, *restuser.API) error
	}{
		{
			name: "CreateUser",
			do: func(ctx context.Context, api *restuser.API) error {
				_, err := api.CreateUser(ctx, someUser)
				return err
			},
		},
		{
			name: "UpdateUser",
			do: func(ctx context.Context, api *restuser.API) error {
				_, err := api.UpdateUser(ctx, someUser)
				return err
			},
		},
		{
			name: "DeleteUser",
			do: func(ctx context.Context, api *restuser.API) error {
				return api.DeleteUser(ctx, someUser.ID)
			},
		},
		{
			name: "GetUser",
			do: func(ctx context.Context, api *restuser.API) error {
				_, err := api.GetUser(ctx, someUser.ID)
				return err
			},
		},
		{
			name: "ListUsers",
			do: func(ctx context.Context, api *restuser.API) error {
				_, err := api.ListUsers(ctx, restuser.ListUsersParams{})
				return err
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const timeout = 10 * time.Millisecond
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				time.Sleep(2 * timeout) // too long for our context
			}))
			defer srv.Close()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			api := restuser.New(restuser.Config{URL: srv.URL})
			err := tc.do(ctx, api)

			assert.Error(t, err)
			assert.True(t, errors.Is(err, context.DeadlineExceeded))
		})
	}
}

type testServerExpectations struct {
	method          string
	url             string
	body            interface{}
	responseStatus  int
	responsePayload interface{}
}

func startTestServer(t *testing.T, expected testServerExpectations) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, expected.method, req.Method)
		assert.Equal(t, expected.url, req.RequestURI)
		if expected.body != nil {
			// build same thing but empty, try to JSON decode it
			gotJSONBody := sameTypeButEmptyElement(expected.body)
			err := json.NewDecoder(req.Body).Decode(gotJSONBody)
			require.NoError(t, err)
			assert.Equal(t, expected.body, gotJSONBody)
		}

		rw.WriteHeader(expected.responseStatus)
		if expected.responsePayload != nil {
			require.NoError(t, json.NewEncoder(rw).Encode(expected.responsePayload))
		}
	}))
}

func sameTypeButEmptyElement(val interface{}) interface{} {
	typ := reflect.TypeOf(val)
	isPtr := typ.Kind() == reflect.Ptr
	if isPtr {
		typ = typ.Elem()
	}
	if isPtr {
		return reflect.New(typ).Interface()
	}
	return reflect.Zero(typ).Interface()
}

type roundTripperFunc func(r *http.Request) (*http.Response, error)

func (r roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return r(request)
}
