package tssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trzsz/ssh_config"
)

func TestForkDefaultUserWithCanonicalConfig(t *testing.T) {
	config, err := ssh_config.DecodeBytes([]byte(`
Host configured
    User deploy
Host canonical.example
    User canonical-user
`))
	require.NoError(t, err)
	originalConfig := userConfig
	userConfig = &tsshConfig{config: &sshConfig{config: config}}
	t.Cleanup(func() { userConfig = originalConfig })

	for _, tt := range []struct {
		name string
		args sshArgs
		want string
	}{
		{"default", sshArgs{Destination: "missing"}, "root"},
		{"explicit", sshArgs{Destination: "configured"}, "deploy"},
		{"canonical", sshArgs{Destination: "short", canonicalDest: "canonical.example"}, "canonical-user"},
		{"canonical-default", sshArgs{Destination: "short", canonicalDest: "missing.example"}, "root"},
		{"original-precedence", sshArgs{Destination: "configured", canonicalDest: "canonical.example"}, "deploy"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getConfig(&tt.args, "User"))
		})
	}
}

func TestForkHostLookupWithArgumentContext(t *testing.T) {
	config, err := ssh_config.DecodeBytes([]byte(`
Host alpha
    HostName alpha.example
    User deploy
    #!! GroupLabels production
Host beta
    HostName beta.example
`))
	require.NoError(t, err)
	originalConfig, originalCallbacks := userConfig, afterLoginFuncs
	userConfig = &tsshConfig{config: &sshConfig{config: config}}
	afterLoginFuncs = nil
	t.Cleanup(func() { userConfig, afterLoginFuncs = originalConfig, originalCallbacks })

	args := &sshArgs{NoSave: true, FileManager: true}
	require.NoError(t, args.Option.UnmarshalText([]byte("Port=2200")))
	hosts := getAllHosts(args)
	require.Len(t, hosts, 2)
	for i, host := range hosts {
		assert.Equal(t, i, host.Index)
		assert.Equal(t, host.Alias, host.Args.Destination)
		assert.Equal(t, host.Alias, host.Args.originalDest)
		assert.True(t, host.Args.NoSave)
		assert.True(t, host.Args.FileManager)
		assert.Equal(t, "2200", host.Port)
	}
	assert.Equal(t, "production", hosts[0].GroupLabels)
	assert.Equal(t, "root", hosts[1].User)
	assert.Equal(t, "alpha", findHostInConfig(args, "alpha.example"))
	assert.Same(t, hosts[0], getExistingHostConfig(args, "alpha"))
	assert.False(t, needsHostUpdate(args, "alpha", "deploy", "alpha.example", "2200", "production"))
	assert.True(t, needsHostUpdate(args, "alpha", "deploy", "alpha.example", "2200", "staging"))
}
