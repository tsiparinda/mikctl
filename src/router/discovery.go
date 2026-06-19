package router

import (
	"errors"
	"fmt"
	"mikctl/src/config"
	"mikctl/src/db"
	"mikctl/src/models"
	"mikctl/src/script"
	"strings"
	"sync"
)

func DiscoveryRouters(routers []models.Router) error {
	fmt.Printf(
		"Discovery %d main routers\n",
		len(routers),
	)
	CheckNeighbors(routers)
	return nil
}

func CheckNeighbors(routers []models.Router) {
	jobs := make(chan models.Router)

	var wg sync.WaitGroup

	for i := 0; i < config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for router := range jobs {
				err := CheckNeighbor(router)
				if err != nil {
					fmt.Printf("[ERR] %s : %v\n", router.Name, err)
				} else {
					fmt.Printf("[OK ] %s\n", router.Name)
				}
			}
		}()
	}
	for _, router := range routers {
		if router.Main {
			jobs <- router
		}
	}
	close(jobs)
	wg.Wait()

}

func CheckNeighbor(parent models.Router) error {
	if !parent.Main {
		return fmt.Errorf("router %s is not main", parent.Name)
	}
	fmt.Printf("[%s] discovery\n", parent.Name)

	//--------------------------------------------------
	// 1. Get neighbor's list
	//--------------------------------------------------

	neighbors, err := GetNeighbors(parent)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] found %d neighbors\n", parent.Name, len(neighbors))

	//--------------------------------------------------
	// 2. Render neighbors
	//--------------------------------------------------

	for _, n := range neighbors {
		if n.Platform != "MikroTik" {
			config.Verbosef("CheckNeighbor: [%s] skip: platform %s\n", n.IP, n.Platform)
			continue
		}
		if n.IP == "" {
			config.Verbosef("CheckNeighbor: mac [%s] skip: no IP, parent: %s\n", n.MAC, parent.Name)
			continue
		}
		if n.Name == "" {
			config.Verbosef("CheckNeighbor: [%s] skip: no Name %s\n", n.IP, n.Name)
			continue
		}
		config.Verbosef("[ParentID %s\n", n.ParentID)
		n.Main = false
		n.SSHUser = config.SSHUser
		n.PasswordID = config.PasswordID
		// Get info about neighbor
		err := GetRouterInfo(&n)
		if err != nil {
			fmt.Printf("[%s] %s : %v\n", parent.Name, n.IP, err)
			//continue
		}
		fMain, err := db.IsMainRouter(n.Name)
		if err != nil {
			fmt.Printf("CheckNeighbor: Error to check IsMainRouter %s \n", n.Name)
			continue
		}
		if fMain {
			n.ParentID = 0
			continue
		}
		// skip CHR and x86
		model := strings.TrimSpace(n.Model)
		if strings.Contains(model, "CHR") || strings.Contains(model, "x86") {
			config.Verbosef("[%s] skip %s - model [%s]\n", parent.Name, n.IP, model)
			continue
		}
		// skip no serial also
		serial := strings.TrimSpace(n.Serial)

		if serial == "" {
			config.Verbosef("[%s] skip adding device %s - serial empty\n", parent.Name, n.IP)
			//continue
		}
		// save info
		err = db.SaveRouterInfo(&n)
		if err != nil {
			fmt.Printf("[%s] SaveRouterInfo device error: %v\n", parent.Name, err)
			continue
		}
	}

	return nil
}

func GetNeighbors(parent models.Router) ([]models.Router, error) {
	cmd := `:foreach n in=[/ip neighbor find] do={ :put ([/ip neighbor get $n address] . ";" . [/ip neighbor get $n mac-address] . ";" . [/ip neighbor get $n identity] . ";" . [/ip neighbor get $n platform] . ";" . [/ip neighbor get $n board] . ";" . [/ip neighbor get $n version])}`
	config.Verbosef("CMD: \n%q\n", cmd)
	out := script.RunRouterScript(
		parent, []string{cmd},
	)
	if out.Error != "" {
		return nil, errors.New(out.Error)
	}
	var result []models.Router
	out.Output = strings.ReplaceAll(out.Output, "\r", "")
	lines := strings.Split(out.Output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ";")
		if len(fields) < 6 {
			continue
		}
		n := models.Router{
			IP:         strings.TrimSpace(fields[0]),
			MAC:        strings.TrimSpace(fields[1]),
			Name:       strings.TrimSpace(fields[2]),
			Platform:   strings.TrimSpace(fields[3]),
			Model:      strings.TrimSpace(fields[4]),
			ROSVersion: strings.TrimSpace(fields[5]),
			ParentID:   parent.ID,
		}
		result = append(result, n)
	}
	config.Verbosef("[%s] neighbors found: %d\n", parent.Name, len(result))
	return result, nil
}
