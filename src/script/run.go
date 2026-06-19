package script

import (
	"bufio"
	"fmt"
	"mikctl/src/config"
	"mikctl/src/db"
	"mikctl/src/models"
	"mikctl/src/ssh"
	"os"
	"strings"
	"sync"
	"time"
)

func RunScript(routers []models.Router, scriptName string, workers int) error {

	fileName := fmt.Sprintf("%s", scriptName)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var commands []string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(
			scanner.Text(),
		)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		commands = append(
			commands,
			line,
		)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Printf("Running %s on %d routers\n", scriptName, len(routers))

	jobs := make(chan models.Router)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {

		wg.Add(1)

		go func() {

			defer wg.Done()

			for r := range jobs {
				result := RunRouterScript(
					r,
					commands,
				)

				result.FinishedAt = time.Now()

				if result.Success {
					fmt.Printf("[OK ] %s\n", result.Router.Name)
				} else {
					fmt.Printf("[ERR] %s %s\n", result.Router.Name, result.Error)
				}
				err := db.SaveScriptResult(result)
				if err != nil {
					config.Verbosef("[SaveScriptResult error] %s %v\n", r.Name, err)
				}
			}
		}()
	}

	for _, r := range routers {
		jobs <- r
	}

	close(jobs)

	wg.Wait()

	return nil
}

func RunRouterScript(r models.Router, commands []string) models.ScriptResult {
	tag := config.GetScriptTag(commands[0])
	result := models.ScriptResult{
		Router:    r,
		Success:   true,
		StartedAt: time.Now(),
		Script:    strings.Join(commands, "\r\n"),
		Tag:       tag,
	}

	var output strings.Builder

	for _, cmd := range commands {

		normcmd := NormalizeScript(cmd)
		config.Verbosef("CMD=[%s]\n", normcmd)

		out, err := ssh.RunSSH(
			r.IP,
			normcmd,
		)
		output.WriteString(out)
		output.WriteString("\n")

		if err != nil {
			result.Success = false
			result.Output = output.String()
			result.Error = fmt.Sprintf("%s : %v", cmd, err)
			return result
		}
	}
	result.Output = output.String()
	return result
}

func NormalizeScript(script string) string {
	script = strings.ReplaceAll(script, "\r", " ")
	script = strings.ReplaceAll(script, "\n", " ")
	script = strings.ReplaceAll(script, "\t", " ")
	return script
}
