package sleep

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const defaultRandomMaxSeconds = 60

type Sleeper interface {
	Sleep(ctx context.Context, in types.Number)
}

type defaultSleeper struct {
	maxDefault int64
}

func NewSleeper(def types.Number) Sleeper {
	m := NumberToInt64(def)
	if m <= 0 {
		m = rand.Int63n(defaultRandomMaxSeconds)
	}
	return &defaultSleeper{maxDefault: m}
}

func (d *defaultSleeper) Sleep(ctx context.Context, in types.Number) {
	m := NumberToInt64(in)
	var random bool
	if m <= 0 {
		random = true
		m = rand.Int63n(d.maxDefault)
	}
	sleepFor := time.Duration(m) * time.Second
	tflog.Info(ctx, fmt.Sprintf("sleeping for %s (random: %t)", sleepFor, random))
	<-time.After(sleepFor)
}

func NumberToInt64(number types.Number) int64 {
	if number.IsUnknown() || number.IsNull() {
		return 0
	}
	s, _ := number.ValueBigFloat().Int64()
	return s
}
