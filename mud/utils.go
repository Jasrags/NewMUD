package mud

import (
	"fmt"
	"io"

	"github.com/i582/cfmt/cmd/cfmt"
)

func CreateEntityRef(area, id string) string {
	return fmt.Sprintf("%s:%s", area, id)
}

func RenderRoom(player *Player, room *Room) string {
	io.WriteString(player.Conn, cfmt.Sprintf("{{%s}}::green|bold\n", player.Room.Title))
	io.WriteString(player.Conn, cfmt.Sprintf("{{%s}}::white\n", player.Room.Description))

	if len(room.Exits) == 0 {
		io.WriteString(player.Conn, cfmt.Sprint("{{There are no exits.}}::red\n"))
	} else {
		io.WriteString(player.Conn, cfmt.Sprint("{{Exits:}}::yellow|bold\n"))
		for direction, _ := range player.Room.Exits {
			io.WriteString(player.Conn, cfmt.Sprintf("{{ - %s}}::yellow\n", direction))
		}
	}

	return fmt.Sprintf("%s\n%s\n", room.Title, room.Description)
}
