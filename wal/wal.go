package wal

type WalRecord struct {
	recordLength uint32
	crc32        uint32
	opType       uint8
	payloadLen   uint32
	payload      []byte
}

type WAL interface {
	AppendEnqueue(payload []byte) error
	AppendDequeue(payload []byte) error // optional for later
	Flush() error                       // --> Explicitly do Fsync
	Sync() error
}

const (
	OPERATION_TYPE_UNKNOWN = 0
	OPERATION_TYPE_ENQUEUE = 1
	OPERATION_TYPE_DEQUEUE = 2
)

// Constants for crc checksum
const (
	IEEE       = 0xedb88320
	Castagnoli = 0x82f63b78
	Koopman    = 0xeb31d82e
)

func (wr *WalRecord) WriteFsync(w *WalManager) error {
	return nil
}
