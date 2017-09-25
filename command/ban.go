package command

import (
	"github.com/olsdavis/goelan/permission"
	"github.com/olsdavis/goelan/server"
	"strings"
	"fmt"
)

type BanCommand struct {}

func (cmd BanCommand) Labels() []string {
	return []string{"ban"}
}

func (cmd BanCommand) MinArgs() int {
	return 1
}

func (cmd BanCommand) RequiredPermission() string {
	return permission.BanPermission
}

func (cmd BanCommand) Help() string {
	return "ban (player) <reason>"
}

func (cmd BanCommand) Description() string {
	return "Bans the given player for the given reason."
}

func (cmd BanCommand) Execute(label string, args []string, sender CommandSender) {
	if ok, ban := server.Get().GetPlayerByName(args[0]); ok {
		var reason string
		if len(args) > 1 {
			reason = strings.Join(args[1:], " ")
		} else {
			reason = "You have been banned."
		}
		server.Get().BanList.AddPlayer(ban.Player.Profile.UUID, reason)
		ban.Disconnect(reason)
	} else {
		sender.SendMessage(fmt.Sprintf("The player %v could not be found.", args[0]))
	}
}
