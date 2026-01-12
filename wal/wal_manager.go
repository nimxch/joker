package wal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nimxch/joker/custom"
)

var ErrEmptyPayload = errors.New("wal: empty payload")

type WalManager struct {
	fd *os.File
}

// Constructor for WalManager
func InitWal(path string) (*WalManager, error) {
	path = filepath.Clean(path) //Sanitize the Path

	if strings.HasSuffix(path, string(os.PathSeparator)) {
		return nil, custom.ErrInvalidFilePath
	}
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

//  Payload Structure
// | recordLen (4) |
// | crc32 (4)     |
// | opType (1)    |
// | payloadLen(4) |
// | payload (N)   |

func (w *WalManager) AppendEnqueue(payload []byte) error {
	// 9 = 4byte[CRC] + 1byte[OpType]+ 4byte [PayloadLen]
	payloadLen := len(payload)
	if payloadLen == 0 {
		return ErrEmptyPayload
	}
	recordLength := uint32(4 + 1 + 4 + payloadLen)

	buf := make([]byte, 4+recordLength)

	// Payload building starts
	offset := 0
	// Put 4 bytes of payload length to buffer
	binary.LittleEndian.PutUint32(buf[offset:], recordLength) //recordLen
	offset += 4
	crcOffset := offset
	offset += 4                          // CRC will be stored between crcOffset and offset
	buf[offset] = OPERATION_TYPE_ENQUEUE // OpType
	// Move offset to 1 byte, as len(opType is 1 byte)
	offset += 1
	binary.LittleEndian.PutUint32(buf[offset:], uint32(payloadLen)) // payload len
	// Move offset to offset
	offset += 4
	copy(buf[offset:], payload)

	crc := crc32.ChecksumIEEE(buf[crcOffset+4:]) // CRC after CRC field
	//put crc offset into the buf
	binary.LittleEndian.PutUint32(buf[crcOffset:], crc)
	// payload building end
	// Write to the wal and fsync
	n, err := w.fd.Write(buf)
	if err != nil {
		return err
	}

	if n != len(buf) {
		return io.ErrShortWrite
	}

	// do the fsync
	if err := w.fd.Sync(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// TODO: Implement Write for dequeue
// TODO: Remove code dups
func (w *WalManager) AppendDequeue(payload []byte) error {
	return nil
}

func (w *WalManager) Flush() error {
	return w.fd.Sync()
}

func GetCrc(payload []byte) uint32 {
	return crc32.Checksum(payload, crc32.MakeTable(Koopman))
}
