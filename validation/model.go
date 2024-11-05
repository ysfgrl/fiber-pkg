package validation

type Error struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	Tag   string      `json:"tag"`
	Param string      `json:"param"`
}
