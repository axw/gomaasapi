// Copyright 2013 Canonical Ltd.  This software is licensed under the
// GNU Affero General Public License version 3 (see the file LICENSE).

package gomaasapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type singleServingServer struct {
	*httptest.Server
	requestContent *string
	requestHeader  *http.Header
}

// newSingleServingServer create a single-serving test http server which will
// return only one response as defined by the passed arguments.
func newSingleServingServer(uri string, response string, code int) *singleServingServer {
	var requestContent string
	var requestHeader http.Header
	var requested bool
	handler := func(writer http.ResponseWriter, request *http.Request) {
		if requested {
			http.Error(writer, "Already requested", http.StatusServiceUnavailable)
		}
		res, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		requestContent = string(res)
		requestHeader = request.Header
		if request.URL.String() != uri {
			http.Error(writer, "404 page not found", http.StatusNotFound)
		} else {
			writer.WriteHeader(code)
			fmt.Fprint(writer, response)
		}
		requested = true
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	return &singleServingServer{server, &requestContent, &requestHeader}
}
