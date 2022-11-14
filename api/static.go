// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

// Static file server.
package api

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/lggruspe/polycloze/basedir"
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
	return http.FileServer(http.FS(sub))
}

// Usage: http.Handle("/public/*", http.StripPrefix("/public/", servePublic()))
func servePublic() http.Handler {
	sub, err := fs.Sub(public, "js/public")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}

// Usage: http.Handle("/svg/*", http.StripPrefix("/svg/", serveSVG()))
func serveSVG() http.Handler {
	sub, err := fs.Sub(dist, "js/dist/svg")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}

// Serve files from data directory.
func serveShare() http.Handler {
	return http.FileServer(http.Dir(basedir.DataDir))
}
