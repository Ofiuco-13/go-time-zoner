package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)
func TestTimeHandler(t *testing.T) {
    req, err := http.NewRequest("GET", "/time", nil)
    if err != nil {
        t.Fatal(err)
    }

    q := req.URL.Query()
    q.Add("tz", "UTC")
    req.URL.RawQuery = q.Encode()

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(timeHandler)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expectedContentType := "application/json"
    if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
        t.Errorf("handler returned wrong content type: got %v want %v",
            contentType, expectedContentType)
    }

    // Test with valid timezone
    q = req.URL.Query()
    q.Add("tz", "America/New_York")
    req.URL.RawQuery = q.Encode()

    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    // Test with invalid timezone
    q = req.URL.Query()
    q.Set("tz", "invalid")
    req.URL.RawQuery = q.Encode()

    rr = httptest.NewRecorder()
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }

    expectedError := "unknown timezone"
    body, err := ioutil.ReadAll(rr.Body)
    if err != nil {
        t.Fatal(err)
    }

    var response map[string]interface{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Fatal(err)
    }

    if errorValue, ok := response["error"]; !ok || errorValue != expectedError {
        t.Errorf("handler returned unexpected error value: got %v want %v",
            errorValue, expectedError)
    }

    if timeValue, ok := response["time"]; ok {
        if timeValue != nil {
            _, ok := timeValue.(string)
            if !ok {
                t.Errorf("handler returned unexpected time value: got %T want string or nil",
                    timeValue)
            }
        }
    } else {
        t.Errorf("handler did not return expected time value")
    }
}

