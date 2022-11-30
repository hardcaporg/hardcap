package client

import (
    "github.com/hardcaporg/hardcap/internal/model"
    "net/rpc"
)

type Registration struct {
    c *rpc.Client
}

type RegistrationsListArgs struct {
    Limit, Offset int
}

type RegistrationsListReply struct {
    List []*model.Registration
}

func RegistrationClient(c *rpc.Client) *Registration {
    return &Registration{c: c}
}

func (obj *Registration) List(args RegistrationsListArgs, reply *RegistrationsListReply) error {
    return obj.c.Call("Service.List", args, reply)
}