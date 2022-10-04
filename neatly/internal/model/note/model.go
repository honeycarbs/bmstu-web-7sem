package note

import (
	"neatly/internal/model/tag"
	"strings"
)

const (
	DEFAULT_NOTE_COLOR = "CFD2CF"
)

type Note struct {
	ID        int       `json:"id" db:"id"`
	Header    string    `json:"header" db:"header"`
	Body      string    `json:"body" db:"body"`
	ShortBody string    `json:"shortBody" db:"short_body"`
	Tags      []tag.Tag `json:"tags" db:"tags"` // []tag.Tag
	Color     string    `json:"color" db:"color"`
}

func (n *Note) GenerateShortBody() {
	var shortLen int
	if len(n.Body) > 500 {
		shortLen = 300
	} else {
		shortLen = len(n.Body)
	}
	n.ShortBody = n.Body[0:shortLen] // TODO: make smart
}

func (n *Note) HasEveryTag(tagNames []string) bool {
	for _, tn := range tagNames {
		if !n.hasSpecificTag(tn) {
			return false
		}
	}
	return true
}

func (n *Note) hasSpecificTag(tagName string) bool {
	for _, t := range n.Tags {
		if strings.Contains(strings.ToLower(t.Name), strings.ToLower(tagName)) {
			return true
		}
	}
	return false
}
