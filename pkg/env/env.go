package env

import "os"

// GetEnv returns value of the name if it's exists, otherwise returns _default value
func GetEnv(name string, _default string) string {
	value, isSet := os.LookupEnv(name)
	if !isSet {
		return _default
	}
	return value
}
