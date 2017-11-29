package trumptracker

type Point struct {
	TimeSeconds int64 `json:"x"`
	Value       int64 `json:"y"`
}

type Series struct {
	Points []Point `json:"values"`
	Name   string  `json:"key"`
}
