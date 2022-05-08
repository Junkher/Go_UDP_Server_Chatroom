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

// func (u *user) err(s *server, err error) {
// 		s.conn.WriteToUDP([]byte("ERR: "+ err.Error() + "\n"), u.remoteAddr)
// }

func (u *user) msg(s *server ,msg string) {
		_, err := s.conn.WriteToUDP([]byte(msg), u.remoteAddr)
		checkError(err)
}
