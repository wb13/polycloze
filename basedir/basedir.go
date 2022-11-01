// Copyright (c) 2022 Levi Gruspe
// License: GNU AGPLv3 or later

package basedir

import (
	"log"
	"os"
	"path"
)

var (
	Home     string
	DataDir  string
	StateDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	Home = home

	DataDir = path.Join(xdgDataHome(), "polycloze")
	if err := initStateDir(); err != nil {
		log.Fatal(err)
	}
}

func xdgDataHome() string {
	val := os.Getenv("XDG_DATA_HOME")
	if val != "" {
		return val
	}
	return path.Join(Home, ".local", "share")
}

func xdgStateHome() string {
	val := os.Getenv("XDG_STATE_HOME")
	if val != "" {
		return val
	}
	return path.Join(Home, ".local", "state")
}

func initStateDir() error {
	StateDir = path.Join(xdgStateHome(), "polycloze")
	users := path.Join(StateDir, "users")

	if err := os.MkdirAll(users, 0o700); err != nil {
		StateDir = ""
		return err
	}
	return nil
}
