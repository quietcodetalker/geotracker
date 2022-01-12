package filter

type Operator uint

const (
	Eq Operator = iota + 1
	Ne
	Lt
	Lte
	Gt
	Gte
)

type Filter struct {
	Field    string
	Operator Operator
	Value    interface{}
}
