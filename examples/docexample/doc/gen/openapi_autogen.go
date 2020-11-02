// Package classification awesome.
//
// Documentation of our awesome API.
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: localhost:1323
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic
//
//    SecurityDefinitions:
//    basic:
//      type: basic
//
// swagger:meta
package doc

import "github.com/kamva/hexa-echo/examples/docexample/api"

// route:begin: hi::say
// swagger:route POST /hi  hiSayParams
// say hi to yourself.
// responses:
//   200: hiSaySuccessResponse

// swagger:parameters hiSayParams
type hiSayParamsWrapper struct {
	// in:body
	Body api.HiRequest `json:"body"`
}

// success response
// swagger:response hiSaySuccessResponse
type hiSayResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	Data api.HiResponse `json:"data"`
    }
}

// route:end: hi::say

