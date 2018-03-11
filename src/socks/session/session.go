//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//---------------------------------------------------------

package session

import (
        "socks"
        "socks/log"
        "socks/context"
        "socks/handshake"
        "socks/request"
)

type Session interface {
        Start()		(error)
}

type SessionV4 struct {
        context *context.Context
}

type SessionV5 struct {
        context *context.Context
}


type SessionContext struct {
        context		*context.Context
        session		Session
}

/*----------------------------------------------------------
    SessionContext
-----------------------------------------------------------*/

func New(contxt *context.Context) (*SessionContext) {
    
    var session Session
    
    switch contxt.Version() {
        case socks.SOCKS_VERSION_V4:
            session = createSessionV4(contxt)
            break
        case socks.SOCKS_VERSION_V5:
            session = createSessionV5(contxt)
            break;
    }
    
    return &SessionContext{ context : contxt, session : session } 
}

func (sc *SessionContext) Start () (error) {
    return sc.session.Start()
}

/*----------------------------------------------------------
    Socks Version 4
-----------------------------------------------------------*/
func createSessionV4(contxt *context.Context) (*SessionV4) {
    
    return &SessionV4{context : contxt}
}

func (session *SessionV4) Start() (error){
    return nil
}

/*----------------------------------------------------------
    Socks Version 5
-----------------------------------------------------------*/
func createSessionV5(contxt *context.Context) (*SessionV5) {
    
    return &SessionV5{context : contxt}
}


func (session *SessionV5) Start() (error) {
    
    // Do the handshake first
    statuscode, err := handshake.New(session.context).Handshake()
    
    // Is there any error?
    if ( err != nil) {
        session.reponse(statuscode)
        log.Errorf("Handshake failed, error: %s\n", err.Error())
        return err
    }
    
    // Send the found method
    session.reponse(statuscode)
    
    // Accept the requests
    request := request.New(session.context)
    status, err := request.Start()
    if (status == false) {
        
        // send the error code back to client, really don't 
        // care of the rest of reply data since we are going
        // to close the connection any way.
        session.reponse(socks.SOCKS_V5_STATUS_SERVER_FAILURE)
        log.Errorf("Process Request failed, error: %s\n", err.Error())
        return err 
    }
    
    // Run the command
    (*request.Command()).Execute()
    
    // Done
    return nil
}

func (session *SessionV5) reponse(statuscode byte) {
    
    // Send response back.
    session.context.Writer().WriteByte(session.context.Version())
    session.context.Writer().WriteByte(statuscode)
    session.context.Writer().Flush()
}