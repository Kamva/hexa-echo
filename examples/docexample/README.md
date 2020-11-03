to generate docs: cd to the main package and 
run following commands:  
```bash
$ cd examples/docexample
$ go run main.go extract # to extract new routes 
$ go run main.go trim # to remove old deleted routes
$ swagger generate spec -o ./swagger.json
$ swagger generate spec -o ./swagger.yamlswagger generate spec -o ./swagger.yaml