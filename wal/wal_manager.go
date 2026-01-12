package wal

import (
	"fmt"
	"os"
	"strings"

	"github.com/nimxch/joker/custom"
)

type WalManager struct {
	fd *os.File
}

// Constructor for WalManager
func InitWal(path string) (*WalManager, error) {
	if strings.HasSuffix(path, "/") {
		return nil, custom.ErrInvalidFilePath
	}
	if strings.HasPrefix(path, "/") {
		// sanitize string
		path, _ = strings.CutPrefix(path, "/")
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	path = fmt.Sprintf("%s/%s", wd, path)
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
