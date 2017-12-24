package protocol

// Represents server list ping response.
type ServerListPing struct {
	Ver  Version       `json:"version"`
	Pl   Players       `json:"players"`
	Desc ChatComponent `json:"description"`
	Fav  string        `json:"favicon,omitempty"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol uint32 `json:"protocol"`
}

type Players struct {
	Max    uint `json:"max"`
	Online uint `json:"online"`
}

// Chat
type MessageMode int

const (
	// chat message (only for players)
	ChatMessageMode MessageMode = iota
	// system message (what you should use)
	DefaultMessageMode
	// action bar message
	ActionBarMode
)

// Animation
const (
	SwingMainArmAnimation        = iota
	TakeDamageAnimation
	LeaveBedAnimation
	SwingOffHandAnimation
	CriticalEffectAnimation
	MagicCriticalEffectAnimation
)
