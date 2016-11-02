package logger

//env var names
const (

	// EnvironmentVariableLogEvents is the log verbosity environment variable.
	EnvironmentVariableLogEvents = "LOG_EVENTS"

	// EnvironmentVariableUseAnsiColors is the env var that controls if we use ansi colors in output.
	EnvironmentVariableUseAnsiColors = "LOG_USE_COLOR"
	// EnvironmentVariableShowTimestamp is the env var that controls if we show timestamps in output.
	EnvironmentVariableShowTimestamp = "LOG_SHOW_TIME"
	// EnvironmentVariableShowLabel is the env var that controls if we show a descriptive label in output.
	EnvironmentVariableShowLabel = "LOG_SHOW_LABEL"
	// EnvironmentVariableLogLabel is the env var that sets the descriptive label in output.
	EnvironmentVariableLogLabel = "LOG_LABEL"

	// EnvironmentVariableLogOutFile is the variable for what file to write to.
	EnvironmentVariableLogOutFile = "LOG_OUT_FILE"
	// EnvironmentVariableLogErrFile is the variable for what file to write to for the error stream.
	EnvironmentVariableLogErrFile = "LOG_ERR_FILE"

	// EnvironmentVariableLogRedirectOutput indicates if we need to redirect output from stdout or stderr to the log files.
	EnvironmentVariableLogRedirectOutput = "LOG_REDIRECT"
)
