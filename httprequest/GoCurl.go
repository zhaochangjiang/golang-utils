package httprequest

/**
 * golang版本的curl请求库
 * Request构造类，用于设置请求参数，发起http请求
 * @author mike <mikemintang@126.com>
 * @blog http://idoubi.cc
 */
import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	neturl "net/url"
)

// Request 构造类
type Request struct {
	cli          *http.Client
	req          *http.Request
	Raw          *http.Request
	Method       string
	URL          string
	Headers      map[string]string
	Cookies      map[string]string
	Queries      map[string]string
	PostDataJSON map[string]interface{}
	PostData     map[string]string
}

//NewRequest 创建一个Request实例
func NewRequest() *Request {
	return &Request{}
}

//SetPostDataJSON 设置输入json的参数
func (req *Request) SetPostDataJSON(postDataJSON map[string]interface{}) *Request {
	req.PostDataJSON = postDataJSON
	return req
}

//SetMethod 设置请求方法
func (req *Request) SetMethod(method string) *Request {
	req.Method = method
	return req
}

//SetURL 设置请求地址
func (req *Request) SetURL(url string) *Request {
	req.URL = url
	return req
}

//SetHeaders 设置请求头
func (req *Request) SetHeaders(headers map[string]string) *Request {
	req.Headers = headers
	return req
}

//setHeaders 将用户自定义请求头添加到http.Request实例上
func (req *Request) setHeaders() error {
	for k, v := range req.Headers {
		req.req.Header.Set(k, v)
	}
	return nil
}

//SetCookies 设置请求cookies
func (req *Request) SetCookies(cookies map[string]string) *Request {
	req.Cookies = cookies
	return req
}

//setCookies 将用户自定义cookies添加到http.Request实例上
func (req *Request) setCookies() error {
	for k, v := range req.Cookies {
		req.req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}
	return nil
}

//SetQueries 设置url查询参数
func (req *Request) SetQueries(queries map[string]string) *Request {
	req.Queries = queries
	return req
}

//setQueries 将用户自定义url查询参数添加到http.Request上
func (req *Request) setQueries() error {
	q := req.req.URL.Query()
	for k, v := range req.Queries {
		q.Add(k, v)
	}
	req.req.URL.RawQuery = q.Encode()
	return nil
}

//SetPostData 设置post请求的提交数据
func (req *Request) SetPostData(postData map[string]string) *Request {
	req.PostData = postData
	return req
}

//Get 发起get请求
func (req *Request) Get() (*Response, error) {
	return req.Send(req.URL, http.MethodGet)
}

//Post 发起post请求
func (req *Request) Post() (*Response, error) {
	return req.Send(req.URL, http.MethodPost)
}

//Run 运行方法
func (req *Request) Run() (*Response, error) {
	return req.Send(req.URL, req.Method)
}

//Send 发起请求
func (req *Request) Send(url string, method string) (*Response, error) {
	// 检测请求url是否填了
	if url == "" {
		return nil, errors.New("Lack of request url")
	}
	// 检测请求方式是否填了
	if method == "" {
		return nil, errors.New("Lack of request method")
	}

	// 初始化Response对象
	response := NewResponse()

	// 初始化http.Client对象
	req.cli = &http.Client{}

	// 加载用户自定义的post数据到http.Request
	var payload io.Reader
	switch method {
	case "POST":
		if req.PostData != nil {
			postDataString := ""
			ind := 0
			for k, v := range req.PostData {
				v = neturl.QueryEscape(v)
				if ind == 0 {
					postDataString += k + "=" + v
				} else {
					postDataString += "&" + k + "=" + v
				}
				ind++
			}
			payload = bytes.NewReader([]byte(postDataString))
		}
		break
	case "POSTJSON":
		if req.PostData != nil {

			jData, err := json.Marshal(req.PostDataJSON)
			if err != nil {
				return nil, err
			}

			payload = bytes.NewReader(jData)
		}
	default:
		payload = nil
		break
	}
	reqs, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req.req = reqs

	req.setHeaders()
	req.setCookies()
	req.setQueries()

	req.Raw = req.req
	resp, err := req.cli.Do(req.req)
	if err != nil {
		return nil, err
	}
	response.Raw = resp

	defer response.Raw.Body.Close()

	response.parseHeaders()
	response.parseBody()

	return response, nil
}
