package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"bytes"
	"flag"
	"golang.org/x/crypto/ssh/knownhosts"
)

func main(){
	//var hostKey ssh.PublicKey
	// An SSH client is represented with a ClientConn.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	
	hostname := flag.String("host", "", "SSH remote host")
	port := flag.String("port", "20", "SSH remote port")

	username := flag.String("user", "ubuntu", "Remote hostname")
	sshPassword := flag.String("ssh-password", "", "Remote SSH password")

	//password := flag.String("user-password", "", "Remote user password")
	cmd := flag.String("cmd", "ls", "Commands to run during ssh session")

	flag.Parse()
	
	keyCallBack, err := knownhosts.New("/home/adonese/.ssh/known_hosts")
	if err != nil {
		log.Printf("Error in knownhosts: %v", err)
	}

	config := &ssh.ClientConfig{
		User: *username,
		Auth: []ssh.AuthMethod{
			ssh.Password(*sshPassword),
		},
		HostKeyCallback: keyCallBack,
	}
	client, err := ssh.Dial("tcp", *hostname + ":" + *port, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(*cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
}
