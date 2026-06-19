package router

import (
	"fmt"
	"mikctl/src/config"
	"mikctl/src/db"
	"mikctl/src/models"
	"mikctl/src/ssh"
	"strings"
	"sync"
)

func UpdateRouters(routers []models.Router) error {
	fmt.Printf("Updating %d routers\n", len(routers))
	RunWorkers(
		routers,
	)
	return nil
}

func RunWorkers(
	routers []models.Router,
) {
	jobs := make(chan *models.Router)
	var wg sync.WaitGroup
	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for router := range jobs {
				err := UpdateRouter(router)
				if err != nil {
					fmt.Printf(
						"[ERR] %s : %v\n",
						router.Name,
						err,
					)
				} else {
					fmt.Printf(
						"[OK ] %s\n",
						router.Name,
					)
				}
			}
		}()
	}
	for i := range routers {
		jobs <- &routers[i]
	}

	close(jobs)

	wg.Wait()
}

func UpdateRouter(router *models.Router) error {

	err := GetRouterInfo(router)
	if err != nil {
		return err
	}

	return db.SaveRouterInfo(router)
}

func GetRouterInfo(router *models.Router) error {
	get := func(cmd string) (string, error) {
		return ssh.RunSSH(router.IP, cmd)
	}
	cmd := `:put ([/system identity get name] . ";" . [/system resource get version] . ";" . [/system resource get board-name])`
	config.Verbosef("GetRouterInfo 1: %q\n", cmd)
	out, err := get(
		cmd,
	)
	if err != nil {
		return fmt.Errorf("1 %v: %s", err, out)
	}
	parts := strings.Split(out, ";")
	if len(parts) != 3 {
		return fmt.Errorf("unexpected response:  %s", out)
	}

	router.Name = strings.TrimSpace(parts[0])
	router.ROSVersion = strings.TrimSpace(parts[1])
	router.Model = strings.TrimSpace(parts[2])

	if !strings.Contains(router.Model, "CHR") &&
		!strings.Contains(router.Model, "x86") {
		serial, err := get(
			":put [/system routerboard get serial-number]",
		)
		if err != nil {
			return fmt.Errorf("%v: %s", err, serial)
		}
		router.Serial = strings.TrimSpace(serial)
	}
	return nil
}
