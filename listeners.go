package logger

// EventListener is a listener for a specific event as given by its flag.
type EventListener func(eventFlag uint64, state ...interface{})

// ErrorListener is a listener for errors.
type ErrorListener func(err error)
