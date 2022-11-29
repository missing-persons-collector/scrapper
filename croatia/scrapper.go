package croatia

import (
	"fmt"
	"github.com/gocolly/colly"
	"missingPersons/contract"
	"missingPersons/kernel"
)

type dataOrError struct {
	data  []contract.ReceiverData
	error error
}

func Start() []contract.CollectedPage {
	baseUrl := "https://nestali.gov.hr"
	url := "/nestale-osobe-403/403?&page=%d"
	element := ".osoba-wrapper .osoba-img"
	pageCollector := make(chan contract.CollectedPage, 0)

	fmt.Println("Croatia: Starting collection...")

	go func() {
		page := 1

		for {
			people := make([][]contract.ReceiverData, 0)

			kernel.Start(
				baseUrl,
				fmt.Sprintf(url, page),
				element,
				func(e *colly.HTMLElement, signal contract.Signal) {
					url := fmt.Sprintf("%s%s", baseUrl, e.Attr("href"))

					internalScraper(url, func(d dataOrError) {
						if d.error != nil {
							pageCollector <- contract.CollectedPage{
								Page:  0,
								Data:  nil,
								Error: d.error,
							}

							close(pageCollector)
							return
						}

						people = append(people, d.data)
					})
				}, func(_ *colly.Response, err error, signal contract.Signal) {
					pageCollector <- contract.CollectedPage{
						Page:  0,
						Data:  nil,
						Error: err,
					}
					close(pageCollector)
				}, func(r *colly.Request, signal contract.Signal) {
				}, func(_ *colly.Response, signal contract.Signal) {
				})

			if len(people) == 0 {
				close(pageCollector)
				return
			}

			pageCollector <- contract.CollectedPage{
				Page: page,
				Data: people,
			}

			page++
		}
	}()

	pages := make([]contract.CollectedPage, 0)
	for page := range pageCollector {
		if page.Error != nil {
			fmt.Printf("Croatia: An error occurred: %s. Stopping collection!", page.Error.Error())

			break
		}

		pages = append(pages, page)
		fmt.Println("Croatia: page collected: ", page.Page)
	}

	fmt.Println("Croatia: scrapping done!")

	return pages
}

func internalScraper(url string, onData func(d dataOrError)) {
	c := colly.NewCollector()
	holder := make([]contract.ReceiverData, 0)

	c.OnHTML(".profile_details_right dl:first-of-type", func(e *colly.HTMLElement) {
		data := contract.ReceiverData{
			Key:   "",
			Value: "",
		}
		fullCollected := false

		e.ForEach("*", func(i int, element *colly.HTMLElement) {
			if element.Name == "dt" {
				data.Key = element.Text
			}

			if element.Name == "dd" {
				data.Value = element.Text
				fullCollected = true
			}

			if fullCollected {
				holder = append(holder, data)
				fullCollected = false
				data = contract.ReceiverData{
					Key:   "",
					Value: "",
				}
			}
		})
	})

	c.OnScraped(func(response *colly.Response) {
		onData(dataOrError{
			data:  holder,
			error: nil,
		})
	})

	if err := c.Visit(url); err != nil {
		onData(dataOrError{
			data:  nil,
			error: err,
		})
	}
}
