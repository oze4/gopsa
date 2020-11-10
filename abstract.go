package gopsa

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

// Constants
const (
	URLBase    = "https://www.psacard.com"
	maxResults = "999999"

	SetOriginal Set = iota
	SetFossil
	SetJungle
)

// Card represents a PSA card
type Card struct {
	Number     string `json:"CardNumber,omitempty"`
	RawName    string `json:"CardName,omitempty"`
	name       string
	identifier string
}

// Name returns the name of the Card
func (c *Card) Name() string {
	return c.name
}

// Identifier returns the identifier (which is used for crafting queries)
func (c *Card) Identifier() string {
	return c.identifier
}

// SetList holds a pokemon set card list
type SetList struct {
	Draw              int     `json:"draw,omitempty"`
	RecordsTotal      int     `json:"recordsTotal,omitempty"`
	RecordsFiltered   int     `json:"recordsFiltered,omitempty"`
	HasCheckListItems bool    `json:"hasCheckListItems,omitempty"`
	Data              []*Card `json:"data,omitempty"`
}

// Set is a collection of cards, eg: `Pokemon Fossil (1st Edition)`
type Set int

// Identifier gets the PSA set identifier for a pokemon set
func (s *Set) Identifier() (string, error) {
	switch *s {
	case SetOriginal:
		return "29137", nil
	case SetFossil:
		return "", errors.New("Not implemented")
	case SetJungle:
		return "", errors.New("Not implemented")
	default:
		return "", errors.New("Invalid Set Identifier")
	}
}

// Name gets the PSA Set name for a pokemon set
func (s *Set) Name() (string, error) {
	switch *s {
	case SetOriginal, SetFossil, SetJungle:
		return "1999+Nintendo+Pokemon+Game", nil
	default:
		return "", errors.New("Invalid Set Name")
	}
}

// Form we use to POST
type requestForm struct {
	Draw         string
	Start        string
	Length       string
	SetID        string
	CategoryName string
	SetName      string
}

func (r *requestForm) ToRequestBody() (io.Reader, error) {
	b, e := json.Marshal(r)
	if e != nil {
		return nil, e
	}
	return bytes.NewBuffer(b), nil
}
