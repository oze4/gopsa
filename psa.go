package gopsa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetSetList gets the cards in a given set
func GetSetList(ctx context.Context, c *http.Client, s Set) (SetList, error) {
	body, err := newRequestBody(s)
	if err != nil {
		return SetList{}, err
	}

	url := URLBase + "/cardfacts/GetSetList"
	resp, err := doNewRequest(ctx, c, http.MethodPost, url, body, 200)
	if err != nil {
		return SetList{}, err
	}

	return readGetSetListResponse(resp)
}

// Copies io.ReadCloser to byte array, then closes the ReadCloser
func copyAndCloseBody(r io.ReadCloser) ([]byte, error) {
	p, e := ioutil.ReadAll(r)
	if e != nil {
		return nil, e
	}
	r.Close()
	return p, nil
}

// Creates POST request body for a given Set
func newRequestForm(s Set) (*requestForm, error) {
	id, err := s.ID()
	if err != nil {
		return nil, err
	}

	return &requestForm{
		Draw:   "0", // There doesn't seem to be an
		Start:  "0", // issue with leaving these two at "0".
		Length: maxResults,
		SetID:  id,
	}, nil
}

// Creates POST body based upon set
func newRequestBody(s Set) (io.Reader, error) {
	form, err := newRequestForm(s)
	if err != nil {
		return nil, err
	}

	body, err := form.ToRequestBody()
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Creates a new http.Request then sends it
func doNewRequest(ctx context.Context, c *http.Client, method string, url string, body io.Reader, expectedStatus int) (io.ReadCloser, error) {
	q, e := http.NewRequest(method, url, body)
	if e != nil {
		return nil, e
	}

	q.Header.Set("Content-Type", "application/json")

	p, e := c.Do(q.WithContext(ctx))
	if e != nil {
		return nil, e
	}

	if p.StatusCode != expectedStatus {
		return nil, fmt.Errorf("Expected StatusCode %d : got %d", expectedStatus, p.StatusCode)
	}

	return p.Body, nil
}

// Reads the response sent to the /GetSetList/ API endpoint
//
// TODO: elaborate more on how data is returned from PSA
// (they return each card name as an html link element) aka `<a>` tag
func readGetSetListResponse(r io.ReadCloser) (SetList, error) {
	b, e := copyAndCloseBody(r)
	if e != nil {
		return SetList{}, e
	}

	var setlist SetList
	if e := json.Unmarshal(b, &setlist); e != nil {
		return SetList{}, e
	}

	for _, card := range setlist.Data {
		z := html.NewTokenizer(strings.NewReader(card.RawName))

		for {
			if z.Next() == html.ErrorToken {
				break
			}

			t := z.Token()
			// Since we are give a single html node, we just need to extract the text from it
			if t.Type == html.TextToken {
				card.name = t.Data
			}

			// Finds an href within a "raw" card name
			for i := 0; i < len(t.Attr); i++ {
				if t.Attr[i].Key == "href" {
					// We are given the card name as an html element in string form,
					// the ID for that card can be found at the end of that URL.
					// This simply pulls out the last query string and holds onto it
					// as the card ID
					hrefStrArr := strings.Split(t.Attr[i].Val, "/")
					index := len(hrefStrArr) - 1
					card.psaIdentifier = hrefStrArr[index]
					// Break out of our loop by making `i` equal to `N < len(t.Attr)`
					i = len(t.Attr)
				}
			}
		}
	}

	return setlist, nil
}
