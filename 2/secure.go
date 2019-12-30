package main

import (
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/nacl/box"
)

type secureReader struct {
	sourceReader          io.Reader
	privateKey, publicKey *[32]byte
}

func (r secureReader) Read(p []byte) (int, error) {
	// We consider that our messages will always be smaller than 32KB
	buffer := make([]byte, 32*1024)
	n, err := r.sourceReader.Read(buffer[:])
	if err != nil && err != io.EOF {
		return n, err
	}

	var nonce [24]byte
	// first 24 bytes of the encrypted data is nonce
	copy(nonce[:], buffer[:24])

	encrypted := buffer[24:n]

	decrypted, ok := box.Open(nil, encrypted, &nonce, r.publicKey, r.privateKey)
	if !ok {
		return 0, errors.New("Can not decrypt data")
	}
	copy(p, decrypted[:])

	return len(decrypted), nil
}

type secureWriter struct {
	sourceWriter          io.Writer
	privateKey, publicKey *[32]byte
}

func (w secureWriter) Write(p []byte) (n int, err error) {
	// We must use a different nonce for each message you encrypt with the same key.
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return n, err
	}

	encrypted := box.Seal(nonce[:], p, &nonce, w.publicKey, w.privateKey)

	return w.sourceWriter.Write(encrypted)
}
