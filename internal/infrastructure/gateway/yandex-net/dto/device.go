package dto

type Devices struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
