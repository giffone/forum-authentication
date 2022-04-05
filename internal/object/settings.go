package object

type Settings struct {
	AllPost bool                     // all posts (home handler) or one post (post handler) to show
	Key     map[string][]interface{} // keys for queries db
	Refers  bool                     // make/or not additional refers to struct if need it
}

func NewSettings() *Settings {
	return &Settings{
		Key: make(map[string][]interface{}, 1),
	}
}

// ClearKey checks if Key already have data.
// Need to clear map for new data, not adding another.
// Only one key must be!
func (s *Settings) ClearKey() *Settings {
	if len(s.Key) > 0 {
		s.Key = make(map[string][]interface{}, 1)
	}
	return s
}
