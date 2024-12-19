package mud_test

import (
	"fmt"
	"testing"

	"github.com/Jasrags/NewMUD/mud"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		text     string
		width    int
		expected string
	}{
		{"", 10, ""},
		{"short text", 10, "short text"},
		{"this is a longer text that needs to be wrapped", 10, "this is a\nlonger text\nthat needs\nto be\nwrapped"},
		{"this is a longer text that needs to be wrapped", 20, "this is a longer text\nthat needs to be\nwrapped"},
		{"thisisaverylongwordthatcannotbewrapped", 10, "thisisaverylongwordthatcannotbewrapped"},
		{"this is a text with multiple spaces", 10, "this is a\ntext with\nmultiple\nspaces"},
		{"this is a text with\nnew lines", 10, "this is a\ntext with\nnew lines"},
	}

	for _, test := range tests {
		t.Run(test.text, func(t *testing.T) {
			result := mud.WrapText(test.text, test.width)
			if result != test.expected {
				t.Errorf("WrapText(%q, %d) = %q; want %q", test.text, test.width, result, test.expected)
			}
		})
	}
}
func TestParseEntityRef(t *testing.T) {
	tests := []struct {
		entityRef    string
		expectedArea string
		expectedID   string
	}{
		{"area1:id1", "area1", "id1"},
		{"area2:id2", "area2", "id2"},
		{"invalidref", "", ""},
		{"area3:", "area3", ""},
		{":id3", "", "id3"},
		{"", "", ""},
	}

	for _, test := range tests {
		t.Run(test.entityRef, func(t *testing.T) {
			area, id := mud.ParseEntityRef(test.entityRef)
			if area != test.expectedArea || id != test.expectedID {
				t.Errorf("ParseEntityRef(%q) = (%q, %q); want (%q, %q)", test.entityRef, area, id, test.expectedArea, test.expectedID)
			}
		})
	}
}
func TestCreateEntityRef(t *testing.T) {
	tests := []struct {
		area     string
		id       string
		expected string
	}{
		{"area1", "id1", "area1:id1"},
		{"area2", "id2", "area2:id2"},
		{"", "id3", ":id3"},
		{"area4", "", "area4:"},
		{"", "", ":"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s:%s", test.area, test.id), func(t *testing.T) {
			result := mud.CreateEntityRef(test.area, test.id)
			if result != test.expected {
				t.Errorf("CreateEntityRef(%q, %q) = %q; want %q", test.area, test.id, result, test.expected)
			}
		})
	}
}
