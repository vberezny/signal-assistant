package assistant

type alert struct {
	name string
}

func NewAlert() *alert {
	return &alert{}
}
