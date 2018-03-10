package context

import (
    "bufio"
    "errors"
    "net"
    "socks"
    "socks/log"
    "socks/config"
)

type Context struct {
    version		byte
    connection 	net.Conn
    reader		*bufio.Reader
    writer		*bufio.Writer
    config		*config.Config
}

func New(conn net.Conn, config *config.Config) (*Context, error) {
    
    var reader *bufio.Reader = bufio.NewReader(conn)
    var writer *bufio.Writer = bufio.NewWriter(conn)
    
    // Read the version
    var version 	byte
    var err		error
    version, err = reader.ReadByte()
   
    // if version is not presented or error, then return nil, and err
    if (err != nil) {
        return nil, err
    }
    
    log.Infof("version: %d\n", version)
    
    // Check version is supported?
    if ((version != socks.SOCKS_VERSION_V4) && (version != socks.SOCKS_VERSION_V5)) {
        return nil, errors.New("Version is not supported")
    }
    
    // Rewind the position of version
    reader.UnreadByte();
    
    // return the Context object
    return &Context {   version		: version,
                        connection 	: conn,
                        reader		: reader,
                        writer		: writer,
                        config		: config }, nil
}

func (context *Context)Connection() (*net.Conn) {
    return &context.connection
}

func (contxt *Context)Version() byte {
    return contxt.version
}

func (context *Context) Reader() (*bufio.Reader) {
    return context.reader
}

func (context *Context) Writer() (*bufio.Writer) {
    return context.writer
}

func (context *Context) LocalAddr() (string) {
    return context.connection.LocalAddr().String()
}

func (context *Context) Config() (*config.Config) {
    return context.config
}