package types

type Result[T any] struct {
	Page      uint64 `json:"page"`
	RPP       uint64 `json:"rpp"`
	Retrieved uint64 `json:"retrieved"`
	Payload   *[]*T  `json:"payload"`
}
