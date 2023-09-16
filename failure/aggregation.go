package failure

type Aggregation struct {
	errors []string
}

func NewAggregation() *Aggregation {
	return &Aggregation{}
}

func (a *Aggregation) Error() string {
	str := ""
	for _, err := range a.errors {
		str += err + "\n"
	}
	return str
}

func (a *Aggregation) Append(err error) {
	a.errors = append(a.errors, err.Error())
}

func (a *Aggregation) Dump() []string {
	return a.errors
}

func (a *Aggregation) Has() bool {
	return len(a.errors) > 0
}
