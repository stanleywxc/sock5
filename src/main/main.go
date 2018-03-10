package main

import (
        "socks/log"
        "socks/config"
)

func main () {
    
    log.Infof("Starting....\n")
    
    // parse args, only support '-f' now
    args := parseArgs()
    
    // initialization
    config := config.Initialize(args.Get("-f"))
    
    // Set log Level and log file path
    log.SetLevel(log.Level(config.Log.Level))
    log.SetOutput(config.Log.Path)
  
    // create a server instance
    server := New(config)
    
    // Start the server   
    if (server.Start() != true) {
        log.Errorf("Statring socks failed\n")
    }
}