# OnionTree CLI utility

Manage OnionTree repository.

## Installation

```
$ go get https://github.com/oniontree-org/go-oniontree/cmd/oniontree
```

## Usage

```
NAME:
   oniontree - Manage OnionTree repository

USAGE:
   oniontree [global options] command [command options] [arguments...]

VERSION:
   0.1

COMMANDS:
   init    Initialize a new repository
   add     Add a new service to the repository
   update  Update a service
   show    Show service's content
   remove  Remove services from the repository
   tag     Tag services
   untag   Untag services
   lint    Lint the repository content

GLOBAL OPTIONS:
   -C value       change directory to (default: ".")
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Examples

### Add a new service

```
$ git clone https://github.com/onionltd/oniontree
$ cd oniontree/
$ oniontree add --name "Dummy Service" \
                --description "In all people I see *myself*, none more and not one a barleycorn less" \
                --url "http://2efafjga32zajfcny.onion" \
                --public-key /tmp/dummyservice-pgp.asc \
                dummyservice
```

### Tag a service

```
$ oniontree tag --name dummy --name test dummyservice
```
