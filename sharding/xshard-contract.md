

### XShard Contract Metadata

* Shard-ID : uint64
* All-Shard : bool
* Invoked-Contract : []Address


### All-Shard contracts

1. Shard-ID == service.ShardID
2. Invoked-Contract == []
3. All-ShardID = true


### Root-Shard contracts

1. Shard-ID == service.ShardID
2. Invoked-Contract == []
3. All-ShardID = false


### Child-Shard contracts:

1. Shard-ID.ParentShard == service.ShardID
2. All-ShardID = false
3. no cycle-dependency


