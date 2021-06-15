package weapp

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	link := URLLink{
		Path: "/",
		Query: "",
		IsExpire: true,
		ExpireType: 1,
		ExpireInterval: 7,
	}

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
