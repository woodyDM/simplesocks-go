package pro

type encrypter interface {
	enc(raw []byte) []byte

	dec(data []byte) []byte
}

type encrypterFactory interface {

	iv() []byte

	ivLen() int

	newInstance() encrypter

}

type caesarFactory struct {
	offset byte
}

func (c caesarFactory) iv() []byte {
	panic("implement me")
}

func (c caesarFactory) ivLen() int {
	panic("implement me")
}

func (c caesarFactory) newInstance() encrypter {
	return &caesarEncrypter{offset:c.offset}
}

type caesarEncrypter struct {
	offset byte
}

func (c *caesarEncrypter) enc(raw []byte) []byte {
	return c.doLoop(raw, func(r []byte, pos int) {
		r[pos] = raw[pos] + c.offset
	})
}

func (c *caesarEncrypter) dec(data []byte) []byte {
	return c.doLoop(data, func(r []byte, pos int) {
		r[pos] = data[pos] - c.offset
	})
}

func (c *caesarEncrypter) doLoop(data []byte, consumer func(r []byte, pos int)) []byte {
	l := len(data)
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		consumer(result, i)
	}
	return result

}
