package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/skida12138/drawffxiv-se/wsproto"
)

// Consume : consume service will consume buff stream and dispatch corresponding service
func Consume(buff []byte, session Session) ([]byte, error) {
	var cmdType string
	var rawArgs wsproto.RawCmdArg
	var hasNext bool
	for {
		if cmdType, rawArgs, buff, hasNext = wsproto.NextCommand(buff); !hasNext {
			return buff, nil
		}
		if err := dispatch(cmdType, rawArgs, session); err != nil {
			return buff, err
		}
	}
}

const (
	cmdTypeStroke  string = "stroke"
	cmdTypePrepare string = "prepare"
	cmdTypeChat    string = "chat"
)

func dispatch(cmdType string, rawArgs wsproto.RawCmdArg, session Session) error {
	switch cmdType {
	case cmdTypeStroke:
		return cmdStroke(rawArgs, session)
	case cmdTypePrepare:
		return cmdPrepare(session)
	case cmdTypeChat:
		return cmdChat(session, rawArgs)
	default:
		return fmt.Errorf("unrecognized command: %s", cmdType)
	}
}

func cmdStroke(args []byte, session Session) error {
	if role, err := session.getRole(); role != rolePainter {
		return fmt.Errorf("A user call `stroke` command in role: %s", role)
	} else if err != nil {
		return nil
	}
	roomID := session.GetRoomID()
	if status := getRoomStatus(roomID); status != roomOnGame {
		return fmt.Errorf("a user called 'stroke' command in room with status: %s", status.String())
	}
	if err := appendRoomGraph(roomID, args); err != nil {
		return err
	}
	msg := fmt.Sprintf("%s:%s;", cmdTypeStroke, string(args))
	writeToAllUsersInRoom(roomID, session, msg)
	return nil
}

func cmdPrepare(session Session) error {
	if role, err := session.getRole(); role != rolePreparing {
		return fmt.Errorf("A user call `stroke` command in role: %s", role)
	} else if err != nil {
		return nil
	}
	roomID := session.GetRoomID()
	if status := getRoomStatus(roomID); status != roomPreparing {
		return fmt.Errorf("a user called 'stroke' command in room with status: %s", status.String())
	}
	userRoomID := session.GetUserRoomID()
	msg := fmt.Sprintf("%s:%d;", cmdTypePrepare, userRoomID)
	writeToAllUsersInRoom(roomID, session, msg)
	return nil
}

func cmdChat(session Session, content []byte) error {
	guessed := false
	roomID := session.GetRoomID()
	if role, err := session.getRole(); err != nil {
		return err
	} else if role == roleGuessing && getRoomStatus(roomID) == roomOnGame {
		decodedContent, _ := url.QueryUnescape(string(content))
		guessed = getRoomAnswer(roomID) == decodedContent
	}
	dispatchBody := map[string]string{
		"content": string(content),
		"userID":  fmt.Sprint(session.GetUserRoomID()),
		"guessed": fmt.Sprint(guessed),
	}
	encodedBody, err := json.Marshal(dispatchBody)
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s:%s;", cmdTypeChat, string(encodedBody))
	writeToAllUsersInRoom(roomID, session, msg)
	return nil
}

func writeToAllUsersInRoom(roomID int32, exception Session, msg string) {
	allWritesFinished := new(sync.WaitGroup)
	roomLock[roomID].RLock()
	for toUser := range activeUsers[roomID] {
		if toUser != exception {
			go func() {
				defer allWritesFinished.Done()
				allWritesFinished.Add(1)
				conns[toUser].SyncWrite(msg)
			}()
		}
	}
	roomLock[roomID].RUnlock()
	allWritesFinished.Wait()
}
