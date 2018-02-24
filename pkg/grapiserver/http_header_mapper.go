package grapiserver

import "net/http"

// HTTPHeaderMappingConfig contains functions for deciding and mapping http header keys.
type HTTPHeaderMappingConfig struct {
	DeciderFunc func(string) bool
	MapperFunc  func(string) string
}

type httpHeaderMapper struct {
	*HTTPHeaderMappingConfig
	mappingCache map[string]string
}

func newHTTPHeaderMapper(c *HTTPHeaderMappingConfig) *httpHeaderMapper {
	return &httpHeaderMapper{
		HTTPHeaderMappingConfig: c,
		mappingCache:            make(map[string]string),
	}
}

func (m *httpHeaderMapper) wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.apply(r)
		next.ServeHTTP(w, r)
	})
}

func (m *httpHeaderMapper) apply(r *http.Request) {
	newHeader := make(map[string]string, len(r.Header))

	for k := range r.Header {
		if newKey, ok := m.mappingCache[k]; ok {
			newHeader[newKey] = r.Header.Get(k)
		} else if m.DeciderFunc(k) {
			newKey := m.MapperFunc(k)
			newHeader[newKey] = r.Header.Get(k)
			m.mappingCache[k] = newKey
		}
	}

	for k, v := range newHeader {
		r.Header.Set(k, v)
	}
}
