package httprequest

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	utils "github.com/zhaochangjiang/golang-utils/myutils"
)

//RequestParamsFormat 公共参数解析
type RequestParamsFormat struct {
	httpRequestContent *http.Request
	RequestParams      *map[string]interface{}
	PostParams         *map[string]interface{}
	GetParams          *map[string]interface{}
}

//Run 初始化导出服务
func (ea *RequestParamsFormat) Run(r *http.Request) *map[string]interface{} {
	ea.httpRequestContent = r
	ea.httpRequestContent.ParseForm()
	ea.httpRequestContent.ParseMultipartForm(32 << 20)
	ea.initParams()
	return ea.RequestParams
}

func (ea *RequestParamsFormat) paramsMaps(k string, v []string) *map[string]interface{} {
	var params = make(map[string]interface{})
	lengthValue := len(v)
	if lengthValue < 1 {
		v[0] = ""
	}

	regex := regexp.MustCompile(`(\[.*\])+$`).FindAllString(k, -1)

	if len(regex) > 0 {
		switch len(regex) {
		case 1:
			count := len([]rune(regex[0]))
			rightString := utils.SubString(k, -count, count)
			if rightString == "[]" {

				paramsK := utils.SubString(k, 0, utils.StringLength(k)-2)

				if lengthValue > 1 {
					params[paramsK] = v
				} else {
					params[paramsK] = v[0]
				}
			} else {
				rightString = strings.TrimRight(strings.TrimLeft(rightString, "["), "]")
				list := strings.Split(rightString, "][")
				length := len(list)
				if length > 1 {
					var res interface{}
					for m := 0; m < length; m++ {
						if m == 0 {
							res = map[string]interface{}{list[length-m-1]: v[0]}
						} else {
							res = map[string]interface{}{list[length-m-1]: res}
						}
					}
					paramsK := utils.SubString(k, 0, utils.StringLength(k)-count)
					params[paramsK] = res
				}
			}
			break
		default:
			panic("the params is not support,please do it. the content is follow:")
		}

	} else {
		params[k] = v[0]
	}
	return &params
}

//orgDataFormat
func (ea *RequestParamsFormat) orgDataFormat(m int, list []string, v []string, res interface{}) interface{} {

	keyVal, err := strconv.Atoi(list[m])
	if nil != err {
		panic(err)
	}

	if list[m] != "0" && keyVal == 0 {
		res = map[string]string{list[m]: v[0]}
	} else {
		res = []string{v[0]}
	}
	return res
}

//ParamsOrganization
func (ea *RequestParamsFormat) initParams() {

	var c = ea.httpRequestContent
	if query := c.URL.Query(); nil != query {
		for k, v := range query {
			params := ea.paramsMaps(k, v)
			ea.GetParams = utils.MapMerge(ea.GetParams, params)
		}
	}
	res, err := json.Marshal(c.PostForm)
	if nil != err {
		panic(err)
	}
	log.Println("the request params:", string(res))
	if nil != c.PostForm {
		for k, v := range c.PostForm {
			params := ea.paramsMaps(k, v)
			ea.PostParams = utils.MapMerge(ea.PostParams, params)
		}
	}

	postP, err := json.Marshal(*ea.PostParams)
	if err != nil {
		panic(postP)
	}
	ea.RequestParams = utils.MapMerge(ea.GetParams, ea.PostParams)
}
