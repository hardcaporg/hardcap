package server

import (
    "github.com/hardcaporg/hardcap/internal/db"
    "github.com/hardcaporg/hardcap/internal/rpc/common"
)

type Registration struct{}

func (s *Registration) List(args common.RegistrationListArgs, reply *common.RegistrationListReply) error {
    var err error
    if args.Limit == 0 {
        args.Limit = 100
    }

    reply.List, _, err = db.Registration.FindByPage(args.Offset, args.Limit)
    if err != nil {
        return err
    }

    return nil
}