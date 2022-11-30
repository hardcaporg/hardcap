package server

import (
    "github.com/hardcaporg/hardcap/internal/db"
    "github.com/hardcaporg/hardcap/internal/model"
)

type Registration struct{}

type RegistrationsListArgs struct {
    Limit, Offset int
}

type RegistrationsListReply struct {
    List []*model.Registration
}

func (s *Registration) List(args RegistrationsListArgs, reply *RegistrationsListReply) error {
    var err error
    if args.Limit == 0 {
        args.Limit = 100
    }

    reply.List, _, err = db.Registration.FindByPage(args.Limit, args.Offset)
    if err != nil {
        return err
    }

    return nil
}