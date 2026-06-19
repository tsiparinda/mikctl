package models

type Neighbor struct {
	ParentRouterID int
	IP             string
	MAC            string
	Name           string
	Platform       string
	Model          string
	ROSVersion     string
}
