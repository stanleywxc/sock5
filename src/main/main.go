package main

import (
    "log"
    "net"
    "sync"
    "socks/context"    
    "socks/session"
)

func main () {
    
    log.Printf("Starting....\n")
  
    listener, err := net.Listen("tcp4", net.JoinHostPort("0.0.0.0" , "9090"))
    
    if (err != nil) {
        log.Printf("Error : %s", err)
    }
   
    for { 
        connection, _ := listener.Accept()
 
        go handleIncoming(connection)
        
    }
}

func handleIncoming(conn net.Conn) {
    
    log.Printf("Incomming: %s, Remote Addr: %s\n", conn.LocalAddr().Network(), conn.RemoteAddr().String())
    
    // create the context 
    contxt, err := context.New(conn)
    
    // Check the context is valid
    if (contxt == nil) {
        log.Printf("error during create Context: %s\n", err.Error())
        conn.Close()
        return;
    }
    
    log.Printf("Context created, version: %d\n", contxt.Version())
        
    // create a session
    var sessionCtxt *session.SessionContext = session.New(contxt)
    
    // start the session
    err = sessionCtxt.Start()
    
    // If there is any error, log it
    if (err != nil) {
        log.Printf("Session create failed: %s\n", err.Error())
    }
    
    // Done
    conn.Close()
}
