package weapp

import (
	"fmt"
	"testing"
)

func TestUrlGetAccessToken(t *testing.T) {
	res, err := GetAccessToken("appid", "secret")
	if err != nil {
		// 处理一般错误信息
		return
	}

	if err := res.GetResponseError(); err !=nil {
		// 处理微信返回错误信息
		return
	}

	fmt.Printf("返回结果: %#v", res)
}

func TestGenerate(t *testing.T) {
	link := URLLink{
		Path: "/path",
		Query: "name=kk&age=29",
		IsExpire: true,
		ExpireType: 1,
		ExpireInterval: 7,
	}

	res, err := link.Generate("access_token")
	if err != nil {
		fmt.Println("----", err)
		return
	}
	fmt.Println(res)

	err =res.GetResponseError()
	if err != nil {
		fmt.Println("++++",err)
		return
	}

	fmt.Println(res)
}
