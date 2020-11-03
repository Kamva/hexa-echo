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

// route:begin: hi::say
// swagger:route GET /hi  hiSayParams
//
// responses:
//   200: hiSaySuccessResponse

// swagger:parameters hiSayParams
type hiSayParamsWrapper struct {
	// in:body
	Body struct{
	    // DOCTODO: place your params body here
	}
}

// success response
// swagger:response hiSaySuccessResponse
type hiSayResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	// DOCTODO: place your body data here
    }
}

// route:end: hi::say

// route:begin: hi::create
// swagger:route POST /hi/{id}/{code}  hiCreateParams
//
// responses:
//   200: hiCreateSuccessResponse

// swagger:parameters hiCreateParams
type hiCreateParamsWrapper struct {
     // in:path
     Id string `json:"id"`
  
     // in:path
     Code string `json:"code"`
  
	// in:body
	Body struct{
	    // DOCTODO: place your params body here
	}
}

// success response
// swagger:response hiCreateSuccessResponse
type hiCreateResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	// DOCTODO: place your body data here
    }
}

// route:end: hi::create

