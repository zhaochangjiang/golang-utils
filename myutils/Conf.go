package myutils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//GetConf 获得配置文件内容
type GetConf struct {
	Conf   *map[string]string
	suffix string
}

//NewGetConf 初始化 NewGetConf
func NewGetConf() *GetConf {
	returnData := &GetConf{}
	config := make(map[string]string)
	returnData.Conf = &config
	return returnData
}

//GetAllConfig 获得读取的配置文件所有的内容
func (gfc *GetConf) GetAllConfig() *map[string]string {
	return gfc.Conf
}

//inArray 判断某个字符串是否在数组内
func (gfc *GetConf) inArray(extString string, ext *[]string) bool {
	flag := false
	for _, v := range *ext {
		if v == extString {
			flag = true
			break
		}
	}
	return flag
}

//获得目录下文件列表
func (gfc *GetConf) getFileList(dir string, ext []string) *[]string {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	fileName := make([]string, 0)
	for _, f := range files {
		fileNameString := f.Name()
		extString := path.Ext(fileNameString)
		if extString == "" {
			continue
		}
		extString = StringSub(extString, 1)
		if gfc.inArray(extString, &ext) == false {
			continue
		}

		fi, err := os.Stat(dir + fileNameString)
		if nil != err {
			panic(err)
		}
		if !fi.IsDir() {
			fileName = append(fileName, fileNameString)
		}
	}
	return &fileName
}

//SetPathString 设置获取配置的前缀
func (gfc *GetConf) SetPathString(suffix string) *GetConf {
	gfc.suffix = suffix
	return gfc
}

//InitConf 获得配置文件信息
func (gfc *GetConf) InitConf(dir string) *GetConf {
	dir = gfc.getCurrentDirectory() + dir
	list := gfc.getFileList(dir, []string{"ini"})
	for _, v := range *list {

		//逐行获取ini文件 非注释内容
		gfc.getFileContentLineByLine(dir + v)
	}
	return gfc
}

//Get 获取指定
func (gfc *GetConf) Get() *map[string]string {
	res := make(map[string]string)
	keyLen := len([]rune(gfc.suffix))
	for k, v := range *gfc.Conf {
		if strings.Index(k, gfc.suffix) == 0 {
			cut := StringSub(k, keyLen)
			if strings.Index(cut, ".") == 0 {
				cut = StringSub(cut, 1)
			}
			res[cut] = v
		}
	}
	return &res
}

//GetFileContent 获得文件内容
func (gfc *GetConf) getFileContentLineByLine(file string) {
	fi, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("Error: %s\n", err))
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {

		a, _, c := br.ReadLine()

		//如果到最后一行了，退出
		if c == io.EOF {
			break
		}

		//去掉注释字符串
		lineString := string(a)
		if strings.Index(lineString, ";") == 0 || strings.TrimSpace(lineString) == "" {
			continue
		}

		stringSplit := strings.Split(lineString, "=")
		if lenSlice := len(stringSplit); lenSlice < 2 {
			stringSplit = append(stringSplit, "")
		}
		stringSplit[0] = strings.TrimSpace(stringSplit[0])
		stringSplit[1] = strings.TrimSpace(stringSplit[1])
		(*(gfc.Conf))[stringSplit[0]] = stringSplit[1]
	}
}

//getCurrentDirectory 获得当前目录
func (gfc *GetConf) getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
