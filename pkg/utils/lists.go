package utils

func UniqueList[T comparable](list []T) []T {
	keys := make(map[T]bool)
	var uniqueList []T

	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniqueList = append(uniqueList, entry)
		}
	}

	return uniqueList
}

func InList[T comparable](list []T, item T) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
