package schedule

import "time"

// CustomTime handles time fields that may be empty strings.
type CustomTime struct {
	time.Time
}

// UnmarshalJSON handles the deserialization of CustomTime.
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == `""` {
		// Empty string case
		return nil
	}

	// Remove the surrounding quotes
	str = str[1 : len(str)-1]

	// Parse the time
	parsedTime, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}

type Schedule struct {
	Data struct {
		Segments []Segment `json:"segments"`
	} `json:"data"`
}

type Segment struct {
	ID             string      `json:"id"`
	StartTime      CustomTime  `json:"start_time"`
	EndTime        *CustomTime `json:"end_time,omitempty"` // Use a pointer to handle optional end_time
	Title          string      `json:"title"`
	CancelledUntil string      `json:"cancelled_until,omitempty"`
	Category       Category    `json:"category"`
	IsRecurring    bool        `json:"is_recurring"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
