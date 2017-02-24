package world

import "../material"

type Block struct {
	location SimpleLocation
	material material.Material
}

func (b *Block) GetLocation() SimpleLocation {
	return b.location
}

func (b *Block) GetMaterial() material.Material {
	return b.material
}
