package common

func InStringList(target string, strArray []string) bool {
	for _, v := range strArray {
		if target == v {
			return true
		}
	}
	return false
}
