package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/boringproxy/boringproxy/httpconnect"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
)

func HTTPCONNECT(network string, url *url.URL, transport *http.Transport) (proxy.Dialer, error) {
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, errors.New("Unsupported scheme: " + url.Scheme)
	}
	d := httpconnect.NewDialer(network, url, transport)
	return d, nil
}

func main() {
	body, err := ioutil.ReadFile("/Users/fipina/.ssh/id_ed25519")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	signer, err := ssh.ParsePrivateKey(body)
	if err != nil {
		log.Fatalf("unable to parse key: %v", err)
	}
	// ssh config
	config := &ssh.ClientConfig{
		User:            "opc",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	//dialer, err := proxy.SOCKS5("tcp", "www-proxy-lon.uk.oracle.com:80", nil, proxy.Direct)
	//dialer, err := proxy.SOCKS5("tcp", "localhost:9898", nil, proxy.Direct)
	//dialer := proxy.FromEnvironment()
	//fmt.Println(dialer)
	tr := http.DefaultTransport.(*http.Transport)
	//tr.Dial = proxy.Direct.Dial
	// if forward != nil {
	// 	tr.Dial = forward.Dial
	// }
	proxyURL, err := url.Parse("http://www-proxy-lon.uk.oracle.com:80")
	if err != nil {
		log.Fatalf("unable to parse key: %v", err)
	}
	dialer, err := HTTPCONNECT("tcp", proxyURL, tr)
	if err != nil {
		log.Fatalf("unable to dial http: %v", err)
	}
	pconn, err := dialer.Dial("tcp", "b.oracole.ga:22")
	if err != nil {
		log.Fatalf("unable to connect ssh: %v", err)
	}

	conn, chans, reqs, err := ssh.NewClientConn(pconn, "b.oracole.ga:22", config)
	if err != nil {
		log.Fatalf("unable to parse key: %v", err)
	}
	client := ssh.NewClient(conn, chans, reqs)
	session, err := client.NewSession()

	// connect ot ssh server
	/*
		conn, err := ssh.Dial("tcp", "b.oracole.ga:22", config)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
	*/

	// Create ssh-session
	// session, err := conn.NewSession()
	if err != nil {
		log.Fatal("Failed to create SSH session", err)
	}
	defer session.Close()
	// Execute remote commands
	combo, err := session.CombinedOutput("whoami; cd /; ls -al")
	if err != nil {
		log.Fatal("remote execution CMD failed", err)
	}
	log.Println("Command Output:", string(combo))
}
