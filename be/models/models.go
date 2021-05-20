package models

type Url struct {
	Uri   string
	Level int
}
type Results struct {
	BrokenLinks []string
	AllLinks    []string
	MainURL     string
}
