package discovery_test

import (
	"github.com/hmuendel/deputyl/discovery"
	"testing"
)

func TestNewVersions(t *testing.T) {
	plainUpsteamVersions := []string{
		"0.1.0", "0.1.1", "0.2.0", "0.2.1", "1.0.0", "1.0.1", "1.1.0", "1.1.1",
	}
	prefixedUpstreamVersions := []string{
		"v0.1.0", "v0.1.1", "v0.2.0", "v0.2.1", "v1.0.0", "v1.0.1", "v1.1.0", "v1.1.1",
	}
	testCases := [...]struct {
		name            string
		version         string
		err             bool
		upsteamVersions []string
		patch           uint64
		minor           uint64
		major           uint64
	}{
		{"normal-semvers-all-newer", "0.1.0", false, plainUpsteamVersions, 1, 2, 4},
		{"normal-semvers-some-newer", "0.3.1", false, plainUpsteamVersions, 0, 0, 4},
		{"normal-semvers-one-newer", "1.1.0", false, plainUpsteamVersions, 1, 0, 0},
		{"prefixed-semvers-all-newer", "0.1.0", false, prefixedUpstreamVersions, 1, 2, 4},
		{"prefixed-semvers-some-newer", "0.3.1", false, prefixedUpstreamVersions, 0, 0, 4},
		{"prefixed-semvers-one-newer", "1.1.0", false, prefixedUpstreamVersions, 1, 0, 0},
		{"normal-semvers-all-newer-v", "v0.1.0", false, plainUpsteamVersions, 1, 2, 4},
		{"normal-semvers-some-newer-v", "v0.3.1", false, plainUpsteamVersions, 0, 0, 4},
		{"normal-semvers-one-newer-v", "v1.1.0", false, plainUpsteamVersions, 1, 0, 0},
		{"prefixed-semvers-all-newer-v", "v0.1.0", false, prefixedUpstreamVersions, 1, 2, 4},
		{"prefixed-semvers-some-newer-v", "v0.3.1", false, prefixedUpstreamVersions, 0, 0, 4},
		{"prefixed-semvers-one-newer-v", "v1.1.0", false, prefixedUpstreamVersions, 1, 0, 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patch, minor, major, err := discovery.NewerVersions(tc.version, tc.upsteamVersions)
			if err == nil && tc.err {
				t.Errorf("%s, should have errored", tc.name)
			}
			if err != nil && !tc.err {
				t.Errorf("%s, should not errored with: %v", tc.name, err)
			}
			if patch != tc.patch {
				t.Errorf("expected %v, got %v", tc.patch, patch)
			}
			if minor != tc.minor {
				t.Errorf("expected %v, got %v", tc.minor, minor)
			}
			if major != tc.major {
				t.Errorf("expected %v, got %v", tc.major, major)
			}
		})
	}

	// }
	// }
}
