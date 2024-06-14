# Introduction

This project is for generating harvester related iPXE boot configs. Please follow steps to generate and use iPXE configs.

1. Replace all variables in `pkg/template/generate.go` with your own.
2. Run `go generate` to generate iPXE configs.
3. Run `go run server.go` to serve files in `./public/harvester` directory. HTTP port is 3333 by default.

