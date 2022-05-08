package backend

type commandID int

const (
	CMD_NICK commandID = iota //0
	CMD_JOIN
	CMD_CREATE
	CMD_ROOMS
	CMD_MSG
	CMD_USERS
	CMD_QUITROOM
	CMD_DISCONNECT
)

type command struct {
	id   commandID
	user *user
	args []string
}
