package authkit

import (
	"fmt"
	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"git.solsynth.dev/hypernet/pusher/pkg/pushkit"
	"github.com/goccy/go-json"
)

func NotifyUser(nx *nex.Conn, userId uint64, notify pushkit.Notification, unsaved ...bool) error {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return fmt.Errorf("failed to get auth service client: %v", err)
	}
	raw, _ := json.Marshal(notify)
	if len(unsaved) == 0 {
		unsaved = append(unsaved, false)
	}
	_, err = proto.NewNotifyServiceClient(conn).NotifyUser(nil, &proto.NotifyUserRequest{
		UserId: userId,
		Notify: &proto.NotifyInfoPayload{
			Unsaved: unsaved[0],
			Data:    raw,
		},
	})
	return err
}

func NotifyUserBatch(nx *nex.Conn, userId []uint64, notify pushkit.Notification, unsaved ...bool) error {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return fmt.Errorf("failed to get auth service client: %v", err)
	}
	raw, _ := json.Marshal(notify)
	if len(unsaved) == 0 {
		unsaved = append(unsaved, false)
	}
	_, err = proto.NewNotifyServiceClient(conn).NotifyUserBatch(nil, &proto.NotifyUserBatchRequest{
		UserId: userId,
		Notify: &proto.NotifyInfoPayload{
			Unsaved: unsaved[0],
			Data:    raw,
		},
	})
	return err
}
