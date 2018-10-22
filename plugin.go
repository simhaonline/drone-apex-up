package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type (
	// Config for the plugin.
	Config struct {
		AccessKey string
		SecretKey string
		Stage     []string
	}

	// Plugin values
	Plugin struct {
		Config Config
	}
)

// Exec executes the plugin.
func (p *Plugin) Exec() error {
	if p.Config.AccessKey != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", p.Config.AccessKey)
	}
	if p.Config.SecretKey != "" {
		os.Setenv("AWS_SECRET_ACCESS_KEY", p.Config.SecretKey)
	}
	commands := []*exec.Cmd{
		p.versionCommand(),
	}

	for _, v := range p.Config.Stage {
		commands = append(commands, p.deployCommand(v))
	}

	for _, cmd := range commands {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = os.Environ()

		trace(cmd)

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Plugin) deployCommand(stage string) *exec.Cmd {
	args := []string{
		"deploy",
		stage,
	}

	return exec.Command(
		"up",
		args...,
	)
}

func (p *Plugin) versionCommand() *exec.Cmd {
	args := []string{
		"version",
	}

	return exec.Command(
		"up",
		args...,
	)
}

func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
