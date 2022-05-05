package backend

import "net"

type room struct {
	name    string
	members map[*net.UDPAddr]*user
}

func (r *room) broadcast(s *server, msg string) {

	for _, c := range r.members {
		// if addr != sender.conn.RemoteAddr() {
		// 		m.msg(msg)
		// }
		// n, err := s.conn.WriteToUDP([]byte(msg), c.remoteAddr())
		c.msg(s, msg)
		// log.Println(c, err)

	}
}