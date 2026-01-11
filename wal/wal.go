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
}

const (
	OPERATION_TYPE_UNKNOWN = 0
	OPERATION_TYPE_ENQUEUE = 1
	OPERATION_TYPE_DEQUEUE = 2
)
