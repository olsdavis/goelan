package world

import "../material"

type Block struct {
	location   SimpleLocation
	material   material.Material
	BlockState byte
}

func (b *Block) GetLocation() SimpleLocation {
	return b.location
}

func (b *Block) GetMaterial() material.Material {
	return b.material
}
