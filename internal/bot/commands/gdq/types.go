package gdq

type GDQResponse struct {
	Event     int    `json:"event"`
	Date      string `json:"date"`
	Comment   string `json:"comment"`
	EventName string `json:"eventName"`
}
