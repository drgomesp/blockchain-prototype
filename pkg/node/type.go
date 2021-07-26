package node

// Type of the rhznode.
type Type int

const (
	TypeFull = Type(iota)
	TypeProducer
	TypeValidator
)
