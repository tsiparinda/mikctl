package models

type Router struct {
	ID         int
	Name       string
	IP         string
	SSHUser    string
	Site       string
	ROSVersion string
	Model      string
	Serial     string
	DeviceID   int
	Main       bool
	ParentID   int
	PasswordID int
	MAC        string
	Platform   string
}
