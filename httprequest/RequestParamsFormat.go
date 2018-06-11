package httprequest

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	utils "github.com/zhaochangjiang/golang-utils/utils"
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
	ea.initParams(r)
	return ea.RequestParams
}

//InitParams
func (ea *RequestParamsFormat) initParams(r *http.Request) {
	ea.httpRequestContent = r
	ea.httpRequestContent.ParseForm()
	ea.paramsOrganization()
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
func (ea *RequestParamsFormat) paramsOrganization() {
	var c = ea.httpRequestContent
	if nil != c.Form {
		for k, v := range c.Form {
			params := ea.paramsMaps(k, v)
			*ea.GetParams = *(utils.MapMerge(ea.GetParams, params))
		}
	}
	if nil != c.PostForm {
		for k, v := range c.PostForm {
			params := ea.paramsMaps(k, v)
			*ea.PostParams = *(utils.MapMerge(ea.PostParams, params))
		}
	}
	ea.RequestParams = utils.MapMerge(ea.GetParams, ea.PostParams)
}
