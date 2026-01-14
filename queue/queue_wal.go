package queue

import (
	"github.com/nimxch/joker/custom"
	"github.com/nimxch/joker/wal"
)

func CommitEnqueue(q *Queue, w wal.WAL, payload []byte) error {
	if w == nil {
		return custom.ErrWalMissing
	}
	if err := w.AppendEnqueue(payload); err != nil {
		return err
	}
	return q.Enqueue(payload)
}

func Sync(w wal.WAL) error {
	if w == nil {
		return custom.ErrWalMissing
	}
	return w.WalReplay()
}
