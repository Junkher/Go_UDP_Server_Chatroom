package backend

import (
	"net"
)

type user struct {
	nick string
	room *room
	remoteAddr *net.UDPAddr
	// 单向write-only
	// commands chan<- command
}

// func (u *user) readInput() {
// 	for {
// 			msg, err := bufio.NewReader(u.conn).ReadString('\n')
// 			if err != nil {
// 				return
// 			}

// 			msg = strings.Trim(msg, "\r\n")

// 			args := strings.Split(msg, " ")
// 			cmd := strings.TrimSpace(args[0])

// 			switch cmd {
// 					case "/nick":
// 							u.commands <- command{
// 								id: CMD_NICK,
// 								user:  u,
// 								args: args,
// 							}
// 				 	case "/join":
// 						u.commands <- command{
// 							id: CMD_JOIN,
// 							user:  u,
// 							args: args,
// 						}
// 					case "/rooms":
// 						u.commands <- command{
// 							id: CMD_ROOMS,
// 							user:  u,
// 							args: args,
// 						}
// 					case "/msg":
// 						u.commands <- command{
// 							id: CMD_MSG,
// 							user:  u,
// 							args: args,
// 						}
// 					case "/quit":
// 						u.commands <- command{
// 							id: CMD_QUIT,
// 							user:  u,
// 							args: args,
// 						}
// 					default:
// 							u.err(fmt.Errorf("unknown command: %s", cmd))
				
// 			}

// 	}
// }

// func (u *user) err(s *server, err error) {
// 		s.conn.WriteToUDP([]byte("ERR: "+ err.Error() + "\n"), u.remoteAddr)
// }

func (u *user) msg(s *server ,msg string) {
		_, err := s.conn.WriteToUDP([]byte(msg), u.remoteAddr)
		checkError(err)
}
