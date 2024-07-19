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
	Type string    `json:"type"`
	Body ChangeMap `json:"body"`
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

type ChangeMap struct {
	// The object's ID
	ID string `json:"id"`
	// The type of the object
	Kind ObjectKind `json:"kind"`
	// Contextual is whether or not this event is only relating
	// to the specific source conditions and not indicative of a
	// genuine creation, deletion, or update to the object
	Contextual bool `json:"contextual,omitempty"`
	// The user who made changes to the object
	Actor UserPartialModel `json:"actor,omitempty"`
	// A list of added fields
	Added []ChangeField `json:"added,omitempty"`
	// A list of updated fields
	Updated []ChangeField `json:"updated,omitempty"`
	// A list of removed fields
	Removed []ChangeField `json:"removed,omitempty"`
	// A list of items pushed to an array
	Pushed []ChangeField `json:"pushed,omitempty"`
	// A list of items pulled from an array
	Pulled []ChangeField `json:"pulled,omitempty"`
	// A full object. Only available during a "create" event
	Object json.RawMessage `json:"object,omitempty"`
}

type ChangeField struct {
	Key      string          `json:"key"`
	Index    *int32          `json:"index"`
	Nested   bool            `json:"nested,omitempty"`
	Type     ChangeFieldType `json:"type"`
	OldValue any             `json:"old_value,omitempty"`
	Value    any             `json:"value"`
}

type ChangeFieldType string

const (
	ChangeFieldTypeString ChangeFieldType = "string"
	ChangeFieldTypeNumber ChangeFieldType = "number"
	ChangeFieldTypeBool   ChangeFieldType = "bool"
	ChangeFieldTypeObject ChangeFieldType = "object"
)

type UserPartialModel struct {
	ID          string                       `json:"id"`
	UserType    UserTypeModel                `json:"type,omitempty" enums:",BOT,SYSTEM"`
	Username    string                       `json:"username"`
	DisplayName string                       `json:"display_name"`
	AvatarURL   string                       `json:"avatar_url,omitempty" extensions:"x-omitempty"`
	Style       UserStyle                    `json:"style"`
	RoleIDs     []string                     `json:"roles,omitempty" extensions:"x-omitempty"`
	Connections []UserConnectionPartialModel `json:"connections,omitempty" extensions:"x-omitempty"`
}

type UserStyle struct {
	Color   int32   `json:"color,omitempty" extensions:"x-omitempty"`
	PaintID *string `json:"paint_id,omitempty" extensions:"x-omitempty"`
}

type UserConnectionPartialModel struct {
	ID string `json:"id"`
	// The service of the connection.
	Platform UserConnectionPlatformModel `json:"platform" enums:"TWITCH,YOUTUBE,DISCORD"`
	// The username of the user on the platform.
	Username string `json:"username"`
	// The display name of the user on the platform.
	DisplayName string `json:"display_name"`
	// The time when the user linked this connection
	LinkedAt int64 `json:"linked_at"`
	// The maximum size of emote sets that may be bound to this connection.
	EmoteCapacity int32 `json:"emote_capacity"`
	// The emote set that is linked to this connection
	EmoteSetID *string `json:"emote_set_id" extensions:"x-nullable"`
}

type UserConnectionPlatformModel string

var (
	UserConnectionPlatformTwitch  UserConnectionPlatformModel = "TWITCH"
	UserConnectionPlatformYouTube UserConnectionPlatformModel = "YOUTUBE"
	UserConnectionPlatformDiscord UserConnectionPlatformModel = "DISCORD"
)

type UserTypeModel string

var (
	UserTypeRegular UserTypeModel = ""
	UserTypeBot     UserTypeModel = "BOT"
	UserTypeSystem  UserTypeModel = "SYSTEM"
)

type ObjectID = string

type ObjectKind int8

const (
	ObjectKindUser        ObjectKind = 1
	ObjectKindEmote       ObjectKind = 2
	ObjectKindEmoteSet    ObjectKind = 3
	ObjectKindRole        ObjectKind = 4
	ObjectKindEntitlement ObjectKind = 5
	ObjectKindBan         ObjectKind = 6
	ObjectKindMessage     ObjectKind = 7
	ObjectKindReport      ObjectKind = 8
	ObjectKindPresence    ObjectKind = 9
	ObjectKindCosmetic    ObjectKind = 10
)
