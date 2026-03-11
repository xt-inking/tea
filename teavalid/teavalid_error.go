package teavalid

type Error interface {
	error
	Key(k string)
}

func NewError(k string, text string) Error {
	return &errorValid{k, text}
}

type errorValid struct {
	k    string
	text string
}

func (e *errorValid) Error() string {
	return "`" + e.k + "` " + e.text
}

func (e *errorValid) Key(k string) {
	e.k = k + "." + e.k
}
