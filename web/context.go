// Copyright 2014 li. All rights reserved.

package web

import (
	//"code.google.com/p/go.net/websocket"
	"bytes"
	"fmt"
	"github.com/roverli/light/log"
	"github.com/roverli/light/mux"
	"github.com/roverli/light/session"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

// A unified view of the request params.
// - Rest path params.
// - URL query string.
// - Form values.
// - File uploads.
// - NOTE: param maps may be nil if there were none.
type Params struct {
	url.Values                                    // A unified view of all the individual param maps below.
	Route      url.Values                         // Parameters extracted from the route,  e.g. /profile/{id}
	Query      url.Values                         // Parameters from the query string, e.g. /list?page=2
	Form       url.Values                         // Parameters from the request body.
	Files      map[string][]*multipart.FileHeader // Files uploaded in a multipart form
	TmpFiles   []*os.File                         // Temp files used during the request.
}

func (p *Params) merge() url.Values {
	num := len(p.Route) + len(p.Query) + len(p.Form)

	if num == 0 {
		return make(url.Values, 0)
	}

	switch num {
	case len(p.Query):
		return p.Query
	case len(p.Route):
		return p.Route
	case len(p.Form):
		return p.Form
	}

	values := make(url.Values, num)
	for k, v := range p.Route {
		values[k] = append(values[k], v...)
	}
	for k, v := range p.Form {
		values[k] = append(values[k], v...)
	}
	for k, v := range p.Query {
		values[k] = append(values[k], v...)
	}

	return values
}

// Wrapper for http request
type Request struct {
	*http.Request
	ContentType     string // eg. "text/html"
	Format          string // eg. "html", "xml", "json", or "txt"
	Locale          string
	AcceptLanguages AcceptLanguages
	//Websocket       *websocket.Conn
}

// Wrapper for http response writer.
type Response struct {
	Status      int
	ContentType string
	Resp        http.ResponseWriter
}

// Wrapper all.
type Context struct {
	Req         *http.Request       // Origin http request
	Resp        http.ResponseWriter // Origin http response writer
	Params      *Params             // All request params
	RouteResult *mux.Result         // The route result
	Session     session.Session     // Http Session
	//	Status      Status              // Handle status
}

// Copy from reveal
// AcceptLanguage is a single language from the Accept-Language HTTP header.
type AcceptLanguage struct {
	Language string
	Quality  float32
}

// AcceptLanguages is collection of sortable AcceptLanguage instances.
type AcceptLanguages []AcceptLanguage

func (al AcceptLanguages) Len() int           { return len(al) }
func (al AcceptLanguages) Swap(i, j int)      { al[i], al[j] = al[j], al[i] }
func (al AcceptLanguages) Less(i, j int) bool { return al[i].Quality > al[j].Quality }
func (al AcceptLanguages) String() string {
	output := bytes.NewBufferString("")
	for i, language := range al {
		output.WriteString(fmt.Sprintf("%s (%1.1f)", language.Language, language.Quality))
		if i != len(al)-1 {
			output.WriteString(", ")
		}
	}
	return output.String()
}

// ResolveAcceptLanguage returns a sorted list of Accept-Language
// header values.
//
// The results are sorted using the quality defined in the header for each
// language range with the most qualified language range as the first
// element in the slice.
//
// See the HTTP header fields specification
// (http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html#sec14.4) for more details.
func ResolveAcceptLanguage(req *http.Request) AcceptLanguages {
	header := req.Header.Get("Accept-Language")
	if header == "" {
		return nil
	}

	acceptLanguageHeaderValues := strings.Split(header, ",")
	acceptLanguages := make(AcceptLanguages, len(acceptLanguageHeaderValues))

	for i, languageRange := range acceptLanguageHeaderValues {
		if qualifiedRange := strings.Split(languageRange, ";q="); len(qualifiedRange) == 2 {
			quality, error := strconv.ParseFloat(qualifiedRange[1], 32)
			if error != nil {
				log.Warn("Detected malformed Accept-Language header quality in '%s', assuming quality is 1", languageRange)
				acceptLanguages[i] = AcceptLanguage{qualifiedRange[0], 1}
			} else {
				acceptLanguages[i] = AcceptLanguage{qualifiedRange[0], float32(quality)}
			}
		} else {
			acceptLanguages[i] = AcceptLanguage{languageRange, 1}
		}
	}

	sort.Sort(acceptLanguages)
	return acceptLanguages
}
