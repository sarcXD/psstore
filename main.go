package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/sarcxd/psstore/constants"
	kv "github.com/sarcxd/psstore/stores/keyValue"
)

type SimpleKV int

var kvI32 *kv.KVStore_num
var kvStr *kv.KVStore_str

/********************************** STRING **************************/

func (s SimpleKV) AddStr(args constants.InsArgs_str, reply *int) error {
	key, value := args.Key, args.Value
	kvStr.RwLock.RLock()
	_, exs := kvStr.Store[key]
	kvStr.RwLock.RUnlock()
	if exs {
		return constants.ErrHandlef(constants.KEY_EXIST, key)
	}
	kvStr.RwLock.Lock()
	kvStr.Store[key] = value
	kvStr.RwLock.Unlock()
	*reply = 1 // SUCCESS
	return nil
}

// if data exists, update data
func (s SimpleKV) UpdateStr(args constants.InsArgs_str, reply *int) error {
	key, value := args.Key, args.Value
	kvStr.RwLock.RLock()
	_, exs := kvStr.Store[key]
	kvStr.RwLock.RUnlock()
	if !exs {
		return constants.ErrHandlef(constants.DNE, key)
	}
	kvStr.RwLock.Lock()
	kvStr.Store[key] = value
	kvStr.RwLock.Unlock()
	*reply = 1 // SUCCESS
	return nil
}

/*
When many entries need to be added, it will be more efficient to have one op that will
add the values in bulk
args: contain list of elements (Key Value Pairs) to add
reply: is an int and list of all elements that were not added
*/
func (s SimpleKV) BulkAddStr(args constants.BulkInsArgs_str, reply *constants.BulkInsReply_str) error {
	insList := args.InsList
	kvStr.RwLock.RLock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		_, prs := kvStr.Store[iVal.Key]
		if prs {
			iVal.Ins = false
		} else {
			iVal.Ins = true
		}
		insList[i] = iVal
	}
	kvStr.RwLock.RUnlock()
	kvStr.RwLock.Lock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		if iVal.Ins {
			kvStr.Store[iVal.Key] = iVal.Value
		} else {
			reject := constants.InsArgs_str{Key: iVal.Key, Value: iVal.Value}
			reply.Rejects = append(reply.Rejects, reject)
		}
	}
	kvStr.RwLock.Unlock()
	return nil
}

func (s SimpleKV) BulkUpdate(args constants.BulkInsArgs_str, reply *constants.BulkInsReply_str) error {
	insList := args.InsList
	kvStr.RwLock.RLock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		_, prs := kvStr.Store[iVal.Key]
		if !prs {
			iVal.Ins = false
		} else {
			iVal.Ins = true
		}
		insList[i] = iVal
	}
	kvStr.RwLock.RUnlock()
	kvStr.RwLock.Lock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		if iVal.Ins {
			kvStr.Store[iVal.Key] = iVal.Value
		} else {
			reject := constants.InsArgs_str{Key: iVal.Key, Value: iVal.Value}
			reply.Rejects = append(reply.Rejects, reject)
		}
	}
	kvStr.RwLock.Unlock()
	return nil
}

func (s SimpleKV) GetStr(key string, reply *string) error {
	kvStr.RwLock.RLock()
	value, prs := kvStr.Store[key]
	kvStr.RwLock.RUnlock()
	if !prs {
		return constants.ErrDne
	}
	*reply = value
	return nil
}

func (s SimpleKV) DeleteStr(key string, reply *int) error {
	kvStr.RwLock.RLock()
	_, prs := kvStr.Store[key]
	kvStr.RwLock.RUnlock()
	if !prs {
		return constants.ErrHandlef(constants.DNE, key)
	}
	kvStr.RwLock.Lock()
	delete(kvStr.Store, key)
	kvStr.RwLock.Unlock()
	*reply = 1
	return nil
}

func (s SimpleKV) ClearStr(args interface{}, reply *int) error {
	if len(kvStr.Store) == 0 {
		return constants.ErrEmptyKv
	}
	// re-initialize KVStore_str
	kvStr = new(kv.KVStore_str)
	kvStr.Store = make(map[string]string)
	*reply = 1
	return nil
}

/********************************** I32 **************************/

/*
Function adds args.Value to args.Key index in the kvStore
*/
func (s SimpleKV) AddI32(args constants.InsArgs_i32, reply *int) error {
	key, value := args.Key, args.Value
	kvI32.RwLock.RLock()
	_, exs := kvI32.Store[key]
	kvI32.RwLock.RUnlock()
	if exs {
		return constants.ErrHandlef(constants.KEY_EXIST, key)
	}
	kvI32.RwLock.Lock()
	kvI32.Store[key] = value
	kvI32.RwLock.Unlock()
	*reply = 1 // SUCCESS
	return nil
}

// if data exists, update data
func (s SimpleKV) UpdateI32(args constants.InsArgs_i32, reply *int) error {
	key, value := args.Key, args.Value
	kvI32.RwLock.RLock()
	_, exs := kvI32.Store[key]
	kvI32.RwLock.RUnlock()
	if !exs {
		return constants.ErrHandlef(constants.DNE, key)
	}
	kvI32.RwLock.Lock()
	kvI32.Store[key] = value
	kvI32.RwLock.Unlock()
	*reply = 1 // SUCCESS
	return nil
}

/*
When many entries need to be added, it will be more efficient to have one op that will
add the values in bulk
args: contain list of elements (Key Value Pairs) to add
reply: is an int and list of all elements that were not added
*/
func (s SimpleKV) BulkAddI32(args constants.BulkInsArgs_i32, reply *constants.BulkInsReply_i32) error {
	insList := args.InsList
	kvI32.RwLock.RLock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		_, prs := kvI32.Store[iVal.Key]
		if prs {
			iVal.Ins = false
		} else {
			iVal.Ins = true
		}
		insList[i] = iVal
	}
	kvI32.RwLock.RUnlock()
	kvI32.RwLock.Lock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		if iVal.Ins {
			kvI32.Store[iVal.Key] = iVal.Value
		} else {
			reject := constants.InsArgs_i32{Key: iVal.Key, Value: iVal.Value}
			reply.Rejects = append(reply.Rejects, reject)
		}
	}
	kvI32.RwLock.Unlock()
	return nil
}

func (s SimpleKV) BulkUpdateI32(args constants.BulkInsArgs_i32, reply *constants.BulkInsReply_i32) error {
	insList := args.InsList
	kvI32.RwLock.RLock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		_, prs := kvI32.Store[iVal.Key]
		if !prs {
			iVal.Ins = false
		} else {
			iVal.Ins = true
		}
		insList[i] = iVal
	}
	kvI32.RwLock.RUnlock()
	kvI32.RwLock.Lock()
	for i := 0; i < len(insList); i++ {
		iVal := insList[i]
		if iVal.Ins {
			kvI32.Store[iVal.Key] = iVal.Value
		} else {
			reject := constants.InsArgs_i32{Key: iVal.Key, Value: iVal.Value}
			reply.Rejects = append(reply.Rejects, reject)
		}
	}
	kvI32.RwLock.Unlock()
	return nil
}

func (s SimpleKV) GetI32(key string, reply *int32) error {
	kvI32.RwLock.RLock()
	value, prs := kvI32.Store[key]
	kvI32.RwLock.RUnlock()
	if !prs {
		return constants.ErrDne
	}
	*reply = value
	return nil
}

func (s SimpleKV) DeleteI32(key string, reply *int) error {
	kvI32.RwLock.RLock()
	_, prs := kvI32.Store[key]
	kvI32.RwLock.RUnlock()
	if !prs {
		return constants.ErrHandlef(constants.DNE, key)
	}
	kvI32.RwLock.Lock()
	delete(kvI32.Store, key)
	kvI32.RwLock.Unlock()
	*reply = 1
	return nil
}

func (s SimpleKV) ClearI32(arg string, reply *int) error {
	if len(kvI32.Store) == 0 {
		return constants.ErrEmptyKv
	}
	// re-initialize KVStore_num
	kvI32 = new(kv.KVStore_num)
	kvI32.Store = make(map[string]int32)
	*reply = 1
	return nil
}

// clearDB - Clears all data stores and re-initializes them

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
	kvI32 = new(kv.KVStore_num)
	kvI32.Store = make(map[string]int32)
	kvStr = new(kv.KVStore_str)
	kvStr.Store = make(map[string]string)

	// populate data stores with data from backups
	restore()
	psstore := new(SimpleKV)
	rpc.Register(psstore)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", address)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	defer backup()
	http.Serve(l, nil)
}
