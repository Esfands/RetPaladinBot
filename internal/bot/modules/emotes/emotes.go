package emotes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/dghubble/sling"
	"github.com/esfands/retpaladinbot/internal/global"
	"github.com/esfands/retpaladinbot/pkg/domain"
	"nhooyr.io/websocket"
)

type EmoteModule struct {
	BTTVGlobalEmotes    []string
	FFZGlobalEmotes     []string
	SevenTVGlobalEmotes []string

	BTTVChannelEmotes    []string
	FFZChannelEmotes     []string
	SevenTVChannelEmotes []string

	wsConn *websocket.Conn
}

// NewEmoteModule creates a new EmoteModule and initializes the global and channel emotes.
func NewEmoteModule(gctx global.Context, channelID string) (*EmoteModule, error) {
	em := &EmoteModule{}
	var err error

	err = em.InitializeGlobalEmotes()
	if err != nil {
		return nil, err
	}

	err = em.InitializeChannelEmotes(channelID)
	if err != nil {
		return nil, err
	}

	fmt.Println(em.GetAllEmotes())

	// Initialize WebSocket client
	err = em.initializeWebSocketClient(gctx)
	if err != nil {
		return nil, err
	}

	// Schedule updates for global and channel emotes
	gctx.Crate().Scheduler.Scheduler().Every(12).Hours().Do(func() {
		err = em.InitializeGlobalEmotes()
		if err != nil {
			slog.Error("error initializing global emotes", "error", err.Error())
		}

		err = em.InitializeChannelEmotes(channelID)
		if err != nil {
			slog.Error("error initializing channel emotes", "error", err.Error())
		}
	})

	return em, nil
}

// InitializeGlobalEmotes initializes the global emotes for BTTV, FFZ, and 7TV.
func (em *EmoteModule) InitializeGlobalEmotes() error {
	if err := em.getBTTVGlobalEmotes(); err != nil {
		return err
	}
	if err := em.getFFZGlobalEmotes(); err != nil {
		return err
	}
	if err := em.getSevenTVGlobalEmotes(); err != nil {
		return err
	}
	return nil
}

// InitializeChannelEmotes initializes the channel emotes for BTTV, FFZ, and 7TV.
func (em *EmoteModule) InitializeChannelEmotes(channelID string) error {
	if err := em.getBTTVChannelEmotes(channelID); err != nil {
		return err
	}
	if err := em.getFFZChannelEmotes(channelID); err != nil {
		return err
	}
	if err := em.getSevenTVChannelEmotes(channelID); err != nil {
		return err
	}
	return nil
}

// GetAllGlobalEmotes returns all global emotes from all providers.
func (em *EmoteModule) GetAllGlobalEmotes() []string {
	return append(append(em.BTTVGlobalEmotes, em.FFZGlobalEmotes...), em.SevenTVGlobalEmotes...)
}

// GetAllChannelEmotes returns all channel emotes from all providers.
func (em *EmoteModule) GetAllChannelEmotes() []string {
	return append(append(em.BTTVChannelEmotes, em.FFZChannelEmotes...), em.SevenTVChannelEmotes...)
}

// GetAllEmotes returns all global and channel emotes from all providers.
func (em *EmoteModule) GetAllEmotes() []string {
	return append(em.GetAllGlobalEmotes(), em.GetAllChannelEmotes()...)
}

// getBTTVGlobalEmotes retrieves the global BTTV emotes.
func (em *EmoteModule) getBTTVGlobalEmotes() error {
	req, err := sling.New().Get("https://api.betterttv.net/3/cached/emotes/global").Request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var emotes []domain.BTTVEmote
	if err := json.Unmarshal(body, &emotes); err != nil {
		return err
	}
	for _, emote := range emotes {
		if !emote.Modifier {
			em.BTTVGlobalEmotes = append(em.BTTVGlobalEmotes, emote.Code)
		}
	}
	return nil
}

// getFFZGlobalEmotes retrieves the global FFZ emotes.
func (em *EmoteModule) getFFZGlobalEmotes() error {
	req, err := sling.New().Get("https://api.frankerfacez.com/v1/set/global").Request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var globalResponse domain.FFZGlobalEmoteResponse
	if err := json.Unmarshal(body, &globalResponse); err != nil {
		return err
	}
	for _, set := range globalResponse.Sets {
		for _, emote := range set.Emoticons {
			if !emote.Modifer {
				em.FFZGlobalEmotes = append(em.FFZGlobalEmotes, emote.Name)
			}
		}
	}
	return nil
}

// getSevenTVGlobalEmotes retrieves the global 7TV emotes.
func (em *EmoteModule) getSevenTVGlobalEmotes() error {
	req, err := sling.New().Get("https://7tv.io/v3/emote-sets/global").Request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var globalResponse domain.SevenTVSet
	if err := json.Unmarshal(body, &globalResponse); err != nil {
		return err
	}
	for _, emote := range globalResponse.Emotes {
		em.SevenTVGlobalEmotes = append(em.SevenTVGlobalEmotes, emote.Name)
	}
	return nil
}

// getBTTVChannelEmotes retrieves the channel BTTV emotes.
func (em *EmoteModule) getBTTVChannelEmotes(channelID string) error {
	req, err := sling.New().Get(fmt.Sprintf("https://api.betterttv.net/3/cached/users/twitch/%s", channelID)).Request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var emoteResponse domain.BTTVChannelEmoteResponse
	if err := json.Unmarshal(body, &emoteResponse); err != nil {
		return err
	}
	for _, emote := range emoteResponse.ChannelEmotes {
		em.BTTVChannelEmotes = append(em.BTTVChannelEmotes, emote.Code)
	}
	for _, sharedEmote := range emoteResponse.SharedEmotes {
		em.BTTVChannelEmotes = append(em.BTTVChannelEmotes, sharedEmote.Code)
	}
	return nil
}

// getFFZChannelEmotes retrieves the channel FFZ emotes.
func (em *EmoteModule) getFFZChannelEmotes(channelID string) error {
	req, err := sling.New().Get(fmt.Sprintf("https://api.frankerfacez.com/v1/room/id/%s", channelID)).Request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var emoteResponse domain.FFZChannelEmoteResponse
	if err := json.Unmarshal(body, &emoteResponse); err != nil {
		return err
	}
	setName := strconv.Itoa(emoteResponse.Room.Set)
	for _, emote := range emoteResponse.Sets[setName].Emoticons {
		em.FFZChannelEmotes = append(em.FFZChannelEmotes, emote.Name)
	}
	return nil
}

// getSevenTVChannelEmotes retrieves the channel 7TV emotes.
func (em *EmoteModule) getSevenTVChannelEmotes(channelID string) error {
	req, err := sling.New().Get(fmt.Sprintf("https://7tv.io/v3/users/twitch/%s", channelID)).Request()
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var emoteResponse domain.SevenTVChannelResponse
	if err := json.Unmarshal(body, &emoteResponse); err != nil {
		return err
	}
	for _, emote := range emoteResponse.EmoteSet.Emotes {
		em.SevenTVChannelEmotes = append(em.SevenTVChannelEmotes, emote.Name)
	}
	return nil
}

// initializeWebSocketClient initializes the WebSocket client and handles incoming messages.
func (em *EmoteModule) initializeWebSocketClient(gctx global.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	wsConn, _, err := websocket.Dial(ctx, "wss://events.7tv.io/v3", nil)
	if err != nil {
		cancel()
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	em.wsConn = wsConn

	subscriptionPayload := Message[SubscribePayload]{
		Op:        OpcodeSubscribe,
		Timestamp: time.Now().Unix(),
		Data: SubscribePayload{
			Type: EventTypeAnyEmoteSet,
			Condition: map[string]string{
				"object_id": "613793270dac665160c56d8f",
			},
		},
	}

	bytes, err := json.Marshal(subscriptionPayload)
	if err != nil {
		cancel()
		return fmt.Errorf("failed to marshal subscription payload: %w", err)
	}

	em.wsConn.Write(ctx, websocket.MessageText, bytes)

	// Close the WebSocket connection when gctx is done
	go func() {
		<-gctx.Done()
		cancel()
		em.wsConn.Close(websocket.StatusNormalClosure, "context done")
	}()

	// Handle incoming messages
	go em.handleWebSocketMessages(ctx)
	return nil
}

// handleWebSocketMessages reads messages from the WebSocket and processes them.
func (em *EmoteModule) handleWebSocketMessages(ctx context.Context) {
	for {
		_, message, err := em.wsConn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				return
			}
			slog.Error("error reading from WebSocket", "error", err.Error())
			return
		}
		em.processWebSocketMessage(ctx, message)
	}
}

// processWebSocketMessage processes a single WebSocket message.
func (em *EmoteModule) processWebSocketMessage(ctx context.Context, message []byte) {
	var commonMsg Message[json.RawMessage]
	if err := json.Unmarshal(message, &commonMsg); err != nil {
		fmt.Println("Error unmarshaling common message:", err)
		return
	}

	switch commonMsg.Op {
	case OpcodeHello:
		var payload HelloPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling hello payload:", err)
			return
		}
		em.handleHelloPayload(payload)
	case OpcodeAck:
		var payload AckPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling ack payload:", err)
			return
		}
		em.handleAckPayload(payload)
	case OpcodeHeartbeat:
		var payload HeartbeatPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling heartbeat payload:", err)
			return
		}
		em.handleHeartbeatPayload(payload)
	case OpcodeReconnect:
		var payload ReconnectPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling reconnect payload:", err)
			return
		}
		em.handleReconnectPayload(ctx, payload)
	case OpcodeResume:
		var payload ResumePayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling resume payload:", err)
			return
		}
		em.handleResumePayload(payload)
	case OpcodeSubscribe:
		var payload SubscribePayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling subscribe payload:", err)
			return
		}
		em.handleSubscribePayload(payload)
	case OpcodeUnsubscribe:
		var payload UnsubscribePayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling unsubscribe payload:", err)
			return
		}
		em.handleUnsubscribePayload(payload)
	case OpcodeDispatch:

		var payload Event
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling dispatch payload:", err)
			return
		}
		em.handleDispatchPayload(payload)
	case OpcodeSignal:
		var payload SignalPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling signal payload:", err)
			return
		}
		em.handleSignalPayload(payload)
	case OpcodeError:
		var payload ErrorPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling error payload:", err)
			return
		}
		em.handleErrorPayload(payload)
	case OpcodeEndOfStream:
		var payload EndOfStreamPayload
		if err := json.Unmarshal(commonMsg.Data, &payload); err != nil {
			fmt.Println("Error unmarshaling end of stream payload:", err)
			return
		}
		em.handleEndOfStreamPayload(payload)
	default:
		fmt.Println("Unknown message opcode:", commonMsg.Op)
	}
}

func (em *EmoteModule) handleHelloPayload(payload HelloPayload) {
	// Handle hello payload
	fmt.Println("Received Hello payload:", payload)
}

func (em *EmoteModule) handleAckPayload(payload AckPayload) {
	// Handle ack payload
	bytes, err := payload.Data.MarshalJSON()
	if err != nil {
		fmt.Println("Error marshaling Ack payload data:", err)
		return
	}

	msg := Message[AckPayload]{}
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		fmt.Println("Error unmarshaling Ack payload data:", err)
		return

	}

	fmt.Println("Received Ack payload:", msg)
}

func (em *EmoteModule) handleHeartbeatPayload(payload HeartbeatPayload) {
	// Handle heartbeat payload
	fmt.Println("Received Heartbeat payload:", payload)
}

func (em *EmoteModule) reconnect(ctx context.Context) error {
	// Close the current WebSocket connection
	if err := em.wsConn.Close(websocket.StatusNormalClosure, "reconnecting"); err != nil {
		fmt.Println("Error closing WebSocket connection:", err)
	}

	// Establish a new WebSocket connection
	wsConn, _, err := websocket.Dial(ctx, "wss://events.7tv.io/v3", nil)
	if err != nil {
		return fmt.Errorf("failed to reconnect to WebSocket: %w", err)
	}
	em.wsConn = wsConn

	// Handle incoming messages with the new connection
	go em.handleWebSocketMessages(ctx)

	return nil
}

func (em *EmoteModule) handleReconnectPayload(ctx context.Context, payload ReconnectPayload) {
	fmt.Println("Received Reconnect payload:", payload)

	// Cancel the current context and create a new one
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Reconnect to the WebSocket
	if err := em.reconnect(ctx); err != nil {
		fmt.Println("Error reconnecting to WebSocket:", err)
	}
}

func (em *EmoteModule) handleResumePayload(payload ResumePayload) {
	// Handle resume payload
	fmt.Println("Received Resume payload:", payload)
}

func (em *EmoteModule) handleSubscribePayload(payload SubscribePayload) {
	// Handle subscribe payload

	fmt.Println("Received Subscribe payload:", payload)
}

func (em *EmoteModule) handleUnsubscribePayload(payload UnsubscribePayload) {
	// Handle unsubscribe payload
	fmt.Println("Received Unsubscribe payload:", payload)
}

func findEmoteIndex(emotes []string, emote string) int {
	for i, e := range emotes {
		if e == emote {
			return i
		}
	}
	return -1
}

func (em *EmoteModule) handleDispatchPayload(payload Event) {
	body := payload.Body

	// Handle emote_set.update event
	if payload.Type == string(EventTypeUpdateEmoteSet) {
		// Emote was added
		if len(body.Pushed) > 0 {
			for _, change := range body.Pushed {
				if change.Key == "emotes" {
					if emote, ok := change.Value.(map[string]interface{}); ok {
						if name, exists := emote["name"].(string); exists {
							fmt.Println("7tv emote added:", name)
							em.SevenTVChannelEmotes = append(em.SevenTVChannelEmotes, name)
						}
					}
				}
			}
		}

		// Emote was removed
		if len(body.Pulled) > 0 || len(body.Removed) > 0 || len(body.Updated) > 0 {
			if len(body.Pulled) > 0 {
				for _, change := range body.Pulled {
					if change.Key == "emotes" {
						if emote, ok := change.OldValue.(map[string]interface{}); ok {
							if name, exists := emote["name"].(string); exists {
								fmt.Println("7tv emote removed:", name)
								// Remove the emote from SevenTVChannelEmotes
								index := findEmoteIndex(em.SevenTVChannelEmotes, name)
								if index != -1 {
									em.SevenTVChannelEmotes = append(em.SevenTVChannelEmotes[:index], em.SevenTVChannelEmotes[index+1:]...)
								}
							}
						}
					}
				}
			}
			if len(body.Removed) > 0 {
				for _, change := range body.Removed {
					if change.Key == "emotes" {
						if emote, ok := change.OldValue.(map[string]interface{}); ok {
							if name, exists := emote["name"].(string); exists {
								fmt.Println("7tv emote removed:", name)
								// Remove the emote from SevenTVChannelEmotes
								index := findEmoteIndex(em.SevenTVChannelEmotes, name)
								if index != -1 {
									em.SevenTVChannelEmotes = append(em.SevenTVChannelEmotes[:index], em.SevenTVChannelEmotes[index+1:]...)
								}
							}
						}
					}
				}
			}
			if len(body.Updated) > 0 {
				for _, change := range body.Updated {
					if change.Key == "emotes" {
						if emote, ok := change.OldValue.(map[string]interface{}); ok {
							if name, exists := emote["name"].(string); exists {
								fmt.Println("7tv emote removed:", name)
								// Remove the emote from SevenTVChannelEmotes
								index := findEmoteIndex(em.SevenTVChannelEmotes, name)
								if index != -1 {
									em.SevenTVChannelEmotes = append(em.SevenTVChannelEmotes[:index], em.SevenTVChannelEmotes[index+1:]...)
								}
							}
						}
					}
				}
			}
		}
	}
}

func (em *EmoteModule) handleSignalPayload(payload SignalPayload) {
	// Handle signal payload
	fmt.Println("Received Signal payload:", payload)
}

func (em *EmoteModule) handleErrorPayload(payload ErrorPayload) {
	// Handle error payload
	fmt.Println("Received Error payload:", payload)
}

func (em *EmoteModule) handleEndOfStreamPayload(payload EndOfStreamPayload) {
	// Handle end of stream payload
	fmt.Println("Received End Of Stream payload:", payload)
}
