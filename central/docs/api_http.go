package docs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Swagger returns an HTTP handler that exposes the swagger.json doc directly.
// It's not a gRPC method because some clients will want to consume this URL directly,
// rather than interpreting a JSON string from inside a response.
func Swagger() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		b, err := swaggerForRequest(req)
		if err != nil {
			w.WriteHeader(500)
			msg := err.Error()
			_, _ = w.Write([]byte(msg))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
}

func swaggerForRequest(req *http.Request) ([]byte, error) {
	b, err := ioutil.ReadFile("/docs/api/v1/swagger.json")
	if err != nil {
		return nil, fmt.Errorf("could not load swagger file: %v", err)
	}

	var swaggerSpec map[string]json.RawMessage
	if err := json.Unmarshal(b, &swaggerSpec); err != nil {
		return nil, fmt.Errorf("could not parse swagger spec: %v", err)
	}

	swaggerSpecOut := make(map[string]interface{}, len(swaggerSpec)+2)
	for k, v := range swaggerSpec {
		swaggerSpecOut[k] = v
	}

	scheme, host := extractSchemeAndHost(req)
	swaggerSpecOut["host"] = host
	swaggerSpecOut["schemes"] = []string{scheme}

	out, err := json.MarshalIndent(swaggerSpecOut, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not marshal swagger spec: %v", err)
	}
	return out, nil
}

func extractSchemeAndHost(req *http.Request) (string, string) {
	forwardedProto := req.Header.Get("X-Forwarded-Proto")
	forwardedHost := req.Header.Get("X-Forwarded-Host")
	if forwardedHost != "" && forwardedProto != "" {
		return strings.ToLower(forwardedProto), forwardedHost
	}

	scheme := req.URL.Scheme
	if scheme == "" {
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	return scheme, host
}
