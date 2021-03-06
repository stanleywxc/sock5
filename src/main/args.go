//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//---------------------------------------------------------

package main

import (
	"fmt"
	"os"
)

const helpMessage = "Usage: main --help | -f conf-file [start | stop]\nWhen 'start' will start it as a daemon process, 'stop' will stop the daemon\nIf no 'start' provided, then it will start as normal process"

// Args holds all the commandline parameters
type Args struct {
	args map[string]string
}

// Get method: returns the commandline key/value
func (arg *Args) Get(key string) string {
	return arg.args[key]
}

func parseArgs() (*Args, string) {

	var msg string = ""
	var arg Args = Args{}

	arg.args = make(map[string]string)

	arg.args["self"] = os.Args[0]

	for i := 1; i < len(os.Args); i++ {

		switch os.Args[i] {
		case "-f":
			if (i + 1) >= len(os.Args) {
				msg = "Param '-f' is provided, but missing '-f' value\nExample usage: 'main -f /tmp/sock5.conf'\nUsing switch '--help' for help info"
				return &arg, msg
			}
			arg.args[os.Args[i]] = os.Args[i+1]
			i++
			break
		case "--help":
			return &arg, helpMessage
		case "start":
			arg.args["cmd"] = "start"
			break
		case "stop":
			arg.args["cmd"] = "stop"
			break
		case "daemon":
			arg.args["cmd"] = "daemon"
			break
		default:
			msg = fmt.Sprintf("Unsupported switch: '%s'\nTo get help info: 'main --help'\n", os.Args[i])
			return &arg, msg
		}
	}

	return &arg, msg
}
