package xplorersbot

func ArrayContainsItem(array []string, item string) bool {
	for _, arrayItem := range array {
		if item == arrayItem {
			return true
		}
	}
	return false
}
