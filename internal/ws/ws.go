package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 📱 The Aaxion Device Profile (Upgraded from a simple bool)
type Client struct {
	Conn       *websocket.Conn
	DeviceName string
}

// 📨 The Standardized Aaxion Envelope for JSON messages
type WSMessage struct {
	Type       string `json:"type"`                 // e.g., "TRACK_ADDED", "COMMAND", "STATE_SYNC"
	SenderID   string `json:"senderId,omitempty"`   // Who sent it
	DeviceName string `json:"deviceName,omitempty"` // Human readable (e.g., "Swapnil's iPhone")
	TargetID   string `json:"targetId,omitempty"`   // Who should receive it (if it's a direct command)
	Payload    any    `json:"payload"`              // The actual data
}

// ws upgrader ( it converts fucking http request to websocket suckers --it sucks ws forever until it gets spited out the server err ) and it also allows all origins (because we don't care about security in this project, i mean it's a local application, who the hell cares?)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	// Map the unique deviceID directly to the Client struct now
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
)

// this shouts loudly to all connected suckers (hey why are you reading this comment? go check the ws handler and see how it works)
// here is  commercial of my  web : just visit https://codershuubinc.com and check out my projects and blog posts, you won't regret it, i promise (you are good to go now , commercial is over)

func Broadcast(message any) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for id, client := range clients {
		err := client.Conn.WriteJSON(message)
		if err != nil {
			log.Printf("🔴 WS Write Error: %v", err)
			client.Conn.Close()
			delete(clients, id)
		}
	}
}

//	if you are still reading this comment, you are either a curious developer or you thinking that im using the AI ,
//
// AI is just a  unpaid intern who writes code for me  (and don't worry i treat him very slaveFully and i pay him with cookies and milk *just emojis  🍪 and 🥛 )
func Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	// 1. Grab identity from the URL (e.g., ws://localhost:8080/ws?deviceId=xyz&deviceName=Desktop)
	deviceID := r.URL.Query().Get("deviceId")
	deviceName := r.URL.Query().Get("deviceName")

	if deviceID == "" {
		deviceID = "anonymous-sucker"
	}
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	//  now you are looking curious
	clientsMu.Lock()
	clients[deviceID] = &Client{Conn: conn, DeviceName: deviceName}
	clientsMu.Unlock()

	log.Printf("🟢 %s connected to Aaxion Sync (ID: %s)", deviceName, deviceID)

	// 📢 Shout that a new device joined so the frontend UI can update instantly
	Broadcast(WSMessage{
		Type:       "DEVICE_JOINED",
		SenderID:   deviceID,
		DeviceName: deviceName,
		Payload:    map[string]string{"status": "online"},
	})

	// ok let me tell you about myself  (on line no 68 )
	go func() {
		defer func() {
			clientsMu.Lock()
			delete(clients, deviceID)
			clientsMu.Unlock()
			conn.Close()
			log.Printf("🔴 %s disconnected", deviceName)

			// Tell everyone else this device died
			Broadcast(WSMessage{Type: "DEVICE_LEFT", SenderID: deviceID})
		}()

		//  you can find me on every social media platform with the same username (codershubinc)
		// instagram : https://www.instagram.com/codershubinc/
		// twitter : https://twitter.com/codershubinc
		// linkedin : https://www.linkedin.com/in/codershubinc/
		// github : https://github.com/codershubinc
		// mail me at : admin@codershubinc.com ===> https://mail.google.com/mail/u/0/?view=cm&fs=1&to=admin@codershubinc.com
		for {
			var msg WSMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}

			// 🧠 THE NEW SWITCHBOARD LOGIC
			if msg.Type == "COMMAND" && msg.TargetID != "" {
				// Sneak the message directly to the Target device only
				clientsMu.Lock()
				if targetClient, exists := clients[msg.TargetID]; exists {
					targetClient.Conn.WriteJSON(msg)
					log.Printf("🔀 Routed COMMAND to %s", targetClient.DeviceName)
				}
				clientsMu.Unlock()
			} else if msg.Type == "STATE_SYNC" {
				// If a device updates its play state, shout it to everyone
				Broadcast(msg)
			}
		}
	}()
}
