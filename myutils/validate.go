package myutils

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

const (

	//ErrorParamsUNSet //参数没设置
	ErrorParamsUNSet = 7110000

	//ErrorParamsEmpty 参数不能为空
	ErrorParamsEmpty = 7110001

	//ErrorParamsCategoryError 数据格式不正确
	ErrorParamsCategoryError = 7110002

	//ErrorParamsCategoryPregError 数据格式不正确
	ErrorParamsCategoryPregError = 7110003

	//ErrorParamsAreaInError 数据区间不正确
	ErrorParamsAreaInError = 7110004

	//ErrorParamsMaxAndMinAreaInError 数据区间不正确
	ErrorParamsMaxAndMinAreaInError = 7110005
)

//ErrorMsg 错误信息描述
var ErrorMsg = make(map[int]string)

//@params language 设置语言包的类型 当前默认zh-cn
func errorMsginit(language string) {
	ErrorMsg[ErrorParamsEmpty] = "%s未设置"
	ErrorMsg[ErrorParamsUNSet] = "参数%s不存在"
	ErrorMsg[ErrorParamsCategoryError] = "%s格式不正确"
	ErrorMsg[ErrorParamsCategoryPregError] = "%s不是%s"
	ErrorMsg[ErrorParamsAreaInError] = "%s不在(%s-%s)"
	ErrorMsg[ErrorParamsMaxAndMinAreaInError] = "%s的校验参数Min和Max不在(%s-%s)范围内"
}

//ValidateRule 校验参数的类型是
type ValidateRule struct {
	Key          string //数据对应的key
	Val          string //数据值
	Set          bool   //判断是否存在，true:如果不存在，则报错
	Empty        bool   //验证是否为空 ，默认false，true:如果为空字符串（数字为0）不报错
	Category     string //数据类型 当前支持 int
	Preg         string //数据类型对应的正则表达式，空字符串不校验
	DefaultValue string //默认值设置 默认空字符串
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
	NotValue string //TODO 参数值永远不可能是的值,此参数用于如果key不存在，给key赋默认值占位，防止程序报错用
}

//ValidateResult 校验结果
type ValidateResult struct {
	Code    int
	Message string
	Data    interface{}
}

//NewValidate 参数校验初始化
func NewValidate(params *map[string]string, rule *[]ValidateRule) *Validate {
	tmp := make([]ValidateResult, 0)
	//NewValidate 初始化验证参数对象
	validate := &Validate{NotValue: "", Rules: rule, Params: params, Language: "zh-cn", Result: &tmp}
	return validate
}

//Run 运行
func (validate *Validate) Run() *[]ValidateResult {
	//初始化提示信息
	errorMsginit(validate.Language)
	for _, v := range *validate.Rules {
		if v.Alias == "" {
			v.Alias = v.Key
		}
		if validateResult := validate.everyValidate(&v); validateResult.Code > 0 {
			*validate.Result = append(*validate.Result, validateResult)
		}
	}
	return validate.Result
}

//SetNotValue TODO 保留功能字段
func (validate *Validate) SetNotValue(notValue string) *Validate {
	validate.NotValue = notValue
	return validate
}

//SetLanguage 设置语言包
func (validate *Validate) SetLanguage(language string) *Validate {
	validate.Language = language
	return validate
}

//判断是否存在
func (validate *Validate) judgeIsSet(v *ValidateRule) bool {
	//判断key是否存在
	if _, ok := (*validate.Params)[v.Key]; ok {
		//如果Params 中存在key ，则将值设置进规则。用于判断
		v.Val = (*validate.Params)[v.Key]
		return true
	}
	return false
}

//needJudgeIsSet
func (validate *Validate) needJudgeIsSet(flagExists bool, v *ValidateRule, validateResult *ValidateResult) {
	//如果需要判断是否存在
	if v.Set == false && flagExists == false {
		validate.getError(ErrorParamsUNSet, validateResult, v)
	}
}

func (validate *Validate) getError(errcode int, vr *ValidateResult, rule *ValidateRule) {

	vr.Code = errcode
	if vr.Message == "" {
		vr.Message = fmt.Sprintf(ErrorMsg[errcode], rule.Alias)
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

//参数区间范围校验
func (validate *Validate) paramsUnit64Area(numberVal uint64, v *ValidateRule, validateResult *ValidateResult) {

	//如果传递的数据范围不在int64 和uint64最大的支持范围的交集内，报错
	if v.Min < 0 || v.Max > 9223372036854775807 {
		validateResult.Message = fmt.Sprintf(ErrorMsg[ErrorParamsMaxAndMinAreaInError], v.Alias, strconv.FormatInt(0, 10), strconv.FormatInt(9223372036854775807, 10))
		validate.getError(ErrorParamsMaxAndMinAreaInError, validateResult, v)
		return
	}
	//如果最下值和最小值不等，则说明一定设置了最大值或最小值，或两个都设置
	//否则最大值等于最小值且不等于0，需要判断
	if v.Min != v.Max || (v.Min == v.Max && v.Min != 0 && uint64(v.Min) != numberVal) {
		if numberVal < uint64(v.Min) || numberVal > uint64(v.Max) {
			validateResult.Message = fmt.Sprintf(ErrorMsg[ErrorParamsAreaInError], v.Alias, strconv.FormatInt(v.Min, 10), strconv.FormatInt(v.Max, 10))
			validate.getError(ErrorParamsCategoryError, validateResult, v)
			return
		}
	}
}

//参数区间范围校验
func (validate *Validate) paramsArea(numberVal int64, v *ValidateRule, validateResult *ValidateResult) {

	//如果最下值和最小值不等，则说明一定设置了最大值或最小值，或两个都设置
	//否则最大值等于最小值且不等于0，需要判断
	if v.Min != v.Max || (v.Min == v.Max && v.Min != 0 && v.Min != numberVal) {
		if numberVal < v.Min || numberVal > v.Max {
			validateResult.Message = fmt.Sprintf(ErrorMsg[ErrorParamsAreaInError], v.Alias, strconv.FormatInt(v.Min, 10), strconv.FormatInt(v.Max, 10))
			validate.getError(ErrorParamsCategoryError, validateResult, v)
			return
		}
	}
}

//验证数据类型
func (validate *Validate) validateCategory(v *ValidateRule, validateResult *ValidateResult) {
	var err error
	switch v.Category {
	case "string": //字符串不验证，但是要写在此处
		validate.validateLength(v, validateResult)
		break
	case "int": //int 如果机器是64位 等于int64，如果机器是32位 等于int32
		var numberVal int
		numberVal, err = strconv.Atoi(v.Val)
		if err != nil {
			validate.getError(ErrorParamsCategoryError, validateResult, v)
		}
		validate.paramsArea(int64(numberVal), v, validateResult)
		break
	case "int8":
		validate.validateNumber(v, -128, 127, validateResult)
		break
	case "int16":
		validate.validateNumber(v, -32768, 32767, validateResult)
		break
	case "int32":
		validate.validateNumber(v, -2147483648, 2147483647, validateResult)
		break
	case "int64":
		validate.validateNumber(v, 0, 0, validateResult)
		break
	case "uint8":
		validate.validateNumber(v, 0, 255, validateResult)
		break
	case "uint16":
		validate.validateNumber(v, 0, 65535, validateResult)
		break
	case "uint32":
		validate.validateNumber(v, 0, 4294967295, validateResult)
		break
	case "uint64":
		var number uint64
		number, err = strconv.ParseUint(v.Val, 10, 64)
		if err != nil {
			validate.getError(ErrorParamsCategoryError, validateResult, v)
		}
		validate.paramsUnit64Area(uint64(number), v, validateResult)
		break
	case "float64":
		var number float64
		number, err = strconv.ParseFloat(v.Val, 64)
		if err != nil {
			validate.getError(ErrorParamsCategoryError, validateResult, v)
		}
		validate.paramsArea(int64(number), v, validateResult)
		break
	case "float32":
		validate.validateNumberFloat(v, -2147483648, 2147483647, validateResult)
		break
	case "bool":
		err = validate.validateBoolean(v)
		if err != nil {
			validate.getError(ErrorParamsCategoryError, validateResult, v)
		}

		break
	case "boolean":
		err = validate.validateBoolean(v)
		if err != nil {
			validate.getError(ErrorParamsCategoryError, validateResult, v)
		}
		break
	default:
		break
	}

	if validateResult.Code > 0 {
		return
	}

}

//validateLeng 校验字符串长度
func (validate *Validate) validateLength(v *ValidateRule, vr *ValidateResult) {
	stringLength := int64(StringLength(v.Val))

	//如果最下值和最小值不等，则说明一定设置了最大值或最小值，或两个都设置
	////否则最大值等于最小值且不等于0，需要判断
	if (v.Min != v.Max) || (v.Min != v.Max && v.Min != 0 && v.Min != stringLength) {
		if stringLength < v.Min || stringLength > v.Max {
			validate.getError(ErrorParamsAreaInError, vr, v)
		}
	}
	//其他条件正常过
	return
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
func (validate *Validate) validateNumber(v *ValidateRule, min int64, max int64, validateResult *ValidateResult) {
	var numberValue int64
	var err error
	numberValue, err = strconv.ParseInt(v.Val, 10, 64)

	//如果数值不在区间int32区间，则报错
	if err != nil || (numberValue < min || numberValue > max) {
		log.Println(err, max)
		validateResult.Message = fmt.Sprintln(ErrorMsg[ErrorParamsCategoryPregError], v.Alias, v.Category)
		// + v.Category
		validate.getError(ErrorParamsCategoryError, validateResult, v)
		return
	}
	var numberVal = int64(numberValue)
	validate.paramsArea(numberVal, v, validateResult)
}

//校验数字
func (validate *Validate) validateNumberFloat(v *ValidateRule, min float64, max float64, validateResult *ValidateResult) {
	var numberValue float64
	var err error

	numberValue, err = strconv.ParseFloat(v.Val, 64)
	//如果数值不在区间int32区间，则报错
	if err != nil || (numberValue < min || numberValue > max) {
		validateResult.Message = fmt.Sprintln(ErrorMsg[ErrorParamsCategoryPregError], v.Alias, v.Category)
		// + v.Category
		validate.getError(ErrorParamsCategoryError, validateResult, v)
		return
	}

	var numberVal = int64(numberValue)
	//如果最下值和最小值不等，则说明一定设置了最大值或最小值，或两个都设置
	//否则最大值等于最小值且不等于0，需要判断
	if v.Min != v.Max || (v.Min == v.Max && v.Min != 0 && v.Min != numberVal) {
		if numberVal < v.Min || numberVal > v.Max {

			validateResult.Message = fmt.Sprintln(ErrorMsg[ErrorParamsAreaInError], v.Alias, strconv.FormatInt(v.Min, 10), strconv.FormatInt(v.Max, 10))
			// + v.Category
			validate.getError(ErrorParamsCategoryError, validateResult, v)
			return

		}
	}

}

//验证是否为空
func (validate *Validate) validateEmpty(v *ValidateRule, validateResult *ValidateResult) {
	//如果要验证为空
	if v.Empty == false {
		validate.getError(ErrorParamsEmpty, validateResult, v)
	}
}
