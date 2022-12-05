package worker

type DataOrError interface {
	Data() interface{}
	Error() error
}

type Worker[T DataOrError, F DataOrError] struct {
	workerNum      int
	producerStream chan T
	consumerStream chan F
	doneStreams    []chan bool
}

func NewWorker[T DataOrError, F DataOrError](workerNum int) *Worker[T, F] {
	doneStreams := make([]chan bool, workerNum)
	for i := 0; i < workerNum; i++ {
		doneStreams[i] = make(chan bool)
	}

	return &Worker[T, F]{
		workerNum:      workerNum,
		producerStream: make(chan T, workerNum),
		consumerStream: make(chan F),
		doneStreams:    doneStreams,
	}
}

func (w *Worker[T, F]) Produce(producer func(producerStream chan T, stop func())) {
	go func() {
		producer(w.producerStream, func() {
			close(w.producerStream)

			for _, doneStream := range w.doneStreams {
				<-doneStream
			}

			close(w.consumerStream)
		})
	}()
}

func (w *Worker[T, F]) Consume(consumer func(val interface{}, stream chan F)) {
	for i := 0; i < w.workerNum; i++ {
		go func(idx int) {
			doneStream := w.doneStreams[idx]

			for val := range w.producerStream {
				consumer(val, w.consumerStream)
			}

			doneStream <- true
		}(i)
	}
}

func (w *Worker[T, F]) Wait(waitFn func(data DataOrError)) {
	for val := range w.consumerStream {
		waitFn(val)
	}
}
