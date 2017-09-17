package util

type ChatCode string

const (
	// Here are defined all the chat codes that can be used
	// for formatting messages. All you need to do is is include
	// them in your message.
	ObfuscatedChat    ChatCode = "§k"
	BoldChat          ChatCode = "§l"
	StrikethroughChat ChatCode = "§m"
	UnderlineChat     ChatCode = "§n"
	ItalicChat        ChatCode = "§o"
	ResetChat         ChatCode = "§r"

	LimeChat        ChatCode = "§a"
	AquaChat        ChatCode = "§b"
	RedChat         ChatCode = "§c"
	LightPurpleChat ChatCode = "§d"
	YellowChat      ChatCode = "§e"
	WhiteChat       ChatCode = "§f"

	BlackChat    ChatCode = "§0"
	DarkBlueChat ChatCode = "§1"
	GreenChat    ChatCode = "§2"
	DarkAquaChat ChatCode = "§3"
	DarkChat     ChatCode = "§4"
	PurpleChat   ChatCode = "§5"
	GoldChat     ChatCode = "§6"
	GrayChat     ChatCode = "§7"
	DarkGrayChat ChatCode = "§8"
	BlueChat     ChatCode = "§9"
)
