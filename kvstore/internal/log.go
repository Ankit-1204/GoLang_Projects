package internal

import (
	"os"
	"sync"
)

type LogFile struct {
	LSN       int64
	Data      interface{}
	Operation string
}

type Wal struct {
	file  *os.File
	entry *[]LogFile
	lock  sync.Mutex
	lsn   int64
}

func (w *Wal) Append(log LogFile) {
	w.lock.Lock()
	defer w.lock.Unlock()
	log.LSN = w.lsn + 1
	w.lsn = log.LSN

}
