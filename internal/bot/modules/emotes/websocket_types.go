package emotes

import (
	"encoding/json"
	"time"
)

type Message[D AnyPayload] struct {
	Op        Opcode `json:"op"`
	Timestamp int64  `json:"t"`
	Data      D      `json:"d"`
	Sequence  uint64 `json:"s,omitempty"`
}

type Opcode uint8

const (
	// Default ops (0-32)
	OpcodeDispatch    Opcode = 0 // R - Server dispatches data to the client
	OpcodeHello       Opcode = 1 // R - Server greets the client
	OpcodeHeartbeat   Opcode = 2 // R - Keep the connection alive
	OpcodeReconnect   Opcode = 4 // R - Server demands that the client reconnects
	OpcodeAck         Opcode = 5 // R - Acknowledgement of an action
	OpcodeError       Opcode = 6 // R - Extra error context in cases where the closing frame is not enough
	OpcodeEndOfStream Opcode = 7 // R - The connection's data stream is ending

	// Commands (33-64)
	OpcodeIdentify    Opcode = 33 // S - Authenticate the session
	OpcodeResume      Opcode = 34 // S - Resume the previous session and receive missed events
	OpcodeSubscribe   Opcode = 35 // S - Subscribe to an event
	OpcodeUnsubscribe Opcode = 36 // S - Unsubscribe from an event
	OpcodeSignal      Opcode = 37 // S - Emit a spectator signal
	OpcodeBridge      Opcode = 38 // S - Send a special command
)

func (op Opcode) String() string {
	switch op {
	case OpcodeDispatch:
		return "DISPATCH"
	case OpcodeHello:
		return "HELLO"
	case OpcodeHeartbeat:
		return "HEARTBEAT"
	case OpcodeReconnect:
		return "RECONNECT"
	case OpcodeAck:
		return "ACK"
	case OpcodeError:
		return "ERROR"
	case OpcodeEndOfStream:
		return "END_OF_STREAM"

	case OpcodeIdentify:
		return "IDENTIFY"
	case OpcodeResume:
		return "RESUME"
	case OpcodeSubscribe:
		return "SUBSCRIBE"
	case OpcodeUnsubscribe:
		return "UNSUBSCRIBE"
	case OpcodeSignal:
		return "SIGNAL"
	case OpcodeBridge:
		return "BRIDGE"
	default:
		return "UNDOCUMENTED_OPERATION"
	}
}

type CloseCode uint16

const (
	CloseCodeServerError           CloseCode = 4000 // an error occured on the server's end
	CloseCodeUnknownOperation      CloseCode = 4001 // the client sent an unexpected opcode
	CloseCodeInvalidPayload        CloseCode = 4002 // the client sent a payload that couldn't be decoded
	CloseCodeAuthFailure           CloseCode = 4003 // the client unsucessfully tried to identify
	CloseCodeAlreadyIdentified     CloseCode = 4004 // the client wanted to identify again
	CloseCodeRateLimit             CloseCode = 4005 // the client is being rate-limited
	CloseCodeRestart               CloseCode = 4006 // the server is restarting and the client should reconnect
	CloseCodeMaintenance           CloseCode = 4007 // the server is in maintenance mode and not accepting connections
	CloseCodeTimeout               CloseCode = 4008 // the client was idle for too long
	CloseCodeAlreadySubscribed     CloseCode = 4009 // the client tried to subscribe to an event twice
	CloseCodeNotSubscribed         CloseCode = 4010 // the client tried to unsubscribe from an event they weren't subscribing to
	CloseCodeInsufficientPrivilege CloseCode = 4011 // the client did something that they did not have permission for
	CloseCodeReconnect             CloseCode = 4012 // the client should reconnect
)

func (c CloseCode) String() string {
	switch c {
	case CloseCodeServerError:
		return "Internal Server Error"
	case CloseCodeUnknownOperation:
		return "Unknown Operation"
	case CloseCodeInvalidPayload:
		return "Invalid Payload"
	case CloseCodeAuthFailure:
		return "Authentication Failed"
	case CloseCodeAlreadyIdentified:
		return "Already identified"
	case CloseCodeRateLimit:
		return "Rate limit reached"
	case CloseCodeRestart:
		return "Server is restarting"
	case CloseCodeMaintenance:
		return "Maintenance Mode"
	case CloseCodeTimeout:
		return "Timeout"
	case CloseCodeAlreadySubscribed:
		return "Already Subscribed"
	case CloseCodeNotSubscribed:
		return "Not Subscribed"
	case CloseCodeInsufficientPrivilege:
		return "Insufficient Privilege"
	case CloseCodeReconnect:
		return "Reconnect"
	default:
		return "Undocumented Closure"
	}
}

type AnyPayload interface {
	json.RawMessage | HelloPayload | AckPayload | HeartbeatPayload | ReconnectPayload | ResumePayload |
		SubscribePayload | UnsubscribePayload | DispatchPayload | SignalPayload |
		ErrorPayload | EndOfStreamPayload
}

type HelloPayload struct {
	HeartbeatInterval uint32                   `json:"heartbeat_interval"`
	SessionID         string                   `json:"session_id"`
	SubscriptionLimit int32                    `json:"subscription_limit"`
	Instance          HelloPayloadInstanceInfo `json:"instance,omitempty"`
}

type HelloPayloadInstanceInfo struct {
	Name       string `json:"name"`
	Population int32  `json:"population"`
}

type AckPayload struct {
	Command string          `json:"command"`
	Data    json.RawMessage `json:"data"`
}

type HeartbeatPayload struct {
	Count uint64 `json:"count"`
}

type ReconnectPayload struct {
	Reason string `json:"reason"`
}

type ResumePayload struct {
	SessionID string `json:"session_id"`
}

type SubscribePayload struct {
	Type      EventType         `json:"type"`
	Condition map[string]string `json:"condition"`
	TTL       time.Duration     `json:"ttl,omitempty"`
}

type UnsubscribePayload struct {
	Type      EventType         `json:"type"`
	Condition map[string]string `json:"condition"`
}

type DispatchPayload struct {
	Type EventType `json:"type"`
	// Hash is a hash of the target object, used for deduping
	Hash *uint32 `json:"hash,omitempty"`
	// A list of subscriptions that match this dispatch
	Matches []uint32 `json:"matches,omitempty"`
	// A list of conditions where at least one must have all its fields match a subscription in order for this dispatch to be delivered
	Conditions []EventCondition `json:"condition,omitempty"`
	// This dispatch is a whisper to a specific session, usually as a response to a command
	Whisper string `json:"whisper,omitempty"`
}

type Event struct {
	Type string `json:"type"`
	Body Body   `json:"body"`
}

type Body struct {
	ID      string   `json:"id"`
	Kind    int      `json:"kind"`
	Actor   Actor    `json:"actor"`
	Pulled  []Pulled `json:"pulled"`
	Pushed  []Pushed `json:"pushed"`
	Removed []Pushed `json:"removed"`
	Updated []Pushed `json:"updated"`
}

type Actor struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Username    string       `json:"username"`
	DisplayName string       `json:"display_name"`
	AvatarURL   string       `json:"avatar_url"`
	Style       Style        `json:"style"`
	Roles       []string     `json:"roles"`
	Connections []Connection `json:"connections"`
}

type Style struct {
	Color   int     `json:"color"`
	PaintID *string `json:"paint_id"`
	BadgeID *string `json:"badge_id"`
	Paint   *string `json:"paint"`
	Badge   *string `json:"badge"`
}

type Connection struct {
	ID            string `json:"id"`
	Platform      string `json:"platform"`
	Username      string `json:"username"`
	DisplayName   string `json:"display_name"`
	LinkedAt      int64  `json:"linked_at"`
	EmoteCapacity int    `json:"emote_capacity"`
	EmoteSetID    string `json:"emote_set_id"`
}

type Pulled struct {
	Key      string  `json:"key"`
	Index    int     `json:"index"`
	Type     string  `json:"type"`
	OldValue Emote   `json:"old_value"`
	Value    *string `json:"value"` // Value is null when an emote is removed
}

type Pushed struct {
	Key      string `json:"key"`
	Index    int    `json:"index"`
	Type     string `json:"type"`
	OldValue Emote  `json:"old_value"`
	Value    Value  `json:"value"`
}

type Emote struct {
	ActorID   string `json:"actor_id"`
	Flags     int    `json:"flags"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

type Value struct {
	ActorID   string `json:"actor_id"`
	Data      Data   `json:"data"`
	Flags     int    `json:"flags"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Timestamp int64  `json:"timestamp"`
}

type Data struct {
	Animated  bool     `json:"animated"`
	Flags     int      `json:"flags"`
	Host      Host     `json:"host"`
	ID        string   `json:"id"`
	Lifecycle int      `json:"lifecycle"`
	Listed    bool     `json:"listed"`
	Name      string   `json:"name"`
	Owner     Owner    `json:"owner"`
	State     []string `json:"state"`
}

type Host struct {
	Files []File `json:"files"`
	URL   string `json:"url"`
}

type File struct {
	Format     string `json:"format"`
	FrameCount int    `json:"frame_count"`
	Height     int    `json:"height"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	StaticName string `json:"static_name"`
	Width      int    `json:"width"`
}

type Owner struct {
	AvatarURL   string `json:"avatar_url"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Style       Style  `json:"style"`
	Username    string `json:"username"`
}

type EventCondition map[string]string

func (evc EventCondition) Set(key string, value string) EventCondition {
	evc[key] = value

	return evc
}

func (evc EventCondition) Match(other EventCondition) bool {
	for k, v := range other {
		if evc[k] != v {
			return false
		}
	}

	return true
}

type SignalPayload struct {
	Sender SignalUser `json:"sender"`
	Host   SignalUser `json:"host"`
}

type SignalUser struct {
	ChannelID   string `json:"channel_id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
}

type ErrorPayload struct {
	Message       string         `json:"message"`
	MessageLocale string         `json:"message_locale,omitempty"`
	Fields        map[string]any `json:"fields"`
}

type EndOfStreamPayload struct {
	Code    CloseCode `json:"code"`
	Message string    `json:"message"`
}

type EventType string

const (
	// System

	EventTypeAnySystem          EventType = "system.*"
	EventTypeSystemAnnouncement EventType = "system.announcement"

	// Emote

	EventTypeAnyEmote    EventType = "emote.*"
	EventTypeCreateEmote EventType = "emote.create"
	EventTypeUpdateEmote EventType = "emote.update"
	EventTypeDeleteEmote EventType = "emote.delete"

	// Emote Set

	EventTypeAnyEmoteSet    EventType = "emote_set.*"
	EventTypeCreateEmoteSet EventType = "emote_set.create"
	EventTypeUpdateEmoteSet EventType = "emote_set.update"
	EventTypeDeleteEmoteSet EventType = "emote_set.delete"

	// User

	EventTypeAnyUser    EventType = "user.*"
	EventTypeCreateUser EventType = "user.create"
	EventTypeUpdateUser EventType = "user.update"
	EventTypeDeleteUser EventType = "user.delete"

	EventTypeAnyEntitlement    EventType = "entitlement.*"
	EventTypeCreateEntitlement EventType = "entitlement.create"
	EventTypeUpdateEntitlement EventType = "entitlement.update"
	EventTypeDeleteEntitlement EventType = "entitlement.delete"

	// Cosmetics

	EventTypeAnyCosmetic    EventType = "cosmetic.*"
	EventTypeCreateCosmetic EventType = "cosmetic.create"
	EventTypeUpdateCosmetic EventType = "cosmetic.update"
	EventTypeDeleteCosmetic EventType = "cosmetic.delete"

	// Special

	EventTypeWhisper EventType = "whisper.self"
)
