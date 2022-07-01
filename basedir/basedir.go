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

func initDataDir() error {
	DataDir = path.Join(xdgDataHome(), "polycloze")
	languages := path.Join(DataDir, "languages")

	if err := os.MkdirAll(languages, 0700); err != nil {
		DataDir = ""
		return err
	}
	return nil
}

func initStateDir() error {
	StateDir = path.Join(xdgStateHome(), "polycloze")
	reviews := path.Join(StateDir, "reviews", "user")
	logs := path.Join(StateDir, "logs", "user")

	if err := os.MkdirAll(reviews, 0700); err != nil {
		StateDir = ""
		return err
	}
	if err := os.MkdirAll(logs, 0700); err != nil {
		StateDir = ""
		return err
	}
	return nil
}

// NOTE This function is different from init.
func Init() error {
	if err := initDataDir(); err != nil {
		return err
	}
	if err := initStateDir(); err != nil {
		return err
	}
	return nil
}
