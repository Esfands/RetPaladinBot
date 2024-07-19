package domain

type BTTVChannelEmoteResponse struct {
	ID            string      `json:"id"`
	Bots          []string    `json:"bots"`
	Avatar        string      `json:"avatar"`
	ChannelEmotes []BTTVEmote `json:"channelEmotes"`
	SharedEmotes  []BTTVEmote `json:"sharedEmotes"`
}

// Represents a BTTV global emote
type BTTVEmote struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	ImageType string `json:"imageType"`
	Animated  bool   `json:"animated"`
	UserID    string `json:"userId,omitempty"`
	User      struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		ProviderID  string `json:"providerId"`
	} `json:"user,omitempty"`
	Modifier bool `json:"modifier"`
}

type FFZChannelEmoteResponse struct {
	Room struct {
		InternalID     int     `json:"_id"`
		TwitchID       int     `json:"twitch_id"`
		YouTubeID      *int    `json:"youtube_id"`
		ID             string  `json:"id"`
		IsGroup        bool    `json:"is_group"`
		DisplayName    string  `json:"display_name"`
		Set            int     `json:"set"`
		ModeratorBadge string  `json:"moderator_badge"`
		VIPBadge       *string `json:"vip_badge"`
		ModURLs        struct {
			One  string `json:"1"`
			Two  string `json:"2"`
			Four string `json:"4"`
		} `json:"mod_urls"`
	} `json:"room"`

	Sets map[string]FFZSet `json:"sets"`
}

type FFZGlobalEmoteResponse struct {
	DefaultSets []int             `json:"default_sets"`
	Sets        map[string]FFZSet `json:"sets"`
}

type FFZSet struct {
	ID        int        `json:"_id"`
	Type      int        `json:"_type"`
	Icon      *string    `json:"icon"`
	Title     string     `json:"title"`
	CSS       *string    `json:"css"`
	Emoticons []FFZEmote `json:"emoticons"`
}

type FFZEmote struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Height        int     `json:"height"`
	Width         int     `json:"width"`
	Public        bool    `json:"public"`
	Hidden        bool    `json:"hidden"`
	Modifer       bool    `json:"modifer"`
	ModifierFlags int     `json:"modifier_flags"`
	Offset        *int    `json:"offset"`
	Margins       *int    `json:"margins"`
	CSS           *string `json:"css"`
	Owner         struct {
		ID          int    `json:"_id"`
		Name        string `json:"name"`
		Displayname string `json:"display_name"`
	} `json:"owner"`
	Artist *string `json:"artist"`
	URLs   struct {
		One  string `json:"1"`
		Two  string `json:"2"`
		Four string `json:"4"`
	} `json:"urls"`
	Status     int    `json:"status"`
	UsageCount int    `json:"usage_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type SevenTVChannelResponse struct {
	ID           string     `json:"id"`
	Platform     string     `json:"platform"`
	Username     string     `json:"username"`
	DisplayName  string     `json:"display_name"`
	LinkedAt     int        `json:"linked_at"`
	EmoteCapcity int        `json:"emote_capacity"`
	EmoteSetID   *int       `json:"emote_set_id"`
	EmoteSet     SevenTVSet `json:"emote_set"`
}

type SevenTVSet struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Flags      int             `json:"flags"`
	Tags       []string        `json:"tags"`
	Immutable  bool            `json:"immutable"`
	Privileged bool            `json:"privileged"`
	Emotes     []SevenTVEmote  `json:"emotes"`
	EmoteCount int             `json:"emote_count"`
	Capacity   int             `json:"capacity"`
	Owner      SevenTVSetOwner `json:"owner"`
}

type SevenTVSetOwner struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Style       struct {
		Color int `json:"color"`
	} `json:"style"`
	Roles []string `json:"roles"`
}

type SevenTVEmote struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Flags     int    `json:"flags"`
	Timestamp int    `json:"timestamp"`
	ActorID   string `json:"actor_id"`
}
