package kernel

import (
	"fmt"
	"github.com/gocolly/colly"
	"missingPersons/contract"
)

type signal struct {
	errorCh chan error
	dataCh  chan [][]contract.ReceiverData
}

func (s signal) Error() chan error {
	return s.errorCh
}

func (s signal) Data() chan [][]contract.ReceiverData {
	return s.dataCh
}

type kernel struct {
	url       string
	element   string
	baseUrl   string
	onHtml    func(e *colly.HTMLElement, signal contract.Signal)
	onError   func(_ *colly.Response, err error, signal contract.Signal)
	onScraped func(_ *colly.Response, signal contract.Signal)
	onRequest func(r *colly.Request, signal contract.Signal)
}

// e.Request.AbsoluteURL(e.Attr("href"))

func (r kernel) Run() {
	c := colly.NewCollector()

	c.OnHTML(r.Element(), func(element *colly.HTMLElement) {
		r.onHtml(element, nil)
	})

	c.OnError(func(response *colly.Response, err error) {
		r.onError(response, err, nil)
	})

	c.OnRequest(func(request *colly.Request) {
		r.onRequest(request, nil)
	})

	c.OnScraped(func(res *colly.Response) {
		r.onScraped(res, nil)
	})

	if err := c.Visit(fmt.Sprintf("%s%s", r.BaseURL(), r.url)); err != nil {
		fmt.Println("error shit")
	}
}

func (r kernel) Element() string {
	return r.element
}

func (r kernel) BaseURL() string {
	return r.baseUrl
}

func initKernel(
	baseUrl string,
	url string,
	element string,
	onHtml func(e *colly.HTMLElement, signal contract.Signal),
	onError func(_ *colly.Response, err error, signal contract.Signal),
	onRequest func(r *colly.Request, signal contract.Signal),
	onScraped func(_ *colly.Response, signal contract.Signal),
) contract.Kernel {
	return kernel{
		url:       url,
		element:   element,
		baseUrl:   baseUrl,
		onHtml:    onHtml,
		onError:   onError,
		onScraped: onScraped,
		onRequest: onRequest,
	}
}

func Start(
	baseUrl string,
	url string,
	element string,
	onHtml func(e *colly.HTMLElement, signal contract.Signal),
	onError func(_ *colly.Response, err error, signal contract.Signal),
	onRequest func(r *colly.Request, signal contract.Signal),
	onScraped func(_ *colly.Response, signal contract.Signal),
) {
	kernel := initKernel(baseUrl, url, element, onHtml, onError, onRequest, onScraped)
	kernel.Run()
}
