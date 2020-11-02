// Package classification awesome.
//
// Documentation of our awesome API.
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//     Host: some-url.com
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

// route:begin: hi::say
// swagger:route POST /hi hi,say hiSayParams
//
// responses:
//   2xx: hiSaySuccessResponse

// swagger:parameters hiSayParams
type hiSayParamsWrapper struct {
	// in:body
	Body struct{
	    // -> place your params body here
	}
}

// success response
// swagger:response hiSaySuccessResponse
type hiSayResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	// -> place your body data here
    }
}

// route:end: hi::say

