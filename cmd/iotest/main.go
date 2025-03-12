package main

import (
	"fmt"
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

	// acct := game.AccountMgr.GetByUsername("Jasrags")
	char := game.CharacterMgr.GetCharacterByName("fred")
	// char.Role = game.CharacterRolePlayer
	char.Role = game.CharacterRoleAdmin
	char.Equipment.Slots["body"] = game.EntityMgr.CreateItemInstanceFromBlueprintID("synth_leather_jacket")

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

	// mob1 := game.EntityMgr.GetMob("ork_thug_lieutenant")

	// mob1 := game.EntityMgr.CreateMobInstanceFromBlueprintID("ork_thug_lieutenant")
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

	char.Room.AddMobInstance(game.EntityMgr.CreateMobInstanceFromBlueprintID("ork_thug_basic"))
	char.Room.AddMobInstance(game.EntityMgr.CreateMobInstanceFromBlueprintID("ork_thug_lieutenant"))

	// item1 := game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock")
	// item2 := game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock")
	// item3 := game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key")

	char.Room.Inventory.Add(game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock"))
	char.Room.Inventory.Add(game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock"))
	char.Room.Inventory.Add(game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key"))

	// item1 := game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock")

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

	// qualityAmbidextrousBP := game.EntityMgr.GetQualityBlueprint("ambidextrous")
	// qualityAmbidextrous := &game.Quality{
	// 	BlueprintID: qualityAmbidextrousBP.ID,
	// 	Blueprint:   qualityAmbidextrousBP,
	// 	Rating:      1,
	// }

	// qualityAmbidextrous := game.EntityMgr.CreateQualityFromBlueprintID("ambidextrous", 1)
	// qualityAllergyBP := game.EntityMgr.GetQualityBlueprint("allergy")
	// qualityAllergy := game.EntityMgr.CreateQualityFromBlueprintID("allergy", 0)
	// BlueprintID: qualityAllergyBP.ID,
	// Blueprint:   qualityAllergyBP,
	// }
	// skillPistols := game.EntityMgr.CreateSkillInstanceFromBlueprintID("pistols", 2, "colt_45")

	inv := game.NewInventory()
	inv.Add(game.EntityMgr.CreateItemInstanceFromBlueprintID("small_rock"))
	inv.Add(game.EntityMgr.CreateItemInstanceFromBlueprintID("jagged_rock"))
	inv.Add(game.EntityMgr.CreateItemInstanceFromBlueprintID("test_key"))

	output := termenv.NewOutput(os.Stdout)
	output.ClearScreen()
	// output.DisableMouse()
	// output.EnableMouse()
	// output.AltScreen()
	// defer output.ExitAltScreen()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Body: %d AV: %d", char.GetBody(), char.GetArmorValue()))
	// sb.WriteString(inv.FormatTable())
	// sb.WriteString(item1.FormatListItem() + game.CRLF)
	// sb.WriteString(item1.FormatDetailed() + game.CRLF)
	// sb.WriteString(skillPistols.FormatListItem() + game.CRLF)
	// sb.WriteString(qualityAmbidextrous.FormatListItem() + game.CRLF)
	// sb.WriteString(qualityAllergy.FormatListItem() + game.CRLF)
	sb.WriteString(game.CRLF)
	// sb.WriteString(skillPistols.FormatDetailed() + game.CRLF)
	// sb.WriteString(qualityAmbidextrous.FormatDetailed() + game.CRLF)
	// sb.WriteString(qualityAllergy.FormatDetailed() + game.CRLF)

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
	// output.WriteString(game.RenderRoom(acct, char, room))
	// output.WriteString(game.RenderMobTable(mob1))
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
