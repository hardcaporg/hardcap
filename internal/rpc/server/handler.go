package server

import (
    "context"
    "github.com/hardcaporg/hardcap/internal/config"
    "github.com/hardcaporg/hardcap/internal/ctxval"
    "net"
    "net/http"
    "net/rpc"
)

func Initialize(ctx context.Context) error {
    log := ctxval.Logger(ctx)

    err := rpc.Register(new(Registration))
    if err != nil {
        return err
    }

    err = rpc.Register(new(Appliance))
    if err != nil {
        return err
    }

    rpc.HandleHTTP()
    l, listenErr := net.Listen("tcp", config.Application.RpcListenAddress)
    if listenErr != nil {
        return listenErr
    }
    go func() {
        log.Info().Msgf("Starting new %s RPC instance", config.Application.RpcListenAddress)

        err := http.Serve(l, nil)
        if err != nil {
            panic(err)
        }
    }()

    return nil
}