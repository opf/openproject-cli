package projects

func (dto *ProjectCollectionDto) convert() []*Project {
	var projects = make([]*Project, len(dto.Embedded.Elements))

	for idx, p := range dto.Embedded.Elements {
		projects[idx] = p.convert()
	}

	return projects
}

func (dto *ProjectDto) convert() *Project {
	return &Project{
		Id:   dto.Id,
		Name: dto.Name,
	}
}
