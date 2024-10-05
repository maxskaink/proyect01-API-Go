package dto

// Error contains the information if some process go wrong
// it is the answere if something doesnt work
type Error struct {
	Message   string `json:"message"`
	Status    int    `json:"status"`
	TypeError string `json:"typeError"`
}
