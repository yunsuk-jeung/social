package main

import (
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {

	app := newTestApplication(t, config{})
	mux := app.mount()

	testToken, _ := app.authenticator.GenerateToken(nil)

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/201", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := exceteRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should allow authenticated requests", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/v1/users/201", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := exceteRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
	})

}
