#### The Go Challenge 2

##### Can you help me secure my company’s data transmission?

Last week, our competitor released a feature that we were working on in secret. We suspect that they are spying on our network. Can you help us prevent our competitor from spying on our network?

#### To get started

Check out the two files **main.go** and **main_test.go** [here](http://web.archive.org/web/20200712092904/https://gist.github.com/creack/333f89f6aec5b789c1a0). These files are the starting point for this challenge.

![update.jpg](http://web.archive.org/web/20200712092904im_/http://golang-challenge.org/images/update.jpg)

**3rd April 2015**: Guillaume has made 2 changes to the test files:

-   The first change is minor and concerns Window user with dual stack. He has changed the `net.Listen("tcp", ":0")` to `net.Listen("tcp", "127.0.0.1:0")`. This shouldn’t change anything for most people.
-   The second is more problematic and would impact people using `io.Copy`. In TestReadWritePing, he has moved the `w.Close()` in the writing goroutine just like in the other tests.

#### Goal of the challenge

In order to prevent our competitor from spying on our network, we are going to write a small system that leverages [NaCl](http://web.archive.org/web/20200712092904/http://nacl.cr.yp.to/) to establish secure communication. **NaCl** is a crypto system that uses a public key for encryption and a private key for decryption.

Your goal is to implement the functions in **main.go** and make it pass the provided tests in **main_test.go**.

#### Steps involved

The first step is going to be able to generate the public and private keys. Next, we want to create an `io.Writer` and `io.Reader` that will allow us to automatically encrypt/decrypt our data.

##### Part 1

Implement the following helpers that will return our NACL Reader / Writer.

```
func NewSecureReader(r io.Reader, priv, pub *[32]byte) io.Reader
func NewSecureWriter(w io.Writer, priv, pub *[32]byte) io.Writer
```

##### Part 2

Now that we can encrypt/decrypt message locally, it would be interesting to do so over the network!

We are going to write a server that will exchange keys with the client in order to establish a secure communication.

In order to be able to encrypt/decrypt, we need to perform a key exchange upon connection.

For the sake of the exercise, performing the key exchange in plain text is acceptable. (In a production system, it would not be acceptable due to MITM risk!)

Unfortunately, everybody has already left for the day, so let’s write a secure echo server so we can test!

In order to test our echo server, we can do:

```
$> ./challenge2 -l 8080&
$> ./challenge2 8080 “hello world”
hello world
```

#### Requirements of the challenge

-   Use the latest version of Go i.e. version 1.4.2
-   Use only standard library and package(s) under `golang.org/x/crypto/nacl`.

#### Hints

-   We consider that our messages will always be smaller than 32KB
-   Most of the elements involved for encryption/decryption have fixed length
-   `encoding/binary` can be useful for the variable parts

#### Further exploration

If you have liked this challenge, you can keep programming **outside** the main challenge. You can:

-   create an actual user interface. It would be nice if both side had a prompt and could send more than one message at a time.
-   Handle multiple client and group chat
-   Add compression

Please keep in mind that cryptography is a complex subject matter, with many very subtle risks and mistakes to make. While it makes for a fun exercise, inventing your own crypto systems for production use is usually not a good idea and should be left to a handful of very talented people.
