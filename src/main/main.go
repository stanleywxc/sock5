//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// socks is a Socks V5 server, which implements SocksV5
// protocol.
//---------------------------------------------------------
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"socks/config"
	"socks/log"
	"strconv"
	"syscall"
)

var pidFile = "/tmp/socks5.pid"

func main() {

	// parse args, only support '-f' now
	args, msg := parseArgs()

	if len(msg) != 0 {
		fmt.Printf("%s\n", msg)
		os.Exit(0)
	}

	switch args.Get("cmd") {
	case "daemon":
		daemonize(args)
		break
	case "start":
		startDaemon(args)
		break
	case "stop":
		stopDaemon(args)
		break
	default:
		start(args)
		break
	}
}

func daemonize(args *Args) {

	// Remove pid file upon receiving signal SIGTERM
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, os.Kill, syscall.SIGTERM)

	// watcher on SIGTERM signal
	go func() {
		signalType := <-channel
		signal.Stop(channel)

		fmt.Printf("Received signal: %v. Exiting...\n", signalType)

		// remove pid file
		os.Remove(pidFile)

		os.Exit(0)
	}()

	// start the server
	if start(args) != true {
		os.Exit(1)
	}
}

func startDaemon(args *Args) {

	// check if daemon already running.
	_, err := os.Stat(pidFile)

	// The pid file has already been created, which means
	// the daemon is running
	if err == nil {
		fmt.Printf("Daemon '%s' has been already running\n", args.Get("self"))
		os.Exit(1)
	}

	fmt.Printf("Daemon is starting ...\n")

	// No daemon is running, trying to start it
	var cmd *exec.Cmd
	if len(args.Get("-f")) > 0 {
		cmd = exec.Command(args.Get("self"), "-f", args.Get("-f"), "daemon")
	} else {
		cmd = exec.Command(args.Get("self"), "daemon")
	}

	//cmd := exec.Command(args.Get("self"), command)
	cmd.Start()

	fmt.Printf("Daemon '%s' started (pid: %v)\n", args.Get("self"), cmd.Process.Pid)

	// Save the pid to file
	savePid(cmd.Process.Pid)

	// successful, exit
	os.Exit(0)
}

func stopDaemon(args *Args) {

	// when received the 'stop', check if pid file exists first
	// if it exists, read pid from file
	// then trying to stop the process.
	_, err := os.Stat(pidFile)

	// pid file exists?
	if err != nil {
		fmt.Printf("'%s' is not running, pid file doesn't exist\n", args.Get("self"))
		os.Exit(1)
	}

	// read the pid from file.
	bytes, err := ioutil.ReadFile(pidFile)
	if err != nil {
		fmt.Printf("'%s' is not running, unable to read pid file\n", args.Get("self"))
		os.Exit(1)
	}

	// get the pid
	pid, err := strconv.Atoi(string(bytes))

	// Any error?
	if err != nil {
		fmt.Printf("Invalid pid in: %s\n", pidFile)
		os.Exit(1)
	}

	// Find the process by pid
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Unable to find process: [%v] with error %v \n", pid, err)
		os.Exit(1)
	}

	// remove PID file
	os.Remove(pidFile)

	fmt.Printf("Shutdown the daemon '%s' with pid: [%v] ...\n", args.Get("self"), pid)

	// kill process and exit immediately
	err = process.Kill()

	// Any error?
	if err != nil {
		fmt.Printf("Unable to shutdown daemon: '%s'(pid: %v) with error %v\n", args.Get("self"), pid, err)
		os.Exit(1)
	}

	fmt.Printf("Daemon (pid: %v) shutdown successfully\n", pid)

	// Exit
	os.Exit(0)
}

func savePid(pid int) {

	file, err := os.Create(pidFile)
	if err != nil {
		log.Errorf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Errorf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	// flush it to disk
	file.Sync()
}

// start the socks server
func start(args *Args) bool {

	// initialization
	config := config.Initialize(args.Get("-f"))

	// Set log Level and log file path
	log.SetLevel(log.Level(config.Log.Level))
	log.SetOutput(config.Log.Path)

	// create a server instance
	server := New(config)

	log.Infof("Socks5 server is starting....\n")

	// Start the server
	if server.Start() != true {
		log.Errorf("Statring socks failed\n")
		return false
	}

	return true
}
