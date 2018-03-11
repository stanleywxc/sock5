//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//---------------------------------------------------------

package handshake

import (
    "bufio"
    "errors"
    "socks"
    "socks/log"
    "socks/authentication"
    "socks/context"
)

type Handshake interface {
    Handshake() (error)
}

type HandshakeV5 struct {
    version			byte
    nmethods			byte
    methods			[]byte
    authenticator	authentication.Authenticator
    context			*context.Context
}

type HandshakeV4 struct {
    version			byte
    nmethods			byte
    methods			[]byte
    authenticator	authentication.Authenticator
    context			*context.Context
}

func New(context *context.Context) (Handshake) {
    
    var handshake Handshake
    switch context.Version() {
        case socks.SOCKS_VERSION_V4:
            handshake = &HandshakeV4{context : context}
            break
        case socks.SOCKS_VERSION_V5:
            handshake = &HandshakeV5{context : context}
            break
    }
    
    return handshake
}

/*----------------------------------------------------------
    HandshakeV5 Implementation
-----------------------------------------------------------*/
func (handshake *HandshakeV5)Handshake() (error) {
    
    var err		error
    
    err = handshake.methodNegotiation()
    if (err != nil) {
        log.Errorf(" Method Negotiation failed, error: %s\n", err.Error())
        return err
    }
        
    _, err = handshake.authentication()
    if (err != nil) {
        // Send authentication successful response
        log.Errorf(" Authentication Negotiation failed, error: %s\n", err.Error())
        return err
    }
        
    return nil    
}

func (handshake *HandshakeV5)methodNegotiation() (error) {    
    
    var err		error
    var reader	*bufio.Reader = handshake.context.Reader()
       
    // Get the version, 1st byte
    handshake.version, err = reader.ReadByte()

    // Error?    
    if err != nil {
        // Send method negotiation error response
        response(socks.SOCKS_AUTH_NOACCEPTABLE, handshake.context)        
        return err
    }

    log.Infof("version : %d\n", handshake.version)
 
    // how many methods client support?
    handshake.nmethods, err = reader.ReadByte()
    
    // can't read the nmethods
    if err != nil {
        
        // Send method negotiation error response
        response(socks.SOCKS_AUTH_NOACCEPTABLE, handshake.context)
        
        // Send back the error code.
        return err
    }
    
    log.Infof("methods: %d\n", handshake.nmethods)

    // Get the methods clients supports    
    handshake.methods = make([]byte, handshake.nmethods)
    
    var index byte
    for index = 0; index < handshake.nmethods; index++ {
        handshake.methods[index], err = reader.ReadByte()

        log.Infof("method[%d]: %d\n", index, handshake.methods[index])        
        if (err != nil) {
            break
        }
    }
    
    // Return false when no methods, and error happened
    // If there are some methods, but error occured, continue.
    if (len(handshake.methods) == 0) {
        
        // Send method negotiation successful response
        response(socks.SOCKS_AUTH_NOACCEPTABLE, handshake.context)
        // Send back the error code
        return errors.New("No acceptable method")
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

    // Send back the error code
    response(found, handshake.context)
    
    // return the 
    return err
}

func (handshake *HandshakeV5)authentication () (bool, error) {
    return handshake.authenticator.Authenticate(handshake.context) 
}

/*----------------------------------------------------------
    HandshakeV4 Implementation
-----------------------------------------------------------*/
func (handshake *HandshakeV4)Handshake() (error) {
    return nil
}

func (handshake *HandshakeV4)authentication () (bool, error) {
    return handshake.authenticator.Authenticate(handshake.context) 
}

/*----------------------------------------------------------
    Public Methods Implementation
-----------------------------------------------------------*/

func response(code byte, context *context.Context) {
    context.Writer().WriteByte(context.Version())
    context.Writer().WriteByte(code)
    context.Writer().Flush()
}