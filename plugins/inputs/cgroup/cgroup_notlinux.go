// +build !linux

package cgroup

import (
	"github.com/MadDogTechnology/telegraf"
)

func (g *CGroup) Gather(acc telegraf.Accumulator) error {
	return nil
}
