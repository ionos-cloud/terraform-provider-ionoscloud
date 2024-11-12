// Copyright 2020 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Package ftps implements a simple FTPS client.
package ftps

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"unicode"
)

// Client FTPS.
type Client struct {
	plain  net.Conn
	secure *tls.Conn

	tc *textproto.Conn

	opt DialOptions
}

// DialOptions for the FTPS client.
type DialOptions struct {
	Host     string
	Port     int // If zero, this will default to 990.
	Username string
	Passowrd string

	// If true, will connect un-encrypted, then upgrade to using AUTH TLS command.
	ExplicitTLS bool

	// If true, will NOT attempt to encrypt.
	InsecureUnencrypted bool

	TLSConfig *tls.Config
}

func joinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.FormatInt(int64(port), 10))
}

// Dial a FTPS server and return a Client.
func Dial(ctx context.Context, opt DialOptions) (*Client, error) {
	port := opt.Port
	if port <= 0 {
		if opt.InsecureUnencrypted {
			port = 21
		} else {
			port = 990
		}
	}
	dialer := &net.Dialer{}
	dialTo := joinHostPort(opt.Host, port)
	conn, err := dialer.DialContext(ctx, "tcp", dialTo)
	if err != nil {
		return nil, fmt.Errorf("ftps: network dial failed: %w", err)
	}

	client := &Client{
		plain: conn,
		opt:   opt,
	}

	if err = client.setup(); err != nil {
		client.plain.Close()
		return nil, fmt.Errorf("ftps: connection setup failed: %w", err)
	}
	return client, nil
}

func (c *Client) setup() error {
	if c.opt.ExplicitTLS {
		c.tc = textproto.NewConn(c.plain)
		if _, err := c.read(220); err != nil {
			return fmt.Errorf("setup init read: %w", err)
		}
		if _, err := c.cmd(234, "AUTH TLS"); err != nil {
			return err
		}
	}

	if !c.opt.InsecureUnencrypted {
		c.secure = tls.Client(c.plain, c.opt.TLSConfig)
		if err := c.secure.Handshake(); err != nil {
			return err
		}
		c.tc = textproto.NewConn(c.secure)
	} else {
		c.tc = textproto.NewConn(c.plain)
	}

	if !c.opt.ExplicitTLS {
		if _, err := c.read(220); err != nil {
			return fmt.Errorf("setup init read: %w", err)
		}
	}

	if _, err := c.cmd(331, "USER %s", c.opt.Username); err != nil {
		return err
	}
	if _, err := c.cmd(230, "PASS %s", c.opt.Passowrd); err != nil {
		return err
	}
	if _, err := c.cmd(200, "TYPE I"); err != nil {
		return err
	}
	if _, err := c.cmd(200, "PBSZ %d", 0); err != nil {
		return err
	}
	if c.opt.InsecureUnencrypted {
		return nil
	}

	if _, err := c.cmd(200, "PROT %s", "P"); err != nil {
		return err
	}
	return nil
}

func (c *Client) read(expectCode int) (string, error) {
	gotCode, message, err := c.tc.ReadResponse(expectCode)
	if err != nil {
		return "", fmt.Errorf("failed to read code, got code %d and message %s: %w", gotCode, message, err)
	}
	return message, nil
}

func (c *Client) cmd(expectedCode int, cmd string, args ...interface{}) (string, error) {
	id, err := c.tc.Cmd(cmd, args...)
	if err != nil {
		return "", fmt.Errorf("cmd %q failed with ID %d: %w", cmd, id, err)
	}

	message, err := c.read(expectedCode)
	if err != nil {
		return "", fmt.Errorf("cmd %q failed read expected code %d with message %q: %w", cmd, expectedCode, message, err)
	}

	return message, nil
}

func (c *Client) data(ctx context.Context, expectedCode int, cmd string, args ...interface{}) (io.ReadWriteCloser, error) {
	message, err := c.cmd(227, "PASV")
	if err != nil {
		return nil, err
	}

	// Expected Message: Entering Passive Mode (x,x,x,x,p1,p2)
	start := strings.Index(message, "(")
	end := strings.LastIndex(message, ")")
	if start < 0 || end < 0 || end < start {
		return nil, fmt.Errorf("invalid PASV response, got %q", message)
	}
	portPartList := strings.Split(message[start+1:end], ",")
	if len(portPartList) < 6 {
		return nil, fmt.Errorf("invalid PASV port response, got %q", portPartList)
	}
	p1, err := strconv.ParseInt(portPartList[4], 10, 16)
	if err != nil {
		return nil, err
	}
	p2, err := strconv.ParseInt(portPartList[5], 10, 16)
	if err != nil {
		return nil, err
	}
	port := int(p1)*256 + int(p2)
	// Ignore the IP address.

	dialer := &net.Dialer{}
	dconn, err := dialer.DialContext(ctx, "tcp", joinHostPort(c.opt.Host, port))
	if err != nil {
		return nil, fmt.Errorf("dial data conn failed: %w", err)
	}

	_, err = c.cmd(expectedCode, cmd, args...)
	if err != nil {
		dconn.Close()
		return nil, err
	}

	if c.opt.InsecureUnencrypted {
		return dconn, nil
	}
	secure := tls.Client(dconn, c.opt.TLSConfig)

	return secure, nil
}

// Close the FTPS client connection.
func (c *Client) Close() error {
	_, qerr := c.cmd(221, "QUIT")
	if c.secure != nil {
		serr := c.secure.Close()
		if serr != nil {
			return serr
		}
		return qerr
	}
	c.plain.Close()
	return qerr
}

// Getwd gets the current working directory.
func (c *Client) Getwd() (dir string, err error) {
	return c.cmd(257, "PWD")
}

// Chdir changes the current working directory.
func (c *Client) Chdir(dir string) error {
	if _, err := c.cmd(250, "CWD %s", dir); err != nil {
		return fmt.Errorf("ftps: Chdir failed: %w", err)
	}
	return nil
}

// Mkdir makes a new directory.
func (c *Client) Mkdir(name string) error {
	if _, err := c.cmd(257, "MKD %s", name); err != nil {
		return fmt.Errorf("ftps: Mkdir failed: %w", err)
	}
	return nil
}

// RemoveFile removes a file.
func (c *Client) RemoveFile(name string) error {
	if _, err := c.cmd(250, "DELE %s", name); err != nil {
		return fmt.Errorf("ftps: RemoveFile failed: %w", err)
	}
	return nil
}

// RemoveDir removes a directory.
func (c *Client) RemoveDir(name string) error {
	if _, err := c.cmd(250, "RMD %s", name); err != nil {
		return fmt.Errorf("ftps: RemoveDir failed: %w", err)
	}
	return nil
}

// File of a directory list.
type File struct {
	Name string
}

// List the contents of the current working directory.
func (c *Client) List(ctx context.Context) ([]File, error) {
	data, err := c.data(ctx, 1, "LIST") // 150
	if err != nil {
		return nil, fmt.Errorf("ftps: failed to List, unable to get data conn: %w", err)
	}
	defer data.Close()

	list := make([]File, 0, 3)

	reader := bufio.NewReader(data)
	for {
		select {
		default:
		case <-ctx.Done():
			return list, fmt.Errorf("ftps: List canceled: %w", ctx.Err())
		}
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		f, err := readLine(line)
		if err != nil {
			return list, fmt.Errorf("ftps: List line parse: %w", err)
		}
		list = append(list, f)
	}
	data.Close()

	_, err = c.read(2) // 226
	if err != nil {
		return list, fmt.Errorf("ftps: List ack failed: %w", err)
	}

	return list, nil
}

func readLine(line string) (File, error) {
	f := File{}
	line = strings.TrimSpace(line)

	filenameIndex := 0
	ct := 0
	inSP := true
	for index, r := range line {
		sp := unicode.IsSpace(r)
		if inSP == sp {
			continue
		}
		inSP = sp
		if sp {
			continue
		}
		ct++

		if ct == 9 {
			filenameIndex = index
			break
		}
	}
	f.Name = line[filenameIndex:]
	return f, nil
}

// Upload the contents of Reader to the file name to the current working directory.
func (c *Client) Upload(ctx context.Context, name string, r io.Reader) error {
	data, err := c.data(ctx, 1, "STOR %s", name) // 150
	if err != nil {
		return fmt.Errorf("upload data: %w", err)
	}
	defer data.Close()

	_, err = io.Copy(data, r)
	if err != nil {
		return fmt.Errorf("upload copy: %w", err)
	}

	if err = data.Close(); err != nil {
		return fmt.Errorf("upload close: %w", err)
	}
	_, err = c.read(2) // 226
	if err != nil {
		return fmt.Errorf("upload read: %w", err)
	}
	return nil
}

// Download the file name from the current working directory to the Writer.
func (c *Client) Download(ctx context.Context, name string, w io.Writer) error {
	data, err := c.data(ctx, 1, "RETR %s", name) // 150
	if err != nil {
		return fmt.Errorf("download data: %w", err)
	}
	defer data.Close()

	_, err = io.Copy(w, data)
	if err != nil {
		return fmt.Errorf("download copy: %w", err)
	}
	data.Close()

	_, err = c.read(2) // 226
	if err != nil {
		return fmt.Errorf("download read: %w", err)
	}
	return nil
}
