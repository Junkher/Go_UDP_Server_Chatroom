package backend

import (
	"fmt"
	"net"
)

type room struct {
	name    string
	members map[*net.UDPAddr]*user
}

func (r *room) broadcast(s *server, msg string) {

	for _, u := range r.members {
		fmt.Println("广播", msg)
		s.sendMsg(u, msg)
		// log.Println(c, err)

	}
}