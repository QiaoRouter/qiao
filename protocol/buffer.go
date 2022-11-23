package protocol

type Buffer struct {
	Octet []byte
}

func (buf *Buffer) Length() uint16 {
	return uint16(len(buf.Octet))
}
