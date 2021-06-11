// Copyright 2020 Nokia
// Licensed under the BSD 3-Clause License.
// SPDX-License-Identifier: BSD-3-Clause

package runtime

import (
	"context"
	"log"
	"time"

	"github.com/srl-labs/containerlab/runtime/containerd"
	"github.com/srl-labs/containerlab/runtime/docker"
	"github.com/srl-labs/containerlab/types"
)

const (
	DockerRuntime     = "docker"
	ContainerdRuntime = "containerd"
)

type ContainerRuntime interface {
	// Set the network management details (generated by the config.go)
	SetMgmtNet(*types.MgmtNet)
	// Create container (bridge) network
	CreateNet(context.Context) error
	// Delete container (bridge) network
	DeleteNet(context.Context) error
	// Pull container image if not present
	PullImageIfRequired(context.Context, string) error
	// Create container
	CreateContainer(context.Context, *types.Node) error
	// Start pre-created container by its name
	StartContainer(context.Context, string) error
	// Stop running container by its name
	StopContainer(context.Context, string, *time.Duration) error
	// List all containers matching labels
	ListContainers(context.Context, []*types.GenericFilter) ([]types.GenericContainer, error)
	// Inspect container (extract its PID)
	ContainerInspect(context.Context, string) (*types.GenericContainer, error)
	// Get a netns path using the pid of a container
	GetNSPath(context.Context, string) (string, error)
	// Executes cmd on container identified with id and returns stdout, stderr bytes and an error
	Exec(context.Context, string, []string) ([]byte, []byte, error)
	// ExecNotWait executes cmd on container identified with id but doesn't wait for output nor attaches stodout/err
	ExecNotWait(context.Context, string, []string) error
	// Delete container by its name
	DeleteContainer(context.Context, *types.GenericContainer) error
}

func NewRuntime(name string, d bool, dur time.Duration, gracefulShutdown bool) ContainerRuntime {
	switch name {
	case DockerRuntime:
		return docker.NewDockerRuntime(d, dur, gracefulShutdown)
	case ContainerdRuntime:
		return containerd.NewContainerdRuntime(d, dur, gracefulShutdown)
	}

	log.Fatalf("Unexpected runtime name: %s", name)
	return nil
}
