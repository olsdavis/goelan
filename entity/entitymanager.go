package entity

type EntityManager struct {
	entities map[int]Entity
}

func (manager EntityManager) AddEntity(entity Entity) {
	manager.entities[len(manager.entities)] = entity
}

type Entity interface {}
