package release

import (
	"context"
	"strings"

	"github.com/invopop/jsonschema"
)

// PendingStrategy is a type for enumerating strategies for handling pending releases.
type PendingStrategy string

const (
	// PendingStrategyRollback rolls back pending release.
	PendingStrategyRollback PendingStrategy = "rollback"
	// PendingStrategyUninstall uninstalls pending release.
	PendingStrategyUninstall PendingStrategy = "uninstall"
)

func (PendingStrategy) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type: "string",
		Enum: []any{
			PendingStrategyRollback,
			PendingStrategyUninstall,
			"",
		},
	}
}

// Helm outputs error into fmt.Errorf which makes errors.Is unusable.
// So we have to find this substring in error string.
const errNoVersionToRollback = "release has no 0 version"

func (rel *config) isPending() (bool, error) {
	status, err := rel.Status()
	if err != nil {
		return false, err
	}

	return status.Info.Status.IsPending(), nil
}

func (rel *config) fixPending(ctx context.Context) error {
	switch rel.PendingReleaseStrategy {
	case PendingStrategyRollback:
		err := rel.Rollback(ctx, 0)

		// If no version to rollback, uninstall the release
		if strings.Contains(err.Error(), errNoVersionToRollback) {
			_, err = rel.Uninstall(ctx)
		}

		return err
	case PendingStrategyUninstall:
		_, err := rel.Uninstall(ctx)

		return err
	default:
		return ErrPendingRelease
	}
}
