package tealog

type handler interface {
	Handle(b *Buffer, r Record)
}
