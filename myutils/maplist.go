package myutils

import (
	"container/list"
)

type Keyer interface {
	GetKey() string
}

//MapList 有序map的实现
type MapList struct {
	dataMap  map[string]*list.Element
	dataList *list.List
}

//NewMapList 实例化一个有序map
func NewMapList() *MapList {
	return &MapList{
		dataMap:  make(map[string]*list.Element),
		dataList: list.New(),
	}
}

//Exists
func (mapList *MapList) Exists(data Keyer) bool {
	_, exists := mapList.dataMap[string(data.GetKey())]
	return exists
}

//Push
func (mapList *MapList) Push(data Keyer) bool {
	if mapList.Exists(data) {
		return false
	}
	elem := mapList.dataList.PushBack(data)
	mapList.dataMap[data.GetKey()] = elem
	return true
}

//Remove
func (mapList *MapList) Remove(data Keyer) {
	if !mapList.Exists(data) {
		return
	}
	mapList.dataList.Remove(mapList.dataMap[data.GetKey()])
	delete(mapList.dataMap, data.GetKey())
}

//Size
func (mapList *MapList) Size() int {
	return mapList.dataList.Len()
}

//Walk
func (mapList *MapList) Walk(cb func(data Keyer)) {
	for elem := mapList.dataList.Front(); elem != nil; elem = elem.Next() {
		cb(elem.Value.(Keyer))
	}
}

//Elements
type Elements struct {
	value string
}

//GetKey 实现interface Keyer方法
func (e Elements) GetKey() string {
	return e.value
}
