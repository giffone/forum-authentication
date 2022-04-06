package dto

type CheckID struct {
	ID  []string
	Who string
}

func NewCheckID(who string, id []string) *CheckID {
	return &CheckID{
		ID:  id,
		Who: who,
	}
}
