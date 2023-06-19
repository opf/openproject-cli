package projects

type Project struct {
	Id   int64
	Name string
}

type ProjectDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type elements struct {
	Elements []*ProjectDto `json:"elements"`
}

type ProjectCollectionDto struct {
	Embedded elements `json:"_embedded"`
	Type     string   `json:"_type"`
}
