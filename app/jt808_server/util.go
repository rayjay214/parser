package main

func calDuration(fileSize int) int {
	quotient := fileSize / 702
	remainder := fileSize % 702

	// 如果余数大于等于除数的一半，向上取整
	if remainder >= 702/2 {
		quotient++
	}
	return quotient
}
