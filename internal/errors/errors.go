package errors

import "errors"

var ErrorEnvVariableRequired = errors.New("environment variable is required")
var ErrorInvalidDSN = errors.New("invalid DSN")
