package queue

// Queue Rules:
//
// 1. If q.head == nil, then q.tail == nil and q.len == 0
// 2. If q.len > 0, then q.head != nil and q.tail != nil
// 3. For any node: 0 <= readOffset <= writeOffset <= MAX_NODE_BYTES
// 4. There are no empty nodes in the list
// 5. Enqueue mutates only tail and writeOffset
// 6. Dequeue mutates only head and readOffset
// 7. q.len == total number of enqueued but not dequeued messages

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
	readOffset  uint32               // Head
	writeOffset uint32               // Tail
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

// Errors
var (
	ErrEntryTooLarge = errors.New("Entry larger than MAX_NODE_SIZE")
	ErrQueueLength   = errors.New("Queue Length is invalid")
	ErrQueueOffset   = errors.New("Offset miss-match")
	ErrEmptyQueue    = errors.New("queue is empty")
)

// Example content
// | len  | data | len | data | len | data
// len : uint32
// payload: size of data
//

func NewNode() *Node {
	return &Node{}
}

func (q *Queue) Enqueue(payload []byte) error {
	payloadSize := len(payload)
	entrySize := payloadSize + LENGTH_BYTES
	// Validate the payload size
	if entrySize > MAX_NODE_BYTES {
		return ErrEntryTooLarge
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
	offSet := int(node.writeOffset)

	// Write length prefix
	binary.LittleEndian.PutUint32(
		node.content[offSet:offSet+LENGTH_BYTES],
		uint32(payloadSize),
	)

	// Write the payload
	copy(
		node.content[offSet+LENGTH_BYTES:offSet+entrySize],
		payload,
	)

	// Update the cursor
	node.writeOffset += uint32(entrySize)
	q.len += 1
	return nil
}

func (q *Queue) Peek() ([]byte, error) {
	// Remove Empty Head
	q.advanceHeadPastEmpty()

	if q.head == nil {
		return nil, ErrEmptyQueue
	}

	node := q.head

	if node.readOffset >= node.writeOffset {
		return nil, ErrQueueOffset
	}

	offSet := node.readOffset

	payLoadLen := binary.LittleEndian.Uint32(
		node.content[offSet : offSet+LENGTH_BYTES],
	)

	// Calculate Payload boundary
	start := offSet + LENGTH_BYTES
	end := start + payLoadLen

	if end > node.writeOffset {
		panic("queue corruption: payload exceeds writeOffset")
	}

	payLoad := make([]byte, payLoadLen)
	copy(payLoad, node.content[start:end])
	return payLoad, nil
}

func (q *Queue) Dequeue() ([]byte, error) {
	// Remove Empty Head
	q.advanceHeadPastEmpty()

	if q.head == nil {
		return nil, ErrEmptyQueue
	}

	node := q.head

	if node.readOffset >= node.writeOffset {
		return nil, ErrQueueOffset
	}

	offSet := node.readOffset

	payLoadLen := binary.LittleEndian.Uint32(
		node.content[offSet : offSet+LENGTH_BYTES],
	)

	// Calculate Payload boundary
	start := offSet + LENGTH_BYTES
	end := start + payLoadLen

	if end > node.writeOffset {
		panic("queue corruption: payload exceeds writeOffset")
	}

	payLoad := make([]byte, payLoadLen)
	copy(payLoad, node.content[start:end])

	// Move the read head to next
	node.readOffset = end
	if node.readOffset == node.writeOffset {
		// node is empty
		q.head = node.next
		if q.head == nil {
			q.tail = nil
		} else {
			q.head.prev = nil
		}
	}
	if q.len == 0 {
		return nil, ErrQueueLength
	}
	q.len -= 1
	return payLoad, nil
}

func (q *Queue) advanceHeadPastEmpty() {
	for q.head != nil && q.head.readOffset == q.head.writeOffset {
		q.head = q.head.next
		if q.head != nil {
			q.head.prev = nil
		} else {
			q.tail = nil
		}
	}

}
