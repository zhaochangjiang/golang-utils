package utils

import (
	"sort"
)

//MapMerge 合并两个Map
func MapMerge(a *map[string]interface{}, b *map[string]interface{}) *map[string]interface{} {
	for k, v := range *b {

		//判断a和b是否有相同的key
		if MapKeyIsSet(k, a) {
			if true == MapIsMapStringInterface(v) && true == MapIsMapStringInterface((*a)[k]) {
				aVal := (*a)[k].(map[string]interface{})
				bVal := v.(map[string]interface{})
				MapMerge(&aVal, &bVal)
			} else {
				(*a)[k] = v
			}
		} else {
			(*a)[k] = v
		}
	}
	return a
}

//MapIsMapStringInterface 判断是否为 map[string]interface类型
func MapIsMapStringInterface(mapContent interface{}) bool {
	res := false
	if _, ok := mapContent.(map[string]interface{}); ok {
		res = true
	}
	return res
}

//MapKeyIsSet 判断map 是否存在key
func MapKeyIsSet(key string, mapValues *map[string]interface{}) bool {
	flag := false
	if _, ok := (*mapValues)[key]; ok {
		flag = true
	}
	return flag
}

//MapSortByKeyString map按照键排序
func MapSortByKeyString(mapPointer *map[string]interface{}) *map[string]interface{} {

	var keys []string
	for k := range *mapPointer {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	res := make(map[string]interface{})
	for _, k := range keys {
		res[k] = (*mapPointer)[k]
	}
	return &res
}

//MapSortByKeyInt map按照键排序
func MapSortByKeyInt(mapPointer *map[int]interface{}) *map[int]interface{} {

	var keys []int
	for k := range *mapPointer {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	res := make(map[int]interface{})
	for _, k := range keys {
		res[k] = (*mapPointer)[k]
	}
	return &res
}
