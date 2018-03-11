//---------------------------------------------------------
// Author: Stanley Wang
// Copyright 2018. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//---------------------------------------------------------

package authentication

import (
        "errors"
        "socks"
        "socks/context"
)

/* RFC 1928
   The client connects to the server, and sends a version
   identifier/method selection message:

                   +----+----------+----------+
                   |VER | NMETHODS | METHODS  |
                   +----+----------+----------+
                   | 1  |    1     | 1 to 255 |
                   +----+----------+----------+

   The VER field is set to X'05' for this version of the protocol.  The
   NMETHODS field contains the number of method identifier octets that
   appear in the METHODS field.

   The server selects from one of the methods given in METHODS, and
   sends a METHOD selection message:

                         +----+--------+
                         |VER | METHOD |
                         +----+--------+
                         | 1  |   1    |
                         +----+--------+

   If the selected METHOD is X'FF', none of the methods listed by the
   client are acceptable, and the client MUST close the connection.
*/

type Identity struct {
    Username		string
    Password		string
}

type Authenticator interface {
    Authenticate (context *context.Context) (bool, error)
}

type NoAuthentication struct {
    
}

type UserPasswordAuthentication struct {
    username		string
    password		string
}

type GssAPIAuthentication struct {
    
}

var AUTHENTICATORS map[byte]Authenticator

/*----------------------------------------------------------
    init called when loading this module
-----------------------------------------------------------*/
func init () {
    
    AUTHENTICATORS = make(map[byte]Authenticator, 3)
    AUTHENTICATORS[socks.SOCKS_AUTH_NOAUTHENTICATION] 	= NewNoAuthentication()
    AUTHENTICATORS[socks.SOCKS_AUTH_GSSAPI] 				= NewGssAPIAuthentication()
    AUTHENTICATORS[socks.SOCKS_AUTH_USERPASSWORD] 		= NewUserPasswordAuthentication()
}

/*----------------------------------------------------------
    NoAuthentication Implementation
-----------------------------------------------------------*/

func NewNoAuthentication() (*NoAuthentication){
    return &NoAuthentication{}
}

func (auth *NoAuthentication) Authenticate(context *context.Context) (bool, error) {
    return true, nil
}

/*----------------------------------------------------------
    UserPasswordAuthentication Implementation
-----------------------------------------------------------*/
func NewUserPasswordAuthentication() (*UserPasswordAuthentication){
    return &UserPasswordAuthentication{}
}

func (auth *UserPasswordAuthentication) Authenticate(context *context.Context) (bool, error) {

    var statuscode byte = socks.SOCKS_AUTH_NOACCEPTABLE
    
    defer response(statuscode, context)
    
    // check of the server config correctly    
    if ((len(context.Config().Auth.Username) == 0) || (len(context.Config().Auth.Password) == 0)) {
        return false, errors.New("Socks server doesn't config authentication correctly")
    }
    
    // Read the version
    _, err := context.Reader().ReadByte()
    if (err != nil) {
        return false, err
    }
    
    length, err := context.Reader().ReadByte()
    if (err != nil) {
        return false, err
    }
    
    bytes := make([]byte, length)
    count, err := context.Reader().Read(bytes)
    if (err != nil) {
        return false, err
    }
    
    username := string(bytes[:count])
    
    // Read password
    length, err = context.Reader().ReadByte()
    if (err != nil) {
        return false, err
    }
    count, err = context.Reader().Read(bytes)
    if (err != nil){
        return false, err        
    }
    password := string(bytes[:count])
    
    // check if the username and password provided
    if ((len(username) == 0) || (len(password) == 0)) {
        return false, errors.New("Authentication failed")
    }
    
    if ((username != context.Config().Auth.Username) || (password != context.Config().Auth.Password)) {
        return false, errors.New("Authentication failed")
    }

    statuscode = socks.SOCKS_V5_STATUS_SUCCESS

    return true, nil
}

/*----------------------------------------------------------
    UserPasswordAuthentication Implementation
-----------------------------------------------------------*/
func NewGssAPIAuthentication() (*GssAPIAuthentication){
    return &GssAPIAuthentication{}
}

func (auth *GssAPIAuthentication) Authenticate(context *context.Context) (bool, error) {
    return true, nil
}

/*----------------------------------------------------------
    public methods Implementation
-----------------------------------------------------------*/

func response(code byte, context *context.Context) {
    context.Writer().WriteByte(context.Version())
    context.Writer().WriteByte(code)
    context.Writer().Flush()
}