package tree

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	Lines    int
	Children []*Node
}
