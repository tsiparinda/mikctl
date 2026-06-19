package config

import (
	"bufio"
	"os"
	"strings"
)

// if the first line of script looks like "#backup2ftp", this line used as tag and this tag is stored in the t.scripts_runs in separate field
func GetScriptTag(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") {
			return strings.TrimSpace(strings.TrimPrefix(line, "#"))
		}
	}

	return ""
}
