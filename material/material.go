package material

var (
	Stone Material = Material{1, "stone"}
	Grass = Material{2, "grass"}
	Dirt = Material{3, "dirt"}
	Cobblestone = Material{4, "cobblestone"}
	WoodPlank = Material{5, "planks"}
	Sapling = Material{6, "sapling"}
	Bedrock = Material{7, "bedrock"}
)

var (
	materialMap map[int]Material = map[int]Material{
		Stone.ID: Stone,
		Grass.ID: Grass,
		Dirt.ID: Dirt,
		Cobblestone.ID: Cobblestone,
		WoodPlank.ID: WoodPlank,
		Sapling.ID: Sapling,
		Bedrock.ID: Bedrock,
	}
)

type Material struct {
	ID   int
	Name string
}

func GetById(id int) Material {
	return materialMap[id]
}
