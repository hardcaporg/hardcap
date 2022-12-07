package server

import (
    "github.com/hardcaporg/hardcap/internal/db"
    "github.com/hardcaporg/hardcap/internal/rpc/common"
)

type System struct{}

func (s *System) List(args common.SystemListArgs, reply *common.SystemListReply) error {
    var err error
    if args.Limit == 0 {
        args.Limit = 100
    }

    reply.List, _, err = db.System.FindByPage(args.Offset, args.Limit)
    if err != nil {
        return err
    }

    return nil
}
