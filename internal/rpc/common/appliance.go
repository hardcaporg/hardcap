package common

import "github.com/hardcaporg/hardcap/internal/model"

type ApplianceRegisterArgs struct {
    Appliance *model.Appliance
}

type ApplianceRegisterReply struct {
    Appliance *model.Appliance
}

type ApplianceListArgs struct {
    Limit, Offset int
}

type ApplianceListReply struct {
    List []*model.Appliance
}
