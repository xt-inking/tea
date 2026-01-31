package tealog

type recordHandler interface {
	HandleRecord(b *Buffer, r Record)
}
