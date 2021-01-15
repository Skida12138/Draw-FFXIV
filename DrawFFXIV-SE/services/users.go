package services

import (
	uuid "github.com/satori/go.uuid"
	"github.com/skida12138/drawffxiv-se/i18n"
	"github.com/skida12138/drawffxiv-se/wsproto"
)

type roomUserRole int

const (
	roleError     roomUserRole = -1
	roleAudience  roomUserRole = 0
	rolePainter   roomUserRole = 1
	roleGuessed   roomUserRole = 2
	rolePreparing roomUserRole = 3
	rolePerpared  roomUserRole = 4
	roleGuessing  roomUserRole = 5
)

func (role roomUserRole) String() string {
	switch role {
	case roleAudience:
		return i18n.Msg("roleAudience")
	case rolePainter:
		return i18n.Msg("rolePainter")
	case roleGuessed:
		return i18n.Msg("roleGuessed")
	case roleGuessing:
		return i18n.Msg("roleGuessing")
	case rolePreparing:
		return i18n.Msg("rolePreparing")
	case rolePerpared:
		return i18n.Msg("rolePerpared")
	default:
		return ""
	}
}

const (
	setSessions            string = "sessions"
	hashSession2NickName   string = "session2nickname"
	hashSession2Room       string = "session2room"
	hashSession2Score      string = "session2score"
	hashSession2UserRoomID string = "session2userroomid"
	hashSession2UserRole   string = "session2userrole"
	hashSession2Rounds     string = "session2rounds"
)

// Session : session type bind with some operations
type Session string

var conns map[Session]*wsproto.Conn

// RegisterNewUser : register a new user
func RegisterNewUser() (Session, error) {
	newUserSession := uuid.NewV4().String()
	err := rdb.SAdd(runCtx, setSessions, newUserSession).Err()
	if err != nil {
		return "", err
	}
	return Session(newUserSession), err
}

// IsValid : check if a sessionID is valid
func (sessionID Session) IsValid() (bool, error) {
	return rdb.SIsMember(runCtx, setSessions, sessionID).Result()
}

// GetNickName : get user's nick name by sessionID
func (sessionID Session) GetNickName() (string, error) {
	return rdb.HGet(runCtx, hashSession2NickName, string(sessionID)).Result()
}

// SetNickName : set nick name for user by sessionID
func (sessionID Session) SetNickName(nickName string) error {
	return rdb.HSet(runCtx, hashSession2NickName, nickName).Err()
}

// GetRoomID : get roomID by sessionID
func (sessionID Session) GetRoomID() int32 {
	var roomID int
	var err error
	if roomID, err = rdb.HGet(runCtx, hashSession2Room, string(sessionID)).Int(); err != nil {
		return -1
	}
	return int32(roomID)
}

// GetUserRoomID : get ID in room by sessionID
func (sessionID Session) GetUserRoomID() int32 {
	var userRoomID int
	var err error
	if userRoomID, err = rdb.HGet(runCtx, hashSession2UserRoomID, string(sessionID)).Int(); err != nil {
		return -1
	}
	return int32(userRoomID)
}

func (sessionID Session) getRounds() int32 {
	var rounds int
	var err error
	if rounds, err = rdb.HGet(runCtx, hashSession2Rounds, string(sessionID)).Int(); err != nil {
		return -1
	}
	return int32(rounds)
}

func (sessionID Session) initUserForRoom(roomID int32, inheritScore bool) error {
	return nil
}

func (sessionID Session) getRole() (roomUserRole, error) {
	role, err := rdb.HGet(runCtx, hashSession2UserRole, string(sessionID)).Int()
	if err != nil {
		return roleError, err
	}
	return roomUserRole(role), nil
}
