package wal

import (
	"fmt"
	"os"
)

type WalManager struct {
	fd *os.File
}

// Constructor for WalManager
func InitWal(path string) (*WalManager, error) {
	fd, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return nil, err
	}
	return &WalManager{
		fd: fd,
	}, nil
}

func (w *WalManager) AppendEnqueue(payload []byte) error {
	fmt.Println("Payload: ", payload)
	fmt.Println("Payload len: ", len(payload))
	return nil
}

func (w *WalManager) AppendDequeue(payload []byte) error {
	return nil
}

func (w *WalManager) Flush() error {
	return nil
}
