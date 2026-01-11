package wal

import "os"

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
	return nil
}
