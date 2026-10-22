package bot

import (
	"strings"
)

// fetchValue returns the part of a "key=value" callback payload after the
// first "=", or "" when there is none.
func fetchValue(payload string) string {
	_, value, _ := strings.Cut(payload, "=")
	return value
}
