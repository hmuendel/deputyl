package discovery

import (
	"errors"
	"strings"
	"time"

	"github.com/blang/semver"
	"github.com/hmuendel/deputyl/depcheck"
	"github.com/hmuendel/deputyl/metrics"
	"github.com/hmuendel/glog"
)

type Artifact struct {
	Name     string
	Metadata map[string]string
}

type Checker interface {
	Check(artifact string) (string, []string, error)
}

type Emitter interface {
	Emit(metrics.Metric)
}

type Discoverer interface {
	Discover() []Artifact
}

type Config struct {
	Interval  time.Duration `desc:"interval for discover and check new versions, defaults to 10s"`
	SkipPre   bool          `desc:"skip upstream pre release versions, defaults to false"`
	SkipBuild bool          `desc:"skip upstream build versions, defaults to false"`
}

var DefaultConfig = Config{
	Interval:  120 * time.Second,
	SkipPre:   false,
	SkipBuild: false,
}

func parseNonSemver(version string) (semver.Version, error) {
	preReleaseSplit := strings.Split(version, "-")
	buildSplit := strings.Split(preReleaseSplit[0], "+")
	versionSplit := strings.Split(buildSplit[0], ".")
	if len(versionSplit) == 2 {
		version = version + ".0"
	}
	if len(version) == 1 {
		version = version + ".0.0"
	}
	buildSplit[0] = version
	version = strings.Join(buildSplit, "+")
	preReleaseSplit[0] = version
	version = strings.Join(preReleaseSplit, "-")
	return semver.Make(version)
}

func NewerVersions(version string, upstreamVersions []string) (patch, minor, major uint64, err error) {
	hasSemvers := false
	version = strings.TrimPrefix(version, "v")
	v, err := semver.Make(version)
	if err != nil {
		if glog.V(3) {
			glog.Errorf("could not parse %v into semver: %v", version, err)
		}
		v, err = parseNonSemver(version)
		if err != nil {
			return 0, 0, 0, err
		}
	}
	for _, upstreamVersion := range upstreamVersions {
		upstreamVersion = strings.TrimPrefix(upstreamVersion, "v")
		uv, err := semver.Make(upstreamVersion)
		if err != nil {
			if glog.V(3) {
				glog.Errorf("could not parse upstream %v into semver: %v", version, err)
			}
			uv, err = parseNonSemver(upstreamVersion)
			if err != nil {
				if glog.V(3) {
					glog.Errorf("could not parse upstream %s, skipping", upstreamVersion)
				}
				continue
			}
		}
		hasSemvers = true
		if DefaultConfig.SkipPre && len(uv.Pre) > 0 {
			if glog.V(5) {
				glog.Infof("skipping pre release version %s", upstreamVersion)
			}
			continue
		}
		if DefaultConfig.SkipBuild && len(uv.Build) > 0 {
			if glog.V(5) {
				glog.Infof("skipping build version %s", upstreamVersion)
			}
			continue
		}
		if uv.LTE(v) {
			continue
		}
		if uv.Major > v.Major {
			major += 1
			continue
		}
		if uv.Minor > v.Minor {
			minor += 1
			continue
		}
		if uv.Patch > v.Patch {
			patch += 1
			continue
		}

	}
	if !hasSemvers {
		return 0, 0, 0, errors.New("none of the upsteam versions seems to be a semver")
	}
	return
}

type DiscoverEmitter struct {
	Checker    Checker
	Emitter    Emitter
	ticker     *time.Ticker
	Discoverer Discoverer
}

func NewDiscoverEmitter() DiscoverEmitter {
	de := DiscoverEmitter{}
	de.Checker = depcheck.NewDockerDepChecker()
	de.Emitter = metrics.NewPrometheusEmitter()
	de.ticker = time.NewTicker(DefaultConfig.Interval)
	de.Discoverer, _ = NewK8sDiscoverer()
	return de
}

func (de *DiscoverEmitter) StartDiscovery() {
	go de.discover()
}

func (de *DiscoverEmitter) discover() {
	for {
		select {
		case <-de.ticker.C:
			if glog.V(7) {
				glog.Infof("begin discovery step")
			}
			pods := de.Discoverer.Discover()
			if glog.V(8) {
				glog.Infof("discovered pods: %#v", pods)
			}
			for _, p := range pods {
				version, upsteamVersions, err := de.Checker.Check(p.Name)
				if err != nil {
					glog.Errorf("Image: %s: %v", p.Name, err)
					continue
				}
				if glog.V(10) {
					glog.Infof("%s.version: %s, %s.tags: %v", p.Name, version, p.Name, upsteamVersions)
				}
				patch, minor, major, err := NewerVersions(version, upsteamVersions)
				if err != nil {
					glog.Errorf("Image: %s: %v", p.Name, err)
					continue
				}
				if glog.V(10) {
					glog.Infof("newer version for %s: patch: %v, minor: %v, major: %v", p.Name, patch, minor, major)
				}
				m := metrics.Metric{
					PatchVersions: float64(patch),
					MinorVersions: float64(minor),
					MajorVersions: float64(major),
					Repository:    p.Name,
					Labels:        p.Metadata,
				}
				de.Emitter.Emit(m)
			}
		}
	}
}
