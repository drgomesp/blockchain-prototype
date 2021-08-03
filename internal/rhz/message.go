package rhz

type Message interface {
	Decode(val interface{}) error
}
