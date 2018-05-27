package avaxo

import (
	"fmt"
	"testing"
)

func TestAllForwardOpts(t *testing.T) {
	t.Skip()
	db := NewDb("host= user= password= dbname=")
	opts, err := db.AllForwardOpts()
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range opts {
		fmt.Println(item)
	}
}
