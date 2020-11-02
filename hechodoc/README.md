Conventions:
---
- You must not import package which contains openapi template file.
  this is because `go-swagger` generates two docs if we import it in our project.
  
- You must create openapi template go file
  initially (e.g., openapi_tpl.go), this is 
  because we just.
  

Notes:
-------
- `raw route name` is the raw route name which we set on route definition in 
  our source code. e.g., `hi::mehran`
- `route name` in the source code is camelcased value of the
  raw route name e.g., `hiMehran`.
  
  
How to genrate docs?
---
- Install swagger
- Place a package named something like `doc`.
- In the `doc` package create a file named something like `openapi_tpl.go`.
  This is your openapi comments file, so all of your
  documentation should be in this file, we will generate
  final go in the `doc/gen/openapi_autogen.go` file.
- In the `doc` package create another package named `gen`
- In the gen package create just an empty file named `openapi_autogen.go`
  __Important__ Set package name of this file same as its parent package name(here `doc`).
- In your main file import `doc/gen` package as blank. (e.g., `_ "github.com/kamva/hexa-echo/examples/docexample/doc/gen"`)
- Now in you source code create three command: `extract`,`trim` and `render`
- In each command in order do `extract template`, `trim template` and `render template`. see the example.

FAQ:
---
- How should I install swagger command? simply install
it as go package: 
```bash
GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger
```

