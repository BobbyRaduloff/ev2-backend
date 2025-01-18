package utils

func RemoveDuplicates[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func SplitArray(arr []string, maxSize int) [][]string {
	if maxSize <= 0 {
		return nil
	}

	var result [][]string
	for i := 0; i < len(arr); i += maxSize {
		end := i + maxSize
		if end > len(arr) {
			end = len(arr)
		}
		result = append(result, arr[i:end])
	}

	return result
}

func ChunkArray[T any](arr []T, maxSize int) [][]T {
	if maxSize <= 0 {
		return nil
	}

	var chunks [][]T
	for maxSize < len(arr) {
		arr, chunks = arr[maxSize:], append(chunks, arr[0:maxSize:maxSize])
	}
	return append(chunks, arr)
}
