package main

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

// Most of the code was copied from https://github.com/mikeflynn/go-alexa/blob/master/skillserver/skillserver.go which
// is under the MIT licence. I added only the cache-feature.

var certCache = cache.New(time.Minute, 10*time.Minute)

// IsValidAlexaRequest handles all the necessary steps to validate that an incoming http.Request has actually come from
// the Alexa service. If an error occurs during the validation process, an http.Error will be written to the provided http.ResponseWriter.
// The required steps for request validation can be found on this page:
// https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service#hosting-a-custom-skill-as-a-web-service
func IsValidAlexaRequest(w http.ResponseWriter, req *http.Request) bool {
	certURL := req.Header.Get("SignatureCertChainUrl")

	if !verifyCertURL(certURL) {
		HTTPError(w, "Invalid cert URL: "+certURL, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	cert := fetchAndVerifyCert(w, certURL)
	if cert == nil {
		return false
	}

	encryptedSignature, err := base64.StdEncoding.DecodeString(req.Header.Get("Signature"))
	if err != nil {
		HTTPError(w, "Could not decode signature: "+err.Error(), "Unauthorized", http.StatusUnauthorized)
		return false
	}

	var bodyBuf bytes.Buffer
	hash := sha1.New()
	_, err = io.Copy(hash, io.TeeReader(req.Body, &bodyBuf))
	if err != nil {
		HTTPError(w, err.Error(), "Internal Error", http.StatusInternalServerError)
		return false
	}

	req.Body = ioutil.NopCloser(&bodyBuf)

	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSignature)
	if err != nil {
		HTTPError(w, "Signature match failed.", "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
		return false
	}

	if !strings.HasPrefix(link.Path, "/echo.api/") {
		return false
	}

	return true
}

func fetchAndVerifyCert(w http.ResponseWriter, certURL string) *x509.Certificate {
	cacheValue, found := certCache.Get("cert")
	if found {
		log.Debug("Found cert in Cache")
		return cacheValue.(*x509.Certificate)
	}

	certContents, err := fetchCert(certURL)
	if err != nil {
		HTTPError(w, err.Error(), "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	certBlock, remaining := pem.Decode(certContents)
	if certBlock == nil {
		HTTPError(w, "Failed to parse certificate PEM.", "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		HTTPError(w, err.Error(), "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(remaining)

	if _, err = cert.Verify(x509.VerifyOptions{DNSName: "echo-api.amazon.com", Roots: roots}); err != nil {
		HTTPError(w, "Amazon certificate invalid: "+err.Error(), "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	certCache.Set("cert", cert, cache.DefaultExpiration)

	return cert
}

func fetchCert(certURL string) ([]byte, error) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, errors.New("Could not download Amazon cert file.")
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, errors.New("Could not read Amazon cert file.")
	}
	return certContents, nil
}

func HTTPError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		log.Info(logMsg)
	}

	http.Error(w, err, errCode)
}
