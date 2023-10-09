package types

type Result[T any] struct {
	Page      int64 `json:"page"`
	RPP       int64 `json:"rpp"`
	Retrieved int64 `json:"retrieved"`
	Payload   []*T  `json:"payload"`
}
