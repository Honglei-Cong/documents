

## Hotel Contract

预订Hotel 房间示例合约


```
roomKey = "Room"

def Main (operation, args):
        if operation == "book":
                account = args[0]
                roomNo = args[1]
                return bookRoom(account, roomNo)
        if operation == "checkout":
                account = args[0]
                roomNo = args[1]
                return checkoutRoom(account, roomNo)
        return False

def bookRoom  (account, roomNo):
        if CheckWitness(account):
                key = concatKey(roomKey, roomNo)
                roomBooked = Get(ctx, key)
                if not roomBooked || roomBooked == False:
                        Put(ctx, key, account)
                        Notify("room booking succeeded")
                        return True
                else:
                        Notify("room had been booked")
                        return False
        else:
                Notify("checkwitness failed")
        return False
```

## Train Ticket Contract

预订火车票示例合约

```
ticketKey = "Ticket"

def Main (operation, args):
        if operation == "buy":
                account = args[0]
                ticketNo = args[1]
                return buyTicket(account, ticketNo)
        return False

def buyTicket (account, ticketNo):
        if CheckWitness(account):
                key = concat(ticketKey, ticketNo)
                ticketSold = Get(ctx, key)
                if not ticketSold || ticketSold == False:
                        Put(ctx, key, account)
                        Notify("ticket buying succeeded")
                        return True
                else:
                        Notify("tiket sold")
                        return False
        else:
                Notify("CheckWitness failed")
                return False
```

## GoTravel Contract

该合约基于Hotel 和 Ticket 两个合约，通过RemoteInvoke调用这两个合约。

如果此合约和Hotel处于同一shard，而且资源没有被占用，RemoteInvoke 将等同于 Invoke。
如果不处于同一shard，remoteInvoke将会触发cross-shard调用，当前执行逻辑将被中断，由Callback中继续处理。


```
RoomContractName = "RoomContract"
TicketContractName = "TicketContract"

def Main (operation, args):
        if operatoin == "goTravel":
                account = args[0]
                roomNo = args[1]
                ticketNo = args[2]
                res = RemoteInvoke(RoomContractName, "book", [account, roomNo])
                if res == False:
                        Notify("room booking failed")
                        return False
                res = RemoteInvoke(TicketContractName, "buy", [account, ticketNo])
                if res == False:
                        Notify("ticket buying failed")
                        return False
                Notify("travel ready")
                return True
        return False

def Callback(contract, operation, result):
        if contract == RoomContractName && operation == "book":
                if result == True:
                        res = RemoteInvoke(TicketContractName, "buy", [account, ticketNo])
                        if res == False:
                                Notify("ticket buying failed")
                                return False
                        Notify("travel ready")
                        return True
                else:
                        Notify("room booking failed")
                        return False

        if contract == TicketContractName && operation == "buy":
                if result == True:
                        Notify("travel ready")
                        return True
                else:
                        Notify("ticket buying failed")
                        return False

        Notify("unknown callback")
        return False
```


