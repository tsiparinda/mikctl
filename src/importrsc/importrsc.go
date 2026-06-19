package importrsc

import (
	"mikctl/src/config"
	"mikctl/src/db"
	"mikctl/src/models"
	"mikctl/src/ssh"
	"path/filepath"
	"sync"
	"time"
)

func ImportRsc(
	routers []models.Router,
	fileName string,
	workers int,
) error {
	jobs := make(chan models.Router)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range jobs {
				result := ImportRscRouter(
					r,
					fileName,
				)
				result.FinishedAt = time.Now()
				if result.Success {
					config.Verbosef("[OK ] %s\n", r.Name)
				} else {
					config.Verbosef("[ERR] %s %s\n", r.Name, result.Error)
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

func ImportRscRouter(r models.Router, fileName string) models.ScriptResult {
	tag := config.GetScriptTag(fileName)
	if tag == "" {
		tag = filepath.Base(fileName)
	}

	result := models.ScriptResult{
		Router:    r,
		Success:   true,
		StartedAt: time.Now(),
		Script:    fileName,
		Tag:       tag,
	}

	err := ssh.CopyFile(
		r.IP,
		fileName,
		"mikctl_import.rsc",
	)
	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}

	out, err := ssh.RunSSH(
		r.IP,
		"/import file-name=mikctl_import.rsc",
	)

	result.Output = out

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		_, _ = ssh.RunSSH(
			r.IP,
			"/file remove mikctl_import.rsc",
		)
	}

	return result
}
