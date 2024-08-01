package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"
	"unsafe"

	"github.com/esfands/retpaladinbot/pkg/domain"
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

// B2S converts byte slice to a string without memory allocation.
// See https://groups.google.com/forum/#!msg/Golang-Nuts/ENgbUzYvCuU/90yGx7GUAgAJ .
//
// Note it may break if string and/or slice header will change
// in the future go versions.
func B2S(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

// S2B converts string to a byte slice without memory allocation.
// Note it may break if string and/or slice header will change
// in the future go versions.
func S2B(s string) (b []byte) {
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
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

// Util - Ternary:
// A golang equivalent to JS Ternary Operator
//
// It takes a condition, and returns a result depending on the outcome
func Ternary[T any](condition bool, whenTrue T, whenFalse T) T {
	if condition {
		return whenTrue
	}

	return whenFalse
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

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ConvertSliceToJSONString converts a slice of strings to a JSON-compatible string array
func ConvertSliceToJSONString(slice []string) string {
	// Quote each string in the slice
	for i, str := range slice {
		slice[i] = fmt.Sprintf(`"%s"`, str)
	}
	// Join the quoted elements with a comma and space, and wrap with square brackets
	return fmt.Sprintf("[%s]", strings.Join(slice, ", "))
}

// ConvertJSONStringToSlice converts a JSON string array to a slice of strings
func ConvertJSONStringToSlice(jsonStr string) ([]string, error) {
	var slice []string
	err := json.Unmarshal([]byte(jsonStr), &slice)
	if err != nil {
		return nil, err
	}
	return slice, nil
}

// ParseBadgePermissions parses the badges of a user and parses the permissions from their Twitch badges
func ParseBadges(badges map[string]int) []domain.Permission {
	var permissions []domain.Permission

	for badge := range badges {
		switch badge {
		case "broadcaster":
			permissions = append(permissions, domain.PermissionBroadcaster)
		case "moderator":
			permissions = append(permissions, domain.PermissionModerator)
		case "vip":
			permissions = append(permissions, domain.PermissionVIP)
		}
	}

	return permissions
}

// ConvertPermissionsToStrings converts a slice of Permission to a slice of string
func ConvertPermissionsToStrings(permissions []domain.Permission) []string {
	var result []string
	for _, permission := range permissions {
		result = append(result, string(permission))
	}
	return result
}

func GetRandomStringFromSlice(slice []string) string {
	// Seed the random number generator
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Get a random index from the slice
	randomIndex := rng.Intn(len(slice))

	// Return the string at the random index
	return slice[randomIndex]
}
