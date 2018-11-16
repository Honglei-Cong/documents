
## MVCC

MVCC uses timestamps(TS) , and incrementing transaction IDs, to achieve transactional consistency.  MVCC ensures a transaction (T) never has to wait to read a database object (P) by maintaining several versions of the object.  Each version of object P has both a Read Timestamp (RTS) and a Write Timestamp (WTS) which lets a particular transaction T\_i read the most recent version of the object which precedes the transaction's Read Timestamp RTS(T\_i).


If transaction T\_i wants to write to object P, and there is also another transaction T\_k happening to the same object, the Read Timestamp RTS(T\_i) must precedes the Read Timestamp RTS(T\_k), i.e., RTS(T\_i) < RTS(T\_k), for the object Write Operation to succeed.

A Write cannot complete if there are other outstanding transactions with an earlier Read Timestamp (RTS) to the same object.

To restate: every object (P) has a Timestmap (TS), however if transaction T\_i wants to write to an object, and the transaction has a Timestmap (TS) that is earlier than the object's current Read Timestamp, TS(T\_i) < RTS(P), then the transaction is aborted and restarted.  Otherwise, T\_i creates a new version of object P and sets the read/write timestamp TS of the new version to the timestamp of the transaction TS = TS(T\_i).



