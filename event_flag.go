package logger

const (
	// EventNone is effectively logging disabled.
	EventNone = uint64(0)
	// EventAll represents every flag being enabled.
	EventAll = ^EventNone
	// EventFatalError enables logging errors
	EventFatalError = 1 << iota
	// EventError enables logging errors
	EventError = 1 << iota
	// EventWarning enables logging for warning messages.
	EventWarning = 1 << iota
	// EventDebug enables logging for debug messages.
	EventDebug = 1 << iota
	// EventInfo enables logging for informational messages.
	EventInfo = 1 << iota

	// EventUserError enables output for user error events.
	EventUserError = 1 << iota

	// EventRequest is a helper event for logging request events.
	EventRequest = 1 << iota
	// EventRequestComplete is a helper event for logging request events with stats.
	EventRequestComplete = 1 << iota
	// EventRequestBody is a helper event for logging incoming post bodies.
	EventRequestBody = 1 << iota

	// EventResponseBody is a helper event for logging response bodies.
	EventResponseBody = 1 << iota
)

// EventFlagAll returns if all the reference bits are set for a given value
func EventFlagAll(reference, value uint64) bool {
	return reference&value == value
}

// EventFlagAny returns if any the reference bits are set for a given value
func EventFlagAny(reference, value uint64) bool {
	return reference&value > 0
}

// EventFlagCombine combines all the values into one flag.
func EventFlagCombine(values ...uint64) uint64 {
	var outputFlag uint64
	for _, value := range values {
		outputFlag = outputFlag | value
	}
	return outputFlag
}
