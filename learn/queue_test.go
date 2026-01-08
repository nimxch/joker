package learn

import (
	"encoding/binary"
	"errors"
)

const (
	MAX_NODE_BYTES = 1024 * 8 // Max 8 KB
	LENGTH_BYTES   = 4        // uint32 length prefix
)

// Node Data structure
type Node struct {
	prev        *Node
	next        *Node
	readOffset  uint32 // Head
	writeOffset uint32 // Tail
	tail        int16
	content     [MAX_NODE_BYTES]byte // First 4 Byte Length, Rest PayLoad
}

// Queue Data Structure
type Queue struct {
	head *Node
	tail *Node
	len  uint32
}

// Wal Interface
type WAL interface {
	AppendEnqueue(payload []byte) error
}

var (
	EntryTooLarge = errors.New("Entry larger than MAX_NODE_SIZE")
)

// Example content
// | len  | data | len | data | len | data
// len : uint32
// payload: size of data
//

func NewNode() *Node {
	return &Node{}
}

func (q *Queue) Enqueue(wal WAL, payload []byte) error {
	payloadSize := len(payload)
	entrySize := payloadSize + LENGTH_BYTES
	// Validate the payload size
	if payloadSize > MAX_NODE_BYTES {
		return EntryTooLarge
	}

	// Write into WAL (durability)
	if err := wal.AppendEnqueue(payload); err != nil {
		return err
	}

	// In-Memory apply
	var node *Node

	if q.tail == nil {
		node = NewNode()
		q.head = node
		q.tail = node
	} else {
		node = q.tail
	}

	// Capacity Check
	if node.writeOffset+uint32(entrySize) > MAX_NODE_BYTES {
		node = NewNode()
		node.prev = q.tail
		q.tail.next = node
		q.tail = node
	}

	// Append Entry (Atomic)
	offSet := node.writeOffset

	// Write length prefix
	binary.LittleEndian.PutUint32(
		node.content[offSet:offSet+LENGTH_BYTES],
		uint32(payloadSize),
	)

	// Write the payload
	copy(
		node.content[offSet+LENGTH_BYTES:offSet+uint32(entrySize)],
		payload,
	)

	// Update the cursor
	node.writeOffset += uint32(entrySize)
	return nil
}
