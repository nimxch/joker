package wal

import (
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"
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

	// Extract the dir path and create it of not exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

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
	// 9 = 4byte[CRC] + 1byte[OptType]+ 4byte [PayloadLen]
	recordLength := uint32(9) + uint32(len(payload))

	walRecord := WalRecord{
		payloadLen:   uint32(len(payload)),
		payload:      payload,
		opType:       uint8(OPERATION_TYPE_ENQUEUE),
		crc32:        GetCrc(payload),
		recordLength: recordLength,
	}

	err := walRecord.WriteFsync(w)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(walRecord)
	return nil
}

func (w *WalManager) AppendDequeue(payload []byte) error {
	return nil
}

func (w *WalManager) Flush() error {
	return nil
}

func GetCrc(payload []byte) uint32 {
	return crc32.Checksum(payload, crc32.MakeTable(Koopman))
}
