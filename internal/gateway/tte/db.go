package tte

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dan-frohlich/tabetopevents/internal/logging"
)

type DB struct {
	path string
	log  logging.Logger
}

func NewDB(log logging.Logger) DB {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Error("could not locate user home dir", "error", err)
		home = "."
	}
	newpath := filepath.Join(home, ".tte_db")
	_ = os.MkdirAll(newpath, os.ModePerm)

	return DB{path: newpath, log: log}
}

func (db DB) mkdir(kind string) {
	_ = os.MkdirAll(db.kindPath(kind), os.ModePerm)
}

func (db DB) kindPath(kind string) (path string) {
	return filepath.Join(db.path, kind)
}

func (db DB) itemPath(id string, kind string, dataType string) (path string) {
	return filepath.Join(db.path, kind, fmt.Sprintf("%s.%s", id, dataType))
}

func (db DB) Store(id string, kind string, dataType string, data []byte) error {
	db.mkdir(kind)
	filePath := db.itemPath(id, kind, dataType)
	db.log.Debug("writing cache", "path", filePath)
	return os.WriteFile(filePath, data, os.FileMode(0644))
}

func (db DB) Read(id string, kind string, dataType string) (data []byte, err error) {
	filePath := db.itemPath(id, kind, dataType)
	db.log.Debug("reading cache", "path", filePath)
	data, err = os.ReadFile(filePath)
	if err != nil {
		err = fmt.Errorf("unable to load %s : %s", filePath, err)
	}
	return data, err
}

func (db DB) CacheAge(id string, kind string, dataType string) (time.Duration, error) {
	fileInfo, err := os.Stat(db.itemPath(id, kind, dataType))
	if err != nil {
		return 0, err
	}
	modTime := fileInfo.ModTime()
	return time.Since(modTime).Truncate(time.Second), nil
}
