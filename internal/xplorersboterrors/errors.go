package xplorersboterrors

type Err int

const (
	InternalServerError Err = 500
)

func (e Err) String() string {
	switch e {
	case InternalServerError:
		return "INTERNAL_SERVER_ERROR"
	default:
		return ""
	}
}
