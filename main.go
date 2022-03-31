package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const IP = "198.18.0.169"

//https://pascal.bach.ch/2015/12/17/from-tcp-to-tls-in-go/
const serverKey = `-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIHg+g2unjA5BkDtXSN9ShN7kbPlbCcqcYdDu+QeV8XWuoAoGCCqGSM49
AwEHoUQDQgAEcZpodWh3SEs5Hh3rrEiu1LZOYSaNIWO34MgRxvqwz1FMpLxNlx0G
cSqrxhPubawptX5MSr02ft32kfOlYbaF5Q==
-----END EC PRIVATE KEY-----
`

const serverCert = `-----BEGIN CERTIFICATE-----
MIIB+TCCAZ+gAwIBAgIJAL05LKXo6PrrMAoGCCqGSM49BAMCMFkxCzAJBgNVBAYT
AkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRn
aXRzIFB0eSBMdGQxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xNTEyMDgxNDAxMTNa
Fw0yNTEyMDUxNDAxMTNaMFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0
YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMM
CWxvY2FsaG9zdDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHGaaHVod0hLOR4d
66xIrtS2TmEmjSFjt+DIEcb6sM9RTKS8TZcdBnEqq8YT7m2sKbV+TEq9Nn7d9pHz
pWG2heWjUDBOMB0GA1UdDgQWBBR0fqrecDJ44D/fiYJiOeBzfoqEijAfBgNVHSME
GDAWgBR0fqrecDJ44D/fiYJiOeBzfoqEijAMBgNVHRMEBTADAQH/MAoGCCqGSM49
BAMCA0gAMEUCIEKzVMF3JqjQjuM2rX7Rx8hancI5KJhwfeKu1xbyR7XaAiEA2UT7
1xOP035EcraRmWPe7tO0LpXgMxlh2VItpc2uc2w=
-----END CERTIFICATE-----
`

const rootCert = `-----BEGIN CERTIFICATE-----
MIIB+TCCAZ+gAwIBAgIJAL05LKXo6PrrMAoGCCqGSM49BAMCMFkxCzAJBgNVBAYT
AkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRn
aXRzIFB0eSBMdGQxEjAQBgNVBAMMCWxvY2FsaG9zdDAeFw0xNTEyMDgxNDAxMTNa
Fw0yNTEyMDUxNDAxMTNaMFkxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0
YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQxEjAQBgNVBAMM
CWxvY2FsaG9zdDBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHGaaHVod0hLOR4d
66xIrtS2TmEmjSFjt+DIEcb6sM9RTKS8TZcdBnEqq8YT7m2sKbV+TEq9Nn7d9pHz
pWG2heWjUDBOMB0GA1UdDgQWBBR0fqrecDJ44D/fiYJiOeBzfoqEijAfBgNVHSME
GDAWgBR0fqrecDJ44D/fiYJiOeBzfoqEijAMBgNVHRMEBTADAQH/MAoGCCqGSM49
BAMCA0gAMEUCIEKzVMF3JqjQjuM2rX7Rx8hancI5KJhwfeKu1xbyR7XaAiEA2UT7
1xOP035EcraRmWPe7tO0LpXgMxlh2VItpc2uc2w=
-----END CERTIFICATE-----
`

func handleConnTls(src *tls.Conn, addr string, stocLog, ctsLog bool) {

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootCert))
	if !ok {
		log.Fatal("failed to parse root certificate")
	}
	fmt.Println("tls- Dialing...")
	config := &tls.Config{ServerName: "www.nrk.no"}
	dst, err := tls.Dial("tcp", addr, config)
	if err != nil {
		panic(err)
	}
	fmt.Println("tls done")
	// server to client
	go func() {
		defer dst.Close()
		defer src.Close()
		if stocLog {
			tee := io.TeeReader(dst, src)
			io.Copy(os.Stdout, tee)

		} else {
			io.Copy(src, dst)
		}
	}()

	// client to server
	go func() {
		defer dst.Close()
		defer src.Close()

		if ctsLog {
			tee := io.TeeReader(src, dst)
			io.Copy(os.Stdout, tee)
		} else {
			io.Copy(dst, src)
		}

	}()

}

func handleConn(src net.Conn, addr string, stocLog, ctsLog bool) {
	dst, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	// server to client
	go func() {
		defer dst.Close()
		defer src.Close()
		if stocLog {
			tee := io.TeeReader(dst, src)
			io.Copy(os.Stdout, tee)

		} else {
			io.Copy(src, dst)
		}
	}()

	// client to server
	go func() {
		defer dst.Close()
		defer src.Close()

		if ctsLog {
			tee := io.TeeReader(src, dst)
			io.Copy(os.Stdout, tee)
		} else {
			io.Copy(dst, src)
		}

	}()

}

func main() {

	go func() {
		norm, err := net.Listen("tcp", "localhost:80")
		if err != nil {
			panic(err)
		}

		for {
			conn, err := norm.Accept()
			if err != nil {
				panic(err)
			}
			handleConn(conn, IP+":80", false, false)
		}
	}()

	go func() {
		cer, err := tls.X509KeyPair([]byte(serverCert), []byte(serverKey))
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{Certificates: []tls.Certificate{cer}}

		ssl, err := net.Listen("tcp", "localhost:443")
		if err != nil {
			panic(err)
		}
		for {
			conn, err := ssl.Accept()
			if err != nil {
				panic(err)
			}

			_ = config
			tconn := tls.Server(conn, config)
			handleConnTls(tconn, IP+":443", false, true)
			// handleConn(conn, IP+":443", false, true)
		}
	}()

	c := make(chan int)
	<-c

}
