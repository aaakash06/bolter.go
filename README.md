#Bolter.go

1) Setup 
```bash
go mod init bolter.new
```

2) Dependecies install
```bash
go get -u github.com/gorilla/mux
```

3) air
```bash
go install github.com/air-verse/air@latest
alias air = $(go env GOPATH)/bin/air 
air init
air
```
