package dto

type Error struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	TypeError string `json:"typeError"`
}
