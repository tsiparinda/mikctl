package router

import (
	"fmt"
	"mikctl/src/models"
	"sort"
)

func ListRouters(routers []models.Router) error {

	fmt.Printf(
		"%-5s %-25s %-15s %-6s %-8s %-8s\n",
		"ID",
		"NAME",
		"IP",
		"MAIN",
		"PARENT",
		"ROS",
	)
	sort.Slice(routers, func(i, j int) bool {
		return routers[i].Name < routers[j].Name
	})

	for _, r := range routers {

		fmt.Printf(
			"%-5d %-25s %-15s %-6t %-5d %-6s \n",
			r.ID,
			r.Name,
			r.IP,
			r.Main,
			r.ParentID,
			r.ROSVersion,
		)
	}

	fmt.Printf("\nRouters: %d\n", len(routers))

	return nil
}
