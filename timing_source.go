package logger

import "time"

// TimingSource is a type that provides a timestamp.
type TimingSource interface {
	UTCNow() time.Time
}

// SystemClock is the an instance of the system clock timing source.
var SystemClock = timingSourceSystemClock{}

// TimingSourceSystemClock is the system clock timing source.
type timingSourceSystemClock struct{}

// UTCNow returns the current time in UTC.
func (t timingSourceSystemClock) UTCNow() time.Time {
	return time.Now().UTC()
}

// Now returns a historical time instance as a timing source.
func Now() TimeInstance {
	return TimeInstance(time.Now())
}

// TimeInstance is the system clock timing source.
type TimeInstance time.Time

// UTCNow returns the current time in UTC.
func (t TimeInstance) UTCNow() time.Time {
	return time.Time(t).UTC()
}
