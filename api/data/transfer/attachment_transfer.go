package transfer

import "noda/failure"

/* Transfers an attachment creation request.  */
type AttachmentCreation struct {
	FileName string `json:"file_name" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
}

func (a *AttachmentCreation) Validate() *failure.Aggregation {
	return validate(a)
}
