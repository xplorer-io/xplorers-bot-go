package common

import (
	"net/url"
)

func ArrayContainsItem(array []string, item string) bool {
	for _, arrayItem := range array {
		if item == arrayItem {
			return true
		}
	}
	return false
}

func ParseQueryString(requestBody string) (url.Values, error) {
	parsedQueryString, err := url.ParseQuery(string(requestBody))
	if err != nil {
		return nil, err
	}
	return parsedQueryString, nil
}
