package util

import "strconv"

func StringToUint(str string) (uint, error) {
	u64, err := StringToUint64(str)
	if err != nil {
		return 0, err
	}
	return uint(u64), nil
}

func StringToUint64(str string) (uint64, error) {
	return strconv.ParseUint(str, 10, 64)
}
