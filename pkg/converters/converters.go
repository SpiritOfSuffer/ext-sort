package converters

import (
	"fmt"
	"strconv"
)

func StringAsInt(s string) int64 {
	num, err := strconv.ParseInt(s[:len(s)-1], 10, 64) // Remove '\n'
	if err != nil {
		panic(fmt.Sprintf("Error converting string to int: %s", err.Error())) //Guarantee that input data will be stringified ints, so panic is ok
	}
	return num
}
