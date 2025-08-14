package main

import (
	"log"
	"os/exec"
)

func updateCommand() {
	cmd := exec.Command("/root/repetition/update-cache-notifier.sh")
	if err := cmd.Run(); err != nil {
		log.Println("Could not execute update command", err)
	}
}
