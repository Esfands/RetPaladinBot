run:
	go run cmd/app/main.go

# Simulate Eventsub events
stream-online:
	twitch event trigger streamup -F http://localhost:3000/v1/twitch/eventsub/ -s 1234567890 -t 38746172

stream-offline:
	twitch event trigger streamdown -F http://localhost:3000/v1/twitch/eventsub/ -s 1234567890 -t 38746172

stream-change:
	twitch event trigger stream-change -F http://localhost:3000/v1/twitch/eventsub/ -s 1234567890 -t 38746172 -v 2