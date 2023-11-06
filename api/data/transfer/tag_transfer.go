package transfer

/* Transfers a tag creation request.  */
type TagCreation struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Color       string `json:"color" validate:"required"`
}

func (t *TagCreation) Validate() error {
	return validate(t)
}

/* Transfers a tag update request.  */
type TagUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}
