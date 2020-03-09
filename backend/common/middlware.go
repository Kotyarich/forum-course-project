package common

import (
	"github.com/dimfeld/httptreemux"
	"net/http"
)

func setCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD, PUT")
	w.Header().Set("Access-Control-Max-Age", "600")
}

func CORSMiddlware(f httptreemux.HandlerFunc) httptreemux.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, ps map[string]string) {
		setCORSHeaders(w)
		f(w, r, ps)
	}
}

func CORSHandler(w http.ResponseWriter, r *http.Request, ps map[string]string) {
	setCORSHeaders(w)
}

