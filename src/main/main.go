package main

import (
    "fmt"
    "net"
    "time"
    "socks/context"    
    "socks/session"
)

func main () {
    
    fmt.Printf("Starting....\n")
    
    listener, err := net.Listen("tcp4", net.JoinHostPort("127.0.0.1" , "9090"))
    
    if (err != nil) {
        fmt.Printf("Error : %s", err)
    }
   
    for { 
        connection, _ := listener.Accept()
        go handleIncoming(connection)
        
        time.Sleep(time.Duration(90000000))
       
    }
}

func handleIncoming(conn net.Conn) {
    
    fmt.Printf("Incomming: %s, Remote Addr: %s\n", conn.LocalAddr().Network(), conn.RemoteAddr().String())
   
    // create the context 
    contxt, err := context.New(conn)
    
    // Check the context is valid
    if (contxt == nil) {
        fmt.Printf("error during create Context: %s\n", err.Error())
        conn.Close()
        return;
    }
    
    fmt.Printf("Context created, version: %d\n", contxt.Version())
        
    // create a session
    var sessionCtxt *session.SessionContext = session.New(contxt)
    
    // start the session
    err = sessionCtxt.Start()
    
    // If there is any error, log it
    if (err != nil) {
        fmt.Printf("Session create failed: %s\n", err.Error())
    }
    
    // Done
    conn.Close()
}
