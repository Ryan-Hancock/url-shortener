package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func GetTestHandler() *Handler {
	db := newDB(":memory:")

	h := Handler{}
	h.initialise(db)

	return &h
}

func newRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}

	return req
}

func serveRequest(h http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h)

	handler.ServeHTTP(rr, req)

	return rr
}

func TestHandler_postURL(t *testing.T) {
	h := GetTestHandler()

	var jsonStr = []byte(`{"url":"http://google.com"}`)
	req := newRequest(t, "POST", "/url", bytes.NewBuffer(jsonStr))
	rr := serveRequest(h.postURL, req)

	t.Run("should return a status created", func(t *testing.T) {
		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}
	})

	t.Run("should return a message", func(t *testing.T) {
		expected := `{"message":`
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("postURL() returned an error: %s", rr.Body.String())
		}
	})
}

func TestHandler_getURL(t *testing.T) {
	h := GetTestHandler()

	t.Run("should return a not found with wrong parameter", func(t *testing.T) {
		req := newRequest(t, "GET", "/somekey", nil)
		rr := serveRequest(h.getURL, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusNotFound)
		}
	})

	t.Run("and we make a get request", func(t *testing.T) {
		key := "randomkey"
		redirect := "http://google.com"
		err := h.storage.Set(key, redirect)
		if err != nil {
			t.Errorf("h.storage.Set() caused an error: %v", err)
		}

		req := newRequest(t, "GET", fmt.Sprintf("/%s", key), nil)
		req = mux.SetURLVars(req, map[string]string{
			"key": key,
		})
		rr := serveRequest(h.getURL, req)

		t.Run("should return a redirect status", func(t *testing.T) {
			if status := rr.Code; status != http.StatusMovedPermanently {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusMovedPermanently)
			}
		})

		t.Run("and should be redirected correctly", func(t *testing.T) {
			loc, err := rr.Result().Location()
			if err != nil {
				t.Errorf("request location had an error: %v", err)
			}

			if str := loc.String(); str != redirect {
				t.Errorf("not corected redirect: got %s wanted %s", str, redirect)
			}
		})
	})
}
