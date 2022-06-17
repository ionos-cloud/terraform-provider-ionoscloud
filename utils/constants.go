package utils

import "time"

// IonosDebug - env variable, set to true to enable debug
const IonosDebug = "IONOS_DEBUG"

// MaxRetries - number of retries in case of rate-limit
const MaxRetries = 999

// MaxWaitTime - waits 4 seconds before retry in case of rate limit
const MaxWaitTime = 4 * time.Second
