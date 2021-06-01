package models

type Tutor struct {
	Id     int
	Title  string
	Author int
}

func (t *Tutor) TableName() string {
	return "tutorials"
}

type Tutorial struct {
	Id       int
	Title    string
	Author   Author
	Comments []Comment
}

type Author struct {
	Name      string
	Tutorials []int
}

type Comment struct {
	Body string
}

func Populate() []Tutorial {
	author := &Author{
		Name:      "Elliot Forbes",
		Tutorials: []int{1},
	}
	tutorial := Tutorial{
		Id:     1,
		Title:  "GO GraphQL Tutorial",
		Author: *author,
		Comments: []Comment{
			{Body: "First Comment"},
		},
	}
	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)
	return tutorials
}
