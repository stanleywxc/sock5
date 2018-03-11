package main

import (
        "fmt"
        "socks/log"
        "socks/config"
)

func main () {
        
    // parse args, only support '-f' now
    args, msg := parseArgs()
    
    if (len(msg) != 0) {
        fmt.Printf("%s\n", msg)
        return 
    }
    
    // initialization
    config := config.Initialize(args.Get("-f"))
    
    // Set log Level and log file path
    log.SetLevel(log.Level(config.Log.Level))
    log.SetOutput(config.Log.Path)
  
    // create a server instance
    server := New(config)
    
    log.Infof("Socks5 server is starting....\n")

    // Start the server   
    if (server.Start() != true) {
        log.Errorf("Statring socks failed\n")
    }
}