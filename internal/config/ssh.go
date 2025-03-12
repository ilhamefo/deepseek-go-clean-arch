package config

import (
	"golang.org/x/crypto/ssh"
)

func NewSshClient(cfg *Config) *ssh.Client {
	sshConfig := &ssh.ClientConfig{
		User:            cfg.SshUsername,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.SshPassword),
		},
	}

	client, err := ssh.Dial("tcp", cfg.SshAddress, sshConfig)
	if client != nil {
		defer client.Close()
	}

	if err != nil {
		panic(err.Error())
	}

	return client
}
