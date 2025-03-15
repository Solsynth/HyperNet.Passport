package authkit

import (
	"context"
	"fmt"

	"git.solsynth.dev/hypernet/nexus/pkg/nex"
	"git.solsynth.dev/hypernet/nexus/pkg/nex/sec"
	"git.solsynth.dev/hypernet/passport/pkg/proto"
	"github.com/gofiber/fiber/v2"
)

func AddEvent(nx *nex.Conn, userId uint, action string, meta map[string]any, ip, ua string) error {
	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return fmt.Errorf("failed to get auth service client: %v", err)
	}
	_, err = proto.NewAuditServiceClient(conn).RecordEvent(context.Background(), &proto.RecordEventRequest{
		UserId:    uint64(userId),
		Action:    action,
		Metadata:  nex.EncodeMap(meta),
		Ip:        ip,
		UserAgent: ua,
	})
	return err
}

func AddEventExt(nx *nex.Conn, action string, meta map[string]any, c *fiber.Ctx) error {
	user, ok := c.Locals("nex_user").(*sec.UserInfo)
	if !ok {
		return fmt.Errorf("failed to get user info, make sure you call this method behind the ContextMiddleware")
	}

	conn, err := nx.GetClientGrpcConn(nex.ServiceTypeAuth)
	if err != nil {
		return fmt.Errorf("failed to get auth service client: %v", err)
	}
	_, err = proto.NewAuditServiceClient(conn).RecordEvent(context.Background(), &proto.RecordEventRequest{
		UserId:    uint64(user.ID),
		Action:    action,
		Metadata:  nex.EncodeMap(meta),
		Ip:        c.IP(),
		UserAgent: c.Get(fiber.HeaderUserAgent),
	})
	return err
}
