package main

import (
	pb "clienttest/pb"
	"clienttest/renderer"
	"flag"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/veandco/go-sdl2/sdl"
	"google.golang.org/protobuf/proto"
)

var addr = flag.String("addr", "172.30.1.55:8080", "http service address")

func main() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	if e := sdl.Init(sdl.INIT_EVERYTHING); e != nil {
		panic(sdl.GetError())
	}
	re := renderer.NewRenderer()
	flag.Parse()
	//log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			handleMessage(message, re)
		}
	}()

	ticker := time.NewTicker(time.Microsecond)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			input(1, c)
			//log.Println(t)
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))

			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return

		}

	}
}

func handleMessage(message []byte, renderer *renderer.Renderer) {
	serverMessage := &pb.ServerMessage{}
	if err := proto.Unmarshal(message, serverMessage); err != nil {
		log.Println("Failed to unmarshal UserInput:", err)
	}
	switch x := serverMessage.Payload.(type) {
	case *pb.ServerMessage_GameState:
		//log.Println(len(serverMessage.GetGameState().GetPlayers()))
		// for _, p := range serverMessage.GetGameState().GetPlayers() {
		// 	if p.PlayerID == 1 {
		// 		//log.Println(p.GetPos().X, p.GetPos().Y)
		// 		renderer.Render(p.GetPos().X, p.GetPos().Y)
		// 	}
		// 	)
		// }
		renderer.Render(serverMessage.GetGameState().GetPlayers())
	default:
		log.Println("Uknown message type:", x)
	}
}

func input(clientId uint32, c *websocket.Conn) {
	var b []byte = make([]byte, 1)
	// var id int64 = 1
	os.Stdin.Read(b)
	s := string(b)
	protoUserMessage := &pb.UserMessage{
		Payload: &pb.UserMessage_UserInput{
			UserInput: &pb.UserInput{
				Move:   true,
				Facing: 0,
			},
		},
	}
	switch s {
	case "A": // up
		protoUserMessage.GetUserInput().Move = true
		protoUserMessage.GetUserInput().Facing = -90
		m, _ := proto.Marshal(protoUserMessage)
		c.WriteMessage(websocket.BinaryMessage, m)

	case "B": //down
		protoUserMessage.GetUserInput().Move = true
		protoUserMessage.GetUserInput().Facing = 90
		m, _ := proto.Marshal(protoUserMessage)
		c.WriteMessage(websocket.BinaryMessage, m)
	case "C": //left
		protoUserMessage.GetUserInput().Move = true
		protoUserMessage.GetUserInput().Facing = 0
		m, _ := proto.Marshal(protoUserMessage)
		c.WriteMessage(websocket.BinaryMessage, m)
	case "D": //right
		protoUserMessage.GetUserInput().Move = true
		protoUserMessage.GetUserInput().Facing = 180
		m, _ := proto.Marshal(protoUserMessage)
		c.WriteMessage(websocket.BinaryMessage, m)
	case "a": // stop
		protoUserMessage.GetUserInput().Move = false
		protoUserMessage.GetUserInput().Facing = 3.1415 //?
		m, _ := proto.Marshal(protoUserMessage)
		c.WriteMessage(websocket.BinaryMessage, m)
	case "l": // player left
		break
		// case "j": //testing UserJoined
		// 	protoUserMessage.Payload = &api.UserMessage_JoinRequest{
		// 		JoinRequest: &api.JoinRequest{
		// 			UserID: []int64{id},
		// 		},
		// 	}
		// 	event.FireEvent(&events.UserJoined{
		// 		ClientID: uint32(id),
		// 		UserName: strconv.Itoa(int(id)),
		// 	})
		// 	id++
	}
}
