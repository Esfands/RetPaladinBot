package subage

type SubageUser struct {
	ID          string `json:"id"`
	Login       string `json:"login"`
	DisplayName string `json:"displayName"`
}

type SubageStreak struct {
	ElapsedDays   int    `json:"elapsedDays"`
	DaysRemaining int    `json:"daysRemaining"`
	Months        int    `json:"months"`
	End           string `json:"end"`
	Start         string `json:"start"`
}

type SubageCumulative struct {
	ElapsedDays   int    `json:"elapsedDays"`
	DaysRemaining int    `json:"daysRemaining"`
	Months        int    `json:"months"`
	End           string `json:"end"`
	Start         string `json:"start"`
}

type SubageGiftMeta struct {
	GiftDate string     `json:"giftDate"`
	Gifter   SubageUser `json:"gifter"`
}

type SubageMeta struct {
	Type     string          `json:"type"`
	Tier     string          `json:"tier"`
	EndsAt   string          `json:"endsAt"`
	RenewsAt *string         `json:"renewsAt"`
	GiftMeta *SubageGiftMeta `json:"giftMeta"`
}

type SubageResponse struct {
	User         SubageUser        `json:"user"`
	Channel      SubageUser        `json:"channel"`
	StatusHidden bool              `json:"statusHidden"`
	FollowedAt   *string           `json:"followedAt"`
	Streak       *SubageStreak     `json:"streak"`
	Cumulative   *SubageCumulative `json:"cumulative"`
	Meta         *SubageMeta       `json:"meta"`
}

type SubageErrorResponse struct {
	StatusCode    int     `json:"statusCode"`
	SentryEventID *string `json:"sentryEventId"`
	RequestID     string  `json:"requestId"`
	Error         struct {
		Message string `json:"message"`
	} `json:"error"`
}
