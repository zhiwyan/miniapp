package weapp

const (
	apiURLLink = "/wxa/generate_urllink"
)

type URLLink struct {
	//通过 URL Link 进入的小程序页面路径，必须是已经发布的小程序存在的页面，不可携带 query 。path 为空时会跳转小程序主页
	Path string `json:"path"`
	//通过 URL Link 进入小程序时的query，最大1024个字符，只支持数字，大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~
	Query string `json:"query"`
	// 生成的 URL Link 类型，到期失效：true，永久有效：false
	IsExpire bool `json:"is_expire,omitempty"`
	// 小程序 URL Link 失效类型，失效时间：0，失效间隔天数：1
	ExpireType int64 `json:"expire_type,omitempty"`
	// 到期失效的 URL Link 的失效时间，为 Unix 时间戳。生成的到期失效 URL Link 在该时间前有效。最长有效期为1年。expire_type 为 0 必填
	ExpireTime int64 `json:"expire_time,omitempty"`
	// 到期失效的URL Link的失效间隔天数。生成的到期失效URL Link在该间隔时间到达前有效。最长间隔天数为365天。expire_type 为 1 必填
	ExpireInterval int64 `json:"expire_interval,omitempty"`
}

type URLLinkResponse struct {
	CommonError
	// 生成的小程序 URL Link
	UrlLink string `json:"url_link"`
}

// 获取小程序 URL Link，适用于短信、邮件、外部网页等拉起小程序的业务场景。
//
// token 微信access_token
func (link *URLLink) Generate(token string) (*URLLinkResponse, error) {
	api := baseURL + apiURLLink
	return link.generate(api, token)
}

func (link *URLLink) generate(api, token string) (*URLLinkResponse, error) {
	uri, err := tokenAPI(api, token)
	if err != nil {
		return nil, err
	}

	res := new(URLLinkResponse)
	err = postJSON(uri, link, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
