package model

import (
    "github.com/stretchr/testify/require"
    "testing"
)

func TestNewSystemIDEmpty(t *testing.T) {
    sid := NewSystemID("")
    require.Equal(t, "NWcliVXI4KmQBs7BuciYsB", sid.Short)
    require.Equal(t, "VhrU4tRmVSKTTuJZkHkrnQSuvlJy0vvmUwBfxUIihdc", sid.Long)
    require.Equal(t, "Sonia Maroni", sid.FriendlyName)
}

func TestNewSystemIDWithMAC(t *testing.T) {
    sid := NewSystemID("", "52:54:00:f4:0e:79")
    require.Equal(t, "s2gwEjO2fm66Vm5DDvd4sD", sid.Short)
    require.Equal(t, "LM67JzFeGWSg3ZrJqmXfkhzdxhEziz03cWJ9CuLNto8", sid.Long)
    require.Equal(t, "Jake Kallevig", sid.FriendlyName)
}

func TestNewSystemIDWithSerial(t *testing.T) {
    sid := NewSystemID("G6BE94500313", "52:54:00:f4:0e:79")
    require.Equal(t, "dSfGY0NWq3yLc83S71zVmD", sid.Short)
    require.Equal(t, "eB8cl9Be8wlDjLard4CywG291pLRLbuk0KuvZllefPzB", sid.Long)
    require.Equal(t, "Celia Porte", sid.FriendlyName)
}