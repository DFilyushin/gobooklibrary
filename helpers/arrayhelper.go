package helpers


func CheckItemInArray(item string, arr *[]string) bool {
	for _, arrItem := range *arr {
		if arrItem == item {
			return true
		}
	}
	return false
}
