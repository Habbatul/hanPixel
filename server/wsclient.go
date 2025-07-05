package server

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/pion/webrtc/v4"
)

var (
	signalConn      *websocket.Conn
	peerConns       = make(map[string]*webrtc.PeerConnection)
	dataChans       = make(map[string]*DataChannel)
	LocalPlayerID   string
	remotePositions = make(map[string]Position)
	positionsMux    sync.Mutex
	connsMux        sync.Mutex

	chatChan = make(chan Chat)
)

type DataChannel struct {
	position *webrtc.DataChannel
	chat     *webrtc.DataChannel
}

type Position struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}

type Chat struct {
	ID       string `json:"id"`
	ChatText string `json:"chatText"`
}

type SignalMessage struct {
	Type  string                     `json:"type,omitempty"`
	From  string                     `json:"from,omitempty"`
	To    string                     `json:"to,omitempty"`
	SDP   *webrtc.SessionDescription `json:"sdp,omitempty"`
	ICE   *webrtc.ICECandidateInit   `json:"ice,omitempty"`
	Peers []string                   `json:"peers,omitempty"`
}

func StartWebRTC() error {
	ctx := context.Background()
	var err error

	signalConn, _, err = websocket.Dial(ctx, "ws://localhost:8080/signal", nil)
	if err != nil {
		return err
	}

	var initMsg SignalMessage
	if err = wsjson.Read(ctx, signalConn, &initMsg); err != nil {
		return err
	}

	LocalPlayerID = initMsg.From

	log.Println(LocalPlayerID)
	log.Println("Assigned LocalPlayerID:", LocalPlayerID)

	go signalingLoop(ctx)

	return nil
}

func signalingLoop(ctx context.Context) {
	for {

		var msg SignalMessage
		if signalConn != nil {
			//log.Println(signalConn)
			if err := wsjson.Read(ctx, signalConn, &msg); err != nil {
				return
			}

		} else {
			log.Println(signalConn)
		}

		log.Println(msg.Type)

		switch msg.Type {
		case "peers":
			for _, id := range msg.Peers {
				if id == LocalPlayerID {
					continue
				}
				createPeerConnection(id, "offer")
				offer, err := peerConns[id].CreateOffer(nil)
				if err != nil {
					return
				}
				peerConns[id].SetLocalDescription(offer)
				sendSignal(SignalMessage{Type: "offer", From: LocalPlayerID, To: id, SDP: &offer})
			}
		case "offer":
			createPeerConnection(msg.From, "answer")
			setRemoteDescription(msg.From, msg.SDP)

			answer, err := peerConns[msg.From].CreateAnswer(nil)
			if err != nil {
				return
			}
			if err = peerConns[msg.From].SetLocalDescription(answer); err != nil {
				return
			}
			sendSignal(SignalMessage{Type: "answer", From: LocalPlayerID, To: msg.From, SDP: &answer})

		case "answer":
			setRemoteDescription(msg.From, msg.SDP)

		case "candidate":
			addIceCandidate(msg.From, msg.ICE)

		}
	}
}

func createPeerConnection(remoteID string, typeMessage string) error {
	connsMux.Lock()
	defer connsMux.Unlock()
	if _, exist := peerConns[remoteID]; exist {
		return nil
	}

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{
			URLs: []string{"stun:stun.l.google.com:19302"},
		}},
	})
	if err != nil {
		return err
	}

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}

		ice := candidate.ToJSON()
		sendSignal(SignalMessage{Type: "candidate", From: LocalPlayerID, To: remoteID, ICE: &ice})
	})

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		switch state {
		case webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateDisconnected, webrtc.PeerConnectionStateClosed:
			removePeer(remoteID)
			log.Println("Koneksi state :", remoteID, " hilang karena :", state.String())
		}
	})

	dataChans[remoteID] = &DataChannel{}

	if typeMessage == "offer" {
		dcPos, err := pc.CreateDataChannel("position", nil)
		if err != nil {
			return err
		}
		setupDataChannelPos(remoteID, dcPos)
		dataChans[remoteID].position = dcPos

		dcChat, err := pc.CreateDataChannel("chat", nil)
		if err != nil {
			return err
		}
		setupDataChannelChat(remoteID, dcChat)
		dataChans[remoteID].chat = dcChat
	} else {
		pc.OnDataChannel(func(dc *webrtc.DataChannel) {
			if dc.Label() == "position" {
				setupDataChannelPos(remoteID, dc)
				dataChans[remoteID].position = dc
			}
			if dc.Label() == "chat" {
				setupDataChannelChat(remoteID, dc)
				dataChans[remoteID].chat = dc
			}
		})
	}

	peerConns[remoteID] = pc
	return nil
}

func sendSignal(msg SignalMessage) {
	ctx := context.Background()
	if signalConn == nil {
		return
	}
	if err := wsjson.Write(ctx, signalConn, msg); err != nil {
		log.Println("error waktu mengirim signal Type:", msg.Type, "; error:", err)
	}
}

func removePeer(remoteID string) {
	connsMux.Lock()
	defer connsMux.Unlock()

	if dc, exist := dataChans[remoteID]; exist {
		if dc.position != nil {
			dc.position.Close()
		}
		if dc.chat != nil {
			dc.chat.Close()
		}
		delete(dataChans, remoteID)
	}
	if pc, exist := peerConns[remoteID]; exist {
		pc.Close()
		delete(peerConns, remoteID)
	}

	positionsMux.Lock()
	delete(remotePositions, remoteID)
	positionsMux.Unlock()
}

func setupDataChannelPos(remoteID string, dc *webrtc.DataChannel) {
	dc.OnOpen(func() {
		log.Println("data channel dengan ", remoteID, " dibuka")
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		var pos Position
		if err := json.Unmarshal(msg.Data, &pos); err != nil {
			log.Println("unmarshal pesan data channel gagal")
			return
		}
		positionsMux.Lock()
		remotePositions[pos.ID] = pos
		positionsMux.Unlock()
	})
}

func setupDataChannelChat(remoteID string, dc *webrtc.DataChannel) {
	dc.OnOpen(func() {
		log.Println("data channel dengan ", remoteID, " dibuka")
	})

	dc.OnMessage(func(msg webrtc.DataChannelMessage) {
		var chat Chat
		if err := json.Unmarshal(msg.Data, &chat); err != nil {
			log.Println("unmarshal pesan data channel gagal")
			return
		}
		if chat.ChatText != "" {
			chatChan <- chat
		}

	})
}

func setRemoteDescription(remoteID string, sdp *webrtc.SessionDescription) {
	connsMux.Lock()
	defer connsMux.Unlock()

	if pc, exist := peerConns[remoteID]; exist {
		if err := pc.SetRemoteDescription(*sdp); err != nil {
			log.Println("error set remote description : ", err)
		}
	}
}

func addIceCandidate(remoteID string, ice *webrtc.ICECandidateInit) {
	connsMux.Lock()
	defer connsMux.Unlock()

	if pc, exist := peerConns[remoteID]; exist {
		if err := pc.AddICECandidate(*ice); err != nil {
			log.Println("error tmbah ice candidate : ", err)
		}
	}
}

func SendPosition(x float64, y float64) {
	pos := &Position{ID: LocalPlayerID, X: x, Y: y}
	data, _ := json.Marshal(pos)
	connsMux.Lock()
	defer connsMux.Unlock()

	for peerID, datChan := range dataChans {
		if datChan.position != nil && datChan.position.ReadyState() == webrtc.DataChannelStateOpen {
			err := datChan.position.Send(data)
			if err != nil {
				log.Println("peer dengan id :", peerID, "gagal mengirim")
			}
		}
	}
}

func GetRemotePositions() map[string]Position {
	positionsMux.Lock()
	defer positionsMux.Unlock()
	positionCpy := make(map[string]Position)
	for key, val := range remotePositions {
		positionCpy[key] = val
	}
	return positionCpy
}

func SendChat(chatText string) {
	chat := &Chat{ID: LocalPlayerID, ChatText: chatText}
	data, _ := json.Marshal(chat)

	connsMux.Lock()
	connsMux.Unlock()

	for _, dataChan := range dataChans {
		if dataChan.chat != nil && dataChan.chat.ReadyState() == webrtc.DataChannelStateOpen {
			err := dataChan.chat.Send(data)
			if err != nil {
				log.Println("error kirim chat :", err)
			}
		}
	}
}

func GetChat(readerChat func(chatID string, chatText string)) {
	select {
	case chat := <-chatChan:
		readerChat(chat.ID, chat.ChatText)
	default:
		return
	}
}

func StartPositionAsyncDelaySender(getPos func() (float64, float64)) {
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			x, y := getPos()
			SendPosition(x, y)
		}
	}()
}
