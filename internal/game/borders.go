package game

import (
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/i582/cfmt/cmd/cfmt"
)

// ------------------------------
// Predefined Border Types
// ------------------------------

type (
	TextAlign  string
	BorderType string
	Border     struct {
		Top          string
		Bottom       string
		Left         string
		Right        string
		TopLeft      string
		TopRight     string
		BottomLeft   string
		BottomRight  string
		MiddleLeft   string
		MiddleRight  string
		Middle       string
		MiddleTop    string
		MiddleBottom string
	}
)

const (
	BorderTypeNone           BorderType = "none"
	BorderTypeNormal         BorderType = "normal"
	BorderTypeRounded        BorderType = "rounded"
	BorderTypeBlock          BorderType = "block"
	BorderTypeOuterHalfBlock BorderType = "outerHalfBlock"
	BorderTypeInnerHalfBlock BorderType = "innerHalfBlock"
	BorderTypeThick          BorderType = "thick"
	BorderTypeDouble         BorderType = "double"
	BorderTypeHidden         BorderType = "hidden"

	TextAlignLeft   TextAlign = "left"
	TextAlignCenter TextAlign = "center"
	TextAlignRight  TextAlign = "right"
)

var (
	noBorder = Border{}

	normalBorder = Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "┌",
		TopRight:     "┐",
		BottomLeft:   "└",
		BottomRight:  "┘",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}

	roundedBorder = Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}

	blockBorder = Border{
		Top:         "█",
		Bottom:      "█",
		Left:        "█",
		Right:       "█",
		TopLeft:     "█",
		TopRight:    "█",
		BottomLeft:  "█",
		BottomRight: "█",
	}

	outerHalfBlockBorder = Border{
		Top:         "▀",
		Bottom:      "▄",
		Left:        "▌",
		Right:       "▐",
		TopLeft:     "▛",
		TopRight:    "▜",
		BottomLeft:  "▙",
		BottomRight: "▟",
	}

	innerHalfBlockBorder = Border{
		Top:         "▄",
		Bottom:      "▀",
		Left:        "▐",
		Right:       "▌",
		TopLeft:     "▗",
		TopRight:    "▖",
		BottomLeft:  "▝",
		BottomRight: "▘",
	}

	thickBorder = Border{
		Top:          "━",
		Bottom:       "━",
		Left:         "┃",
		Right:        "┃",
		TopLeft:      "┏",
		TopRight:     "┓",
		BottomLeft:   "┗",
		BottomRight:  "┛",
		MiddleLeft:   "┣",
		MiddleRight:  "┫",
		Middle:       "╋",
		MiddleTop:    "┳",
		MiddleBottom: "┻",
	}

	doubleBorder = Border{
		Top:          "═",
		Bottom:       "═",
		Left:         "║",
		Right:        "║",
		TopLeft:      "╔",
		TopRight:     "╗",
		BottomLeft:   "╚",
		BottomRight:  "╝",
		MiddleLeft:   "╠",
		MiddleRight:  "╣",
		Middle:       "╬",
		MiddleTop:    "╦",
		MiddleBottom: "╩",
	}

	hiddenBorder = Border{
		Top:          " ",
		Bottom:       " ",
		Left:         " ",
		Right:        " ",
		TopLeft:      " ",
		TopRight:     " ",
		BottomLeft:   " ",
		BottomRight:  " ",
		MiddleLeft:   " ",
		MiddleRight:  " ",
		Middle:       " ",
		MiddleTop:    " ",
		MiddleBottom: " ",
	}
)

// getBorder returns the Border struct corresponding to the given borderType string.
// If borderType is unknown or empty, it defaults to "normal".
func getBorder(borderType BorderType) Border {
	switch borderType {
	case BorderTypeNone:
		return noBorder
	case BorderTypeRounded:
		return roundedBorder
	case BorderTypeBlock:
		return blockBorder
	case BorderTypeOuterHalfBlock:
		return outerHalfBlockBorder
	case BorderTypeInnerHalfBlock:
		return innerHalfBlockBorder
	case BorderTypeThick:
		return thickBorder
	case BorderTypeDouble:
		return doubleBorder
	case BorderTypeHidden:
		return hiddenBorder
	default:
		return normalBorder
	}
}

// ------------------------------
// Wrap Options Struct and Defaults
// ------------------------------

type WrapOptions struct {
	BorderType    BorderType
	TextWidth     int    // The width for wrapping the text (excluding padding).
	PaddingTop    int    // Number of empty lines to add at the top inside the border.
	PaddingBottom int    // Number of empty lines to add at the bottom inside the border.
	PaddingLeft   int    // Number of spaces to add to the left of each text line.
	PaddingRight  int    // Number of spaces to add to the right of each text line.
	BorderColor   string // The color (and style) for the border (e.g. "white|bold").
	Alignment     TextAlign
}

// Default options if none are provided.
var defaultWrapOptions = WrapOptions{
	BorderType:    BorderTypeNormal,
	TextWidth:     80,
	PaddingTop:    0,
	PaddingBottom: 0,
	PaddingLeft:   0,
	PaddingRight:  0,
	BorderColor:   "white|bold",
	Alignment:     TextAlignLeft,
}

// ------------------------------
// Helper Functions
// ------------------------------

// padRight pads a string with spaces on the right until it reaches the specified rune length.
// func padRight(str string, length int) string {
// 	runes := []rune(str)
// 	if len(runes) >= length {
// 		return str
// 	}
// 	return str + strings.Repeat(" ", length-len(runes))
// }

// alignText aligns the input text within a field of width `width` according to alignment:
// "left", "center", or "right". If the text is shorter than width, extra spaces are added.
func alignText(text string, width int, alignment TextAlign) string {
	runes := []rune(text)
	length := len(runes)
	if length >= width {
		return text
	}
	spaces := width - length
	slog.Info("spaces",
		slog.String("text", text),
		slog.Int("spaces", spaces))
	switch alignment {
	case TextAlignRight:
		return strings.Repeat(" ", spaces) + text
	case TextAlignCenter:
		leftSpaces := spaces / 2
		rightSpaces := spaces - leftSpaces
		return strings.Repeat(" ", leftSpaces) + text + strings.Repeat(" ", rightSpaces)
	default:
		return text + strings.Repeat(" ", spaces)
	}
}

// wrapText splits and wraps the input text into lines that do not exceed maxWidth runes.
// It wraps on spaces, and if a single word is too long, it breaks that word.
func wrapText(text string, maxWidth int) []string {
	var result []string
	paragraphs := strings.Split(text, "\n")
	for _, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			result = append(result, "")
			continue
		}

		currentLine := ""
		for _, word := range words {
			// Break word into chunks if too long.
			for utf8.RuneCountInString(word) > maxWidth {
				if currentLine != "" {
					result = append(result, currentLine)
					currentLine = ""
				}
				runes := []rune(word)
				result = append(result, string(runes[:maxWidth]))
				word = string(runes[maxWidth:])
			}

			// Add word to current line if it fits.
			if currentLine == "" {
				currentLine = word
			} else if utf8.RuneCountInString(currentLine)+1+utf8.RuneCountInString(word) <= maxWidth {
				currentLine += " " + word
			} else {
				result = append(result, currentLine)
				currentLine = word
			}
		}
		if currentLine != "" {
			result = append(result, currentLine)
		}
	}
	return result
}

// ------------------------------
// Main Function: WrapTextInBorder
// ------------------------------

// WrapTextInBorder wraps the provided text to a specified width and returns the text inside a border.
// The options parameter allows customization (border style via a string, text width, padding, border color, and alignment).
// If options is nil, default values are used.
func WrapTextInBorder(text string, options *WrapOptions) string {
	// Use default options if none provided.
	if options == nil {
		options = &defaultWrapOptions
	}
	// Ensure BorderColor and Alignment have defaults if empty.
	if options.BorderColor == "" {
		options.BorderColor = defaultWrapOptions.BorderColor
	}
	if options.Alignment == "" {
		options.Alignment = defaultWrapOptions.Alignment
	}
	if options.BorderType == "" {
		options.BorderType = defaultWrapOptions.BorderType
	}

	// Choose the border based on BorderType.
	chosenBorder := getBorder(options.BorderType)

	// Calculate effective inner width: text width plus left/right padding.
	effectiveWidth := options.PaddingLeft + options.TextWidth + options.PaddingRight
	slog.Info("effectiveWidth", slog.Int("effectiveWidth", effectiveWidth))
	// Wrap the text using the specified TextWidth.
	lines := wrapText(text, options.TextWidth)

	var result strings.Builder
	// Build top and bottom borders.
	topBorder := chosenBorder.TopLeft + strings.Repeat(chosenBorder.Top, effectiveWidth) + chosenBorder.TopRight
	bottomBorder := chosenBorder.BottomLeft + strings.Repeat(chosenBorder.Bottom, effectiveWidth) + chosenBorder.BottomRight

	// Apply the border color using cfmt.Sprintf.
	coloredTopBorder := cfmt.Sprintf("{{%s}}::%s", topBorder, options.BorderColor)
	coloredBottomBorder := cfmt.Sprintf("{{%s}}::%s", bottomBorder, options.BorderColor)

	result.WriteString(coloredTopBorder + "\n")

	// Build an empty line (for padding) with colored side borders.
	emptyLine := cfmt.Sprintf("{{%s}}::%s", chosenBorder.Left, options.BorderColor) +
		strings.Repeat(" ", effectiveWidth) +
		cfmt.Sprintf("{{%s}}::%s", chosenBorder.Right, options.BorderColor)

	// Add top padding lines.
	for i := 0; i < options.PaddingTop; i++ {
		result.WriteString(emptyLine + "\n")
	}

	// Process each wrapped text line.
	for _, line := range lines {
		// Align the text within the specified TextWidth.
		alignedText := alignText(line, options.TextWidth, options.Alignment)
		// Add left and right padding.
		paddedText := strings.Repeat(" ", options.PaddingLeft) + alignedText + strings.Repeat(" ", options.PaddingRight)
		// Surround with colored side borders.
		coloredLine := cfmt.Sprintf("{{%s}}::%s", chosenBorder.Left, options.BorderColor) +
			paddedText + cfmt.Sprintf("{{%s}}::%s", chosenBorder.Right, options.BorderColor)
		result.WriteString(coloredLine + "\n")
	}

	// Add bottom padding lines.
	for i := 0; i < options.PaddingBottom; i++ {
		result.WriteString(emptyLine + "\n")
	}

	result.WriteString(coloredBottomBorder)

	return result.String()
}
