package backend

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	clients map[string]*user
	commands chan command
	conn *net.UDPConn
}

func checkError(err error) {
	if err != nil {
		log.Println("Error:", err)
	}
}

func NewServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Startup(port string) {

	log.Println("启动")


	// listener, err := net.Listen("tcp", ":" + port)
	udpAddr, err := net.ResolveUDPAddr("udp", ":" + port)
	// if err != nil {
	// 	log.Fatalf("启动失败 :%s", err.Error())
	// }
	checkError(err)

	// defer listener.Close()
	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)
	defer s.conn.Close()

	log.Printf("start server on :" + port)
	go s.handleMsg()
	for{
			 s.receiveMsg()
	}
	
}

func (s *server) receiveMsg() {

				buf := make([]byte, 1024)
				n, remoteAddr, err := s.conn.ReadFromUDP(buf)
				checkError(err)
				msg := string(buf[0:n]) 
				args := strings.Split(msg, " ")
				cmd := strings.TrimSpace(args[0])
				u := &user{
					remoteAddr: remoteAddr,
				}
				s.clients[remoteAddr.String()] = u

				switch cmd {
						case "/nick":
								s.commands <- command{
									id: CMD_NICK,
									user:  u,
									args: args,
								}
						case "/join":
								s.commands <- command{
								id: CMD_JOIN,
								user:  u,
								args: args,
							}
						case "/rooms":
								s.commands <- command{
								id: CMD_ROOMS,
								user:  u,
								args: args,
							}
						case "/msg":
								s.commands <- command{
								id: CMD_MSG,
								user:  u,
								args: args,
							}
						case "/quit":
								s.commands <- command{
								id: CMD_QUIT,
								user:  u,
								args: args,
							}
						default:
								s.sendMsg(u, "未知命令")
					
				}

}

func (s *server) sendMsg(u *user, msg string) {
	_, err := s.conn.WriteToUDP([]byte(msg), u.remoteAddr)
	checkError(err)
}


func (s *server) handleMsg() {
		for cmd := range s.commands {
				switch cmd.id {
				case CMD_NICK:
					s.nick(cmd.user, cmd.args);
				case CMD_JOIN:
					s.join(cmd.user, cmd.args);
				case CMD_ROOMS:
					s.listRooms(cmd.user, cmd.args);
				case CMD_MSG:
					s.msg(cmd.user, cmd.args);
				case CMD_QUIT:
					s.quit(cmd.user, cmd.args);
				}
		}
}



func (s *server) nick(u *user, args []string) {
		u.nick = args[1]
		s.sendMsg(u, fmt.Sprintf("allright, I will call you %s", u.nick))
}

func (s *server) join(u *user, args []string) {
		roomName := args[1]

		r, ok := s.rooms[roomName]
		if !ok {
			r = &room{
				name: roomName,
				members: make(map[*net.UDPAddr]*user),
			}
		}
		s.rooms[roomName] = r
		r.members[u.remoteAddr] = u

		s.quitRoom(u)

		u.room =r

		r.broadcast(s, u.nick + " 进入房间")
		// u.msg("welcome to " + r.name)
}

func (s *server) listRooms(u *user, args []string) {
		var  rooms []string
		for name := range s.rooms {
				rooms = append(rooms, name)
		}

		u.msg(s, "available rooms are: " +  strings.Join(rooms, ", "))

}

func (s *server) msg(u *user, args []string) {
		if u.room == nil {
				// u.err(errors.New("you must join one room"))			
				s.sendMsg(u, "你必须先进入房间")
				return
		}
		u.room.broadcast(s, u.nick+":" + strings.Join(args[1:], " "))
}

func (s *server) quit(u *user, args []string) {

		log.Printf("user has disconnected: %s", u.remoteAddr.String())

		s.quitRoom(u)
		u.msg(s,"see you~")
}

func (s *server) quitRoom(u *user) {
	if u.room != nil {
			delete(u.room.members, u.remoteAddr)
			u.room.broadcast(s, fmt.Sprintf("%s has left the room", u.nick))
		}
}
