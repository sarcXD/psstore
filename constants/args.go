package constants

type InsArgs struct {
	Key   string
	Value int32
}

/* When performing bulk inserts struct contains
a Key, Value to identify Key and Value pairs to insert
Ins is a conditional that will be computed on each KVPair
to determine if an element should be inserted (added/updated)
Ins cases:
1. if op is ADD, key is present, Ins = false
2. if op is Update, key is not present, Ins = false
3. if op is Update, key is present, Value is same, Ins = false
*/
type CheckedIns struct {
	Key   string
	Value int32
	Ins   bool
}

type BulkInsArgs struct {
	InsList    []CheckedIns
	Compressed bool
}

type BulkInsReply struct {
	Rejects []InsArgs
}

type GetArgs struct {
	Key string
}
