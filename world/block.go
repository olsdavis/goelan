package world

import (
	"bytes"
	"github.com/olsdavis/goelan/material"
)

type Block struct {
	location   *Location3i
	material   material.Material
	BlockState byte
}

func NewBlock(loc *Location3i, mat material.Material, state byte) *Block {
	return &Block{
		location:   loc,
		material:   mat,
		BlockState: state,
	}
}

func (b *Block) GetLocation() *Location3i {
	return b.location
}

func (b *Block) GetMaterial() material.Material {
	return b.material
}

func (b *Block) writeBlock(buffer *bytes.Buffer) {

}
