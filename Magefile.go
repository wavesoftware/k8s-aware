//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/otiai10/copy"

	// mage:import
	"github.com/wavesoftware/go-magetasks"
	"github.com/wavesoftware/go-magetasks/config"
	"github.com/wavesoftware/go-magetasks/pkg/artifact"
	"github.com/wavesoftware/go-magetasks/pkg/artifact/platform"
	"github.com/wavesoftware/go-magetasks/pkg/checks"
	"github.com/wavesoftware/go-magetasks/pkg/files"
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
		Context: context.WithValue(context.TODO(),
			key{"image"}, artifact.ImageReferenceOf(im)),
	})
}

// Deploy will deploy service on a cluster.
// goland:noinspection GoUnusedExportedFunction
// nolint:deadcode
func Deploy() error {
	mg.Deps(magetasks.Publish)
	files.EnsureBuildDir()
	source := path.Join(files.ProjectDir(), "deploy")
	target := path.Join(files.BuildDir(), "deploy")
	err := copy.Copy(source, target)
	if err != nil {
		return err // nolint:wrapcheck
	}
	return withinDirectory(target, deployWithKubectl)
}

func deployWithKubectl() error {
	resolver, ok := config.Actual().Context.Value(key{"image"}).(config.Resolver)
	if !ok {
		return errNoResolver
	}
	err := replaceImage(resolver)
	if err != nil {
		return err
	}
	err = sh.RunV("kubectl", "apply", "-f", ".")
	if err != nil {
		return err // nolint:wrapcheck
	}
	return nil
}

func replaceImage(resolver config.Resolver) error {
	input, err := ioutil.ReadFile("100-service.yaml")
	if err != nil {
		return err // nolint:wrapcheck
	}

	contents := strings.ReplaceAll(string(input),
		"ghcr.io/wavesoftware/k8s-aware/service:latest",
		resolver())

	err = ioutil.WriteFile("100-service.yaml", []byte(contents), rwUser)
	return err // nolint:wrapcheck
}

func withinDirectory(path string, fn func() error) error {
	dir, err := os.Getwd()
	if err != nil {
		return err // nolint:wrapcheck
	}
	err = os.Chdir(path)
	if err != nil {
		return err // nolint:wrapcheck
	}
	defer func() {
		_ = os.Chdir(dir)
	}()
	return fn()
}

const rwUser = 0o600

var errNoResolver = errors.New("no resolver")

type key struct {
	name string
}
