/*
This file will contain the data models. Until refactored only Headers will be here to be included everywhere else.
*/

package profitbricks

import "net/http"

type Header struct {
	Headers *http.Header
}

// GetHeader to be interfaceable
func (h *Header) GetHeader() *http.Header {
	return h.Headers
}

// SetHeader to be interfaceable
func (h *Header) SetHeader(header *http.Header) {
	h.Headers = header
}

// Get returns the actual value for given header key
func (h *Header) Get(key string) string {
	return h.Headers.Get(key)
}
