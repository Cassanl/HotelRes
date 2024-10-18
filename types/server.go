package types

type Filter map[string]interface{}

const (
	UserKey int = iota
)

type ResponseKind string

const (
	ErrorResp ResponseKind = "Error"
)

type GenericResponse struct {
	Kind   ResponseKind
	Status int
}

// type Constraint func(any) error

// type Validator struct {
// 	model       any
// 	Constraints map[string]Constraint
// }

// func (v *Validator) Validate() map[string]error {
// 	errs := map[string]error{}
// 	for key, constraint := range v.Constraints {
// 		if err := constraint(v.model); err != nil {
// 			errs[key] = err
// 		}
// 	}
// 	return errs
// }
