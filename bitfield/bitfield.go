package bitfield

type Bitfield []byte

func (bf Bitfield) HasPiece(idx int) bool {
    byteIdx := idx/8
    offset := idx%8
    return (bf[byteIdx] >> (7 - offset) & 1) != 0
}

func (bf Bitfield) SetPiece(idx int) {
    byteIdx := idx/8
    offset := idx%8
    bf[byteIdx] |= 1 << (7 - offset)
}
