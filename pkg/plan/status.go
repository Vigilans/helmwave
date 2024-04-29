package plan

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"

	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Status renders status table for list of releases names.
func (p *Plan) Status(names ...string) error {
	return status(p.body.Releases, names)
}

func status(all []release.Config, names []string) error {
	r := all

	if len(names) > 0 {
		sort.Strings(names)

		r = slices.DeleteFunc(r, func(rel release.Config) bool {
			return !slices.Contains(names, rel.Uniq().String())
		})
	}

	for _, rel := range r {
		l := rel.Logger()

		s, err := rel.Status()
		if err != nil {
			l.Errorf("Failed to get status: %v", err)

			continue
		}

		labels, err := json.Marshal(s.Labels)
		if err != nil {
			l.Errorf("Failed to get labels: %v", err)
		}

		values, err := json.Marshal(s.Config)
		if err != nil {
			l.Errorf("Failed to get values: %v", err)
		}

		log.WithFields(log.Fields{
			"name":          s.Name,
			"namespace":     s.Namespace,
			"chart":         fmt.Sprintf("%s-%s", s.Chart.Name(), s.Chart.Metadata.Version),
			"last deployed": s.Info.LastDeployed,
			"status":        s.Info.Status,
			"revision":      s.Version,
		}).Infof("General status of %s", rel.Uniq())

		log.WithFields(log.Fields{
			"notes":         s.Info.Notes,
			"labels":        string(labels),
			"chart sources": s.Chart.Metadata.Sources,
			"values":        string(values),
		}).Debugf("Debug status of %s", rel.Uniq())

		log.WithFields(log.Fields{
			"hooks":    s.Hooks,
			"manifest": s.Manifest,
		}).Tracef("Superdebug status of %s", rel.Uniq())
	}

	return nil
}
