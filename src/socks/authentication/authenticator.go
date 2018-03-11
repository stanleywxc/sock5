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

type Authenticator interface {
    Authenticate () (bool, error)
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

func (auth *NoAuthentication) Authenticate() (bool, error) {
    return true, nil
}

/*----------------------------------------------------------
    UserPasswordAuthentication Implementation
-----------------------------------------------------------*/
func NewUserPasswordAuthentication() (*UserPasswordAuthentication){
    return &UserPasswordAuthentication{}
}

func (auth *UserPasswordAuthentication) Authenticate() (bool, error) {
    
    if (len(auth.username) == 0) {
        return false, errors.New("Authentication failed")
    }
    
    if (len(auth.username) == 0) {
        return false, errors.New("Authentication failed")
    }
    
    return true, nil
}

func (auth *UserPasswordAuthentication) SetUsername(username string) {
    auth.username = username
}

func (auth *UserPasswordAuthentication) SetPassword(password string) {
    auth.password = password
}

/*----------------------------------------------------------
    UserPasswordAuthentication Implementation
-----------------------------------------------------------*/
func NewGssAPIAuthentication() (*GssAPIAuthentication){
    return &GssAPIAuthentication{}
}

func (auth *GssAPIAuthentication) Authenticate() (bool, error) {
    return true, nil
}
