package authorities

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	// TODO: make sure these don't expire
	TLS_SELFSIG_KEY  = "-----BEGIN RSA PRIVATE KEY-----\nMIICXAIBAAKBgQC3j47UVL/c0sqCtvEyNUra56BTK+saal/mp6VnjZWbNM1ZRfH1\nVi8yYWqictyhUQDMbEvO3bDBWdwCBrcuZoXtvHaQBLgMPt4h3+WjfmUSN+H+o4Rx\nPizes25OJSDXWoFX6SfZDG+8w+AsA7I9DRpSxx24uK4Yhb6Lb/emjOaKKwIDAQAB\nAoGBAK5LpQ6eznOKv/nwZgQLsGE9cPGokRvLB/bHFvrF6NmwAJCVJtOTG0uWTp+j\nuzV96ekxp6XswQpHHe3anJN1jBJ+igEUuVoF61tV2M/sZwXuGolyAgFDJSosqPgy\n0oijbuJ+4tQPCHrg0jKsKIkTdom5oQ2OeTaB/UmRJ/cWNW7pAkEA3dcquLexnGes\ncmdo9DCRGHLVKB72XU0UuL8QInZh1T4yQjLzbUUtKPwaM/Nk0XMgF9epy2RmeL7b\nyk5OAAh5twJBANPTbMoPV3+uyuRP/GhedtNJAQIB7qdNqgYPB4cRE6xN+It9xuNI\n+qIuKGMEEkRVvCkC4ff+NwNm66uYeE3kAy0CQFrekTxa2mEDwoqWO1KTNkv6db/a\ndvYe5dcLRHOpZEeyE2o0bqwawvXf1mfjUi/NZZ7+kymiNatOGr/StXakAh8CQH1k\nO5MFL+uigfJTMvxpZve90H3qvOaGv+4kOXWH81hdM5MHWpOy4MLehgPPJi0Tf3Xb\ngf52mwRFiZ6jfBvHrOUCQGeTzEin78yQGqaxjNvFPFfxnoZrx4sPa/dk3Qquj8eA\nWFnr32OyZNtsjUviUg2/n1fK2bUfzi0LO50X9bVlOPQ=\n-----END RSA PRIVATE KEY-----"
	TLS_SELFSIG_CERT = "-----BEGIN CERTIFICATE-----\nMIIB8TCCAVqgAwIBAgIJAIyPpV4fZ+/QMA0GCSqGSIb3DQEBCwUAMBsxGTAXBgNV\nBAMMEHRlc3Qtc2VydmVyLWNlcnQwHhcNMTcwODA2MDU1NDQxWhcNMTcwOTA1MDU1\nNDQxWjAbMRkwFwYDVQQDDBB0ZXN0LXNlcnZlci1jZXJ0MIGfMA0GCSqGSIb3DQEB\nAQUAA4GNADCBiQKBgQC3j47UVL/c0sqCtvEyNUra56BTK+saal/mp6VnjZWbNM1Z\nRfH1Vi8yYWqictyhUQDMbEvO3bDBWdwCBrcuZoXtvHaQBLgMPt4h3+WjfmUSN+H+\no4RxPizes25OJSDXWoFX6SfZDG+8w+AsA7I9DRpSxx24uK4Yhb6Lb/emjOaKKwID\nAQABoz0wOzAJBgNVHRMEAjAAMB0GA1UdJQQWMBQGCCsGAQUFBwMCBggrBgEFBQcD\nATAPBgNVHREECDAGhwR/AAABMA0GCSqGSIb3DQEBCwUAA4GBAI8M+jT5T+ajBsYv\nUHuTB+wkL32qdTU9GB0eobGU2EU+zh8fvnYGyrEsgNP3dSnIx77bpqjZD8Jo7iWP\nZvqtgOKRPaPXgvmgWFZgTLZjFk/YQ05jihEH3wXhnezP0Pb1kkCSLGhXWuB7V00U\nbSt/cILW+uJSxAZLY0jjCJbzy7HT\n-----END CERTIFICATE-----"

	TLS_CLIENT_KEY        = "-----BEGIN RSA PRIVATE KEY-----\nMIICWwIBAAKBgQC0q6RPYtP88n+LWlup97hWb2I3bIwWiIqPR6bsUU6sB5T/mier\nx9kReFSu4346GMyv4rVzarLueipvMPcXP++LZ+sH0NQUwD2uSPe15EgRcoEEPNvV\nxsNJMvQfEfjv+1RHHPHMYJV9MJzXFRj52oyx3xK+jDG4Sm1ThI70fwJHYwIDAQAB\nAoGAClAl7/YnPbAmAbFlvB0M47o19A35LSwcJLOlXqYBhKZmJfUJwK+Gv42L3/PS\nd8SEoqGhU/ZKQnyswW4dHLGkncr+RAAQ5UGRUHr7wsP1c+9kZpkaj1hmyLQvEbL6\nLPFxvno6AGxbURznIBu1hCQUu0aS/QZYfpaYrjo/9N3dg7ECQQDe4HAsUMYah+3b\nGu2q2oTqFOdLU+LA7ZloX338uIXbXCwiZz43b40uNqoXYXRZQQB7qT+zwseqDXWZ\npmWjBTeZAkEAz4VtH9Ug511V7idjlOe0k1kois4ydfvurniUoBtDE6xKD6dR/EZ6\nf5yCVfM0GZAq+BgomYKEBTklo1EuUMYkWwJAF7M0GnJIbp/PukHlzgpIof+xDMCR\n10Qs0P1+jzYr/cSSaOIjqo9xKt3jPnM9hRQ1cfDwdjQbOUkPHVSlcC1o2QJAekup\nWZ8ievbYUzdHSlOaaVObvuFxf3Ju4McS35/xUcCxDLSQblmii13SuZBP3djGWdry\n4jS2VNWuxqZq4xNCDQJAEZ7djTVtLEghjof27CuXMkopZZ4RhYTsAbZAwMBNBhds\npQQS+O5pIVDD8ou3QfifB6G5OmZr0PaKld/99H52LA==\n-----END RSA PRIVATE KEY-----"
	TLS_CLIENT_CERT       = "-----BEGIN CERTIFICATE-----\nMIICHTCCAQUCCQDMWwmePGAAbTANBgkqhkiG9w0BAQsFADAPMQ0wCwYDVQQDDAR0\nZXN0MB4XDTE3MDgwNjA1NDIzMloXDTE3MDkwNTA1NDIzMlowFjEUMBIGA1UEAwwL\nY2xpZW50LXRlc3QwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBALSrpE9i0/zy\nf4taW6n3uFZvYjdsjBaIio9HpuxRTqwHlP+aJ6vH2RF4VK7jfjoYzK/itXNqsu56\nKm8w9xc/74tn6wfQ1BTAPa5I97XkSBFygQQ829XGw0ky9B8R+O/7VEcc8cxglX0w\nnNcVGPnajLHfEr6MMbhKbVOEjvR/AkdjAgMBAAEwDQYJKoZIhvcNAQELBQADggEB\nAL9ean4hhASkXiDMhDN0bTamflo6FlA8P9z+YD+OgyVVI+shEIkgMdB1NojIgzUB\nshPwhEQZrL+sDIKRYiA22tSBbgsbo4dkF5rpU+7RduKTl+q68SBue1TCtOCOimVX\nm/1VboYzfarRbLauHIVbiei4iX8OduIKi5i88yTwZgg7W69PcF8EUlg2brh7gMTa\n9Rn58H788kMokYcoAEZ+LVZdk8PNjYNkRZ3lUYvJbt9Ytdkone2WwyK8nChIV9Lg\n8yaSyFcqH4jkRVZrkdiMqqqbxKh2q0b+tLf0ixdRPeUxyDEikgAGaZ9XhV4+rqWa\nbgtUHwsSXRDGAkwMC3e2wtk=\n-----END CERTIFICATE-----"
	TLS_CLIENT_ENCODED_CN = "client-test"

	TLS_CLIENT_EXPIRED  = "-----BEGIN CERTIFICATE-----\nMIICHTCCAQUCCQDMWwmePGAAcjANBgkqhkiG9w0BAQsFADAPMQ0wCwYDVQQDDAR0\nZXN0MB4XDTE3MDczMDE2MzY0M1oXDTE3MDczMTE2MzY0M1owFjEUMBIGA1UEAwwL\nY2xpZW50LXRlc3QwgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBALSrpE9i0/zy\nf4taW6n3uFZvYjdsjBaIio9HpuxRTqwHlP+aJ6vH2RF4VK7jfjoYzK/itXNqsu56\nKm8w9xc/74tn6wfQ1BTAPa5I97XkSBFygQQ829XGw0ky9B8R+O/7VEcc8cxglX0w\nnNcVGPnajLHfEr6MMbhKbVOEjvR/AkdjAgMBAAEwDQYJKoZIhvcNAQELBQADggEB\nAH33moyxmEfUQk366sKzqMszCRGFMi4hoCFICih2FQI+rLhRLjqgwp+nPJaMxOuK\n3r4+hS3J4dRePdJNCyL26Vc9Sa6Qc825UwLMb5y6PJ7vCqXhlxsRMm69WntKpIyt\nGyTm22GSay8VE4aE24bhwP6SBFm/0hf9l60t7u9UtVYB7duoeNbwzntpG0F5HzMD\nbh5jlUDDoXrBJg/ywgwRsg6zrEh3S/Eijgc2RIXSvbefW3qvyV5H0bxR5ZADN7RB\nk6NFAwMGtFLqwDEIkGuooCkcSt8noptBwjygIALtihGWI2+mghlgpuNjXBG0bqwy\nVNf+DJUbQeTv82e/D1rAO1M=\n-----END CERTIFICATE-----"
	TLS_CLIENT_UNISSUED = "-----BEGIN CERTIFICATE-----\nMIICITCCAQkCCQDMWwmePGAAczANBgkqhkiG9w0BAQsFADAPMQ0wCwYDVQQDDAR0\nZXN0MCIYDzI1MDAwMTAxMDUwMDAwWhgPMjUwMDAxMDIwNTAwMDBaMBYxFDASBgNV\nBAMMC2NsaWVudC10ZXN0MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0q6RP\nYtP88n+LWlup97hWb2I3bIwWiIqPR6bsUU6sB5T/mierx9kReFSu4346GMyv4rVz\narLueipvMPcXP++LZ+sH0NQUwD2uSPe15EgRcoEEPNvVxsNJMvQfEfjv+1RHHPHM\nYJV9MJzXFRj52oyx3xK+jDG4Sm1ThI70fwJHYwIDAQABMA0GCSqGSIb3DQEBCwUA\nA4IBAQAWf1TLwaYLA7xR4Zqfr5Kei5mYYhlYq1eJ5y5j4+TI9XIdUPa5f3mrYZtu\nDqAmQmbqIXYM0YRSKgR/6BnMwEITUYeV1Ahmg631bWxdKdOvxxvtjiXDQnzyYXZl\naAL0URPQ1urnJziLF/SKh3j3HTCeuqpPYp0K7fB/A8EGcVDgy7sBQZzSxTyBiynf\np0S285QIgkoz3rLo3CUrlIyadmkSfmKdfsi5DHloJzbbarSRHbRb2xi1XVfV2ECU\n7HQRdZr9jBwDZGCFQHPIySZzfdUmi6vaMNaDKguBj25LLNKLI1xmJH+64Sc2q23H\n5OpXFz6VRggtqANEMzFvfs5NF+jY\n-----END CERTIFICATE-----"

	TLS_CLIENT_CSR = "-----BEGIN CERTIFICATE REQUEST-----\nMIIBVTCBvwIBADAWMRQwEgYDVQQDDAtjbGllbnQtdGVzdDCBnzANBgkqhkiG9w0B\nAQEFAAOBjQAwgYkCgYEAtKukT2LT/PJ/i1pbqfe4Vm9iN2yMFoiKj0em7FFOrAeU\n/5onq8fZEXhUruN+OhjMr+K1c2qy7noqbzD3Fz/vi2frB9DUFMA9rkj3teRIEXKB\nBDzb1cbDSTL0HxH47/tURxzxzGCVfTCc1xUY+dqMsd8SvowxuEptU4SO9H8CR2MC\nAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4GBALCOKX+QHmNLGrrSCWB8p2iMuS+aPOcW\nYI9c1VaaTSQ43HOjF1smvGIa1iicM2L5zTBOEG36kI+sKFDOF2cXclhQF1WfLcxC\nIi/JSV+W7hbS6zWvJOnmoi15hzvVa1MRk8HZH+TpiMxO5uqQdDiEkV1sJ50v0ZtR\nTMuSBjdmmJ1t\n-----END CERTIFICATE REQUEST-----"

	TLS_SERVER_CSR     = "-----BEGIN CERTIFICATE REQUEST-----\nMIIBVTCBvwIBADAWMRQwEgYDVQQDDAtzZXJ2ZXItdGVzdDCBnzANBgkqhkiG9w0B\nAQEFAAOBjQAwgYkCgYEAvIS94nPwFm2nmv1EdvnzZozLhz8KGeJeJjMTuMAz/Sjh\nTdpxDMwmwYfzLhElhnwdRpebW/vlfnmKEUN+NbKi44KPYPRKASTsjJ5czpZS0pSd\nsRJ4CdrVMcFMCZU1+8pIiDVXf0wV1g2Z1DnRTtFz2lCcmk3CROlQLDqjXh5rJZsC\nAwEAAaAAMA0GCSqGSIb3DQEBCwUAA4GBAGIi05JNB8IMiYAZyx+WtclzB7RryA3J\nMVfMF1LNsa8oX7r0o8bDR6ZlGIG2QIjXfg5XkGCyAhwe2LS0vi1BG4TVBDfPbZUG\nocqPeGRnrElpGkjP68uuzoCfKzbl9/fb8tmXF8gWC+4U01ySCMzzVqbvUAvdRhpv\naf4s0elNnfIF\n-----END CERTIFICATE REQUEST-----"
	TLS_SERVER_CSR_KEY = "-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQC8hL3ic/AWbaea/UR2+fNmjMuHPwoZ4l4mMxO4wDP9KOFN2nEM\nzCbBh/MuESWGfB1Gl5tb++V+eYoRQ341sqLjgo9g9EoBJOyMnlzOllLSlJ2xEngJ\n2tUxwUwJlTX7ykiINVd/TBXWDZnUOdFO0XPaUJyaTcJE6VAsOqNeHmslmwIDAQAB\nAoGBAJpyNg8c5Qm69KGp9Tl0NdFCcExxfRkPuAqbtgCalv1FTLC30f6mEupjMvv+\n2DuB24pGEuYdTTt833ydSv07SO6rBtOrky8RJLjyG35gX5cEMGbBOY4prJtkF7A9\nfCuVoU3grnfK+H/pkBA324qS0LdIyLyFolhC6Z2s0+/wup3hAkEA6kwF1xsSvWln\n0mHJTxuvsO6pdIwQnhib7EpsQtu5FFow/GoLguk0oNZ8yJo/c/YxNUSE8Ei5xceJ\nkmk1I8oQbQJBAM37KBEKQ0WkUZBScXQoeRg/ipgG8aZ+zoxPKJ243tLIRp7YZRMd\nNPQ60Uq/nwV6aIq0MZldjmTiNmnRFxYoGScCQEeocLDVaunbbxF9LuCvCxEGLXWj\n0QFJbYbbTDk2kmiTHSBWHqcvRBVdBFUmN/nzdiFgj3geBhNjb8QDwWjsyBECQQCz\nqjrMgjMjb6xlQkQLmbHmYKY27PeizkDDSdiXmkVhfad7riybe4JQ9WzQ0kuWr6q4\nMWyB5YTqohXsPqDwhpFJAkA+XDOC8WxKbllZBkcnoa8AqPwTmmGsoMWtKi0gO3iJ\nr8AbSq2TQei4ZE4eaOj/sNBGbElJX7NZi7MPnDRIy9Ic\n-----END RSA PRIVATE KEY-----"
)

func launchFakeTLSServer(t *testing.T, handler http.Handler, clientCA *TLSAuthority, clientCert tls.Certificate) (func(), func(string) []byte) {
	falseauthority, err := LoadTLSAuthority([]byte(TLS_SELFSIG_KEY), []byte(TLS_SELFSIG_CERT))
	if err != nil {
		t.Fatal(err)
	}
	return launchFakeTLSServerExt(t, handler, clientCA, falseauthority.(*TLSAuthority), clientCert)
}

func launchFakeTLSServerExt(t *testing.T, handler http.Handler, clientCA *TLSAuthority, serverCertAsCA *TLSAuthority, clientCert tls.Certificate) (func(), func(string) []byte) {
	cert := serverCertAsCA.ToHTTPSCert()

	pool := clientCA.ToCertPool()
	// include an invalid authority for kicks (and better testing that Verify also checks things)
	pool.AddCert(serverCertAsCA.cert)

	server := http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", 50900+rand.Intn(100)),
		Handler: handler,
		TLSConfig: &tls.Config{
			ClientAuth:   tls.VerifyClientCertIfGiven,
			ClientCAs:    pool,
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"http/1.1", "h2"},
		},
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		tlsListener := tls.NewListener(ln, server.TLSConfig)
		err := server.Serve(tlsListener)
		if err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	stop := func() {
		err := server.Shutdown(context.Background())
		if err != nil {
			t.Error(err)
		}
		err = server.Close()
		if err != nil {
			t.Error(err)
		}
	}

	certs := []tls.Certificate{}
	if clientCert.PrivateKey != nil {
		certs = []tls.Certificate{clientCert}
	}

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: certs,
			RootCAs:      serverCertAsCA.ToCertPool(),
		},
	}}

	request := func(url string) []byte {
		request, err := http.NewRequest("GET", "https://"+server.Addr+"/"+url, nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := client.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	return stop, request
}

func TestTLSAuthority_VerifyWorks(t *testing.T) {
	authority := getTLSAuthority(t) // uses cert from tls_parse_test.go
	clientCert, err := tls.X509KeyPair([]byte(TLS_CLIENT_CERT), []byte(TLS_CLIENT_KEY))
	if err != nil {
		t.Fatal(err)
	}
	stop, request := launchFakeTLSServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		principal, err := authority.Verify(request)
		if err != nil {
			t.Error(err)
		} else {
			writer.Write([]byte(principal))
		}
	}), authority, clientCert)
	defer stop()

	principal := request("")
	if string(principal) != TLS_CLIENT_ENCODED_CN {
		t.Error("Mismatch on encoded Common Name")
	}
}

func TestTLSAuthority_VerifyWrongCert(t *testing.T) {
	authority := getTLSAuthority(t) // uses cert from tls_parse_test.go
	clientCert, err := tls.X509KeyPair([]byte(TLS_SELFSIG_CERT), []byte(TLS_SELFSIG_KEY))
	if err != nil {
		t.Fatal(err)
	}
	stop, request := launchFakeTLSServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := authority.Verify(request)
		if err == nil {
			t.Error("Expected error from authority")
		} else {
			writer.Write([]byte(err.Error()))
		}
	}), authority, clientCert)
	defer stop()

	errtext := string(request(""))
	if !strings.Contains(errtext, "Certificate not valid under this authority") {
		t.Errorf("Expected failure of cert, not: %s", errtext)
	}
}

func TestTLSAuthority_VerifyNoCert(t *testing.T) {
	authority := getTLSAuthority(t) // uses cert from tls_parse_test.go
	stop, request := launchFakeTLSServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, err := authority.Verify(request)
		if err == nil {
			t.Error("Expected error from authority")
		} else {
			writer.Write([]byte(err.Error()))
		}
	}), authority, tls.Certificate{})
	defer stop()

	errtext := string(request(""))
	if !strings.Contains(errtext, "must be present") {
		t.Errorf("Expected failure of cert, not: %s", errtext)
	}
}

func TestExpiredTLSAuthority(t *testing.T) {
	authority := getTLSAuthority(t) // uses cert from tls_parse_test.go
	clientCert, err := tls.X509KeyPair([]byte(TLS_CLIENT_CERT), []byte(TLS_CLIENT_KEY))
	if err != nil {
		t.Fatal(err)
	}
	expiredCertBlock, err := loadSinglePEMBlock([]byte(TLS_CLIENT_EXPIRED), []string{"CERTIFICATE"})
	if err != nil {
		t.Fatal(err)
	}
	expiredCert, err := x509.ParseCertificate(expiredCertBlock)
	if err != nil {
		t.Fatal(err)
	}
	stop, request := launchFakeTLSServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// the HTTPS server will reject this itself, but we're checking if *we* handle it correctly
		// so we substitute this here
		request.TLS.VerifiedChains[0][0] = expiredCert
		_, err := authority.Verify(request)
		if err == nil {
			t.Error("Expected verification failure with expired cert.")
		} else {
			writer.Write([]byte(err.Error()))
		}
	}), authority, clientCert)
	defer stop()

	errstr := request("")
	if string(errstr) != "Certificate for /CN=client-test has expired" {
		t.Errorf("Mismatch on expected error; did not expect %s", errstr)
	}
}

func TestUnissuedTLSAuthority(t *testing.T) {
	authority := getTLSAuthority(t) // uses cert from tls_parse_test.go
	clientCert, err := tls.X509KeyPair([]byte(TLS_CLIENT_CERT), []byte(TLS_CLIENT_KEY))
	if err != nil {
		t.Fatal(err)
	}
	unissuedCertBlock, err := loadSinglePEMBlock([]byte(TLS_CLIENT_UNISSUED), []string{"CERTIFICATE"})
	if err != nil {
		t.Fatal(err)
	}
	unissuedCert, err := x509.ParseCertificate(unissuedCertBlock)
	if err != nil {
		t.Fatal(err)
	}
	stop, request := launchFakeTLSServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		// the HTTPS server will reject this itself, but we're checking if *we* handle it correctly
		// so we substitute this here
		request.TLS.VerifiedChains[0][0] = unissuedCert
		_, err := authority.Verify(request)
		if err == nil {
			t.Error("Expected verification failure with expired cert.")
		} else {
			writer.Write([]byte(err.Error()))
		}
	}), authority, clientCert)
	defer stop()

	errstr := request("")
	if string(errstr) != "Certificate for /CN=client-test is not yet valid" {
		t.Errorf("Mismatch on expected error; did not expect %s", errstr)
	}
}

func checkClientCertValidity(t *testing.T, clientCertData []byte, key []byte) (string, error) {
	authority := getTLSAuthority(t) // uses cert from tls_parse_test.go
	clientCert, err := tls.X509KeyPair(clientCertData, key)
	if err != nil {
		t.Fatal(err)
	}
	stop, request := launchFakeTLSServer(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		principal, err := authority.Verify(request)
		if err != nil {
			writer.Write([]byte("E" + principal))
		} else {
			writer.Write([]byte("P" + principal))
		}
	}), authority, clientCert)
	defer stop()

	principal := string(request(""))
	if principal[0] == 'E' {
		return "", errors.New(principal[1:])
	} else {
		return principal[1:], nil
	}
}

func TestTLSAuthority_Sign_UserCert(t *testing.T) {
	a := getTLSAuthority(t)
	cert, err := a.Sign(TLS_CLIENT_CSR, false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if err != nil {
		t.Fatal(err)
	}
	princ, err := checkClientCertValidity(t, []byte(cert), []byte(TLS_CLIENT_KEY))
	if err != nil {
		t.Errorf("Generated certificate did not validate: %s", err)
	} else if princ != "common-name-tc" {
		t.Error("Invalid principal found in generated certificate")
	}
}

func TestTLSAuthority_Sign_HostCert(t *testing.T) {
	a := getTLSAuthority(t)

	serverCA, err := LoadTLSAuthority([]byte(TLS_SELFSIG_KEY), []byte(TLS_SELFSIG_CERT))
	if err != nil {
		t.Fatalf("Could not create TLS authority: %s", err)
	}

	serverCert, err := serverCA.(*TLSAuthority).Sign(TLS_SERVER_CSR, true, time.Hour, "server-common-name", []string{"127.0.0.1"})
	if err != nil {
		t.Fatal(err)
	}

	serverAuthority, err := LoadTLSAuthority([]byte(TLS_SERVER_CSR_KEY), []byte(serverCert))
	if err != nil {
		t.Fatal(err)
	}

	clientCert, err := tls.X509KeyPair([]byte(TLS_CLIENT_CERT), []byte(TLS_CLIENT_KEY))
	if err != nil {
		t.Fatal(err)
	}
	stop, request := launchFakeTLSServerExt(t, http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		principal, err := a.Verify(request)
		if err != nil {
			t.Error(err)
		} else {
			writer.Write([]byte(principal))
		}
	}), a, serverAuthority.(*TLSAuthority), clientCert)
	defer stop()

	principal := string(request(""))
	if principal != TLS_CLIENT_ENCODED_CN {
		t.Error("Did not find expected principal")
	}
}

func signAndLoad(t *testing.T, csr string, ishost bool, duration time.Duration, commonname string, names []string) *x509.Certificate {
	a := getTLSAuthority(t)
	certpem, err := a.Sign(csr, ishost, duration, commonname, names)
	if err != nil {
		t.Fatal(err)
	}
	certdata, err := loadSinglePEMBlock([]byte(certpem), []string{"CERTIFICATE"})
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(certdata)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

func TestTLSAuthority_Sign_DistinctSerials(t *testing.T) {
	c1 := signAndLoad(t, TLS_CLIENT_CSR, false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	c2 := signAndLoad(t, TLS_CLIENT_CSR, false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if c1.SerialNumber.Cmp(c2.SerialNumber) == 0 {
		t.Error("Serial numbers are the same between certificates")
	}
}

func TestTLSAuthority_Sign_IssuedNow(t *testing.T) {
	c1 := signAndLoad(t, TLS_CLIENT_CSR, false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	delta := time.Now().Sub(c1.NotBefore)
	if math.Abs(delta.Seconds()) > 1 {
		t.Error("Certificate was not just issued")
	}
}

func TestTLSAuthority_Sign_CorrectDuration(t *testing.T) {
	for _, duration := range []time.Duration{time.Second, time.Second * 8, time.Minute, time.Minute * 16, time.Hour, time.Hour * 24, time.Hour * 10000} {
		c1 := signAndLoad(t, TLS_CLIENT_CSR, false, duration, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
		delta := c1.NotAfter.Sub(c1.NotBefore)
		if delta != duration {
			t.Errorf("Certificate did not have the correct duration: %v instead of %v", delta, duration)
		}
	}
}

func TestTLSAuthority_Sign_CheckKeyUsage(t *testing.T) {
	c1 := signAndLoad(t, TLS_CLIENT_CSR, false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if c1.KeyUsage != x509.KeyUsageDigitalSignature {
		t.Error("Incorrect client auth usage")
	}
	if len(c1.ExtKeyUsage) != 1 || c1.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth {
		t.Error("Incorrect client auth usage")
	}
	c2 := signAndLoad(t, TLS_CLIENT_CSR, true, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if c2.KeyUsage != x509.KeyUsageDigitalSignature {
		t.Error("Incorrect server auth usage")
	}
	if len(c2.ExtKeyUsage) != 2 || c2.ExtKeyUsage[0] != x509.ExtKeyUsageClientAuth || c2.ExtKeyUsage[1] != x509.ExtKeyUsageServerAuth {
		t.Errorf("Incorrect server auth usage: %v", c2.ExtKeyUsage)
	}
}

func TestTLSAuthority_Sign_CheckBasicConstraints(t *testing.T) {
	for _, ishost := range []bool{false, true} {
		c1 := signAndLoad(t, TLS_CLIENT_CSR, ishost, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
		if !c1.BasicConstraintsValid || c1.IsCA {
			t.Error("Did not properly specify that this certificate is not a CA.")
		}
	}
}

func TestTLSAuthority_Sign_CheckSubject(t *testing.T) {
	for _, ishost := range []bool{false, true} {
		for _, commonname := range []string{"", "testcommon", "hello world", "I HAVE SPACES AND CAPITAL LETTERS AND PUNCTUATION ..,..a,.u*&("} {
			c1 := signAndLoad(t, TLS_CLIENT_CSR, ishost, time.Hour, commonname, []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
			if c1.Subject.CommonName != commonname {
				t.Errorf("Mismatch in common name of generated certificate: expected %s but got %s", commonname, c1.Subject.CommonName)
			}
		}
	}
}

func TestTLSAuthority_Sign_CheckNames(t *testing.T) {
	for _, ishost := range []bool{false, true} {
		c1 := signAndLoad(t, TLS_CLIENT_CSR, ishost, time.Hour, "test-common-tcs", []string{"dns1.mit.edu", "18.181.123.45", "dns2.mit.edu", "18.181.123.78", "fe80::4:6"})
		if len(c1.DNSNames) != 2 {
			t.Errorf("Wrong number of DNS names: %v", c1.DNSNames)
		} else {
			if c1.DNSNames[0] != "dns1.mit.edu" {
				t.Error("Wrong dns name")
			}
			if c1.DNSNames[1] != "dns2.mit.edu" {
				t.Error("Wrong dns name")
			}
		}
		if len(c1.IPAddresses) != 3 {
			t.Errorf("Wrong number of IP addresses: %v", c1.IPAddresses)
		} else {
			if !c1.IPAddresses[0].Equal(net.IPv4(18, 181, 123, 45)) {
				t.Error("Wrong IP address")
			}
			if !c1.IPAddresses[1].Equal(net.IPv4(18, 181, 123, 78)) {
				t.Error("Wrong IP address")
			}
			if !c1.IPAddresses[2].Equal(net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 6}) {
				t.Error("Wrong IP address")
			}
		}
	}
}

func TestTLSAuthority_Sign_MalformedPEM(t *testing.T) {
	a := getTLSAuthority(t)
	_, err := a.Sign(strings.Replace("I'm literally not a PEM file", "Z", "Y", -1), false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if err == nil {
		t.Error("Expected error while signing malformed CSR")
	} else if !strings.Contains(err.Error(), "PEM header") {
		t.Errorf("Unexpected error message -- expected 'PEM header' error: %s", err)
	}
}

func TestTLSAuthority_Sign_MalformedCSRBody(t *testing.T) {
	a := getTLSAuthority(t)
	_, err := a.Sign(strings.Replace(TLS_CLIENT_CSR, "WMRQw", "Y", -1), false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if err == nil {
		t.Error("Expected error while signing malformed CSR")
	} else if !strings.Contains(err.Error(), "asn1") {
		t.Errorf("Unexpected error message -- expected asn1 error: %s", err)
	}
}

func TestTLSAuthority_Sign_MalformedCSR(t *testing.T) {
	a := getTLSAuthority(t)
	_, err := a.Sign(strings.Replace(TLS_CLIENT_CSR, "Z", "Y", -1), false, time.Hour, "common-name-tc", []string{"dns1.mit.edu", "18.181.123.456", "dns2.mit.edu", "18.181.123.789"})
	if err == nil {
		t.Error("Expected error while signing malformed CSR")
	} else if !strings.Contains(err.Error(), "verification error") {
		t.Errorf("Unexpected error message -- expected verification error: %s", err)
	}
}
