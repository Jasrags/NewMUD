package main

import (
	"os"
	"strings"

	"github.com/Jasrags/NewMUD/internal/game"
	"github.com/i582/cfmt/cmd/cfmt"
	"github.com/muesli/termenv"
)

func main() {

	gs := game.NewGameServer()
	gs.SetupConfig()
	gs.SetupLogger()
	gs.OutputConfig()

	game.EntityMgr.LoadDataFiles()
	game.AccountMgr.LoadDataFiles()
	game.CharacterMgr.LoadDataFiles()
	game.GameTimeMgr = game.NewGameTime()

	acct := game.AccountMgr.GetByUsername("Jasrags")
	char := game.CharacterMgr.GetCharacterByName("fred")
	// char.Role = game.CharacterRolePlayer
	char.Role = game.CharacterRoleAdmin
	room := game.EntityMgr.GetRoom("main_place_arcade_1f")

	char.MoveToRoom(room)

	player1 := game.NewCharacter()
	player1.Name = "Player1"
	player1.Room = room
	player1.MetatypeID = "Human"
	char.Room.AddCharacter(player1)

	player2 := game.NewCharacter()
	player2.Name = "Player2"
	player2.Room = room
	player2.MetatypeID = "Ork"
	char.Room.AddCharacter(player2)

	// mob1 := game.EntityMgr.GetMob("orc")
	// if mob1 == nil {
	// 	panic("Mob not found")
	// }
	// mob1.ID = "ork1"
	// mob2 := game.EntityMgr.GetMob("goblin")
	// mob2.ID = "goblin1"
	// mob1 := game.NewMob()
	// mob1.Name = "Mob1"
	// mob2 := game.NewMob()
	// mob2.Name = "Mob2"

	char.Room.AddMob(game.EntityMgr.GetMob("thug_basic"))
	char.Room.AddMob(game.EntityMgr.GetMob("thug_lieutenant"))

	// item1 := game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock")
	// item2 := game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock")
	// item3 := game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key")

	char.Room.Inventory.AddItem(game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock"))
	char.Room.Inventory.AddItem(game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock"))
	char.Room.Inventory.AddItem(game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key"))

	game.GameTimeMgr.Minutes = 90
	char.Prompt = "{{time}} {{date}} {{>}}::white"

	// 	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
	// A
	// AB
	// ABC
	// {{A}}::red|bold
	// {{AB}}::red|bold
	// {{ABC}}::red|bold
	// Nulla {{convallis}}::yellow egestas rhoncus. Donec facilisis fermentum sem, ac viverra ante luctus vel.
	// Phasellus ultrices nulla quis nibh. {{Quisque}}::cyan|bold a lectus. Donec consectetuer ligula vulputate sem tristique cursus.`

	// Example 2: Use custom options.
	// customOptions := &game.WrapOptions{
	// 	BorderType: game.BorderTypeDouble,
	// 	TextWidth:  50,
	// 	// PaddingTop:    1,
	// 	// PaddingBottom: 1,
	// 	// PaddingLeft:   1,
	// 	// PaddingRight:  1,
	// 	BorderColor: "red|bold",
	// 	Alignment:   game.TextAlignLeft,
	// }
	// borderStyle := lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.RoundedBorder()).
	// 	BorderForeground(lipgloss.Color("63")).
	// Padding(0, 1, 0, 1)

	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()
	// output.DisableMouse()
	// output.EnableMouse()
	// output.AltScreen()
	// defer output.ExitAltScreen()
	var sb strings.Builder
	// output.WriteString(stripCfmt("{{convallis}}::yellow"))
	// output.WriteString(game.CRLF)
	// output.WriteString(stripCfmt("{{Quisque}}::cyan|bold"))
	// output.WriteString(game.CRLF)
	// output.WriteString(game.WrapTextInBorder(text, nil))
	// output.WriteString(cfmt.Sprint(game.WrapTextInBorder(text, nil)))
	// output.WriteString(game.CRLF)
	// output.WriteString(game.WrapTextInBorder(text, customOptions))
	// sb.WriteString(game.CRLF)
	// sb.WriteString(game.RenderCharacterTable(char))
	// sb.WriteString(borderStyle.Width(80).Render(cfmt.Sprint(text)))
	sb.WriteString(game.CRLF)
	output.WriteString(game.RenderRoom(acct, char, room))
	// output.WriteString(game.RenderCharacterTable(char))
	// output.WriteString(game.RenderPromptMenu("Main Menu", []string{"Enter Game", "Create Character", "Change Password", "Quit"}))
	// output.WriteString(game.RenderPrompt(char))
	cfmt.Print(sb.String())
}

// Use a regex that removes any {{...}} with optional :: formatting.
// var cfmtRegex = regexp.MustCompile(`\{\{[^}]*\}\}(?:::[^}]*)?`)
// var cfmtRegex = regexp.MustCompile(`\{\{(.*?)\}\}(?:::[^}]+)?`)

// func stripCfmt(s string) string {
// 	return cfmtRegex.ReplaceAllString(s, "$1")
// }

// func visibleLength(s string) int {
// 	return len([]rune(stripCfmt(s)))
// }

// func main() {
// 	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
// Nulla convallis egestas rhoncus. Donec facilisis fermentum sem, ac viverra ante luctus vel.
// Phasellus ultrices nulla quis nibh. Quisque a lectus. Donec consectetuer ligula vulputate sem tristique cursus.`

// 	// Wrap the text in a border with an inner width of 40 characters.
// 	output := WrapTextInBorder(text, roundedBorder, 40)
// 	fmt.Println(output)
// }
