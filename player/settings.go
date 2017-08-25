package player

type ChatMode int
type Hand int

const (
	EnabledChatMode      ChatMode = iota
	CommandsOnlyChatMode
	HiddenChatMode

	LeftHand  Hand = iota
	RightHand
)

type ClientSettings struct {
	Locale             string
	ViewDistance       byte
	ChatMode           ChatMode
	ColorsEnabled      bool
	DisplayedSkinParts uint8
	MainHand           Hand
}

// IsLeftHanded returns true if the client's main hand is the left one.
func (settings *ClientSettings) IsLeftHanded() bool {
	return settings.MainHand == LeftHand
}

// IsRightHanded returns true if the client's main hand is the right one.
func (settings *ClientSettings) IsRightHanded() bool {
	return settings.MainHand == RightHand
}
