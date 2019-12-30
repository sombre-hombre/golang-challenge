package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/nacl/box"
)

// NewSecureReader instantiates a new SecureReader
func NewSecureReader(r io.Reader, priv, pub *[32]byte) io.Reader {
	return &secureReader{r, priv, pub}
}

// NewSecureWriter instantiates a new SecureWriter
func NewSecureWriter(w io.Writer, priv, pub *[32]byte) io.Writer {
	return &secureWriter{w, priv, pub}
}

// Dial generates a private/public key pair,
// connects to the server, perform the handshake
// and return a reader/writer.
func Dial(addr string) (io.ReadWriteCloser, error) {
	// Generating client keys
	clientPublicKey, clientPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// sending public key
	_, err = conn.Write(clientPublicKey[:])
	if err != nil {
		return nil, err
	}
	// reading server's key
	serverKey := [32]byte{}
	_, err = conn.Read(serverKey[:])
	if err != nil {
		return nil, err
	}

	reader := NewSecureReader(conn, clientPrivateKey, &serverKey)
	writer := NewSecureWriter(conn, clientPrivateKey, &serverKey)

	return &secureConnection{conn, reader, writer}, nil
}

type secureConnection struct {
	c net.Conn
	io.Reader
	io.Writer
}

func (c secureConnection) Close() error {
	return c.c.Close()
}

// Serve starts a secure echo server on the given listener.
func Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		// reading client's public key
		clientKey := [32]byte{}
		n, err := conn.Read(clientKey[:])
		if err != nil {
			return errors.New("Client public key expected")
		}

		// generating server's keys
		pubKey, privKey, err := box.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}
		// sending server's public key
		n, err = conn.Write(pubKey[:])
		if err != nil || n != 32 {
			return errors.New("Can not send public key")
		}

		sr := NewSecureReader(conn, privKey, &clientKey)
		sw := NewSecureWriter(conn, privKey, &clientKey)
		// Echo all incoming data.
		io.Copy(sw, sr)
		// Shut down the connection.
		conn.Close()

		// // Handle the connection in a new goroutine.
		// // The loop then returns to accepting, so that
		// // multiple connections may be served concurrently.
		// go func(c net.Conn) {
		// 	sr := NewSecureReader(c, privKey, &clientKey)
		// 	sw := NewSecureWriter(c, privKey, pubKey)
		// 	// Echo all incoming data.
		// 	io.Copy(sw, sr)
		// 	// Shut down the connection.
		// 	c.Close()
		// }(conn)
	}
}

func main() {
	port := flag.Int("l", 0, "Listen mode. Specify port")
	flag.Parse()

	// Server mode
	if *port != 0 {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Fatal(err)
		}
		defer l.Close()
		log.Fatal(Serve(l))
	}

	// Client mode
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <port> <message>", os.Args[0])
	}
	conn, err := Dial("localhost:" + os.Args[1])
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Write([]byte(os.Args[2])); err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, len(os.Args[2]))
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", buf[:n])
}
