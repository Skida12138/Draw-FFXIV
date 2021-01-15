package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/skida12138/drawffxiv-se/i18n"
	"github.com/skida12138/drawffxiv-se/utils"
	"github.com/skida12138/drawffxiv-se/wsproto"
)

const (
	setRooms        string = "rooms"
	hashRoom2Graph  string = "room2graph"
	hashRoom2Status string = "room2status"
	hashRoom2Rounds string = "room2rounds"
	hashRoom2Answer string = "room2answer"
)

// RoomStatus : current status of room
type RoomStatus int

const (
	roomError     RoomStatus = -1
	roomOnGame    RoomStatus = 0
	roomPreparing RoomStatus = 1
	roomPaused    RoomStatus = 2
)

func (status RoomStatus) String() string {
	switch status {
	case roomOnGame:
		return i18n.Msg("roomOnGame")
	case roomPreparing:
		return i18n.Msg("roomPreparing")
	case roomPaused:
		return i18n.Msg("roomPaused")
	default:
		return i18n.Msg("roomError")
	}
}

var (
	roomLock      map[int32]*sync.RWMutex
	activeUsers   map[int32]map[Session]struct{}
	inactiveUsers map[int32]map[Session]struct{}
)

var lastLostConnTime map[int32]time.Time

// UserEnterRoom : initialize user's basic infomations for room
func UserEnterRoom(userSession Session, roomID int32) error {
	if !isRoomOpen(roomID) {
		if err := createRoom(roomID); err != nil {
			return err
		}
		return userSession.initUserForRoom(roomID, false)
	}
	if userSession.GetRoomID() != roomID {
		return userSession.initUserForRoom(roomID, false)
	}
	if userSession.getRounds() != getRoomRounds(roomID) {
		return userSession.initUserForRoom(roomID, true)
	}
	return nil
}

func isRoomOpen(roomID int32) bool {
	return rdb.SIsMember(runCtx, setRooms, fmt.Sprint(roomID)).Val()
}

func createRoom(roomID int32) error {
	if err := utils.Try(
		func() error { return rdb.SAdd(runCtx, setRooms, roomID).Err() },
		func() error { return rdb.HSet(runCtx, hashRoom2Status, roomID, roomPreparing).Err() },
		func() error { return rdb.HSet(runCtx, hashRoom2Rounds, roomID, 0).Err() },
		func() error { return rdb.HSet(runCtx, hashRoom2Graph, roomID, "[]").Err() },
	); err != nil {
		return err
	}
	go func() {
		for {
			time.Sleep(time.Minute)
			checkLostUser(roomID)
			if checkRoomExpired(roomID) {
				break
			}
		}
	}()
	return nil
}

func getRoomStatus(roomID int32) RoomStatus {
	if status, err := rdb.HGet(runCtx, hashRoom2Status, fmt.Sprint(roomID)).Int(); err == nil {
		return RoomStatus(status)
	}
	return roomError
}

func getRoomRounds(roomID int32) int32 {
	if rounds, err := rdb.HGet(runCtx, hashRoom2Rounds, fmt.Sprint(roomID)).Int(); err == nil {
		return int32(rounds)
	}
	return -1
}

func checkLostUser(roomID int32) {
	newInactiveUsers := make([]Session, 0)
	roomLock[roomID].RLock()
	for user := range activeUsers[roomID] {
		if err := conns[user].SyncWrite(wsproto.BeatMsg); err != nil {
			newInactiveUsers = append(newInactiveUsers, user)
			conns[user].Replace(nil)
		}
	}
	roomLock[roomID].RUnlock()
	if len(newInactiveUsers[roomID]) > 0 {
		lastLostConnTime[roomID] = time.Now()
	}
	roomLock[roomID].Lock()
	for _, user := range newInactiveUsers {
		delete(activeUsers[roomID], user)
		inactiveUsers[roomID][user] = struct{}{}
	}
	roomLock[roomID].Unlock()
}

func checkRoomExpired(roomID int32) bool {
	roomLock[roomID].RLock()
	if len(activeUsers[roomID]) == 0 {
		elapsTime := lastLostConnTime[roomID].Add(2 * time.Minute)
		if time.Now().After(elapsTime) {
			roomLock[roomID].RUnlock()
			deleteRoom(roomID)
			return true
		}
	}
	return false
}

func deleteRoom(roomID int32) {
	rdb.SRem(runCtx, setRooms, roomID)
	rdb.HDel(runCtx, hashRoom2Status, fmt.Sprint(roomID))
	rdb.HDel(runCtx, hashRoom2Rounds, fmt.Sprint(roomID))
	rdb.HDel(runCtx, hashRoom2Graph, fmt.Sprint(roomID))
}

func appendRoomGraph(roomID int32, newStroke []byte) error {
	oldGraph := rdb.HGet(runCtx, hashRoom2Graph, fmt.Sprint(roomID)).String()
	if len(oldGraph) > 0 {
		oldGraph += ","
	}
	return rdb.HSet(runCtx, hashRoom2Graph, roomID, oldGraph+string(newStroke)).Err()
}

func getRoomAnswer(roomID int32) string {
	return rdb.HGet(runCtx, hashRoom2Answer, fmt.Sprint(roomID)).String()
}
