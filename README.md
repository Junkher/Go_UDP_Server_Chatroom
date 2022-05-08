使用go编写的简单UDP服务器，实现实时的多人聊天


## Startup code

```go
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
```

通过`net.ResolveUDPAddr`获取对应IP:port的`*net.UDPAddr`，作为`net.ListenUDP`的参数获取UDPconn，然后就可以通过`s.conn`读取UDP数据包以及向指定的remodteAddr发送UDP数据包。

 `receiveMsg`接收数据，通过`chan`传递给`handleMsg`实现具体的处理逻辑。

## Commands Overview

- /nick \<nickname>
- /join \<roomName>
- /create \<roomName>
- /rooms 
- /msg \<content>
- /users 
- /quit
- /disconnect

## Commands Detail


### CMD_NICK

`/nick <nickname`

修改用户的昵称为`<nickname>`

### CMD_JOIN

`/join <roomName>`

若用户不在房间且该房间名`roomName`存在，则将用户加入房间`roomName`，并向在房间中的用户广播用户进入的信息，最后通知
不在房间的用户更新房间信息列表，在该房间的用户更新用户列表

### CMD_CREATE

`/create <roomName>`

若房间名`roomName`不存在则创建房间

### CMD_ROOMS

`/rooms`

向用户列出当前的房间信息列表，格式为`Lobby: roomName@房roomUserNum,...`

### CMD_MSG

`/msg <content>`

向用户房间内的所有用户发送`user.nick: content`，


### CMD_USERS

`/users`

向用户列出用户所在房间的用户名，格式为`Users: nickname1,nickname2,... `

### CMD_QUIT

`/quit`

用户退出当前所在的房间，若用户退出后该房间人数为0，则将该房间从房间列表删除。
并且通知不在房间的用户更新房间信息，在该房间的用户更新房间用户信息。

### CMD_DISCONNECT

`/disconnect`

将用户从服务器用户列表清除