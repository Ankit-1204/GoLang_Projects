package internal

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Oper string

type LogFile struct {
	LSN       int64
	Key       string
	Data      interface{}
	Operation Oper
}

const OpSet Oper = "SET"
const OpDelete Oper = "DELETE"

type Wal struct {
	dir  string
	file *os.File
	lock sync.Mutex
	lsn  int64
}

func CreateLogFile(dir string) (*Wal, error) {
	if err := os.Mkdir(dir, 0755); err != nil {
		return nil, err
	}
	w := &Wal{dir: dir}
	if err := w.NewFile(); err != nil {
		return nil, err
	}
	return w, nil

}
func (w *Wal) NewFile() error {
	if w.file != nil {
		w.file.Close()
	}
	filename := filepath.Join(w.dir, fmt.Sprintf("wal-%06d.log", w.lsn))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = f
	return nil
}

func (w *Wal) Append(op Oper, key string, value interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	log := LogFile{
		LSN:       w.lsn + 1,
		Operation: op,
		Data:      value,
		Key:       key,
	}
	w.lsn = log.LSN

	encoder := gob.NewEncoder(w.file)
	if err := encoder.Encode(log); err != nil {
		return err
	}
	// Actually commiting to the real file and not to a buffer
	if err := w.file.Sync(); err != nil {
		return err
	}
	return nil
}

func (w *Wal) Recover(apply func(LogFile)) error {
	files, err := filepath.Glob(filepath.Join(w.dir, "wal-*.log"))
	if err != nil {
		return err
	}
	for _, f := range files {
		file, err := os.Open(f)
		if err != nil {
			return err
		}
		defer file.Close()
		decoder := gob.NewDecoder(file)
		for {
			var rec LogFile
			if err := decoder.Decode(&rec); err != nil {
				return err
			}
			apply(rec)
			if rec.LSN > w.lsn {
				w.lsn = rec.LSN
			}
		}
	}
	return nil
}

func (w *Wal) Close() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.file != nil {
		w.file.Close()
	}
	return nil
}
