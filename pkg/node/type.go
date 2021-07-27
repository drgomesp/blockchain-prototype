package node

// Type of the node.
type Type int

const (
	TypeFull = Type(iota)
	TypeProducer
	TypeValidator
)

func (t Type) String() string {
	switch t {
	case TypeFull:
		return "full"
	case TypeProducer:
		return "producer"
	case TypeValidator:
		return "validator"
	}

	return ""
}
