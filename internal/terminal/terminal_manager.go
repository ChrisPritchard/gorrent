package terminal

import (
	"fmt"
	"strings"
)

const (
	// ANSI escape codes
	escape     = "\x1b"
	clearLine  = escape + "[2K"
	cursorUp   = escape + "[1A"
	cursorHide = escape + "[?25l"
	cursorShow = escape + "[?25h"
)

func ToggleCursor(show bool) {
	if show {
		fmt.Print(cursorShow)
	} else {
		fmt.Print(cursorHide)
	}
}

func Clear(line_count int) {
	fmt.Print(cursorHide)
	for range line_count {
		fmt.Print(cursorUp, clearLine)
	}
}

func Render(lines []string) {
	for _, l := range lines {
		fmt.Println(l)
	}
}

func ProgressBar(current, max, line_length int, suffix string) (string, error) {
	if current < 0 {
		return "", fmt.Errorf("current can not be less than 0")
	}
	if max < 0 {
		return "", fmt.Errorf("max can not be less than 0")
	}
	if current > max {
		return "", fmt.Errorf("current must be less than or equal to max")
	}

	required_suffix := len(suffix)
	if required_suffix != 0 {
		required_suffix++ // space gap
	}
	available := line_length - required_suffix
	if available < 3 { // '[ ]' is the smallest bar
		return "", fmt.Errorf("line_length is not sufficient to present the progress bar")
	}

	segments := available - 2
	progress_per_segment := max / segments
	current_progress := current / progress_per_segment

	var line strings.Builder
	line.WriteString("[")

	for range current_progress {
		line.WriteString("█")
	}
	for range segments - current_progress {
		line.WriteString(" ")
	}
	line.WriteString("]")

	if len(suffix) > 0 {
		line.WriteString(" " + suffix)
	}

	return line.String(), nil
}
