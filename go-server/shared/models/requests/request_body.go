package requests

type RequestBody[T any, U any] struct {
	Payload T `json:"payload"`
	Options U `json:"options,omitempty"`
}
