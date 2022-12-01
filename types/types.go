package types

type ReceiverData struct {
	ID    string
	Key   string
	Value string
}

type DataOrError struct {
	Data  []ReceiverData
	Error error
}

type Kernel interface {
	Run()
	Element() string
	BaseURL() string
}

type Signal interface {
	Error() chan error
	Data() chan [][]ReceiverData
}

type CollectedPage struct {
	Page  int
	Data  [][]ReceiverData
	Error error
}

type Information struct {
	UpdatedCount int
	CreatedCount int
	DeletedCount int
}
