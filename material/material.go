package material

var materialMap map[int]Material

func init() {
	materialMap = make(map[int]Material)
	materialMap[0] = Material{0, "dirt"}
}

type Material struct {
	ID   int
	Name string
}

func GetById(id int) Material {
	return materialMap[id]
}
