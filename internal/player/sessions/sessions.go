package sessions

type MediaSession interface {
	Init() error
	GetSession() interface{}
	GetPlayer() interface{}
	SetPlayer(player interface{}) error
	Close() error
}

var MediaSessions map[string]MediaSession = map[string]MediaSession{
	"mpris": &DBusSession{conn: nil},
}
