package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/unrolled/secure"
)

func main() {
	fs := http.FileServer(http.Dir("assets/css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.Handle("/", middleware(index))
	fmt.Println("Starting server..")

	http.ListenAndServe(":8080", nil)
}

func middleware(next http.HandlerFunc) http.Handler {
	isProd := os.Getenv("RIDDLEBIT_PRODUCTION")
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:            []string{"localhost", "riddleb.it"},                                                         // AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
		AllowedHostsAreRegex:    false,                                                                                       // AllowedHostsAreRegex determines, if the provided AllowedHosts slice contains valid regular expressions. Default is false.
		HostsProxyHeaders:       []string{"X-Forwarded-Hosts"},                                                               // HostsProxyHeaders is a set of header keys that may hold a proxied hostname value for the request.
		SSLRedirect:             true,                                                                                        // If SSLRedirect is set to true, then only allow HTTPS requests. Default is false.
		SSLTemporaryRedirect:    false,                                                                                       // If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
		SSLHost:                 "riddleb.it",                                                                                // SSLHost is the host name that is used to redirect HTTP requests to HTTPS. Default is "", which indicates to use the same host.
		SSLHostFunc:             nil,                                                                                         // SSLHostFunc is a function pointer, the return value of the function is the host name that has same functionality as `SSHost`. Default is nil. If SSLHostFunc is nil, the `SSLHost` option will be used.
		SSLProxyHeaders:         map[string]string{"X-Forwarded-Proto": "https"},                                             // SSLProxyHeaders is set of header keys with associated values that would indicate a valid HTTPS request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
		STSSeconds:              31536000,                                                                                    // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    true,                                                                                        // If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
		STSPreload:              true,                                                                                        // If STSPreload is set to true, the `preload` flag will be appended to the Strict-Transport-Security header. Default is false.
		ForceSTSHeader:          true,                                                                                        // STS header is only included when the connection is HTTPS. If you want to force it to always be added, set to true. `IsDevelopment` still overrides this. Default is false.
		FrameDeny:               true,                                                                                        // If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
		CustomFrameOptionsValue: "SAMEORIGIN",                                                                                // CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option. Default is "".
		ContentTypeNosniff:      true,                                                                                        // If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
		BrowserXssFilter:        true,                                                                                        // If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
		ContentSecurityPolicy:   "default-src 'self'; script-src 'self'; connect-src 'self';style-src 'self';img-src https:", // ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "". Passing a template string will replace `$NONCE` with a dynamic nonce value of 16 bytes for each request which can be later retrieved using the Nonce function.
		ReferrerPolicy:          "same-origin",                                                                               // ReferrerPolicy allows the Referrer-Policy header with the value to be set with a custom value. Default is "".
		FeaturePolicy:           "vibrate 'none';",                                                                           // FeaturePolicy allows the Feature-Policy header with the value to be set with a custom value. Default is "".

		IsDevelopment: isProd != "1", // This will cause the AllowedHosts, SSLRedirect, and STSSeconds/STSIncludeSubdomains options to be ignored during development. When deploying to production, be sure to set this to false.
	})

	return secureMiddleware.Handler(next)
}

func index(rw http.ResponseWriter, _ *http.Request) {
	fmt.Println("Serving index")
	tmpl, err := template.ParseFiles("assets/html/index.html")
	if err != nil {

	}

	err = tmpl.Execute(rw, struct{ Title string }{Title: "My projects"})
	if err != nil {
		panic(err.Error())
	}
}
