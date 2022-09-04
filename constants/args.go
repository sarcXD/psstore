package constants

type InsArgs_i32 struct {
	Key   string
	Value int32
}
type InsArgs_str struct {
	Key   string
	Value string
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
type CheckedIns_i32 struct {
	Key   string
	Value int32
	Ins   bool
}
type CheckedIns_str struct {
	Key   string
	Value string
	Ins   bool
}
type BulkInsArgs_i32 struct {
	InsList    []CheckedIns_i32
	Compressed bool
}
type BulkInsArgs_str struct {
	InsList    []CheckedIns_str
	Compressed bool
}

type BulkInsReply_i32 struct {
	Rejects []InsArgs_i32
}
type BulkInsReply_str struct {
	Rejects []InsArgs_str
}

type GetArgs struct {
	Key string
}
