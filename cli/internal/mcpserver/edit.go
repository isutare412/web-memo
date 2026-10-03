package mcpserver

import (
	"errors"
	"fmt"
	"strings"
)

// applyEdit replaces oldStr with newStr in content. oldStr must match exactly
// once unless replaceAll is set, so an edit never lands on the wrong spot.
func applyEdit(content, oldStr, newStr string, replaceAll bool) (string, error) {
	if oldStr == "" {
		return "", errors.New("old_string must not be empty")
	}
	n := strings.Count(content, oldStr)
	switch {
	case n == 0:
		return "", errors.New("old_string not found in memo content (0 matches)")
	case n > 1 && !replaceAll:
		return "", fmt.Errorf("old_string matches %d places in memo content; "+
			"add surrounding context to make it unique or set replace_all", n)
	}
	return strings.ReplaceAll(content, oldStr, newStr), nil
}
