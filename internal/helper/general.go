package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	DATE_LAYOUT = "2006/01/02"
)

func PrettyPrint(b ...interface{}) {
	for _, i := range b {
		s, err := json.MarshalIndent(i, "", "\t")
		if err != nil {
			fmt.Print(err.Error())
		}
		fmt.Print(string(s) + "\n")
	}
}

func StartDateParser(str string) (*time.Time, error) {
	res, err := time.Parse(DATE_LAYOUT, str)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func EndDateParser(str string) (*time.Time, error) {
	res, err := time.Parse(DATE_LAYOUT, str)
	if err != nil {
		return nil, err
	}

	res = res.Add(24 * time.Hour).Add(-1 * time.Second)

	return &res, nil
}

func NormalizeString(input string) string {
	input = strings.TrimSpace(input)
	input = strings.ToUpper(input)
	input = regexp.MustCompile(`[^A-Z0-9\s-]`).ReplaceAllString(input, "")
	input = strings.ReplaceAll(input, " ", "_")
	return input
}
