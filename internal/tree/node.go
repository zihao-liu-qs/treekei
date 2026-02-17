package tree

type Node struct {
	Name     string
	Path     string
	IsDir    bool
	IsBinary bool
	Lines    int
	Children []*Node
}
