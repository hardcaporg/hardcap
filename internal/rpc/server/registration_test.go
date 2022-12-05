package server

import (
    "github.com/hardcaporg/hardcap/internal/db"
    "github.com/stretchr/testify/require"
    "sync"
    "testing"
)

var setupOnce = sync.Once{}

func setup() {
    db.Initialize()
}

func TestList(t *testing.T) {
    setupOnce.Do(setup)
    
    r := Registration{}
    reply := &ApplianceListReply{}
    err := r.List(ApplianceListArgs{
        Limit:  0,
        Offset: 0,
    }, reply)
    require.Nil(t, err)
    require.Empty(t, reply.List)
}