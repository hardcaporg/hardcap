package plugin

import (
    "context"
    "errors"
    "fmt"
    "github.com/hardcaporg/hardcap/internal/ctxval"
    "github.com/lzap/cborpc/cmd"
    "github.com/rs/zerolog"
    "os/exec"
)

type AppliancePlugin struct {
    proc *cmd.Command
}

func logger(ctx context.Context) *zerolog.Logger {
    logger := ctxval.Logger(ctx).With().Str("plugin", "appliance").Logger()
    return &logger
}

func StartAppliance(ctx context.Context, commandName string, commandArgs ...string) (*AppliancePlugin, error) {
    log := logger(ctx)

    log.Debug().Msgf("Starting plugin %s %v", commandName, commandArgs)
    proc, err := cmd.NewCommand(WithPluginLogger(ctx, log), commandName, commandArgs...)
    if err != nil {
        return nil, fmt.Errorf("cannot initialize plugin: %w", err)
    }
    err = proc.Start(ctx)
    if err != nil {
        return nil, fmt.Errorf("cannot start plugin: %w", err)
    }

    return &AppliancePlugin{
        proc: proc,
    }, nil
}

type ApplianceMultiplyArgs struct {
    A, B int
}

type ApplianceMultiplyReply struct {
    C int
}

func (plugin *AppliancePlugin) Multiply(ctx context.Context, args *ApplianceMultiplyArgs) (*ApplianceMultiplyReply, error) {
    reply := &ApplianceMultiplyReply{}

    err := plugin.proc.Call(ctx, "Arith.Multiply", args, &reply)
    if err != nil {
        return nil, err
    }

    return reply, nil
}

type ApplianceEnlistArgs struct {
    URL string
    NamePattern string
}

type EnlistedSystem struct {
    Name string
    UID string
    SerialNumber string
    MACs []string
}

type ApplianceEnlistReply struct {
    Systems []EnlistedSystem
}

func (plugin *AppliancePlugin) Enlist(ctx context.Context, args *ApplianceEnlistArgs) (*ApplianceEnlistReply, error) {
    reply := &ApplianceEnlistReply{}

    err := plugin.proc.Call(ctx, "Appliance.Enlist", args, &reply)
    if err != nil {
        return nil, err
    }

    return reply, nil
}

func (plugin *AppliancePlugin) Stop(ctx context.Context) error {
    log := logger(ctx)

    log.Debug().Msg("Stopping plugin")
    err := plugin.proc.Stop(ctx)
    var exitErr *exec.ExitError
    if errors.As(err, &exitErr) {
        log.Info().Msgf("Plugin exited with code %d", exitErr.ProcessState.ExitCode())
    } else if err != nil {
        return fmt.Errorf("cannot stop plugin: %w", err)
    }

    return nil
}
