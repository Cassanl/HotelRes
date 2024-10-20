package types

type Filter map[string]interface{}

const (
	UserKey int = iota
)

type ResponseKind string

const (
	ErrorResp ResponseKind = "Error"
	OkResp    ResponseKind = "OK"
)

type GenericResponse struct {
	Kind ResponseKind
	Msg  string
}
