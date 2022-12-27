// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package api

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

// Content hash of various files
var fileHashes map[string]string

// Returns hash of file contents.
func hashFile(file fs.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %v", err)
	}
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func cachedHashFile(filename string, file fs.File) (string, error) {
	// Uses `filename` instead of `file.Stat()`, because `file` may be in a
	// subdirectory.
	name := filename
	hash, ok := fileHashes[name]
	if ok {
		return hash, nil
	}
	hash, err := hashFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to hash file: %v", err)
	}
	fileHashes[name] = hash
	return hash, nil
}

// Computes hashes of various files.
func computeHashes() error {
	fileHashes = make(map[string]string)

	// Hash files in js/dist/.
	_ = fs.WalkDir(dist, "js/dist", func(path string, d fs.DirEntry, _ error) error {
		file, err := dist.Open(path)
		if err != nil {
			// Ignore error.
			return nil
		}
		defer file.Close()
		_, _ = cachedHashFile(path, file)
		return nil
	})

	// Hash some public files.
	sub, err := fs.Sub(public, "js/public")
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
		_, _ = cachedHashFile(filepath.Join("js", "public", path), file)
		return nil
	})
	return nil
}
