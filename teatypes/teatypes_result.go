package teatypes

type Result[Value any] struct {
	Value Value
	Error error
}
