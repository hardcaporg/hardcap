package client

import (
    "github.com/hardcaporg/hardcap/internal/rpc/common"
    "net/rpc"
)

type Registration struct {
    c *rpc.Client
}

func RegistrationClient(c *rpc.Client) *Registration {
    return &Registration{c: c}
}

func (obj *Registration) List(args common.RegistrationListArgs, reply *common.RegistrationListReply) error {
    return obj.c.Call("System.List", args, reply)
}