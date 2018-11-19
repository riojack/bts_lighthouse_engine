package logic

// ReadWriteBridge provides a communication interface between game logic and I/O.
//
// When the game needs to read from or write to a data stream, it should not be aware of the specific data stream kind (socket, HTTP, etc).
type ReadWriteBridge interface {
	Read() (value string, err error)
	WriteString(value string) (err error)
}
