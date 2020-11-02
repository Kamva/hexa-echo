{{.BeginRouteVal}}
// swagger:route {{hold .Name "Method"}} {{hold .Name "Path"}} {{hold .Name "TagsString"}} {{hold .Name "ParamsId"}}
//
// responses:
//   2xx: {{hold .Name "SuccessRespId"}}

// swagger:parameters {{hold .Name "ParamsId"}}
type {{.Name}}ParamsWrapper struct {
	// in:body
	Body struct{
	    // DOCTODO: place your params body here
	}
}

// success response
// swagger:response {{hold .Name "SuccessRespId"}}
type {{.Name}}ResponseWrapper struct {
	// in:body
	Body struct{
	    // response code
    	Code string `json:"code"`
    	// DOCTODO: place your body data here
    }
}

{{.EndRouteVal}}

