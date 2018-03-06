package handshake

import (
    "fmt"
    "bufio"
    "errors"
    "socks"
    "socks/authentication"
    "socks/context"
)

type Handshake struct {
    version			byte
    nmethods			byte
    methods			[]byte
    authenticator	authentication.Authenticator
    context			*context.Context
}

func New(context *context.Context) (*Handshake) {
    return &Handshake {context : context}
}

func (handshake *Handshake)Handshake() (byte, error) {
    
    var code 	byte
    var err		error
    
    code, err = handshake.methodNegotiation()
    if (err != nil) {
        fmt.Printf(" Method Negotiation failed, error: %s\n", err.Error())
        return code, err
    }
    
    _, err = handshake.authentication()
    if (err != nil) {
        fmt.Printf(" Authentication Negotiation failed, error: %s\n", err.Error())
        return socks.SOCKS_AUTH_NOACCEPTABLE, err
    }
    
    return code, nil    
}

func (handshake *Handshake)methodNegotiation() (byte, error) {    
    
    var err		error
    var reader	*bufio.Reader = handshake.context.Reader()
       
    // Get the version, 1st byte
    handshake.version, err = reader.ReadByte()

    // Error?    
    if err != nil {
        return socks.SOCKS_AUTH_NOACCEPTABLE, err
    }

    fmt.Printf("version : %d\n", handshake.version)
 
    // how many methods client support?
    handshake.nmethods, err = reader.ReadByte()
    
    // can't read the nmethods
    if err != nil {
        
        // Send back the error code.
        return socks.SOCKS_AUTH_NOACCEPTABLE, err
    }
    
    fmt.Printf("nmethods: %d\n", handshake.nmethods)

    // Get the methods clients supports    
    handshake.methods = make([]byte, handshake.nmethods)
    
    var index byte
    for index = 0; index < handshake.nmethods; index++ {
        handshake.methods[index], err = reader.ReadByte()

        fmt.Printf("method[%d]: %d\n", index, handshake.methods[index])        
        if (err != nil) {
            break
        }
    }
    
    // Return false when no methods, and error happened
    // If there are some methods, but error occured, continue.
    if (len(handshake.methods) == 0) {
        
        // Send back the error code
        return socks.SOCKS_AUTH_NOACCEPTABLE, errors.New("No acceptable method")
    }
    
    // Find out if we support the methods.
    var found byte = socks.SOCKS_AUTH_NOACCEPTABLE
    
    for _, method := range handshake.methods {
        
        var authenticator authentication.Authenticator
        
        // Trying to find one:
        authenticator = authentication.AUTHENTICATORS[method]
        if (authenticator != nil) {
            // set the authenticator
            handshake.authenticator 	= authenticator
            found					= method
            break
        }
    }
 
    err = nil
    if (found == socks.SOCKS_AUTH_NOACCEPTABLE) {
        err = errors.New("No acceptbale method")
    }
    // return the 
    return found, err
}

func (handshake *Handshake)authentication () (bool, error) {
    
    return handshake.authenticator.Authenticate() 
}

