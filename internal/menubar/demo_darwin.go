//go:build darwin

package menubar

import (
	_ "embed"

	"github.com/haritabh17/theirtime/internal/config"
	"github.com/haritabh17/theirtime/internal/team"
	"github.com/haritabh17/theirtime/internal/timeformat"
)

//go:embed demo/bob.png
var demoBobPNG []byte

//go:embed demo/ann.png
var demoAnnPNG []byte

const (
	demoBobID = "demo-bob"
	demoAnnID = "demo-ann"
)

func demoConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		cfg = &config.Config{}
	}
	config.ApplyDefaults(cfg)

	// Demo always showcases every display field regardless of saved prefs.
	showAvatar, showName, showTime := true, true, true
	cfg.ShowAvatar = &showAvatar
	cfg.ShowName = &showName
	cfg.ShowTime = &showTime
	cfg.Format24h = false
	cfg.TimePrecision = timeformat.PrecisionMinutes

	return cfg
}

func demoMembers() []team.MemberTime {
	return []team.MemberTime{
		{Label: "bob", ID: demoBobID, TZ: "America/Los_Angeles", DisplayName: "bob"},
		{Label: "ann", ID: demoAnnID, TZ: "America/New_York", DisplayName: "ann"},
	}
}

func demoAvatars() map[string][]byte {
	return map[string][]byte{
		demoBobID: append([]byte(nil), demoBobPNG...),
		demoAnnID: append([]byte(nil), demoAnnPNG...),
	}
}
