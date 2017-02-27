package player

type GameMode int

const (
	SurvivalMode GameMode = iota
	CreativeMode
	AdventureMode
	SpectatorMode
)
