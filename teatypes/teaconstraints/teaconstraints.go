package teaconstraints

type Integer interface {
	Signed | Unsigned
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type String interface {
	~string
}
