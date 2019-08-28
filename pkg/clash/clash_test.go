package clash

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClan(t *testing.T) {

	t.Run("unauthorized", func(t *testing.T) {
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

		client := NewClient("TAG", "TOKEN", testServer.URL)

		_, err := client.GetClan()
		log.Printf("%v", err)

	})

}
