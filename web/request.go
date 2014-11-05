// Copyright 2014 li. All rights reserved.

package web

import (
	"net/http"
	"strings"
)

// Get the http request content type.
// If none is specified, returns "text/html" by default.
func ResolveContentType(req *http.Request) string {
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		return "text/html"
	}
	// From "multipart/form-data; boundary=--" to "multipart/form-data"
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

// ResolveFormat maps the request's Accept MIME type to "html", "xml", "json", or "txt".
// Returning a default of "html" when Accept header cannot be mapped to a value above.
func ResolveFormat(req *http.Request) string {
	accept := req.Header.Get("accept")

	switch {
	case accept == "",
		strings.HasPrefix(accept, "*/*"),
		strings.Contains(accept, "application/xhtml"),
		strings.Contains(accept, "text/html"):
		return "html"

	case strings.Contains(accept, "application/xml"),
		strings.Contains(accept, "text/xml"):
		return "xml"

	case strings.Contains(accept, "text/plain"):
		return "txt"

	case strings.Contains(accept, "application/json"),
		strings.Contains(accept, "text/javascript"):
		return "json"
	}

	return "html"
}

// Write the header (for now, just the status code).
// The status may be set directly by the application (c.Response.Status = 501).
// if it isn't, then fall back to the provided status code.
// func (resp *Response) WriteHeader(defaultStatusCode int, defaultContentType string) {
// 	if resp.Status == 0 {
// 		resp.Status = defaultStatusCode
// 	}
// 	if resp.ContentType == "" {
// 		resp.ContentType = defaultContentType
// 	}
// 	resp.Out.Header().Set("Content-Type", resp.ContentType)
// 	resp.Out.WriteHeader(resp.Status)
// }
