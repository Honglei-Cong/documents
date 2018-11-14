
## Transactions in RocksDB


RocksDB supports transactions when using a TransactionDB.

Transactions have a simple BEGIN/COMMIT/ROLLBACK api and allow applications to modify their data concurrently while letting RocksDB handle the conflict checking.  

RocksDB supports both pessimistic and optimistic concurrent control.

RocksDB provides Atomicity by default when writing multiple keys via **WriteBatch**.

Transactions provide a way to gurantee that a batch of writes will only be written if there are no conflicts.  Simliar to a WriteBatch, no other threads can see the changes in a transaction until it has been committed.


### TransactionDB

When using a TransactionDB, all keys that are written are locked internally by RocksDB to perform confilict detection.  If a key cannot be locked, the operation will return an error.

A TransactionDB will do conflict checking for all write operations, including writes performed outside of a Transaction.


### OptimisticTransactionDB

Optimistic Transactions do not take any locks when preparing writes.  Instead, they rely on doing conflict-detection at commit time to validate that no other writers have modified the keys being written by the current transaction.  If there is a conflict with another write, the commit will return an error and no keys will be written.

OptimisticTransactionDB may be more performant than a TransactionDB for workloads that have many non-transactional writes and few transactions.


```
        OptimisticTransactionDB* txn_db;
        Status s = OptimisticTransactionDB::Open(options, path, &txn_db);
        DB* db = txn_db->GetBaseDB();
        
        OptimisticTransaction* txn = txn_db->BeginTransaction(write_options, txn_options);
        txn->Put(“key”, “value”);
        txn->Delete(“key2”);
        txn->Merge(“key3”, “value”);
        s = txn->Commit();
        delete txn;
```

### Snapshot

If you want to guarantee that no else has written a key since the start of the transaction, this can be accomplished by calling SetSnapshot() after creating the transaction.

In the following example, if this were a TransactionDB, the Put() would have failed.  If this were a OptimisticTransactionDB, the Commit() would failed.

```
        txn = txn_db->BeginTransaction(write_options);
        txn->SetSnapshot();
        
        // Write to key1 OUTSIDE of the transaction
        db->Put(write_options, “key1”, “value0”);
        
        // Write to key1 IN transaction
        s = txn->Put(“key1”, “value1”);
        s = txn->Commit();
        // Transaction will NOT commit since key1 was written outside of this transaction after SetSnapshot() was called (even though this write
        // occurred before this key was written in this transaction).
```

### Repeatable Read

You can achieve *repeatable reads* when reading through a transaction by setting a Snapshot in the ReadOptions.

```
        read_options.snapshot = txn->GetSnapshot();
        Status s = txn->GetForUpdate(read_options, “key1”, &value);
        …
        s = txn->GetForUpdate(read_options, “key1”, &value);
        db->ReleaseSnapshot(read_options.snapshot);
```


### Guarding against Read-Write Conflicts

GetForUpdate() will ensure that no other writer modifies any keys that were read by this transaction.

```
        // Start a transaction 
        txn = txn_db->BeginTransaction(write_options);
        
        // Read key1 in this transaction
        Status s = txn->GetForUpdate(read_options, “key1”, &value);
        
        // Write to key1 OUTSIDE of the transactio,     <---- if TransactionDB: waiting until timeout or txn commit
        s = db->Put(write_options, “key1”, “value0”);   <---- if OptimisticTransactionDB, txn commit will be failed.
```

If this transaction was created by a TransactionDB, the Put would either timeout or block until the transaction commits or aborts. If this transaction were created by an OptimisticTransactionDB(), then the Put would succeed, but the transaction would not succeed if txn->Commit() were called.

```
        // Repeat the previous example but just do a Get() instead of a GetForUpdate()
        txn = txn_db->BeginTransaction(write_options);
        
        // Read key1 in this transaction
        Status s = txn->Get(read_options, “key1”, &value);
        
        // Write to key1 OUTSIDE of the transaction
        s = db->Put(write_options, “key1”, “value0”);
        
        // No conflict since transactions only do conflict checking for keys read using GetForUpdate().
        s = txn->Commit();
```






