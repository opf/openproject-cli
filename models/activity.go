package models

type Activity struct {
	Id        uint64
	Comment   string
	Details   []*string
	Version   uint64
	CreatedAt string
	UpdatedAt string
	UserId    uint64
}
