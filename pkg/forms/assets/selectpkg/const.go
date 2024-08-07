package selectpkg

const (
	DefaultBackOption   = "Back"
	DefaultSubmitOption = "Submit"
)

type QuitType int

const (
	QuitTypeSubmit QuitType = iota + 1
	QuitTypeBack
	QuitTypeInterrupt
)

func (t QuitType) IsEmpty() bool {
	return t == 0
}
