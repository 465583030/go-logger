package logger

// BitFlagAll returns if all the reference bits are set for a given value
func BitFlagAll(reference, value uint64) bool {
	return reference&value == value
}

// BitFlagAny returns if any the reference bits are set for a given value
func BitFlagAny(reference, value uint64) bool {
	return reference&value > 0
}

// BitFlagCombine combines all the values into one flag.
func BitFlagCombine(values ...uint64) uint64 {
	var outputFlag uint64
	for _, value := range values {
		outputFlag = outputFlag | value
	}
	return outputFlag
}
