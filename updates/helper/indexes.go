package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/safing/portbase/updater"
)

// Release Channel Configuration Keys.
const (
	ReleaseChannelKey     = "core/releaseChannel"
	ReleaseChannelJSONKey = "core.releaseChannel"
)

// Release Channels.
const (
	ReleaseChannelStable  = "stable"
	ReleaseChannelBeta    = "beta"
	ReleaseChannelStaging = "staging"
	ReleaseChannelSupport = "support"
)

// SetIndexes sets the update registry indexes and also configures the registry
// to use pre-releases based on the channel.
func SetIndexes(registry *updater.ResourceRegistry, releaseChannel string, deleteUnusedIndexes bool) (warning error) {
	usePreReleases := false

	// Be reminded that the order is important, as indexes added later will
	// override the current release from earlier indexes.

	// Reset indexes before adding them (again).
	registry.ResetIndexes()

	// Always add the stable index as a base.
	registry.AddIndex(updater.Index{
		Path: ReleaseChannelStable + ".json",
	})

	// Add beta index if in beta or staging channel.
	indexPath := ReleaseChannelBeta + ".json"
	if releaseChannel == ReleaseChannelBeta ||
		releaseChannel == ReleaseChannelStaging ||
		(releaseChannel == "" && indexExists(registry, indexPath)) {
		registry.AddIndex(updater.Index{
			Path:       indexPath,
			PreRelease: true,
		})
		usePreReleases = true
	} else if deleteUnusedIndexes {
		err := deleteIndex(registry, indexPath)
		if err != nil {
			warning = fmt.Errorf("failed to delete unused index %s: %w", indexPath, err)
		}
	}

	// Add staging index if in staging channel.
	indexPath = ReleaseChannelStaging + ".json"
	if releaseChannel == ReleaseChannelStaging ||
		(releaseChannel == "" && indexExists(registry, indexPath)) {
		registry.AddIndex(updater.Index{
			Path:       indexPath,
			PreRelease: true,
		})
		usePreReleases = true
	} else if deleteUnusedIndexes {
		err := deleteIndex(registry, indexPath)
		if err != nil {
			warning = fmt.Errorf("failed to delete unused index %s: %w", indexPath, err)
		}
	}

	// Add support index if in support channel.
	indexPath = ReleaseChannelSupport + ".json"
	if releaseChannel == ReleaseChannelSupport ||
		(releaseChannel == "" && indexExists(registry, indexPath)) {
		registry.AddIndex(updater.Index{
			Path: indexPath,
		})
	} else if deleteUnusedIndexes {
		err := deleteIndex(registry, indexPath)
		if err != nil {
			warning = fmt.Errorf("failed to delete unused index %s: %w", indexPath, err)
		}
	}

	// Add the intel index last, as it updates the fastest and should not be
	// crippled by other faulty indexes. It can only specify versions for its
	// scope anyway.
	registry.AddIndex(updater.Index{
		Path: "all/intel/intel.json",
	})

	// Set pre-release usage.
	registry.SetUsePreReleases(usePreReleases)

	return warning
}

func indexExists(registry *updater.ResourceRegistry, indexPath string) bool {
	_, err := os.Stat(filepath.Join(registry.StorageDir().Path, indexPath))
	return err == nil
}

func deleteIndex(registry *updater.ResourceRegistry, indexPath string) error {
	err := os.Remove(filepath.Join(registry.StorageDir().Path, indexPath))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
