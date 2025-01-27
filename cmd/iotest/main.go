package main

import (
	"fmt"
	"strings"

	"github.com/Jasrags/NewMUD/internal/game"
)

var ()

func main() {
	gs := game.NewGameServer()
	gs.SetupConfig()
	gs.SetupLogger()

	// game.EntityMgr = game.NewEntityManager()
	game.EntityMgr.LoadDataFiles()
	// game.AccountMgr = game.NewAccountManager()
	game.AccountMgr.LoadDataFiles()
	// game.CharacterMgr = game.NewCharacterManager()
	game.CharacterMgr.LoadDataFiles()

	acct := game.AccountMgr.GetByUsername("Jasrags")
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

	var output strings.Builder
	output.WriteString(game.RenderRoom(acct, char, room))
	fmt.Print(output.String())
}
