package router

import (
	"fmt"
	"mikctl/src/models"
)

func CountRouters(routers []models.Router) error {
	fmt.Printf("Routers: %d\n", len(routers))
	return nil
}
