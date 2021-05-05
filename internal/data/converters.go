package data

// ConvertBoolToInt converts a boolean to it's metric/inventory representation
func ConvertBoolToInt(val bool) (returnval int) {
	returnval = 0
	if val {
		returnval = 1
	}
	return
}
