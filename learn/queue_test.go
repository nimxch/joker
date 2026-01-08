package learn

const (
	MAX_PACKET = 1024 * 8 // Max 8 Byte
)

type Node struct {
	prev        *Node
	next        *Node
	readOffset  int32 // Head
	writeOffset int32 // Tail
	tail        int16
	content     []byte // First 2 Byte Length, Rest PayLoad
}
