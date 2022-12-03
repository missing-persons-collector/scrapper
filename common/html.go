package common

import (
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
	"io/ioutil"
	"missingPersons/httpClient"
	"strings"
)

func GetListing(url string, query string) ([]*html.Node, error) {
	pageHtml, err := GetHtml(url)

	if err != nil {
		return nil, err
	}

	listing, err := QueryList(pageHtml, query)

	if err != nil {
		return nil, err
	}

	return listing, nil
}

func QueryList(pageHtml *html.Node, query string) ([]*html.Node, error) {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return nil, err
	}

	node := cascadia.QueryAll(pageHtml, sel)

	return node, nil
}

func Query(pageHtml *html.Node, query string) (*html.Node, error) {
	sel, err := cascadia.Parse(query)
	if err != nil {
		return nil, err
	}

	node := cascadia.Query(pageHtml, sel)

	return node, nil
}

func GetHtml(url string) (*html.Node, error) {
	response, err := httpClient.SendRequest(url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(body)))

	if err != nil {
		return nil, err
	}

	return doc, nil
}

func GetAttr(attr string, attributes []html.Attribute) string {
	for _, a := range attributes {
		if a.Key == attr {
			return a.Val
		}
	}

	return ""
}
