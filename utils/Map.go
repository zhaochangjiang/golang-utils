package utils

//MapMerge 合并两个Map
func MapMerge(a *map[string]interface{}, b *map[string]interface{}) *map[string]interface{} {
	for k, v := range *b {
		aValue := *a

		//判断a和b是否有相同的key
		if MapKeyIsSet(k, aValue) {
			if true == MapIsMapStringInterface(v) && true == MapIsMapStringInterface(aValue[k]) {
				aVal := aValue[k].(map[string]interface{})
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
func MapKeyIsSet(key string, mapPointer interface{}) bool {
	flag := false
	mapValues := mapPointer.(map[string]interface{})
	if _, ok := mapValues[key]; ok {
		flag = true
	}
	return flag
}
