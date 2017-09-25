package player

type GameMode byte

const (
	SurvivalMode GameMode = iota
	CreativeMode
	AdventureMode
	SpectatorMode
)
