package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/nozomi-nishinohara/jwt_validation/handler"
	"github.com/nozomi-nishinohara/jwt_validation/token"
	"github.com/stretchr/testify/assert"
)

func getIDToken() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("client_id=%s", os.Getenv("CLIENT_ID")))
	builder.WriteString("&")
	builder.WriteString(fmt.Sprintf("client_secret=%s", os.Getenv("CLIENT_SECRET")))
	builder.WriteString("&")
	builder.WriteString(fmt.Sprintf("username=%s", os.Getenv("USERNAME")))
	builder.WriteString("&")
	builder.WriteString(fmt.Sprintf("password=%s", os.Getenv("PASSWORD")))
	builder.WriteString("&")
	builder.WriteString(fmt.Sprintf("scope=%s", os.Getenv("SCOPE")))
	builder.WriteString("&")
	builder.WriteString(fmt.Sprintf("grant_type=%s", os.Getenv("GRANT_TYPE")))
	realm := os.Getenv("REALM")
	if realm != "" {
		builder.WriteString("&")
		builder.WriteString(fmt.Sprintf("realm=%s", realm))
	}
	req, err := http.Post(os.Getenv("URL"), "application/x-www-form-urlencoded", strings.NewReader(builder.String()))
	if req != nil {
		defer req.Body.Close()
	}
	if err != nil {
		log.Fatalln(err)
		return ""
	}
	m := map[string]interface{}{}
	buf, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(buf, &m)
	if val, ok := m["id_token"]; ok {
		return val.(string)
	}
	return ""
}

func getHandler(w http.ResponseWriter, req *http.Request) {
	value := token.FromContext(req.Context())
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	buf, _ := json.Marshal(value)
	w.Write(buf)
}

func TestJWTHandler(t *testing.T) {
	token := getIDToken()
	router := httprouter.New()
	h := handler.New(nil)
	router.HandlerFunc(http.MethodGet, "/user", h.JWTValidationMiddleware(getHandler))
	req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	m := map[string]interface{}{}
	json.NewDecoder(rec.Body).Decode(&m)
	assert.Equal(t, "auth0|5fe0067678238b0071965f44", m["sub"].(string))
}

func TestJWTRouter(t *testing.T) {
	token := getIDToken()
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/user", getHandler)
	h := handler.New(nil)
	handl := h.JWTValidation(router)
	req := httptest.NewRequest(http.MethodGet, "/user", bytes.NewBufferString(""))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	rec := httptest.NewRecorder()
	handl.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	m := map[string]interface{}{}
	json.NewDecoder(rec.Body).Decode(&m)
	assert.Equal(t, "auth0|5fe0067678238b0071965f44", m["sub"].(string))
}
