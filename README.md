# cfjwt

CFJWT is a JWT verifier for Cloudflare Access.

It's built to act as Middleware for your services.

## Usage

1. Gather your Policy AUD (from the Access Policy in the dashboard) and Auth Domain (like `https://iamtheyammer.cloudflareaccess.com`)
2. `import github.com/iamtheyammer/cfjwt`
3. On the top of your file, add the following:
```go
var jwtVerifier = cfjwt.Verifier{
    PolicyAUD: "yourpolicyaud",
    AuthDomain: "https://yourauthdomain.cloudflareaccess.com",
}.Init()
```

Now, you can use:
- `jwtVerifier.Verify(jwt)` - returns a bool indicating the validity of the JWT
- `jwtVerifier.Middleware(YourRouter())`, like
```go
http.Handle("/pathwithaccess", jwtVerifier.Middleware(YourRouter()))
```
- `ok := verifier.HandlerMiddleware(w, r)`
  - This is meant to be run inside of a route handler, allowing easy compatibility with other frameworks like httprouter.
  - `ok` represents whether the middleware has written a response, meaning that if ok is true, you should simply return from your handler
  - `w` is your `http.ResponseWriter`
  - `r` is your `*http.Request`
  
## Full Examples

### Using normal middleware

```go
package main

import (
    "context"
    "fmt"
    "github.com/iamtheyammer/cfjwt"
    "net/http"
)

var jwtVerifier = cfjwt.Verifier{
    PolicyAUD: "yourpolicyaud",
    AuthDomain: "https://yourauthdomain.cloudflareaccess.com",
}.Init()

func MainHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("welcome"))
    })
}

func main() {
    http.Handle("/", jwtVerifier.Middleware(MainHandler()))
    http.ListenAndServe(":3000", nil)
}
```

### Using router middleware

```go

package main

import (
    "context"
    "fmt"
    "github.com/iamtheyammer/cfjwt"
    "net/http"
)

var jwtVerifier = cfjwt.Verifier{
    PolicyAUD: "yourpolicyaud",
    AuthDomain: "https://yourauthdomain.cloudflareaccess.com",
}.Init()

func MainHandler(w http.RespnseWriter, r *http.Request) {
    if ok := jwtVerifier.HandlerMiddleware(w, r); ok {
        return
    }

    w.Write([]byte("Welcome"))
    return
}

func main() {
    http.Handle("/", MainHandler)
    http.ListenAndServe(":3000", nil)
}
```