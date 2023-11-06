package transfer

/* Transfers an attachment creation request.  */
type AttachmentCreation struct {
	FileName string `json:"file_name" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
}

func (a *AttachmentCreation) Validate() error {
	return validate(a)
}
