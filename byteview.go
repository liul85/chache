package chache

//ByteView encapsulate cache object
type ByteView struct {
	b []byte
}

//Len returns length of byteview
func (v ByteView) Len() int64 {
	return int64(len(v.b))
}

// ByteSlice returns copy of ByteView
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

//String convert byteview to string
func (v ByteView) String() string {
	return string(v.b)
}
