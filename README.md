#### kecho (kitty echo) contain middlewares,handlers,... for the echo.

#### Install
```
go get github.com/Kamva/kecho
```

##### Middlewares
* log: set new log handler as context log that contains:
    - request id in eac log record.
    - users data in each log record.

* transltion: Set new translator in context that localized with
users accept-languages and then fallback and default languages.


##### Handlers
* error handler: handle kitty errors.
    
##### Middleware dependencies:
* `kecho.CurrentUser` middleware requires
    - `kecho.JWT` middleware (load `JWT` middleware before `CurrentUser`).
* `kecho.KittyContext` middleware requires 
    - echo `middleware.RequestID`
    - kitty `kecho.CorrelationID`
    - kitty `kecho.CurrentUser` middleware.
    

#### Todo:
- [ ] Map echo errors (see errors list in `echo.go:263`) to kitty error with translation.
- [ ] Tests
- [ ] Add badges to readme.
- [ ] CI 

