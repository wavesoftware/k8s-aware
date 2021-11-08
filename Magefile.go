//go:build mage
// +build mage

package main

import (

	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
	"github.com/wavesoftware/go-magetasks/pkg/artifact/platform"
	"github.com/wavesoftware/go-magetasks/pkg/checks"
	"github.com/wavesoftware/go-magetasks/pkg/git"
	"github.com/wavesoftware/k8s-aware/pkg/metadata"
)

// Default target is set to binary.
//goland:noinspection GoUnusedGlobalVariable
var Default = magetasks.Build // nolint:deadcode,gochecknoglobals

func init() { //nolint:gochecknoinits
	im := artifact.Image{
		Metadata:      config.Metadata{Name: "service"},
		Architectures: []platform.Architecture{platform.AMD64},
	}
	magetasks.Configure(config.Config{
		Version: &config.Version{
			Path:     metadata.VersionPath(),
			Resolver: git.NewVersionResolver(),
		},
		Artifacts: []config.Artifact{im},
		Checks:    []config.Task{checks.GolangCiLint()},
	})
}
