package cfjwt

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc"
)

type Verifier struct {
	PolicyAUD  string
	AuthDomain string
	ctx        context.Context
}

var (
	ctx        = context.TODO()
	authDomain = ""

	config   = &oidc.Config{}
	keySet   = oidc.NewRemoteKeySet(ctx, "")
	verifier = oidc.NewVerifier(authDomain, keySet, config)
)

func (v Verifier) getCertsURL() string {
	return fmt.Sprintf("%s/cdn-cgi/access/certs", v.AuthDomain)
}

// Init populates package global variables.
func (v Verifier) Init() Verifier {
	authDomain = v.AuthDomain
	config.ClientID = v.PolicyAUD
	keySet = oidc.NewRemoteKeySet(ctx, v.getCertsURL())
	verifier = oidc.NewVerifier(authDomain, keySet, config)

	return v
}

// Verify returns true if the jwt is valid and false if the jwt is invalid.
func (v Verifier) Verify(jwt string) bool {
	// Verify the access token
	_, err := verifier.Verify(ctx, jwt)
	if err != nil {
		fmt.Print(err.Error())
		return false
	}

	return true
}

// Middleware is a middleware to verify a CF Access token
func (v Verifier) Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Make sure that the incoming request has our token header
		// Could also look in the cookies for CF_AUTHORIZATION
		accessJWT := r.Header.Get("Cf-Access-Jwt-Assertion")

		if len(accessJWT) < 1 {
			c, err := r.Cookie("CF_Authorization")
			if err == nil {
				accessJWT = c.Value
			}
		}

		if len(accessJWT) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("No token on the request"))
			return
		}

		// Verify the access token
		ctx = r.Context()
		if ok := v.Verify(accessJWT); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("Invalid token")))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// HandlerMiddleware should be called from a request handler. It returns if true if it has sent a response.
func (v Verifier) HandlerMiddleware(w http.ResponseWriter, r *http.Request) bool {
	// Make sure that the incoming request has our token header
	// Could also look in the cookies for CF_AUTHORIZATION
	accessJWT := r.Header.Get("Cf-Access-Jwt-Assertion")

	if len(accessJWT) < 1 {
		c, err := r.Cookie("CF_Authorization")
		if err == nil {
			accessJWT = c.Value
		}
	}

	if len(accessJWT) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No token on the request"))
		return true
	}

	// Verify the access token
	ctx = r.Context()
	if ok := v.Verify(accessJWT); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(fmt.Sprintf("Invalid token")))
		return true
	}

	return false
}
