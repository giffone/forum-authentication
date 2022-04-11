package dto

type CheckID struct {
	Atoi     bool        // if true - must use IDString
	IDString string      // will use if Atoi = true
	ID       interface{} // will use if Atoi = false
	Who      string
}

func NewCheckID(who string, id interface{}) *CheckID {
	return &CheckID{
		ID:  id,
		Who: who,
	}
}

func NewCheckIDAtoi(who string, id string) *CheckID {
	return &CheckID{
		IDString: id,
		Who:      who,
		Atoi:     true,
	}
}
