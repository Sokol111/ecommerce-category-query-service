package categorylist

type CategoryListViewDTO struct {
	Categories []CategoryDTO
}

type CategoryDTO struct {
	ID   string
	Name string
}
