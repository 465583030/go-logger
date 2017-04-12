package logger

// SyncAgent is an agent that fires events synchronously.
// It wraps a regular agent.
type SyncAgent struct {
	a *Agent
}

// Infof logs an informational message to the output stream.
func (sa *SyncAgent) Infof(format string, args ...interface{}) {
	if sa == nil {
		return
	}
	sa.WriteEventf(EventInfo, ColorWhite, format, args...)
}

// Debugf logs a debug message to the output stream.
func (sa *SyncAgent) Debugf(format string, args ...interface{}) {
	if sa == nil {
		return
	}
	sa.WriteEventf(EventDebug, ColorLightYellow, format, args...)
}

// WriteEventf writes to the standard output and triggers events.
func (sa *SyncAgent) WriteEventf(event EventFlag, color AnsiColorCode, format string, args ...interface{}) {
	if sa == nil {
		return
	}
	if sa.a == nil {
		return
	}
	if sa.a.IsEnabled(event) {
		sa.a.write(append([]interface{}{TimeNow(), event, ColorLightYellow, format}, args...)...)

		if sa.a.HasListener(event) {
			sa.a.triggerListeners(append([]interface{}{TimeNow(), event, format}, args...)...)
		}
	}
}

// WriteErrorEventf writes to the error output and triggers events.
func (sa *SyncAgent) WriteErrorEventf(event EventFlag, color AnsiColorCode, format string, args ...interface{}) {
	if sa == nil {
		return
	}
	if sa.a == nil {
		return
	}
	if sa.a.IsEnabled(event) {
		sa.a.writeError(append([]interface{}{TimeNow(), event, ColorLightYellow, format}, args...)...)

		if sa.a.HasListener(event) {
			sa.a.triggerListeners(append([]interface{}{TimeNow(), event, format}, args...)...)
		}
	}
}

// ErrorEventWithState writes an error and triggers events with a given state.
func (sa *SyncAgent) ErrorEventWithState(event EventFlag, color AnsiColorCode, err error, state ...interface{}) error {
	if sa == nil {
		return err
	}
	if sa.a == nil {
		return err
	}
	if err != nil {
		if sa.a.IsEnabled(event) {
			sa.a.writeError(TimeNow(), event, ColorLightYellow, "%+v", err)
			if sa.a.HasListener(event) {
				sa.a.triggerListeners(append([]interface{}{TimeNow(), event, err}, state...)...)
			}
		}
	}
	return err
}
