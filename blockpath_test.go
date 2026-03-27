package traefik_blockpath

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc    string
		regexps []string
		expErr  bool
	}{
		{
			desc:    "should return no error",
			regexps: []string{`^/foo/(.*)`},
			expErr:  false,
		},
		{
			desc:    "should return an error",
			regexps: []string{"*"},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := &Config{
				Regex:      test.regexps,
				StatusCode: http.StatusForbidden,
			}

			_, err := New(context.Background(), nil, cfg, "name")
			if test.expErr && err == nil {
				t.Errorf("expected error on bad regexp format")
			}
		})
	}
}

func TestServeHTTP(t *testing.T) {
	tests := []struct {
		desc          string
		methods       []string
		regexps       []string
		reqMethod     string
		reqPath       string
		statusCode    int
		expNextCall   bool
		expStatusCode int
	}{
		{
			desc:          "should return forbidden status",
			regexps:       []string{"/test"},
			reqMethod:     http.MethodGet,
			reqPath:       "/test",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{"/test", "/toto"},
			reqMethod:     http.MethodGet,
			reqPath:       "/toto",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexps:       []string{"/test", "/toto"},
			reqMethod:     http.MethodGet,
			reqPath:       "/plop",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status",
			reqMethod:     http.MethodGet,
			reqPath:       "/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return forbidden status",
			regexps:       []string{`^/bar(.*)`},
			reqMethod:     http.MethodPost,
			reqPath:       "/bar/foo",
			expNextCall:   false,
			expStatusCode: http.StatusForbidden,
		},
		{
			desc:          "should return ok status",
			regexps:       []string{`^/bar(.*)`},
			reqMethod:     http.MethodGet,
			reqPath:       "/foo/bar",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return defined statuscode not found",
			regexps:       []string{`^/bar(.*)`},
			reqMethod:     http.MethodGet,
			reqPath:       "/bar/foo",
			statusCode:    http.StatusNotFound,
			expNextCall:   false,
			expStatusCode: http.StatusNotFound,
		},
		{
			desc:          "should return ok status with not matching and statuscode set",
			regexps:       []string{"/test", "/toto"},
			reqMethod:     http.MethodPost,
			reqPath:       "/plop",
			statusCode:    http.StatusNotFound,
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status with no regex and statuscode set",
			reqMethod:     http.MethodGet,
			reqPath:       "/test",
			statusCode:    http.StatusForbidden,
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return not found status for defined methods",
			methods:       []string{http.MethodGet, http.MethodPost},
			regexps:       []string{"/test"},
			reqMethod:     http.MethodGet,
			reqPath:       "/test",
			statusCode:    http.StatusNotFound,
			expNextCall:   false,
			expStatusCode: http.StatusNotFound,
		},
		{
			desc:          "should return ok status with not matching method",
			methods:       []string{http.MethodGet},
			regexps:       []string{"/test", "/toto"},
			reqMethod:     http.MethodPost,
			reqPath:       "/toto",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
		{
			desc:          "should return ok status with not matching path",
			methods:       []string{http.MethodPost},
			regexps:       []string{"/toto"},
			reqMethod:     http.MethodPost,
			reqPath:       "/test",
			expNextCall:   true,
			expStatusCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			cfg := CreateConfig()

			cfg.Regex = test.regexps
			if test.statusCode != 0 {
				cfg.StatusCode = test.statusCode
			}

			if test.methods != nil {
				cfg.Methods = test.methods
			}

			nextCall := false
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				nextCall = true
			})

			handler, err := New(context.Background(), next, cfg, "blockpath")
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()

			url := "http://localhost" + test.reqPath
			req := httptest.NewRequest(test.reqMethod, url, http.NoBody)

			handler.ServeHTTP(recorder, req)

			if nextCall != test.expNextCall {
				t.Errorf("next handler should not be called")
			}

			if recorder.Result().StatusCode != test.expStatusCode {
				t.Errorf("got status code %d, want %d", recorder.Code, test.expStatusCode)
			}
		})
	}
}
