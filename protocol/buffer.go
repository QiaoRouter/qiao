package protocol

type Buffer struct {
	Octet []byte
}

func (buf *Buffer) Length() int {
	return len(buf.Octet)
}
