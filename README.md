# WordSSH

Play wordle via ssh and earn points

## How to run the server

Clone the project and cd into the directory

Run `go mod tidy` to fetch the dependecies
Run `go build` to build the project

You will have a binary called wordssh.
Run that to start the server.

It will run on localhost:23234 by default. To change that, you can edit lines 28 and 29 of main.go

## How to connect to the server

Run `ssh <server ip> -p <server port>`
You must have a SSH private/public key pair 
To generate one you can run `ssh-keygen` and follow the steps.
