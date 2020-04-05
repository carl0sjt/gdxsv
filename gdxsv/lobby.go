package main

type Lobby struct {
	ID         uint16
	Rule       *Rule
	Users      map[string]*AppPeer
	Rooms      map[uint16]*Room
	EntryUsers []string
}

func NewLobby(lobbyID uint16) *Lobby {
	lobby := &Lobby{
		ID:         lobbyID,
		Rule:       NewRule(),
		Users:      make(map[string]*AppPeer),
		Rooms:      make(map[uint16]*Room),
		EntryUsers: make([]string, 0),
	}
	for i := 1; i <= maxRoomCount; i++ {
		roomID := uint16(i)
		lobby.Rooms[roomID] = NewRoom(lobbyID, roomID)
	}
	return lobby
}

func (l *Lobby) Enter(u *AppPeer) {
	l.Users[u.UserID] = u
}

func (l *Lobby) Exit(userID string) {
	_, ok := l.Users[userID]
	if ok {
		delete(l.Users, userID)
		for i, id := range l.EntryUsers {
			if id == userID {
				l.EntryUsers = append(l.EntryUsers[:i], l.EntryUsers[i+1:]...)
				break
			}
		}
	}
}

func (l *Lobby) Entry(u *AppPeer) {
	l.EntryUsers = append(l.EntryUsers, u.UserID)
}

func (l *Lobby) EntryCancel(u *AppPeer) {
	for i, id := range l.EntryUsers {
		if id == u.UserID {
			l.EntryUsers = append(l.EntryUsers[:i], l.EntryUsers[i+1:]...)
			break
		}
	}
}

func (l *Lobby) GetUserCountBySide() (uint16, uint16) {
	a := uint16(0)
	b := uint16(0)
	for _, u := range l.Users {
		switch u.Entry {
		case EntryRenpo:
			a++
		case EntryZeon:
			b++
		}
	}
	return a, b
}

func (l *Lobby) GetLobbyMatchEntryUserCount() (uint16, uint16) {
	a := uint16(0)
	b := uint16(0)
	for _, id := range l.EntryUsers {
		u, ok := l.Users[id]
		if ok {
			switch u.Entry {
			case EntryRenpo:
				a++
			case EntryZeon:
				b++
			}
		}
	}
	return a, b
}

func (l *Lobby) CanBattleStart() bool {
	a, b := l.GetLobbyMatchEntryUserCount()
	return 2 <= a && 2 <= b
}

func (l *Lobby) PickReadyToBattleUsers() []*AppPeer {
	a := uint16(0)
	b := uint16(0)
	battleUsers := []*AppPeer{}
	for _, id := range l.EntryUsers {
		u, ok := l.Users[id]
		if ok {
			switch u.Entry {
			case EntryRenpo:
				if a < 2 {
					battleUsers = append(battleUsers, u)
				}
				a++
			case EntryZeon:
				if b < 2 {
					battleUsers = append(battleUsers, u)
				}
				b++
			}
		}
	}
	for _, u := range battleUsers {
		l.EntryCancel(u)
	}
	return battleUsers
}
