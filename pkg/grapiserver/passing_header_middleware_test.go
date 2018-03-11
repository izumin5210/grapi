package grapiserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_passingHeaderMiddleware(t *testing.T) {
	type Case struct {
		test    string
		decider PassedHeaderDeciderFunc
		in      http.Header
		out     http.Header
	}

	cases := []Case{
		{
			test:    "passing 1 header",
			decider: func(k string) bool { return strings.HasPrefix(k, "X-Debug-") },
			in: http.Header{
				"X-Debug-User-Id": []string{"100"},
				"X-User-Id":       []string{"100"},
			},
			out: http.Header{
				"X-Debug-User-Id":               []string{"100"},
				"Grpc-Metadata-X-Debug-User-Id": []string{"100"},
				"X-User-Id":                     []string{"100"},
			},
		},
	}

	getDefaultHeader := func() http.Header {
		return http.Header{
			"Accept-Encoding": []string{"gzip"},
			"User-Agent":      []string{"Go-http-client/1.1"},
		}
	}

	for _, c := range cases {
		t.Run(c.test, func(t *testing.T) {
			var wantHeader, gotHeader http.Header
			wantHeader = getDefaultHeader()
			for k, v := range c.out {
				wantHeader.Set(k, v[0])
			}

			wrap := createPassingHeaderMiddleware(c.decider)
			s := httptest.NewServer(wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotHeader = r.Header
				w.WriteHeader(200)
			})))
			defer s.Close()

			req, _ := http.NewRequest("GET", s.URL, nil)
			req.Header = c.in

			(&http.Client{}).Do(req)
			if diff := cmp.Diff(gotHeader, wantHeader); diff != "" {
				t.Errorf("Received header differs: (-got +want)\n%s", diff)
			}

			(&http.Client{}).Do(req)
			if diff := cmp.Diff(gotHeader, wantHeader); diff != "" {
				t.Errorf("Received header differs: (-got +want)\n%s", diff)
			}
		})
	}
}
