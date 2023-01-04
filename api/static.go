// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Static file server.
package api

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/polycloze/polycloze/auth"
	"github.com/polycloze/polycloze/basedir"
	"github.com/polycloze/polycloze/sessions"
)

//go:embed js/dist
var dist embed.FS

//go:embed js/public
var public embed.FS

// Usage: http.Handle("/dist/*", http.StripPrefix("/dist/", serveDist()))
func serveDist() http.Handler {
	sub, err := fs.Sub(dist, "js/dist")
	if err != nil {
		panic(err)
	}
	return cacheUntilBusted(http.FileServer(http.FS(sub)))
}

// Usage: http.Handle("/public/*", http.StripPrefix("/public/", servePublic()))
func servePublic() http.Handler {
	sub, err := fs.Sub(public, "js/public")
	if err != nil {
		panic(err)
	}
	return cacheUntilBusted(http.FileServer(http.FS(sub)))
}

// Usage: http.Handle("/svg/*", http.StripPrefix("/svg/", serveSVG()))
func serveSVG() http.Handler {
	sub, err := fs.Sub(public, "js/public/svg")
	if err != nil {
		panic(err)
	}
	return cacheForever(http.FileServer(http.FS(sub)))
}

// Caches responses for a long time.
func cacheForever(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3153600, immutable")
		next.ServeHTTP(w, r)
	})
}

// Caches responses that contain search params until cache gets busted.
// Bust cache by changing the search params to previously unused params.
// Requests without search params won't set caching instructions.
func cacheUntilBusted(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		if values.Has("v") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		next.ServeHTTP(w, r)
	})
}

// Sets ETag header to data version found in `$DATA_DIR/polycloze/version.txt`.
func versioned(next http.Handler) http.HandlerFunc {
	etag := fmt.Sprintf(`"%s"`, dataVersion)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("ETag", etag)
		next.ServeHTTP(w, r)
	})
}

// Serve files from data directory.
func serveShare() http.Handler {
	return versioned(cacheUntilBusted(http.FileServer(http.Dir(basedir.DataDir))))
}

func serveLanguagesJSON() http.HandlerFunc {
	name := filepath.Join(basedir.StateDir, "languages.json")
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
	return cacheUntilBusted(http.HandlerFunc(handler))
}

func serveCoursesJSON() http.HandlerFunc {
	name := filepath.Join(basedir.StateDir, "courses.json")
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
	return cacheUntilBusted(http.HandlerFunc(handler))
}

func serveUserData(w http.ResponseWriter, r *http.Request) {
	// Page redirects to itself recursively without this check...
	if r.URL.Path == "" || r.URL.Path == "/" {
		http.NotFound(w, r)
		return
	}

	// Check if user is signed in.
	db := auth.GetDB(r)
	s, err := sessions.ResumeSession(db, w, r)
	if err != nil || !s.IsSignedIn() {
		http.NotFound(w, r)
		return
	}

	userID := s.Data["userID"].(int)
	name := filepath.Join(basedir.User(userID), r.URL.Path)
	http.ServeFile(w, r, name)
}
