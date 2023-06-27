package dtos

type LinkDto struct {
	Href   string `json:"href,omitempty"`
	Title  string `json:"title,omitempty"`
	Method string `json:"method,omitempty"`
}
