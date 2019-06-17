

XShard Contract Metadata

* Shard-ID : uint64
* All-Shard : bool
* Invoked-Contract : []Address


For All-Shard contracts:

1. Shard-ID == service.ShardID
2. Invoked-Contract == []
3. All-ShardID = true


For Root-Shard contracts, same as all-shard contracts

1. Shard-ID == service.ShardID
2. Invoked-Contract == []
3. All-ShardID = false


For Child-Shard contracts:

1. Shard-ID.ParentShard == service.ShardID
2. All-ShardID = false
3. no cycle-dependency


