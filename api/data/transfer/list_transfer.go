package transfer

type ListCreation struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
