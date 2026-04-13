package models

type Comment struct {
	ID       uint
	Body     string
	Approved bool
	PostID   uint
}
