package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Jasrags/NewMUD/internal/game"
	"github.com/Jasrags/NewMUD/internal/game/templates"
	"github.com/Masterminds/sprig/v3"
	"github.com/gliderlabs/ssh"
	"github.com/muesli/termenv"
)

var ()

func main() {
	gs := game.NewGameServer()
	gs.SetupConfig()
	gs.SetupLogger()
	gs.OutputConfig()

	game.EntityMgr.LoadDataFiles()
	game.AccountMgr.LoadDataFiles()
	game.CharacterMgr.LoadDataFiles()

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

	rtd := RoomTemplateData{
		Title:       "The Void",
		ID:          "the_void",
		Description: "You don't think that you are not floating in nothing.",
		Characters: []Character{
			{Name: "Player1", Race: "Human"},
			{Name: "Player2", Race: "Elf"},
		},
		Mobs: []Mob{
			{Name: "Orc", Disposition: "Hostile"},
			{Name: "Goblin", Disposition: "Neutral"},
		},
		Items: []Item{
			{Name: "Jagged Rock"},
			{Name: "Test Key"},
			{Name: "Small Rock"},
		},
		Exits: []Exit{
			{Direction: "north", Description: "an open doorway leading to Limbo"},
		},
	}

	o := termenv.NewOutput(os.Stdout)
	// f := o.TemplateFuncs()
	tpl, err := template.New("room.tmpl").
		Funcs(o.TemplateFuncs()).
		Funcs(sprig.FuncMap()).
		Funcs(templates.FuncMap).
		ParseFiles("_data/templates/room.tmpl")
	if err != nil {
		panic(err)
	}

	// https://github.com/muesli/termenv/blob/master/examples/color-chart/color-chart.png

	// tpl.ParseFiles("_data/templates/room.tmpl")
	// tpl, err := LoadTemplate("room.tmpl")
	// if err != nil {
	// panic(err)
	// }

	var output strings.Builder
	if err := tpl.Execute(&output, rtd); err != nil {
		panic(err)
	}
	// output.WriteString(game.RenderRoom(acct, char, room))
	// output.WriteString(game.RenderCharacterTable(char))
	fmt.Print(output.String())
}

// TemplateCache stores parsed templates for reuse
var TemplateCache = map[string]*template.Template{}

// LoadTemplate dynamically loads and parses a template file.
func LoadTemplate(templateName string) (*template.Template, error) {
	// Check if the template is already cached
	if tmpl, ok := TemplateCache[templateName]; ok {
		return tmpl, nil
	}

	// Build the template path
	templatePath := filepath.Join("_data/templates", templateName)

	output := termenv.NewOutput(os.Stdout)

	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %w", err)
	}

	tmpl.
		Funcs(sprig.FuncMap()).
		Funcs(output.TemplateFuncs()).
		Funcs(templates.FuncMap)

	// tmpl.
	// 	Funcs(sprig.FuncMap()).
	// 	Funcs(output.TemplateFuncs()).
	// 	Funcs(templates.FuncMap)

	// Cache the template for future use
	TemplateCache[templateName] = tmpl
	return tmpl, nil
}

func RenderRoomDynamic(s ssh.Session, room RoomTemplateData, templateName string) error {
	// Dynamically load the template
	tmpl, err := LoadTemplate(templateName)
	if err != nil {
		slog.Error("Failed to load room template", slog.Any("error", err))
		io.WriteString(s, "An error occurred while displaying the room.\n")
		return err
	}

	// Render the room data using the template

	if err := tmpl.Execute(s, room); err != nil {
		slog.Error("Failed to render room template", slog.Any("error", err))
		io.WriteString(s, "An error occurred while displaying the room.\n")
		return err
	}

	return nil
}

type RoomTemplateData struct {
	Title       string
	ID          string
	Description string
	Characters  []Character
	Mobs        []Mob
	Items       []Item
	Exits       []Exit
}

type Character struct {
	Name string
	Race string
}

type Mob struct {
	Name        string
	Disposition string
}

type Item struct {
	Name string
}

type Exit struct {
	Direction   string
	Description string
}
