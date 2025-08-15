package bot

import (
	"strings"
)

func fetchValue(payload string) string {
	return strings.Split(payload, "=")[1]
}
