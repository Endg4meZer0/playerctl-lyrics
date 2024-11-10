// The playerctl implementation in lrcsnc doesn't require sessions
//
// This file exists purely for satisfying the session interface
// needed to initialize the player listener in the first place
//
// Playerctl support is deprecated as of 0.3.0 and right now exists purely for compatibility purposes

package sessions

type PlayerctlSession struct{}

func (s *PlayerctlSession) Init() error {
	return nil
}

func (s *PlayerctlSession) GetSession() interface{} {
	return nil
}

func (s *PlayerctlSession) GetPlayer() interface{} {
	return nil
}

func (s *PlayerctlSession) SetPlayer(player interface{}) error {
	return nil
}

func (s *PlayerctlSession) Close() error {
	return nil
}
