package types

type Map map[string]interface{}

const (
	UserKey int = iota
)

type ResponseKind string

type GenericResponse struct {
	Status int
	Msg    string
}
