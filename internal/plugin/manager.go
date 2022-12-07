package plugin

import (
    "context"
    "fmt"
    "github.com/hardcaporg/hardcap/internal/ctxval"
)

var (
    Appliance *AppliancePlugin
)

func StartAll(ctx context.Context) error {
    var err error

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
