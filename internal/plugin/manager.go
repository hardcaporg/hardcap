package plugin

import (
    "context"
    "fmt"
    "github.com/hardcaporg/hardcap/internal/ctxval"
    "os"
)

var (
    Appliance *AppliancePlugin
)

func StartAll(ctx context.Context) error {
    var err error

    err = os.Setenv("PLUGIN", "1")
    if err != nil {
        return fmt.Errorf("cannot set plugin env variable: %w", err)
    }

    Appliance, err = StartAppliance(ctx, "python3", "plugins/appliance-libvirt/appliance.py")
    if err != nil {
        return fmt.Errorf("cannot start plugin: %w", err)
    }

    return nil
}

func StopAll(ctx context.Context) {
    logger := ctxval.Logger(ctx)

    if Appliance != nil {
        err := Appliance.Stop(ctx)
        if err != nil {
            logger.Warn().Err(err).Msg("Error when stopping plugin")
        }
    }
}
