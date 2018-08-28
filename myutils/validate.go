package myutils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	//ErrorParamsUNSet //参数没设置
	ErrorParamsUNSet = "7110000|参数%s不存在"

	//ErrorParamsEmpty 参数不能为空
	ErrorParamsEmpty = "7110001| %s不能为空"

	//ErrorParamsCategoryError 数据格式不正确
	ErrorParamsCategoryError = "7110002| %s格式不正确"

	//ErrorParamsCategoryPregError 数据格式不正确
	ErrorParamsCategoryPregError = "7110003| %s格式不正确"
)

//ValidateRule 校验参数的类型是
type ValidateRule struct {
	Key          string //数据对应的key
	Val          string //数据值
	Set          bool   //判断是否存在，true:如果不存在，则报错
	Empty        bool   //验证是否为空 ，默认false，true:如果为空字符串（数字为0）不报错
	Category     string //数据类型 当前支持 int
	Preg         string //数据类型对应的正则表达式，空字符串不校验
	DefaultValue string //默认值设置
	Min          int64  //最小值或字符串最小长度
	Max          int64  //最大值或字符串最大长度

	//字段别名 如果不为"",则默认会用本内容替换错误提示中的字段描述内容
	//如: Key:"username"    "username不能为空"   Alias:"用户名" 则返回信息 "用户名不能为空"
	Alias string
}

//Validate 校验结构体
type Validate struct {
	Params   *map[string]string
	Rules    *[]ValidateRule
	Result   *[]ValidateResult
	Language string
	NotValue string //参数值永远不可能是的值,此参数用于如果key不存在，给key赋默认值占位，防止程序报错用
}

//ValidateResult 校验结果
type ValidateResult struct {
	Code    int
	Message string
	Data    interface{}
}

//NewValidate 参数校验初始化
func NewValidate(params *map[string]string, rule *[]ValidateRule) *Validate {
	//NewValidate 初始化验证参数对象
	validate := &Validate{NotValue: ""}
	validate.Rules = rule
	validate.Params = params
	tmp := make([]ValidateResult, 0)
	validate.Result = &tmp
	validate.Language = "zh-cn"
	return validate
}

//Run 运行
func (validate *Validate) Run() *[]ValidateResult {
	for _, v := range *validate.Rules {
		validateResult := validate.everyValidate(&v)
		if validateResult.Code > 0 {
			*validate.Result = append(*validate.Result, validateResult)
		}
	}
	return validate.Result
}

//判断是否存在
func (validate *Validate) judgeIsSet(v *ValidateRule) bool {
	//判断key是否存在
	if _, ok := (*validate.Params)[v.Key]; ok {
		//如果Params 中存在key ，则将值设置进规则。用于判断
		v.Val = (*validate.Params)[v.Key]
		return true
	}
	v.Val = validate.NotValue
	return false
}
func (validate *Validate) needJudgeIsSet(flagExists bool, v *ValidateRule, validateResult *ValidateResult) {
	//如果需要判断是否存在
	if v.Set == false && flagExists == false {
		validate.getError(ErrorParamsUNSet, validateResult, v)
	}
}

func (validate *Validate) getError(err string, vr *ValidateResult, rule *ValidateRule) {
	v := strings.Split(err, "|")
	var erro error
	vr.Code, erro = strconv.Atoi(v[0])
	if erro != nil {
		panic(err)
	}
	if rule.Alias == "" {
		vr.Message = fmt.Sprintf(v[1], rule.Key)
	} else {
		vr.Message = fmt.Sprintf(v[1], rule.Alias)
	}
}

//everyValidate
func (validate *Validate) everyValidate(v *ValidateRule) ValidateResult {

	validateResult := new(ValidateResult)
	validateResult.Data = v.Key
	//判断返回是否存在的结果,此处的返回结果 参数后边还会使用，顾放置在此处
	if flagExists := validate.judgeIsSet(v); flagExists == true {
		*validateResult = validate.ruleValidate(flagExists, v, validateResult)
		(*validate.Params)[v.Key] = v.Val
	} else { //如果值不存在，则设置默认值
		(*validate.Params)[v.Key] = v.DefaultValue
	}

	return *validateResult
}

//满足条件的规则验证
func (validate *Validate) ruleValidate(flagExists bool, v *ValidateRule, validateResult *ValidateResult) ValidateResult {
	//检查Set项
	validate.needJudgeIsSet(flagExists, v, validateResult)
	if validateResult.Code > 0 {
		return *validateResult
	}

	//如果数据存在，则往下验证。
	//	if flagExists == true {
	validate.validateEmpty(v, validateResult)
	if validateResult.Code > 0 {
		return *validateResult
	}

	//检验类型
	validate.validateCategory(v, validateResult)
	if validateResult.Code > 0 {
		return *validateResult
	}

	//检验类型
	validate.validatePreg(v, validateResult)
	if validateResult.Code > 0 {
		return *validateResult
	}
	return *validateResult
}

func (validate *Validate) validatePreg(v *ValidateRule, validateResult *ValidateResult) {

	//如果验证正则表达式的数值不为空字符串
	if v.Preg != "" {
		match, err := regexp.MatchString(v.Preg, v.Val)
		if match != true || err != nil {
			validate.getError(ErrorParamsCategoryPregError, validateResult, v)
		}
	}
}

//验证数据类型
func (validate *Validate) validateCategory(v *ValidateRule, validateResult *ValidateResult) {
	var err error
	switch v.Category {
	case "string": //字符串不验证，但是要写在此处
		err = validate.validateLength(v)
		break
	case "int": //int 如果机器是64位 等于int64，如果机器是32位 等于int32
		_, err = strconv.Atoi(v.Val)
		break
	case "int8":
		err = validate.validateNumber(v, -128, 127)
		break
	case "int16":
		err = validate.validateNumber(v, -32768, 32767)
		break
	case "int32":
		err = validate.validateNumber(v, -2147483648, 2147483647)
		break
	case "int64":
		_, err = strconv.ParseInt(v.Val, 10, 64)
		break
	case "uint8":
		err = validate.validateNumber(v, 0, 255)
		break
	case "uint16":
		err = validate.validateNumber(v, 0, 65535)
		break
	case "uint32":
		err = validate.validateNumber(v, 0, 4294967295)
		break
	case "uint64":
		_, err = strconv.ParseUint(v.Val, 10, 64)
		break
	case "float64":
		_, err = strconv.ParseFloat(v.Val, 64)
		break
	case "float32":
		err = validate.validateNumberFloat(v, -2147483648, 2147483647)
		break
	case "bool":
		err = validate.validateBoolean(v)
		break
	case "boolean":
		err = validate.validateBoolean(v)
		break
	default:
		break
	}
	if err != nil {
		validate.getError(ErrorParamsCategoryError, validateResult, v)
	}
}

//validateLeng 校验字符串长度
func (validate *Validate) validateLength(v *ValidateRule) error {
	var err error
	stringLength := int64(StringLength(v.Val))
	if stringLength < v.Min || stringLength > v.Max {
		err = errors.New("length is error")
	}
	return err
}

//校验bool型数据
func (validate *Validate) validateBoolean(v *ValidateRule) error {
	var err error
	if v.Val == "" || v.Val == "0" || v.Val == "false" || v.Val == "null" || v.Val == "nil" {
		v.Val = "false"
	} else if v.Val == "true" || v.Val == "1" {
		v.Val = "true"
	} else {
		err = errors.New("is not " + v.Category)
	}
	return err
}

//校验数字
func (validate *Validate) validateNumber(v *ValidateRule, min int64, max int64) error {
	var numberValue int64
	var err error
	numberValue, err = strconv.ParseInt(v.Val, 10, 64)
	if err != nil {
		return err
	}
	//如果数值不在区间int32区间，则报错
	if numberValue < min || numberValue > max {
		err = errors.New("is not " + v.Category)
	}
	return err
}

//校验数字
func (validate *Validate) validateNumberFloat(v *ValidateRule, min float64, max float64) error {
	var numberValue float64
	var err error
	numberValue, err = strconv.ParseFloat(v.Val, 64)
	if err != nil {
		return err
	}
	//如果数值不在区间int32区间，则报错
	if numberValue < min || numberValue > max {

		err = errors.New("is not " + v.Category)
	}
	return err
}

//验证是否为空
func (validate *Validate) validateEmpty(v *ValidateRule, validateResult *ValidateResult) {
	//如果要验证为空
	if v.Empty == false {
		validate.getError(ErrorParamsEmpty, validateResult, v)
	}
}
