# psstore
A simple in-memory datastore

## Usage
Built using `go1.17.12 linux/amd64`  
### server
1. start server with `go run .`  
2. specify port with `go run . --port 8080`  
### client
The rpc uses standard rpc format, whilst the format should be 
language independant, it isn't at the moment. The argument and reply
types are in the constants package of the psstore module.  
*See the rpc-client/main.go file to get an idea of usage*

## Features
See Milestones.md for an idea of future features  

### Developed
**Basic Ops**
* Key Value Map (Hashtable)
  - Add
  - Get
  - Update
  - Delete
  - Clear da
  - Bulk Add
  - Bulk Update
