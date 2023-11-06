package transfer

/* Transfers a step creation request.  */
type StepCreation struct {
	Description string `json:"description" validate:"required"`
}

func (s *StepCreation) Validate() error {
	return validate(s)
}

/* Transfers a step update request.  */
type StepUpdate struct {
	Description string `json:"description"`
}
