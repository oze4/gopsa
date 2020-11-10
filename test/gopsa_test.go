//
// TODO (MAJOR): tests
//
package gopsa

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/oze4/gopsa"
)

func TestDefault(t *testing.T) {
	sl, err := gopsa.GetSetList(context.Background(), http.DefaultClient, gopsa.SetOriginal)
	if err != nil {
		fmt.Println("ERR :", err)
	}

	for _, d := range sl.Data {
		fmt.Println(d.Name(), "\t\t\t", d.PSAIdentifier())
		fmt.Println()
	}
}
