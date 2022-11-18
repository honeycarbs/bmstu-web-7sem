package model

import (
	"strings"
	"time"
)

const (
	DefaultNoteColor = "CFD2CF"
	shortBodyLen     = 255
)

type Note struct {
	ID        int       `json:"id" db:"id"`
	Header    string    `json:"header" db:"header"`
	Body      string    `json:"body" db:"body"`
	ShortBody string    `json:"shortBody" db:"short_body"`
	Tags      []Tag     `json:"tags" db:"tags"`
	Color     string    `json:"color" db:"color"`
	Edited    time.Time `json:"edited"`
}

func (n *Note) GenerateShortBody() {
	if len(n.Body) < shortBodyLen {
		n.ShortBody = n.Body
	} else {
		n.ShortBody = truncate(n.Body, shortBodyLen)
	}
}

func (n *Note) HasEveryTag(tagNames []string) bool {
	for _, tn := range tagNames {
		if !n.HasSpecificTag(tn) {
			return false
		}
	}
	return true
}

func (n *Note) HasSpecificTag(tagName string) bool {
	for _, t := range n.Tags {
		if strings.Contains(strings.ToLower(t.Label), strings.ToLower(tagName)) {
			return true
		}
	}
	return false
}

func truncate(text string, width int) string {
	r := []rune(text)
	trunc := r[:width]
	return string(trunc)
}
