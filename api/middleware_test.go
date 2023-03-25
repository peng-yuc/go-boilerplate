package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/walkccc/go-boilerplate/token"
)

func TestAutuMidderware(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "NoAuthorization",
			setupAuth: func(request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "UnsupportedAuthorizationType",
			setupAuth: func(request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "unsupported_type", "username", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, "", "username", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "username", -time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "OK",
			setupAuth: func(request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "username", time.Minute)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			server := newTestServer(t, nil)

			authUrl := "/auth"
			server.router.GET(
				authUrl,
				authMiddleware(server.tokenMaker),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, authUrl, nil)
			require.NoError(t, err)

			testCase.setupAuth(request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(recorder)
		})
	}
}

func addAuthorization(
	t *testing.T,
	request *http.Request,
	tokenMaker token.Maker,
	authorizationType string,
	username string,
	duration time.Duration,
) {
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}
