package internal

import (
	"encoding/gob"
	"os"
	"sync"
)

type LogFile struct {
	LSN       int64
	Data      interface{}
	Operation string
}

type Wal struct {
	file *os.File
	lock sync.Mutex
	lsn  int64
}

func (w *Wal) Append(log LogFile) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	log.LSN = w.lsn + 1
	w.lsn = log.LSN

	encoder := gob.NewEncoder(w.file)
	if err := encoder.Encode(log); err != nil {
		return err
	}
	// Actually commiting to the real file and not to a buffer
	return w.file.Sync()
}
