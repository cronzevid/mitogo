### Mitogo

#### How to run

1. `git pull git@github.com:cronzevid/mitogo.git`
2. `cd mitogo`
3. `go mod tidy`

#### How to use

1. Put desired go code into file in `payloads/` directory
2. Run it on remote: `go run . -host 10.10.0.1 -file ./payloads/demo.go`
