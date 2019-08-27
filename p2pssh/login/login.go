package login

import (
	"errors"
	"strings"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/wanyvic/p2pssh/api"
)

// func Connect(auth *api.UserAuth) {
// 	client.New(context.BackGround(), auth)
// }
func ParseAuth(str string) (auth api.UserAuth, err error) {
	if str == "" {
		return auth, errors.New("no args")
	}
	s := strings.Split(str, "@")
	if len(s) == 2 {
		auth.User = s[0]
		auth.NodeID, err = peer.IDB58Decode(s[1])
		if err != nil {
			return auth, err
		}
	} else if len(s) == 1 {
		auth.NodeID, err = peer.IDB58Decode(s[1])
		if err != nil {
			return auth, err
		}
	} else {
		return auth, errors.New("args errors")
	}
	return auth, nil
}
