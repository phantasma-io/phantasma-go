package types

type Timestamp struct {
	Value uint32
}

func NewTimestamp(value uint32) *Timestamp {
	t := Timestamp{}
	t.Value = value

	return &t
}
