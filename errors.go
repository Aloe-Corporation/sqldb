package sqldb

// ParamError is return when a param is at it default value.
// A call to Error() return the invalid param name.
type ParamError struct {
	name string
}

func (e ParamError) Error() string {
	return e.name + " param is at it default value"
}
