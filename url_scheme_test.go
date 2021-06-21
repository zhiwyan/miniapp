package weapp

import (
	"fmt"
	"testing"
)

func TestURLSchemeGenerate(t *testing.T) {
	scheme := URLScheme{
		SchemedInfo:  &SchemedInfo{
			Path:  "mock/path",
			Query:  "",
		},
		IsExpire: false,
	}

	res, err := scheme.Generate("acces-token")
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

