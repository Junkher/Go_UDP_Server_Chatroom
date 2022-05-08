package backend

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room//所有房间
	clients map[string]*user//所有用户
	commands chan command//命令channal
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
		clients:  make(map[string]*user),
		commands: make(chan command),
	}
}

func (s *server) Startup(port string) {

	log.Println("启动")

	udpAddr, err := net.ResolveUDPAddr("udp", ":" + port)
	checkError(err)

	s.conn, err = net.ListenUDP("udp", udpAddr)
	checkError(err)
	defer s.conn.Close()
	
	//初始化第一个房间
	s.rooms["UDPyyds"] = &room{
		name: "UDPyyds",
		members: make(map[*net.UDPAddr]*user),
	}

	log.Printf("start server on :" + port)
	go s.handleMsg()
	for{
			 s.receiveMsg()
	}
	
}


func (s *server) receiveMsg() {

				buf := make([]byte, 1024)
				n, remoteAddr, err := s.conn.ReadFromUDP(buf)
				log.Println("收到啦！")
				checkError(err)
				msg := strings.Trim(string(buf[0:n]), "\r\n") 
				args := strings.Split(msg, " ")
				cmd := strings.TrimSpace(args[0])
				fmt.Println("地址:",remoteAddr)
				fmt.Println("完整信息:", msg)
				fmt.Println("命令:", cmd)
				fmt.Println("参数个数:", len(args) )

				u, ok := s.clients[remoteAddr.String()]
				if !ok {
					fmt.Println("新用户:", remoteAddr)
				  	u = &user{
							nick: remoteAddr.String(),
						remoteAddr: remoteAddr,
			  	}
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
						case "/create":
								s.commands <- command{
								id: CMD_CREATE,
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
						case "/users":
							s.commands <- command{
							id: CMD_USERS,
							user:  u,
							args: args,
						}
						case "/quit":
								s.commands <- command{
								id: CMD_QUITROOM,
								user:  u,
								args: args,
							}
						case "/disconnect":
							s.commands <- command{
								id: CMD_DISCONNECT,
								user:  u,
								args: args,
							}
						case "hi":
							  s.sendMsg(u, "hi")
						default:
								s.sendMsg(u, "未知命令")
					
				}

}

// 通过WriteToUDP向指定地址发送信息
func (s *server) sendMsg(u *user, msg string) {
	_, err := s.conn.WriteToUDP([]byte(msg), u.remoteAddr)
	checkError(err)
}


// 根据命令类型，调用对应的函数处理
func (s *server) handleMsg() {
		for cmd := range s.commands {
				switch cmd.id {
				case CMD_NICK:
					s.nick(cmd.user, cmd.args);
				case CMD_CREATE:
					s.create(cmd.user, cmd.args);
				case CMD_JOIN:
					s.join(cmd.user, cmd.args);
				case CMD_ROOMS:
					s.listRooms(cmd.user, cmd.args);
				case CMD_MSG:
					s.msg(cmd.user, cmd.args);
				case CMD_USERS:
					s.listUsers(cmd.user, cmd.args)
				case CMD_QUITROOM:
					s.quit(cmd.user, cmd.args);
				case CMD_DISCONNECT:
					s.disconnect(cmd.user, cmd.args)
				}
		}
}

// 通知函数
func (s *server) notice(r *room) {
	for _, c := range s.clients { 
			if c.room == nil {
			fmt.Println("通知不在房间的人刷新roomlist")
			s.sendMsg(c, "Fresh")
			}
			if c.room == r {
				fmt.Println("通知在房间的人刷新userlist")
				s.sendMsg(c, "FreshUserList")
			}

	}
}

//修改用户昵称
func (s *server) nick(u *user, args []string) {
		u.nick = args[1]
		s.sendMsg(u, fmt.Sprintf("好的，你的昵称是 %s", u.nick))
}

//用户进入房间
func (s *server) join(u *user, args []string) {
		roomName := args[1]
		fmt.Println(u.nick + " try to enter " + roomName)
		//判断用户是否已在房间
		if (u.room != nil) {
			fmt.Println(u.nick, "已在房间", u.room.name)
			s.sendMsg(u, "你已在房间" + u.room.name)
			return
		}
		//判断房间是否存在
		r, ok := s.rooms[roomName]
		if !ok {
		  fmt.Println(roomName, "不存在")
			s.sendMsg(u, roomName + "不存在")	
			return
		}
		u.room = r
		r.members[u.remoteAddr] = u

		fmt.Println(roomName, "存在")
		r.broadcast(s, u.nick + " 进入房间")
		// u.msg("welcome to " + r.name)
		// 通知不在房间的用户更新房间信息，在该房间的用户更新用户信息
		s.notice(r)
}

//用户创建房间，房间不能重名
func  (s *server) create(u *user, args []string) {
		roomName := args[1]
		_, ok := s.rooms[roomName]
		if !ok {
			s.rooms[roomName] = &room{
				name: roomName,
				members: make(map[*net.UDPAddr]*user),
			}
			fmt.Println(roomName, "成功创建")
			s.sendMsg(u, "Create: "+ roomName + " ok")
			// s.notice()
			//因为后续会进入该房间，因此无须在此通知
		} else {
			fmt.Println(roomName, "已存在")
			s.sendMsg(u, "Create: "+ roomName +" fail")
		}
		// u.room = r
		// r.broadcast(s, u.nick + " 进入房间")
}

//获取当前房间列表及房间人数
func (s *server) listRooms(u *user, args []string) {
		var  rooms []string
		for _, r := range s.rooms {
				name_num := r.name + "@" + fmt.Sprint(len(r.members))
				rooms = append(rooms, name_num)
		}
    // 发送房间名数组字符串
		u.msg(s,  "Lobby: " + strings.Join(rooms, ","))

}

// 获取用户所在房间的用户列表
func (s *server) listUsers(u *user, args []string) {
		var users []string
		for _, u := range u.room.members {
			  users = append(users, u.nick)
		}
		//发送用户列表
		s.sendMsg(u, "Users: " + strings.Join(users, ","))
}

//用户在房间发送信息
func (s *server) msg(u *user, args []string) {
		if u.room == nil {	
				s.sendMsg(u, "你必须先进入房间")
				return
		}
		u.room.broadcast(s, "#" + u.nick+": " + strings.Join(args[1:], " "))
}

//用户申请从服务器退出
func (s *server) disconnect(u *user, args []string) {
  	log.Println(u.remoteAddr.String(), "断开连接" )
		s.sendMsg(u, "再见")
	  delete(s.clients, u.remoteAddr.String())
}

//用户离开房间
func (s *server) quit(u *user, args []string) {

		// log.Println(u.remoteAddr.String(), "断开连接" )
		if u.room == nil {
			return
		}
		r := u.room
		delete(u.room.members, u.remoteAddr)
		if len(u.room.members) == 0 && u.room.name != "UDPyyds"{
			delete(s.rooms, u.room.name)
		} else {
			u.room.broadcast(s, fmt.Sprintf("%s 离开了房间", u.nick))
		}
		//将用户的房间置空
		u.room = nil
		// 通知不在房间的用户更新房间信息，在该房间的用户更新用户信息
		s.notice(r)
		// s.quitRoom(u)
		// u.msg(s,"see you~")
}

// func (s *server) quitRoom(u *user) {
// 	if u.room != nil {
// 			delete(u.room.members, u.remoteAddr)
// 			u.room.broadcast(s, fmt.Sprintf("%s 离开了房间", u.nick))
// 		}
// }
