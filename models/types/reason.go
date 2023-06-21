package types

type Reason string

const (
	Mentioned Reason = "mentioned"
	Assigned  Reason = "assigned"
	None      Reason = ""
)

func ParseReason(reason string) Reason {
	switch reason {
	case string(Mentioned):
		return Mentioned
	case string(Assigned):
		return Assigned
	default:
		return None
	}
}
