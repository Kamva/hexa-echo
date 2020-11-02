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
// swagger:route {{.hiSay.Method}} {{.hiSay.Path}} {{.hiSay.TagsString}} {{.hiSay.ParamsId}}
// say hi to yourself.
// responses:
//   200: {{.hiSay.SuccessRespId}}

// swagger:parameters {{.hiSay.ParamsId}}
type hiSayParamsWrapper struct {
	// in:body
	Body api.HiRequest `json:"body"`
}

// success response
// swagger:response {{.hiSay.SuccessRespId}}
type hiSayResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	Data api.HiResponse `json:"data"`
    }
}

// route:end: hi::say

