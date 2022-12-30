// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// Content hash of various files
// fileHashes: file URL -> hash.
var fileHashes map[string]string

// Returns hash of file contents.
func hashFile(file fs.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// `path` is a URL path rather than a file path.
func cachedHashFile(path string, file fs.File) (string, error) {
	hash, ok := fileHashes[path]
	if ok {
		return hash, nil
	}
	hash, err := hashFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}
	fileHashes[path] = hash
	return hash, nil
}

// Computes hashes of various files.
func computeHashes() error {
	fileHashes = make(map[string]string)

	// Hash files in js/dist/.
	sub, err := fs.Sub(dist, "js/dist")
	if err != nil {
		// This shouldn't happen.
		panic(err)
	}
	_ = fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, _ error) error {
		file, err := sub.Open(path)
		if err != nil {
			// Ignore error.
			return nil
		}
		defer file.Close()

		_, _ = cachedHashFile(filepath.Join("/dist", path), file)
		return nil
	})

	// Hash some public files.
	sub, err = fs.Sub(public, "js/public")
	if err != nil {
		// This shouldn't happen.
		panic(err)
	}
	_ = fs.WalkDir(sub, ".", func(path string, d fs.DirEntry, _ error) error {
		// Ignore fonts and icons.
		if strings.HasPrefix(path, "fonts") || strings.HasPrefix(path, "svg") {
			return nil
		}

		file, err := sub.Open(path)
		if err != nil {
			// Ignore error.
			return nil
		}
		defer file.Close()
		_, _ = cachedHashFile(filepath.Join("/public", path), file)
		return nil
	})
	return nil
}

// Returns URL of file with the content hash as its version.
// If the file hasn't been hashed, simply returns the input URL.
func versionedURL(url string) string {
	hash, ok := fileHashes[url]
	if ok {
		return url + "?v=" + hash
	}
	return url
}
