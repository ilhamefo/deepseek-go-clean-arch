package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gofrs/uuid"
)

const (
	DATE_LAYOUT = "2006/01/02"
)

type buffer struct {
	r         []byte
	runeBytes [utf8.UTFMax]byte
}

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

func GenerateUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return uuid.String()
}

func Underscore(s string) string {
	b := buffer{
		r: make([]byte, 0, len(s)),
	}
	var m rune
	var w bool
	for _, ch := range s {
		if unicode.IsUpper(ch) {
			handleUppercase(&b, &m, &w, ch)
		} else {
			handleLowercase(&b, &m, &w, ch)
		}
	}
	if m != 0 {
		if !w {
			b.indent()
		}
		b.write(m)
	}

	return string(b.r)
}

func handleUppercase(b *buffer, m *rune, w *bool, ch rune) {
	if *m != 0 {
		if !*w {
			b.indent()
			*w = true
		}
		b.write(*m)
	}
	*m = unicode.ToLower(ch)
}

func handleLowercase(b *buffer, m *rune, w *bool, ch rune) {
	if *m != 0 {
		b.indent()
		b.write(*m)
		*m = 0
		*w = false
	}
	b.write(ch)
}

func (b *buffer) indent() {
	if len(b.r) > 0 {
		b.r = append(b.r, '_')
	}
}

func (b *buffer) write(r rune) {
	if r < utf8.RuneSelf {
		b.r = append(b.r, byte(r))
		return
	}
	n := utf8.EncodeRune(b.runeBytes[0:], r)
	b.r = append(b.r, b.runeBytes[0:n]...)
}
