package main

import (
	"fmt"
	"strings"

	"github.com/Jasrags/NewMUD/internal/game"
)

var (
	char = &game.Character{
		UserID: "test_user",
		Role:   game.CharacterRolePlayer,
		GameEntity: game.GameEntity{
			Name:            "Street Samurai",
			Title:           "Street Samurai",
			ID:              "ID",
			Metatype:        "Ork",
			Age:             25,
			Sex:             "Male",
			Height:          180,
			Weight:          80,
			Ethnicity:       "White",
			StreetCred:      2,
			Notoriety:       2,
			PublicAwareness: 2,
			Karma:           2,
			TotalKarma:      5,
			Description:     "A street samurai character",
			Attributes: game.Attributes{
				Body:            game.Attribute[int]{Name: "Body", Base: 7},
				Agility:         game.Attribute[int]{Name: "Agility", Base: 6},
				Reaction:        game.Attribute[int]{Name: "Reaction", Base: 5, Delta: 2},
				Strength:        game.Attribute[int]{Name: "Strength", Base: 5},
				Willpower:       game.Attribute[int]{Name: "Willpower", Base: 3},
				Logic:           game.Attribute[int]{Name: "Logic", Base: 2},
				Intuition:       game.Attribute[int]{Name: "Intuition", Base: 3},
				Charisma:        game.Attribute[int]{Name: "Charisma", Base: 2},
				Essence:         game.Attribute[float64]{Name: "Essence", Base: 6.0, Delta: -5.12},
				Magic:           game.Attribute[int]{Name: "Magic"},
				Resonance:       game.Attribute[int]{Name: "Resonance"},
				Initiative:      game.Attribute[int]{Name: "Initiative"},
				InitiativeDice:  game.Attribute[int]{Name: "Initiative Dice", Base: 1},
				Composure:       game.Attribute[int]{Name: "Composure"},
				JudgeIntentions: game.Attribute[int]{Name: "Judge Intentions"},
				Memory:          game.Attribute[int]{Name: "Memory"},
				Lift:            game.Attribute[int]{Name: "Lift"},
				Carry:           game.Attribute[int]{Name: "Carry"},
				Walk:            game.Attribute[int]{Name: "Walk"},
				Run:             game.Attribute[int]{Name: "Run"},
				Swim:            game.Attribute[int]{Name: "Swim"},
			},
			PhysicalDamage: game.PhysicalDamage{
				Current:  0,
				Max:      10,
				Overflow: 0,
			},
			StunDamage: game.StunDamage{
				Current: 0,
				Max:     10,
			},
			Edge: game.Edge{
				Max:       5,
				Available: 5,
			},
			Equipment: map[string]*game.Item{},
		},
	}
)

func main() {
	var output strings.Builder
	output.WriteString(game.RenderCharacterTable(char))
	fmt.Print(output.String())
}
