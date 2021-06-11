package version

import (
	log "github.com/sirupsen/logrus"
)

func CheckVersion(a, b string) {
	if a != b {
		log.Warn("⚠️ Unsupported version ", b)
		log.Debug("🌊 HelmWave version ", a)
	}
}

var Version = "dev"
