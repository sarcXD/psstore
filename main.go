package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"

	"github.com/sarcxd/psstore/constants"
)

type safeKVStore struct {
	store  map[string]int32
	rwLock sync.RWMutex
}

type SimpleStore int

var kvStore *safeKVStore

/*
Function adds args.Value to args.Key index in the kvStore
*/
func (s SimpleStore) Add(args constants.InsArgs, reply *int) error {
	key, value := args.Key, args.Value
	kvStore.rwLock.RLock()
	_, exs := kvStore.store[key]
	kvStore.rwLock.RUnlock()
	if exs {
		return constants.ErrHandlef(constants.KEY_EXIST, key)
	}
	kvStore.rwLock.Lock()
	kvStore.store[key] = value
	kvStore.rwLock.Unlock()
	*reply = 1 // SUCCESS
	return nil
}

/*
When many entries need to be added, it will be more efficient to have one op that will
add the values in bulk
args: contain list of elements (Key Value Pairs) to add
reply: is an int and list of all elements that were not added
*/
func (s SimpleStore) BulkAdd(args constants.BulkInsArgs, reply *constants.BulkInsReply) error {
	insList := args.InsList
	kvStore.rwLock.RLock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		_, prs := kvStore.store[iVal.Key]
		if prs {
			iVal.Ins = false
		} else {
			iVal.Ins = true
		}
		insList[i] = iVal
	}
	kvStore.rwLock.RUnlock()
	kvStore.rwLock.Lock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		if iVal.Ins {
			kvStore.store[iVal.Key] = iVal.Value
		} else {
			reject := constants.InsArgs{Key: iVal.Key, Value: iVal.Value}
			reply.Rejects = append(reply.Rejects, reject)
		}
	}
	kvStore.rwLock.Unlock()
	return nil
}

// if data exists, update data
func (s SimpleStore) Update(args constants.InsArgs, reply *int) error {
	key, value := args.Key, args.Value
	kvStore.rwLock.RLock()
	_, exs := kvStore.store[key]
	kvStore.rwLock.RUnlock()
	if !exs {
		return constants.ErrHandlef(constants.DNE, key)
	}
	kvStore.rwLock.Lock()
	kvStore.store[key] = value
	kvStore.rwLock.Unlock()
	*reply = 1 // SUCCESS
	return nil
}

func (s SimpleStore) BulkUpdate(args constants.BulkInsArgs, reply *constants.BulkInsReply) error {
	insList := args.InsList
	kvStore.rwLock.RLock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		_, prs := kvStore.store[iVal.Key]
		if !prs {
			iVal.Ins = false
		} else {
			iVal.Ins = true
		}
		insList[i] = iVal
	}
	kvStore.rwLock.RUnlock()
	kvStore.rwLock.Lock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		if iVal.Ins {
			kvStore.store[iVal.Key] = iVal.Value
		} else {
			reject := constants.InsArgs{Key: iVal.Key, Value: iVal.Value}
			reply.Rejects = append(reply.Rejects, reject)
		}
	}
	kvStore.rwLock.Unlock()
	return nil
}

func (s SimpleStore) Get(args constants.GetArgs, reply *int32) error {
	kvStore.rwLock.RLock()
	value, prs := kvStore.store[args.Key]
	kvStore.rwLock.RUnlock()
	if !prs {
		return constants.ErrHandle(constants.DNE)
	}
	*reply = value
	return nil
}

func (s SimpleStore) Delete(args constants.GetArgs, reply *int) error {
	kvStore.rwLock.RLock()
	_, prs := kvStore.store[args.Key]
	kvStore.rwLock.RUnlock()
	if !prs {
		return constants.ErrHandlef(constants.DNE, args.Key)
	}
	kvStore.rwLock.Lock()
	delete(kvStore.store, args.Key)
	kvStore.rwLock.Unlock()
	*reply = 1
	return nil
}

func (s SimpleStore) Clear(args constants.GetArgs, reply *int) error {
	if len(kvStore.store) == 0 {
		return constants.ErrHandle(constants.EMPTY_KV)
	}
	// re-initialize kvstore
	kvStore = new(safeKVStore)
	kvStore.store = make(map[string]int32)
	*reply = 1
	return nil
}

// ********************!Storage handling

// restore program data from backup
func restore() {}

// make program backups before server exits
func backup() {}

// *************************************

/*
go run main.go [args]
args:
--port xyz | port number to run service on
*/
func main() {
	args := os.Args
	ip := "localhost"
	port := "8020"
	if len(args) > 0 {
		i := 0
		for i < len(args) {
			switch args[i] {
			case "--port":
				i++
				port = args[i]
			}
			i++
		}
	}
	address := fmt.Sprintf("%s:%s", ip, port)
	fmt.Printf("server running on %s\n", address)
	// Initialize all data stores
	kvStore = new(safeKVStore)
	kvStore.store = make(map[string]int32)
	// populate data stores with data from backups
	restore()
	psstore := new(SimpleStore)
	rpc.Register(psstore)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer backup()
	http.Serve(l, nil)
}
