package utils

import "strconv"

// Int32Pointer Get int32 pointer
func Int32Pointer(i int32) *int32 {
	return &i
}

// StrToUint String to Uint parser
func StrToUint(value string) (uint, error) {
	u64, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	result := uint(u64)
	return result, nil
}
