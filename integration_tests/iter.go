package integration_tests

import (
	"context"
	"os/exec"
	"testing"

	"beautifulthings/server"
	"beautifulthings/store"

	"github.com/stretchr/testify/require"
)

const addr = "localhost:8080"

func startRestServer(t *testing.T) func() {
	cancel, err := server.ServeRest(context.Background(), addr, store.NewInMemoryServer())
	require.NoError(t, err)
	return cancel
}

func runDockerServer(t *testing.T) func() {
	cmd := exec.Command("docker", "build", "-t", "beautifulthings:latest", "..")
	err := cmd.Run()
	require.NoError(t, err)

	cmd = exec.Command("docker", "run", "--rm", "beautifulthings")
	err = cmd.Start()
	require.NoError(t, err)
	cancel := func() {
		cmd.Process.Kill()
	}
	return cancel
}

type serverBuilderFunc func(t *testing.T) (server.Server, func())

type serverTest struct {
	tag string
	f   serverBuilderFunc
}

var serverBuilders = []serverTest{
	{
		"InMemory",
		func(_ *testing.T) (server.Server, func()) { return server.New(store.NewInMemoryServer()), nil },
	},
	{
		"Rest",
		func(t *testing.T) (server.Server, func()) {
			cancel := startRestServer(t)
			return server.NewRemoteRest("http://" + addr), cancel
		},
	},
	// {
	// 	"DockerRest",
	// 	func(t *testing.T) (server.Server, func()) {
	// 		cancel := runDockerServer(t)
	// 		return server.NewRemoteRest("http://" + addr), cancel
	// 	},
	// },
}

type serverIter struct {
	current int
}

func (si *serverIter) Next(t *testing.T) (server.Server, func(), bool) {
	if si.current >= len(serverBuilders) {
		return nil, nil, false
	}
	s, cancel := serverBuilders[si.current].f(t)
	si.current += 1
	return s, cancel, true
}

func Run(t *testing.T, test func(*testing.T, server.Server)) {
	s := server.New(store.NewInMemoryServer())

	it := &serverIter{}
	s, cancel, next := it.Next(t)
	for next {
		t.Run(t.Name()+serverBuilders[it.current-1].tag, func(t *testing.T) { test(t, s) })
		if cancel != nil {
			cancel()
		}
		s, cancel, next = it.Next(t)
	}
}
