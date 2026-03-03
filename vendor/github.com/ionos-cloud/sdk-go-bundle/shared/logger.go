/*
 * IONOS Shared Libraries
 */

package shared

import (
	"log"
	"os"
	"strings"
)

// creates default logger and gets logLevel from environment
func init() {
	NewSdkLogger()
	SdkLogLevel = getLogLevelFromEnv()
}

type LogLevel uint

func (l *LogLevel) Get() LogLevel {
	if l != nil {
		return *l
	}
	return Off
}

// Satisfies returns true if this LogLevel is at least high enough for v
func (l *LogLevel) Satisfies(v LogLevel) bool {
	return l.Get() >= v
}

const (
	Off LogLevel = 0x100 * iota
	Debug
	// Trace We recommend you only set this field for debugging purposes.
	// Disable it in your production environments because it can log sensitive data.
	// It logs the full request and response without encryption, even for an HTTPS call.
	// Verbose request and response logging can also significantly impact your application's performance.
	Trace
)

var LogLevelMap = map[string]LogLevel{
	"off":   Off,
	"debug": Debug,
	"trace": Trace,
}

// getLogLevelFromEnv - gets LogLevel type from env variable IONOS_LOGLEVEL
// returns Off if an invalid log level is encountered
func getLogLevelFromEnv() LogLevel {
	strLogLevel := "off"

	if logLevelFromEnv, isSet := os.LookupEnv(IonosLogLevelEnvVar); isSet {
		strLogLevel = logLevelFromEnv
	}

	logLevel, ok := LogLevelMap[strings.ToLower(strLogLevel)]
	if !ok {
		log.Printf("Cannot set logLevel for value: %s, setting loglevel to Off", strLogLevel)
	}
	return logLevel
}

var SdkLogger Logger
var SdkLogLevel LogLevel

type Logger interface {
	Printf(format string, args ...interface{})
}

func NewSdkLogger() {
	SdkLogger = log.New(os.Stderr, "IONOSLOG ", log.LstdFlags)
}

// LogDebug logs a message at Debug level. It is a no-op when the SDK log
// level is below Debug.
func LogDebug(format string, args ...interface{}) {
	if SdkLogLevel.Satisfies(Debug) {
		SdkLogger.Printf(format, args...)
	}
}

// LogTrace logs a message at Trace level. It is a no-op when the SDK log
// level is below Trace.
func LogTrace(format string, args ...interface{}) {
	if SdkLogLevel.Satisfies(Trace) {
		SdkLogger.Printf(format, args...)
	}
}
