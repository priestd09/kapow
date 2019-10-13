package data

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetRequestMethodReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	value, err := getRequestMethod(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "GET" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "GET", value)
	}
}

func TestGetRequestHostReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	value, err := getRequestHost(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "www.example.com" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "www.example.com", value)
	}
}

func TestGetRequestPathReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	value, err := getRequestPath(req)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "/this/is/a/test" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "/this/is/a/test", value)
	}
}

func TestSetResponseStatusSetsCorrectValue(t *testing.T) {
	res := httptest.NewRecorder()

	setResponseStatus(res, 500)
	if val := res.Result().StatusCode; val != 500 {
		t.Errorf("Unexpected value. Expected: %d, got: %d", 500, val)
	}
}

func TestGetRequestHeaderReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")

	value, err := getRequestHeader(req, "A-Header")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestHeaderReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")

	if _, err := getRequestHeader(req, "Other-Header"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestSetResponseHeaderSetsCorrectValue(t *testing.T) {
	res := httptest.NewRecorder()

	setResponseHeader(res, "A-Header", "With-Value")
	if val := res.Result().Header.Get("A-Header"); val != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", val)
	}
}

func TestGetRequestCookieReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	value, err := getRequestCookie(req, "A-Cookie")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestCookieReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	if _, err := getRequestCookie(req, "Other-Cookie"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestSetResponseCookieSetsCorrectValue(t *testing.T) {
	res := httptest.NewRecorder()

	setResponseCookie(res, "A-Cookie", "With-Value")
	cookies := res.Result().Cookies()
	val := ""
	for _, v := range cookies {
		if v.Name == "A-Cookie" {
			val = v.Value
		}
	}
	if val != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", val)
	}
}

func TestGetRequestParamReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	value, err := getRequestParam(req, "with")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "params" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "params", value)
	}
}

func TestGetRequestParamReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	if _, err := getRequestParam(req, "Other-Param"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestGetRequestFormReturnsCorrectValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.PostForm = url.Values{}
	req.PostForm.Set("A-Field", "With-Value")
	req.PostForm.Set("Another-Field", "With-AnotherValue")

	value, err := getRequestForm(req, "A-Field")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestFormReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.PostForm = url.Values{}
	req.PostForm.Set("A-Field", "With-Value")
	req.PostForm.Set("Another-Field", "With-AnotherValue")

	if _, err := getRequestForm(req, "Other-Field"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestGetRequestFileNameReturnsCorrectValue(t *testing.T) {
	t.Skip("****** WIP ******")
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.MultipartForm = &multipart.Form{}
	req.MultipartForm.File = make(map[string][]*multipart.FileHeader)
	req.MultipartForm.File["A-File"] = append(make([]*multipart.FileHeader, 1), &multipart.FileHeader{Filename: "A-file.txt"})

	value, err := getRequestFileName(req, "A-File")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestGetRequestFileNameReturnsErrorWhenNotExists(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.MultipartForm = &multipart.Form{}
	req.MultipartForm.File = make(map[string][]*multipart.FileHeader)
	req.MultipartForm.File["A-File"] = append(make([]*multipart.FileHeader, 1), &multipart.FileHeader{Filename: "A-file.txt"})

	if _, err := getRequestFileName(req, "Other-Field"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestCopyRequestFileReturnsOK(t *testing.T) {
	t.Skip("****** WIP ******")
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.MultipartForm = &multipart.Form{}
	req.MultipartForm.File = make(map[string][]*multipart.FileHeader)
	req.MultipartForm.File["A-File"] = append(make([]*multipart.FileHeader, 1), &multipart.FileHeader{Filename: "A-file.txt"})

	result := strings.Builder{}
	err := copyRequestFile(req, "A-File", &result)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value := result.String(); value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func TestCopyRequestFileReturnsErrorWhenNotExists(t *testing.T) {
	t.Skip("****** WIP ******")
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", nil)

	req.Header.Add("A-Header", "With-Value")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	req.MultipartForm = &multipart.Form{}
	req.MultipartForm.File = make(map[string][]*multipart.FileHeader)
	req.MultipartForm.File["A-File"] = append(make([]*multipart.FileHeader, 1), &multipart.FileHeader{Filename: "A-file.txt"})

	result := strings.Builder{}
	err := copyRequestFile(req, "A-File", &result)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value := result.String(); value != "With-Value" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "With-Value", value)
	}
}

func generateTargetRequestForMatch() *http.Request {
	var targetRequest *http.Request

	h := mux.NewRouter()
	h.HandleFunc("/a/{foo}", func(res http.ResponseWriter, req *http.Request) { targetRequest = req }).Methods("GET")
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/a/bar", nil))

	return targetRequest
}

func TestGetRequestMatchReturnsCorrectValue(t *testing.T) {
	req := generateTargetRequestForMatch()

	value, err := getRequestMatch(req, "foo")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value != "bar" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "bar", value)
	}
}

func TestGetRequestMatchReturnsErrorWhenNotExists(t *testing.T) {
	req := generateTargetRequestForMatch()

	if _, err := getRequestMatch(req, "bar"); err == nil {
		t.Errorf("Expected error but no error returned")
	}
}

func TestCopyFromRequestBodyReturnsOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "http://www.example.com/this/is/a/test?with=params", strings.NewReader("This is a body content for testing purposes"))

	req.Header.Add("A-Header", "With-Value")
	req.AddCookie(&http.Cookie{Name: "A-Cookie", Value: "With-Value"})

	result := strings.Builder{}
	if err := copyFromRequestBody(req, &result); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if value := result.String(); value != "This is a body content for testing purposes" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "This is a body content for testing purposes", value)
	}
}

func TestCopyToResponseBodyReturnsOK(t *testing.T) {
	res := httptest.NewRecorder()

	if err := copyToResponseBody(res, strings.NewReader("This is a body content for testing purposes")); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Result().Body)
	if err != nil {
		t.Errorf("Unexpected error while reading result body: %+v", err)
	}

	if value := string(bodyBytes); value != "This is a body content for testing purposes" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "This is a body content for testing purposes", value)
	}
}

func TestCopyToResponseStreamReturnsOK(t *testing.T) {
	res := httptest.NewRecorder()

	if err := copyToResponseStream(res, strings.NewReader("This is a body content for testing purposes")); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	bodyBytes, err := ioutil.ReadAll(res.Result().Body)
	if err != nil {
		t.Errorf("Unexpected error while reading result body: %+v", err)
	}

	if value := string(bodyBytes); value != "This is a body content for testing purposes" {
		t.Errorf("Unexpected value. Expected: %s, got: %s", "This is a body content for testing purposes", value)
	}
}

func TestRouterIsWellConfigured(t *testing.T) {
	testCases := []struct {
		pattern, method string
		handler         func(http.ResponseWriter, *http.Request)
		mustMatch       bool
		vars            []struct{ k, v string }
	}{
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/method", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "method"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/host", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "host"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/path", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "path"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/matches/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "matches/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/params/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "params/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "headers/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/cookies/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "cookies/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/form/name", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "form/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/files/name/filename", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "files/name/filename"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/files/name/content", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "files/name/content"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/body", http.MethodGet, readRequestResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "body"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodPost, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodPut, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/request/headers/name", http.MethodDelete, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodGet, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodPost, nil, false, []struct{ k, v string }{}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/status", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "status"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "headers/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/cookies/name", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "cookies/name"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/body", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "body"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/stream", http.MethodPut, writeResponseResources, true, []struct{ k, v string }{{"handler_id", "HANDLER_YYYYYYYYYYYYYYYY"}, {"resource_path", "stream"}}},
		{"/handlers/HANDLER_YYYYYYYYYYYYYYYY/response/headers/name", http.MethodDelete, nil, false, []struct{ k, v string }{}},
	}

	r := configRouter()

	for _, tc := range testCases {
		rm := mux.RouteMatch{}
		rq, _ := http.NewRequest(tc.method, tc.pattern, nil)
		if matched := r.Match(rq, &rm); tc.mustMatch != matched {
			t.Errorf("Route mismatch: Expected: %+v\n\t\t\t\t\t\t got: %+v", tc, rm)
		} else {
			if tc.mustMatch {
				// Check for Handler match.
				realHandler := reflect.ValueOf(rm.Handler).Pointer()
				expectedHandler := reflect.ValueOf(tc.handler).Pointer()
				if realHandler != expectedHandler {
					t.Errorf("Handler mismatch. Expected: %X, got: %X", expectedHandler, realHandler)
				}

				// Check for variables
				for _, v := range tc.vars {
					if value, exists := rm.Vars[v.k]; !exists {
						t.Errorf("Variable not present: %s", v.k)
					} else if v.v != value {
						t.Errorf("Variable value mismatch. Expected: %s, got: %s", v.v, value)
					}
				}
			}
		}
	}
}