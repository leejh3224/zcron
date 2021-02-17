package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
}

func TestGetNextNSchedules_BasicUsage(t *testing.T) {
	u, _ := url.Parse("/next-schedules")
	q := u.Query()
	q.Set("expression", "0/5 23,0-11 * * ? *")
	q.Set("timezoneOffset", "+09:00")
	q.Set("limit", "5")
	u.RawQuery = q.Encode()

	s := newServer()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", u.String(), nil)
	s.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	var res getNextNSchedulesResponse
	json.Unmarshal(body, &res)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "0/5 23,0-11 * * ? *", res.Expression)
	assert.Equal(t, len(res.NextSchedules), res.Limit)
	assert.Equal(t, "+09:00", res.TimezoneOffset)
}
