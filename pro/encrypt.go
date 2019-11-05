package pro

type encrypter interface {
	enc(raw []byte) []byte

	dec(data []byte) []byte
}

type encFactory interface {
	iv() []byte

	ivLen() int

	newEncrypter() encrypter
}

type caesarFactory struct {
	offset byte
}

func (c caesarFactory) iv() []byte {
	return []byte{100}
}

func (c caesarFactory) ivLen() int {
	return 1
}

func (c caesarFactory) newEncrypter() encrypter {
	return &caesarEncrypter{offset: c.offset}
}

type caesarEncrypter struct {
	offset byte
}

func (c *caesarEncrypter) enc(raw []byte) []byte {
	return c.doLoop(len(raw), func(r []byte, pos int) {
		r[pos] = raw[pos] + c.offset
	})
}

func (c *caesarEncrypter) dec(data []byte) []byte {
	return c.doLoop(len(data), func(r []byte, pos int) {
		r[pos] = data[pos] - c.offset
	})
}

func (c *caesarEncrypter) doLoop(l int, consumer func(r []byte, pos int)) []byte {

	result := make([]byte, l)
	for i := 0; i < l; i++ {
		consumer(result, i)
	}
	return result

}
