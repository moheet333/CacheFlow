package server

import (
	root "CacheFlow/cmd"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CachedData struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Status  int               `json:"status"`
}

func (s *Server) UniversalHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getHandler(w, r)
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

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	originUrl,_ := url.JoinPath(root.Newflag.Origin , r.URL.Path)
	cacheKey := "cache>" + originUrl

	cachedResponse, err := s.db.GetCache(ctx, cacheKey)
	if err == nil && cachedResponse != "" {
		// cache hit
		var cachedData CachedData
		err := json.Unmarshal([]byte(cachedResponse), &cachedData)
		if err != nil {
			InternalServerError(w)
			return
		}
		for key, value := range cachedData.Headers {
			w.Header().Set(key, value)
		}
		w.Header().Set("Cache-Status", "HIT")
		w.WriteHeader(cachedData.Status)
		w.Write([]byte(cachedData.Body))
		return
	}
	// cache miss
	currentUrl := root.Newflag.Origin + r.URL.Path
	resp, err := http.Get(currentUrl)
	if err != nil {
		BadGatewayError(w, err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		InternalServerError(w)
		return
	}

	responseData := CachedData{
		Headers: make(map[string]string),
		Body:    string(bodyBytes),
		Status:  resp.StatusCode,
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
			responseData.Headers[key] = value
		}
	}

	jsonData, _ := json.Marshal(responseData)
	err = s.db.SetCache(ctx, cacheKey, string(jsonData), 10*time.Minute)
	if err != nil {	
		InternalServerError(w)
		return
	}
	w.Header().Set("Cache-Status", "MISS")
	w.WriteHeader(resp.StatusCode)
	w.Write(bodyBytes)
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
	forwardHandler(w, r, currentUrl)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	currentUrl := root.Newflag.Origin + r.URL.Path
	forwardHandler(w, r, currentUrl)
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
		BadGatewayError(w, err)
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
