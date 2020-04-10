package bytepal

// BoolToBinary converts a boolean to it's Unary operator
func BoolToBinary(value bool) uint {
	if value {
		return 1
	}
	return 0
}

// ByteToBool takes a byte value and turns it into a boolean. If the value is zero it is false.
func ByteToBool(value uint8) bool {
	return value != 0
}