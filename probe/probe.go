package probe

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/csduarte/mattermost-probe/config"
	"github.com/csduarte/mattermost-probe/mattermost"
	"github.com/csduarte/mattermost-probe/metrics"
	"github.com/pkg/errors"
)

// Probe represents a basic probe
type Probe interface {
	String() string
	Setup() error
	Start() error
}

// NewProbes takes configs, report channel and the client objects and creates each probe.
func NewProbes(cfg config.Config, rc chan metrics.Report, c1, c2 *mattermost.Client) []Probe {
	probes := []Probe{}
	if cfg.BroadcastProbe.Enabled {
		p := NewBroadcastProbe(cfg.BroadcastProbe, c1, c2)
		p.ReportChannel = rc
		probes = append(probes, p)
	}

	if cfg.ChannelJoinProbe.Enabled {
		p := NewChannelJoinProbe(cfg.ChannelJoinProbe, c1)
		probes = append(probes, p)
	}

	if cfg.PingProbe.Enabled {
		p := NewPingProbe(cfg.PingProbe, c1)
		probes = append(probes, p)
	}

	if cfg.SearchProbe.Enabled {
		p := NewSearchProbe(cfg.SearchProbe, c1)
		p.ReportChannel = rc
		probes = append(probes, p)
	}

	if cfg.LoginProbe.Enabled {
		p1 := NewLoginProbe(cfg.LoginProbe, c1, cfg.UserA)
		p2 := NewLoginProbe(cfg.LoginProbe, c2, cfg.UserB)
		probes = append(probes, p1, p2)
	}

	if cfg.APIPingProbe.Enabled {
		p := NewAPIPingProbe(cfg.APIPingProbe, c1)
		probes = append(probes, p)
	}

	return probes
}

// SetupProbes simply calls setup on each of the probes that were initialized
func SetupProbes(probes []Probe, log *logrus.Logger) error {
	for _, p := range probes {
		log.Infof("Setting up probe: %s", p.String())
		if err := p.Setup(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Could not start probe %s", p.String()))
		}
	}
	return nil
}

// StartProbes will call the start method of each probe, blocks app start up if any fail
func StartProbes(probes []Probe, log *logrus.Logger) error {
	for _, p := range probes {
		log.Infof("Starting probe: %s", p.String())
		if err := p.Start(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Could not start probe %s", p.String()))
		}
	}
	return nil
}
