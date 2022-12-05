package croatia

import (
	"fmt"
	"missingPersons/common"
	"missingPersons/download"
	"missingPersons/worker"
)

type downloadData struct {
	url      string
	customId string
}

type pathData struct {
	path     string
	customId string
}

type dataOrError struct {
	data  downloadData
	error error
}

func (c dataOrError) Data() interface{} {
	return c.data
}

func (c dataOrError) Error() error {
	return c.error
}

type pathOrError struct {
	data  pathData
	error error
}

func (c pathOrError) Data() interface{} {
	return c.data
}

func (c pathOrError) Error() error {
	return c.error
}

func preDownloadImages(people []common.RawPerson) (map[string]string, error) {
	downloadCache := make(map[string]string)
	for _, p := range people {
		id, err := createPersonId(p)

		if err != nil {
			return nil, err
		}

		downloadCache[id] = ""
	}

	imageSaver := download.NewFsImageSaver()

	w := worker.NewWorker[dataOrError, pathOrError](10)

	w.Produce(func(producerStream chan dataOrError, stopFn func()) {
		for _, p := range people {
			id, err := createPersonId(p)

			if err != nil {
				producerStream <- dataOrError{error: err}

				continue
			}

			producerStream <- dataOrError{
				data: downloadData{
					url:      p.ImageURL,
					customId: id,
				},
				error: nil,
			}
		}

		stopFn()
	})

	w.Consume(func(val interface{}, consumerStream chan pathOrError) {
		data := val.(dataOrError)

		path, err := imageSaver.Save(data.data.url, data.data.customId)

		if err != nil {
			fmt.Printf("Cannot fetch/save image to filesystem: %s\n", err.Error())
			consumerStream <- pathOrError{error: err}

			return
		}

		consumerStream <- pathOrError{data: pathData{
			path:     path,
			customId: data.data.customId,
		}}
	})

	w.Wait(func(data worker.DataOrError) {
		d := data.(pathOrError)

		if d.error != nil {
			fmt.Printf("Cannot fetch/save image to filesystem: %s\n", d.error.Error())

			return
		}

		downloadCache[d.data.customId] = d.data.path
	})

	return downloadCache, nil
}
