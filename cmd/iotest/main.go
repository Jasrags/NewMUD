package main

import (
	"strings"

	"github.com/Jasrags/NewMUD/internal/game"
	"github.com/i582/cfmt/cmd/cfmt"
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

	// acct := game.AccountMgr.GetByUsername("Jasrags")
	char := game.CharacterMgr.GetCharacterByName("Jasrags")
	room := game.EntityMgr.GetRoom("limbo")

	char.MoveToRoom(room)

	player1 := game.NewCharacter()
	player1.Name = "Player1"
	player1.Room = room
	player1.Metatype = "Human"

	player2 := game.NewCharacter()
	player2.Name = "Player2"
	player2.Room = room
	player2.Metatype = "Ork"

	char.Room.AddCharacter(player1)
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

	char.Room.AddMob(game.EntityMgr.GetMob("orc"))
	char.Room.AddMob(game.EntityMgr.GetMob("goblin"))

	// item1 := game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock")
	// item2 := game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock")
	// item3 := game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key")

	char.Room.Inventory.AddItem(game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock"))
	char.Room.Inventory.AddItem(game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock"))
	char.Room.Inventory.AddItem(game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key"))

	game.GameTimeMgr.Minutes = 90
	char.Prompt = "{{time}} {{date}} {{>}}::white"

	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
    Nulla convallis egestas rhoncus. Donec facilisis fermentum sem, ac viverra ante luctus vel.
    Phasellus ultrices nulla quis nibh. Quisque a lectus. Donec consectetuer ligula vulputate sem tristique cursus.`

	// Example 2: Use custom options.
	customOptions := &game.WrapOptions{
		BorderType:    "rounded", // You can define your own border here.
		TextWidth:     50,
		PaddingTop:    1,
		PaddingBottom: 1,
		PaddingLeft:   1,
		PaddingRight:  1,
		BorderColor:   "red|bold",
		Alignment:     "center",
	}

	var output strings.Builder
	output.WriteString(game.WrapTextInBorder(text, nil))
	output.WriteString(game.CRLF)
	output.WriteString(game.WrapTextInBorder(text, customOptions))
	output.WriteString(game.CRLF)
	// output.WriteString(game.RenderRoom(acct, char, room))
	// output.WriteString(game.RenderCharacterTable(char))
	// output.WriteString(game.RenderPromptMenu("Main Menu", []string{"Enter Game", "Create Character", "Change Password", "Quit"}))
	// output.WriteString(game.RenderPrompt(char))
	cfmt.Print(output.String())
}

// func main() {
// 	text := `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
// Nulla convallis egestas rhoncus. Donec facilisis fermentum sem, ac viverra ante luctus vel.
// Phasellus ultrices nulla quis nibh. Quisque a lectus. Donec consectetuer ligula vulputate sem tristique cursus.`

// 	// Wrap the text in a border with an inner width of 40 characters.
// 	output := WrapTextInBorder(text, roundedBorder, 40)
// 	fmt.Println(output)
// }
