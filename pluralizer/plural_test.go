// pluralizer, translated from https://www.last-outpost.com/LO/pubcode/
// Ported to Go from pluralizer.c by Dr Pogi (drpogi@icloud.com)
package pluralizer_test

import (
	"fmt"
	"testing"

	"github.com/Jasrags/NewMUD/pluralizer"
	"github.com/stretchr/testify/assert"
)

func TestIntToWords(t *testing.T) {
	tests := []struct {
		num      int
		expected string
	}{
		{0, "zero"},
		{1, "one"},
		{-1, "minus one"},
		{7, "seven"},
		{10, "ten"},
		{12, "twelve"},
		{17, "seventeen"},
		{-29, "minus twenty-nine"},
		{99, "ninety-nine"},
		{100, "100"},
		{1234, "1234"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("IntToWords(%d)", test.num), func(t *testing.T) {
			assert.Equal(t, pluralizer.IntToWords(test.num), test.expected)
		})
	}
}

func TestPluralizeNoun(t *testing.T) {
	tests := []struct {
		singular string
		count    int
		expected string
	}{
		{"", 2, ""},
		{"foo", 1, "foo"},
		{"fish", 1, "fish"},
		{"fish", 2, "fish"},
		{"person", 1, "person"},
		{"person", 2, "people"},
		{"boss", 7, "bosses"},
		{"brush", 7, "brushes"},
		{"punch", 2, "punches"},
		{"fox", 2, "foxes"},
		{"avocado", 2, "avocadoes"},
		{"fez", 2, "fezes"},
		{"life", 3, "lives"},
		{"loaf", 2, "loaves"},
		{"bluff", 2, "bluffs"},
		{"entity", 2, "entities"},
		{"tray", 2, "trays"},
		{"virus", 2, "viri"},
		{"terminus", 2, "termini"},
		{"ellipsis", 2, "ellipses"},
		{"onion", 2, "onions"},
		{"gas", 2, "gases"},
		{"Excalibur", 2, "Excaliburs"},
		{"GLAMDRING", 2, "GLAMDRINGs"},
		{"apex", 2, "apexes"},
		{"beau", 2, "beaux"},
		{"quiz", 2, "quizzes"},
		{"elf", 2, "elves"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("PluralizeNoun(%s, %d)", test.singular, test.count), func(t *testing.T) {
			assert.Equal(t, pluralizer.PluralizeNoun(test.singular, test.count), test.expected)
		})
	}
}

func TestPluralizeNounPhrase(t *testing.T) {
	tests := []struct {
		phrase   string
		count    int
		expected string
	}{
		{"a short sword", 0, "zero short swords"},
		{"a short sword", 1, "one short sword"},
		{"a short sword", 3, "three short swords"},
		{"a bag of holding", 0, "zero bags of holding"},
		{"a bag of holding", 1, "one bag of holding"},
		{"a bag of holding", 11, "eleven bags of holding"},
		{"one loaf of crusty bread", 0, "zero loaves of crusty bread"},
		{"one loaf of crusty bread", 1, "one loaf of crusty bread"},
		{"one loaf of crusty bread", 42, "forty-two loaves of crusty bread"},
		{"an Excalibur", 1, "one Excalibur"},
		{"an Excalibur", 7, "seven Excaliburs"},
		{"THE GLAMDRING", 1, "one GLAMDRING"},
		{"THE GLAMDRING", 2, "two GLAMDRINGs"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("PluralizeNounPhrase(%s, %d)", test.phrase, test.count), func(t *testing.T) {
			assert.Equal(t, pluralizer.PluralizeNounPhrase(test.phrase, test.count), test.expected)
		})
	}
}

func TestPluralizeVerb(t *testing.T) {
	tests := []struct {
		verb     string
		expected string
	}{
		{"", ""},
		{"has", "have"},
		{"isn't", "aren't"},
		{"relaxes", "relax"},
		{"blesses", "bless"},
		{"bashes", "bash"},
		{"wrenches", "wrench"},
		{"fuzzes", "fuzz"},
		{"bloodies", "bloody"},
		{"parries", "parry"},
		{"assays", "assay"},
		{"moans", "moan"},
		{"test", "test"},
		{"flies", "fly"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("PluralizeVerb(%s)", test.verb), func(t *testing.T) {
			assert.Equal(t, pluralizer.PluralizeVerb(test.verb), test.expected)
		})
	}
}
