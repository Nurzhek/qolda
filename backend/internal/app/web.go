package app

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// rootHandler is the monolith on one port: /api/* -> grpc-gateway, /files/* -> photos, /docs/* -> Swagger UI,
// everything else -> SPA
func rootHandler(api http.Handler, uploadDir, webDir string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	mux.Handle("/files/", http.StripPrefix("/files/", noDirListing(http.FileServer(http.Dir(uploadDir)))))
	// Swagger UI (docs/index.html) on the app port, so "Try it out" hits /api on the same origin
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	if webDir != "" {
		mux.Handle("/", spaHandler(webDir))
	}
	return mux
}

func noDirListing(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		h.ServeHTTP(w, r)
	})
}

// spaHandler serves existing files as is and index.html for any other path (history mode)
func spaHandler(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+r.URL.Path)))
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			if strings.HasPrefix(r.URL.Path, "/assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fs.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		http.ServeFile(w, r, filepath.Join(dir, "index.html"))
	})
}
