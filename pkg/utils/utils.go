package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
)

func GetTarget(user twitch.User, context []string) string {
	var tagged string
	if len(context) > 0 {
		tagged = context[0]
	}

	if tagged == "" {
		tagged = user.Name
	}

	tagged = strings.TrimPrefix(tagged, "@")

	return strings.ToLower(tagged)
}

// IsEmptyValue uses reflection to determine if a value is empty.
func IsEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// TimeDifference returns the difference between two times in a human readable format
//
// If abbreviate is true, the output will be abbreviated
func TimeDifference(start, end time.Time, abbreviate bool) string {
	duration := end.Sub(start)
	if duration < 0 {
		_ = -duration
		start, end = end, start
	}

	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	days := end.Day() - start.Day()

	if days < 0 {
		months--
		daysInMonth := time.Date(end.Year(), end.Month()-1, 0, 0, 0, 0, 0, time.UTC).Day()
		days += daysInMonth
	}

	if months < 0 {
		years--
		months += 12
	}

	daysDuration := end.AddDate(-years, -months, -days).Sub(start)
	hours := int(daysDuration.Hours())
	minutes := int(daysDuration.Minutes()) % 60
	seconds := int(daysDuration.Seconds()) % 60

	labels := []string{"year", "month", "day", "hour", "minute", "second"}
	if abbreviate {
		labels = []string{"y", "mo", "d", "hr", "min", "sec"}
	}

	times := []int{years, months, days, hours, minutes, seconds}
	var parts []string

	// Track the count of non-zero time units we've added
	count := 0
	for i := 0; i < len(times); i++ {
		if times[i] > 0 {
			label := labels[i]
			if times[i] > 1 {
				label += "s"
			}
			parts = append(parts, fmt.Sprintf("%d %s", times[i], label))
			count++
		}
		// We stop adding time units once we've added 3 non-zero time units
		if count == 3 {
			break
		}
	}

	if len(parts) == 0 {
		return fmt.Sprintf("0 %s", labels[len(labels)-1]+"s")
	}

	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		parts = parts[:len(parts)-1]
		parts = append(parts, "and "+lastPart)
	}

	return strings.Join(parts, ", ")
}
