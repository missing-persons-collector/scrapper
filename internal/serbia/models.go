package serbia

import (
	"golang.org/x/net/html"
	"missingPersons/common"
)

type nodeOrError struct {
	container *html.Node
	node      *html.Node
	error     error
}

type personOrError struct {
	person common.RawPerson
	error  error
}

func (n nodeOrError) Data() interface{} {
	return n.node
}

func (n nodeOrError) Error() error {
	return n.error
}

func (n personOrError) Data() interface{} {
	return n.person
}

func (n personOrError) Error() error {
	return n.error
}
