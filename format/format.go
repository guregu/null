package format

import "time"

var timeFormat = time.RFC3339Nano

// SetTimeFormat sets the time format used for both null and zero times
// Not thread safe (in the interest of speed) so make sure to call only once
// does not check if f is a valid format or not
func SetTimeFormat(f string) {
	timeFormat = f
}

// GetTimeFormat gets the format for the class
// Not thread safe
func GetTimeFormat() string {
	return timeFormat
}
