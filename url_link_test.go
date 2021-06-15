package weapp

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	link := URLLink{}

	res, err := link.Generate("acces-token")
	if err != nil {
		fmt.Println(err)
		return
	}

	err =res.GetResponseError()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(res)
}
