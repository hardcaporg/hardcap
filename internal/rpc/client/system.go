package client

import (
    "github.com/hardcaporg/hardcap/internal/rpc/common"
    "net/rpc"
)

type System struct {
    c *rpc.Client
}

func SystemClient(c *rpc.Client) *System {
    return &System{c: c}
}

func (obj *System) List(args common.SystemListArgs, reply *common.SystemListReply) error {
    return obj.c.Call("System.List", args, reply)
}
