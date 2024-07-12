package song

type Response struct {
	Recenttracks Recenttracks `json:"recenttracks"`
}

type Recenttracks struct {
	Track []Track `json:"track"`
	Attr  Attr    `json:"@attr"`
}

type Track struct {
	Artist     Artist      `json:"artist"`
	Streamable string      `json:"streamable"`
	Image      []ImageSize `json:"image"`
	MBid       string      `json:"mbid"`
	Album      Album       `json:"album"`
	Name       string      `json:"name"`
	URL        string      `json:"url"`
	Date       Date        `json:"date"`
}

type Artist struct {
	MBid string `json:"mbid"`
	Text string `json:"#text"`
}

type ImageSize struct {
	Size string `json:"size"`
	Text string `json:"#text"`
}

type Album struct {
	MBid string `json:"mbid"`
	Text string `json:"#text"`
}

type Date struct {
	Uts  string `json:"uts"`
	Text string `json:"#text"`
}

type Attr struct {
	User       string `json:"user"`
	TotalPages string `json:"totalPages"`
	Page       string `json:"page"`
	PerPage    string `json:"perPage"`
	Total      string `json:"total"`
}
