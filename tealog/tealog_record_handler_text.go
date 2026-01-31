package tealog

type recordHandlerText struct{}

func NewRecordHandlerText() recordHandlerText {
	return recordHandlerText{}
}

var _ recordHandler = recordHandlerText{}

func (recordHandlerText) HandleRecord(b *Buffer, r Record) {
	const layout = "2006-01-02T15:04:05.000Z07:00"
	b.BufferWriteString(r.Time.Format(layout))
	b.BufferWriteByte(' ')
	b.BufferWriteString(r.Level)
	b.BufferWriteByte(' ')
	b.BufferWriteString(r.Msg)
	b.BufferWriteByte('\n')
}
