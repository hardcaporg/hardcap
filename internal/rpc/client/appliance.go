package client

import (
    "github.com/hardcaporg/hardcap/internal/rpc/common"
    "net/rpc"
)

type Appliance struct {
    c *rpc.Client
}

func ApplianceClient(c *rpc.Client) *Appliance {
    return &Appliance{c: c}
}

func (obj *Appliance) List(args common.ApplianceListArgs, reply *common.ApplianceListReply) error {
    return obj.c.Call("Appliance.List", args, reply)
}

func (obj *Appliance) Register(args common.ApplianceRegisterArgs, reply *common.ApplianceRegisterReply) error {
    return obj.c.Call("Appliance.Register", args, reply)
}