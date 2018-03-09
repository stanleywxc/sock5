package main

import (
    "log"
    "net"
    "socks/context"    
    "socks/session"
)

func main () {
    
    log.Printf("Starting....\n")
  
    listener, err := net.Listen("tcp", net.JoinHostPort("" , "9090"))
    
    if (err != nil) {
        log.Printf("Error : %s", err)
    }
   
    for {
        
        // Accept incoming connections
        connection, err := listener.Accept()
        
        if (err != nil) {
            
            log.Printf("Error in accepting incoming connection: %s\n", err.Error())
            continue
        }
        
        // Is it a valid connection?
        if (connection == nil) {
            log.Printf("Error in accepting incoming connection: connection is nil\n")            
            continue
        }
 
        // Handle the incoming connections.
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
