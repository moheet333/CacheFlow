package server

import (
	root "CacheFlow/cmd"
	"io"
	"net/http"
)

func UniversalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	case http.MethodPatch:
		patchHandler(w, r)
	case http.MethodPut:
		putHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	default:
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a get method"))
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	currentUrl := root.Newflag.Origin + r.URL.Path
	forwardHandler(w, r, currentUrl)
}

func putHandler(w http.ResponseWriter, r *http.Request) {
	currentUrl := root.Newflag.Origin + r.URL.Path
	forwardHandler(w, r, currentUrl)
}

func patchHandler(w http.ResponseWriter, r *http.Request) {
	currentUrl := root.Newflag.Origin + r.URL.Path
	forwardHandler(w,r,currentUrl)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	currentUrl := root.Newflag.Origin + r.URL.Path
	forwardHandler(w,r,currentUrl)
}

func forwardHandler(w http.ResponseWriter, r *http.Request, currentUrl string) {
	req, err := http.NewRequest(r.Method, currentUrl, r.Body)
	if err != nil {
		// internal server error
		InternalServerError(w)
	}

	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// bad gateway error
		BadGatewayError(w,err)
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// internal server error
		InternalServerError(w)
	}
}
