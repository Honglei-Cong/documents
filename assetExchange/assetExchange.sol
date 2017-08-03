import "AssetLedger.sol";

// TODO describe rules of exchange.
// In the AF market all users benefit from a "last look" facility to prevent 
// them falling foul of mechanical arb sniping by those exploiting the way
// the blockchain and smart contract data lags the real markets being tracked.
// The system proceeds through a series of batching steps, which in practice
// might correspond to some number of blocks (for example, each batching step
// corresponds to a block). Orders collected in batching step B_i are initially
// in a "pending" state although visible within the book. Orders then remain
// pending in step B_i+1 so that submitters may cancel them before they go
// live in step B_i+2 (the assumption is that submitters can see the final
// state of B_i, calculate what their orders may execute against in step
// B_i+2, and then submit cancellation transactions for inclusion in step
// B_i+1 if execution will not be favorable). The computational cost of
// submitting transactions will restrict the amomunt of "phantom liquidity"
// traders create to distort the market and discover prices. 
//
contract AssetExchange {
  // === Pile of orders ==
  // Orders are bids
  byte constant BID = 0x01;
  // Orders are asks
  byte constant ASK = 0x02;
    // === Order matching rules ==
  // SLO matching
  // Simple limit order: Only matches at specified price or better. Order
  // remains in book until either entire amount requested has been filled/sold
  // (which may occur as result of unlimited partial match executions over any
  // number of blocks) or order is canclled by owner.
  uint64 constant SIMPLE_LIMIT = 0x0100000000000000;
  // LWE matching
  // Limit order with expiry: Same as a normal limit order, but with an expiry
  // date. If order has not been cancelled or filled by the expiry date, it is
  // removed from the book automatically. The expiry date can be set to the
  // current block to create what is effectively a "kill or fill" order.
  uint64 constant LIMIT_WITH_EXPIRY = 0x0200000000000000;
  // MTD matching
  // Market order to DAO's expected best price, then DAO's best price: the
  // order first behaves like a kill or fill limit order set at DAO's expected
  // best price, and then any remaining units are exchanged with the DAO for
  // whatever it determines the best price should be. To the maximum possible
  // extent, the order executes upon submission to provide functionality for
  // "liquid wallet" applications. Since the DAO cannot reliably determine a
  // best price without first hedging its position with liquidity providers,
  // whose pricing it cannot control, it settles the trade for only a portion of
  // the required units immediately and it is the responsibility of Dapps to
  // handle this limitation. For example, if a MTD ask order to sell 100
  // D-FACEBOOK is submitted, and say (i) 10 D-FACEBOOK are sold in limit order
  // fashion, then (ii) the DAO might immediately only pay 85% of the expected
  // value of the remaining 90 D-FACEBOOK tokens, later returning additional
  // currency to the order owner to complete the settlement once it has been
  // able to determine what its best price was by acquiring hedges. Note that if
  // the DAO was forced to pre-commit to a best price it would be trivial for
  // arbitrage traders to exploit it. Final note, this order does not sit in the
  // book, it executes and the DAO forwards change to the owner.
  uint64 constant MARKET_TO_DAO = 0x0300000000000000;
  // LTM matching
  // Limit order with expiry, but if/when expires order converts into an MTD
  // order to fill remaining units.
  uint64 constant LIMIT_THEN_MTD = 0x0400000000000000;
  // === Pointers ==
  // Nil/unassigned
  uint64 constant NIL = 0xffffffffffffffff;
  
  // An order in a pile
  // NOTE: using 3*256 bits since EVM charges by whole word
  struct Order {
    uint128 price;    // price raised 10^18
    uint128 amount;   // amount raised 10^18
    uint32 submitted; // block batch when submitted, determines if pending
    uint64 orderSeqNo;// order sequence number, used resolve price
    uint128 config;   // 16 bytes of config e.g. flags, expiry
    address account;  // owner of order
    uint64 above;     // order above in pile
    uint64 below;     // order below in pile
  }

  // A pile of limit orders
  struct Pile {
    uint64 top;       // the top order in the pile
    uint64 size;      // the size of the pile
    byte orderType;   // type of orders stored: BID|ASK
  }

  // Order book containing orders that limit match. Note this can include
  // LTM orders that have not yet switched into MTD mode.
  struct OrderBook {
    Pile bids;        // pile of bids in limit matching mode
    Pile asks;        // pile of asks in limit matching mode
  }
  OrderBook book;
  
    // MTD orders where the owner has been provided an advance on the final
  // settlement and is now waiting for full settlement once the DAO has
  // determined its best price.
  struct MtdBalance {
    uint64 id;        // id of original order
    uint128 requested;// original amount asset or currency "given" to DAO
    uint128 advanced; // amount immediately advanced to owner order
  }

  // Pool of MTD balances waiting to be settled
  struct MtdPool {
    MtdBalance[] bids;
    MtdBalance[] asks;
  }

  // Id of asset type exchanged in this market instance
  string public assetTypeId;

  // Id of settlement currency
  string public xCurrencyId;

  // Asset Ledger for this asset
  AssetLedger public currencyLedger;
  AssetLedger public assetLedger;

  // What the DAO has given as an indicative sell price
  uint128 public daoIndicativeAskPrice;

  // What the DAO has given as an indicative buy price
  uint128 public daoIndicativeBidPrice;

  // Proportion of DAO order that can be advanced
  uint128 public mtdAdvance = uint128(75);// 75%

  // Current block batch index
  uint32 public batchIdx;
  
  // Store pending order id queue
  uint64[] pendingOrders;
  uint64 pendingStart = 0;
  
  // Store orders
  mapping(uint64 => Order) orders;

	function AssetExchange(string _assetTypeId) {
	  assetTypeId = _assetTypeId;
	  xCurrencyId = "XBIT";
	  currencyLedger = AssetLedger(xCurrencyId, address(0));
	  batchIdx = uint32(block.number);
	  book.bids = Pile({top: NIL, size: 0, orderType: BID});
    book.asks = Pile({top: NIL, size: 0, orderType: ASK});
	}
	
	// Increment the order batch index
	function incOrderBatchIndex() {
	  // DEBUGGING ONLY. Will use block.number
	  batchIdx++;
	}
	
	// The current order batch.
	function orderBatchIndex() returns(uint32) {
	  // DEBUGGING ONLY. Will return uint32(block.number)
	  return batchIdx;
	}
	
  // Return a copy of the book
  function marketOrderBook() returns (uint64[5][], uint64[5][], uint64[5][]) {
    return (pendingQueue(), orderPile(BID), orderPile(ASK));
  }
  
  // Return a copy of the pending order queue
  function pendingQueue()
    private
    returns (uint64[5][])
  {
    uint pendingCount = pendingOrders.length-pendingStart;
    uint64[5][] memory copy;
    if (pendingCount == 0)
      return copy;
    copy = new uint64[5][](pendingCount);
    for (var i=0; i<pendingCount; i++) {
      uint64 id = pendingOrders[pendingStart+i];
      Order o = orders[id];
      copy[i] = copyOrderFields(id, o);
    }
    return copy;
  }
  
  // Return a copy of an order pile
  function orderPile(byte orderType) 
    private
    returns (uint64[5][])
  {
    Pile pile = orderType == BID ? book.bids : book.asks;
    uint64[5][] memory copy;
    if (pile.size == 0)
      return copy;
    copy = new uint64[5][](pile.size);
    uint64 id = pile.top;
    uint i = 0;
    while (id != NIL) {
      Order o = orders[id];
      copy[i] = copyOrderFields(id, o);
      id = o.below;
      i++;
    }
    return copy;
  }
  
  function copyOrderFields(uint64 id, Order memory o)
    private
    returns(uint64[5])
  {
    uint64[5] memory copy;
    copy[0] = id;
    copy[1] = uint64(o.price);
    copy[2] = uint64(o.amount);
    copy[3] = o.submitted;
    copy[4] = uint64(o.account);
    return copy;
  }
  
  // Add a bid order to the book. Returns id of new order, any amount of
  // asset immediately purchased, any amount of currency remaining in a bid
  // (if no currency remains in a bid, then bid has completed)
  function addBid(
    uint64 id,
    address account,
    uint128 amount,
    uint128 price,
    uint64 matching,
    uint64 config)
  {
    // Validate call
    if (amount == 0 || price == 0 || orders[id].price != 0)
      throw;
    // Transfer payment from sender to exchange
    uint128 spend = price*amount;
    currencyLedger.transferToDAO(spend, account);
    // Record pending order
    recordPendingOrder(account, id, price, amount, config);
  }

  // Add a ask order to the book. Returns id of new order, any amount of
  // currency immediately received, and any amount of asset still for sale (if
  // no asset remains for sale, then order complete).
  function addAsk(
    uint64 id,
    address account,
    uint128 amount,
    uint128 price,
    uint64 matching,
    uint64 config)
  {
    // Vaidate call
    if (amount == 0 || price == 0 || orders[id].price != 0)
      throw;
    // Transfer amount of asset to be sold
    assetLedger.transferToDAO(amount, account);
    // Record pending order
    recordPendingOrder(account, id, price, amount, config);
  }
  
  // Process live orders that "cross the book" by creating trades. We rely on
  // this method to be called by somebody so that pending trades that have tran-
  // sitioned into live trades get processed. In practice we might expect
  // committed participants such as liquidity providers to make the calls
  // but anyone can call the function if they want to make sure that it executes.
  // TODO: modify Ethereum so that a block "initiailization" script can be
  // registered by transactions, with execution to be paid from the balance of a
  // contract (here such as an exchange or other DAO contract). Once such a
  // block initializer has been registered, in any subsequent block, it will be
  // called *before* any transactions from the mempool are applied.
  function executeLimitOrders() {
    while (pendingStart < pendingOrders.length) {
      // Get next id in pending queue
      uint64 id = pendingOrders[pendingStart];
      // Continue if empty slot/deleted pending order
      if (id == 0) {
        pendingStart++;
        continue;
      }
      // Is this pending order ready to become live?
      Order o = orders[id];
      if (o.submitted+2 > batchIdx)
        return; // => no following pending order ready go live either.
      // Remove from pending list
      delete pendingOrders[pendingStart];
      pendingStart++;
      // Insert this order into the pile
      byte orderType = id & 0x8000000000000000 == 0 ? BID : ASK;
      Pile pile = orderType == BID ? book.bids : book.asks;
      if (!insertOrder(pile, o, id))
        // If not to of pile, no trade possible so continue
        continue;
      // New order added to top of file, so process possible trades
      processTradeExecutions();
    }
  }
  
  function processTradeExecutions() 
    private
  {
    var (bId, bid) = (book.bids.top, orders[book.bids.top]);
    var (aId, ask) = (book.asks.top, orders[book.asks.top]);
    while (bId != NIL && aId != NIL) {
      // Check more trades to be done given current bid & ask piles
      if (bid.price < ask.price)
        break; // nope
      // Yes, at what price? Age decides who crosses who
      uint128 tradePrice = bid.orderSeqNo < ask.orderSeqNo ? bid.price : ask.price;
      // Calculate quantity traded and adjust orders
      uint128 bidMaxSpend = bid.price*bid.amount;
      uint128 maxAmount = bidMaxSpend/tradePrice;
      uint128 bought;
      uint128 received;
      if (maxAmount > ask.amount) {
        bought = ask.amount;
        bid.amount = (bidMaxSpend-bought*tradePrice)/bid.price;
        ask.amount = 0;
      } else if (maxAmount < ask.amount) {
        bought = maxAmount;
        bid.amount = 0;
        ask.amount -= bought;
      } else {
        bought = ask.amount;
        bid.amount = 0;
        ask.amount = 0;
      }
      received = bought*tradePrice;
      // Update currency and asset ledgers
      currencyLedger.transferFromDAO(ask.account, received);
      assetLedger.transferFromDAO(bid.account, bought);
      // Delete spent orders and move down piles
      if (bid.amount <= 0) {
        uint64 bFilled = bId;
        (bId, bid) = getOrderBelow(bid);
        deleteOrder(bFilled);
      }
      if (ask.amount <= 0) {
        uint64 aFilled = aId;
        (aId, ask) = getOrderBelow(ask);
        deleteOrder(aFilled);
      }
    }
  }
  
  function getOrderBelow(Order memory o)
    private
    returns(uint64, Order storage)
  {
    return (o.below, orders[o.below]);
  }
  
  function recordPendingOrder(
    address account,
    uint64 id,
    uint128 price,
    uint128 amount,
    uint128 config
  )
    private
  {
    // Make sure id not in use (e.g. client error, attempt subvert, etc)
    Order used = orders[id];
    if (used.price != 0)
      throw;
    // Add to queue of orders in pending state
    pendingOrders.push(id);
    // Create order object
    Order memory o = Order({price:price, amount:amount, submitted:batchIdx,
      orderSeqNo:uint64(pendingOrders.length-1), config:config, account:account,
      above: NIL, below: NIL});
    // Store order but do not link into pile sequence yet
    orders[id] = o;
  }
  
  // Delete order of specified type. Returns if order was found.
  function deleteOrder(uint64 id) {
    Order o = orders[id];
    if (o.price == 0)
      throw; // nothing there
    // Pending order?
    if (o.orderSeqNo >= pendingStart) {
      // delete pending order/set value to zero
      delete pendingOrders[o.orderSeqNo];
      // then remove from storage
      delete orders[id];
      return;
    }
    // Live order
    byte orderType = id & 0x8000000000000000 == 0 ? BID : ASK;
    Pile pile = orderType == BID ? book.bids : book.asks;
    pile.size--;
    if (o.above == NIL)
      pile.top = o.below;
    else
      orders[o.above].below = o.below;
    if (o.below != NIL)
      orders[o.below].above = o.above;
    // remove from storage
    delete orders[id];
  }
  
  // // Add a MTD bid order that will spend specified amount. Actual amount bought
  // // to be determined by market limit orders and DAO. The assets immediately
  // // returned are minimally spend/daoIndicativeAskPrice * mtdAdvance/100.
  // // Returns the id of the order, the amount of asset immediately received,
  // // and whether the order is completed or sitting in the MTD pool.
  // function addMtdBid(
  //   address account,
  //   uint128 spend
  // )
  //   returns (uint64, uint128, bool)
  // {
  //   // Transfer payment from sender to exchange
  //   currencyLedger.transferToDAO(spend, account);
  //   // Create id for notional order
  //   uint64 id = nextOrderId(BID, MARKET_TO_DAO);
  //   // Execute MTD bid
  //   var (received, finished) = execMtdBid(account, id, spend);

  //   return (id, received, finished);
  // }

  // function execMtdBid(
  //   address account,
  //   uint64 id,
  //   uint128 spend
  // )
  //   private
  //   returns (uint128, bool)
  // {
  //   // Take everything we can from limit orders to DAO's indicative price
  //   var (bought, unbought) = buyLimitOrdersToPrice(spend, daoIndicativeAskPrice);
  //   // Advance user assets from change
  //   if (unbought != 0) {
  //     // The mtdAdvance coefficient is adjusted by the DAO to account for
  //     // volatility, and determines how much of an asset it will advance. For
  //     // example, if mtdAdvance is 75% this implies the DAO thinks the risk of
  //     // the asset rising in value by 25% in ~20s is sufficiently low that it
  //     // can advance 75% of the amount of assets it expects to eventually
  //     // provide.
  //     uint128 advanced = (unbought*mtdAdvance)/(daoIndicativeAskPrice*100);
  //     // Add the MTD order to the pool so that user receives the additional
  //     // assets/difference once the DAO has established its best price.
  //     bidToMtdPool(account, id, unbought, advanced); // DAO processes pool periodically
  //   }
  //   // Immediately transfer to the user those assets that have been bought on
  //   // the exchange and advanced by the DAO.
  //   assetLedger.transferFromDAO(account, bought+advanced);

  //   return (bought+advanced, unbought == 0);
  // }

  // // Add a MTD ask order that will sell the specified amount. The actual amount
  // // received to be detemined by market limit orders and DAO. The payment
  // // immediately received is minimally
  // // amount*daoIndicativeAskPrice * mtdAdvance/100. Returns the id of the
  // // order, the payment immediately received and whether the order has
  // // completed or is sitting in MTD pool.
  // function addMtdAsk(
  //   address account,
  //   uint128 amount
  // )
  //   returns (uint64, uint128, bool)
  // {
  //   // Transfer amount being sold by sender
  //   assetLedger.transferToDAO(amount, account);
  //   // Create id for notional order
  //   uint64 id = nextOrderId(ASK, MARKET_TO_DAO);
  //   // Execute MTD ask
  //   var (received, finished) = execMtdAsk(account, id, amount);

  //   return (id, received, finished);
  // }

  // function execMtdAsk(
  //   address account,
  //   uint64 id,
  //   uint128 amount
  // )
  //   private
  //   returns (uint128, bool)
  // {
  //   // Take everything we can from limit orders to DAO's indicative price
  //   var (received, unsold) = sellLimitOrdersToPrice(amount, daoIndicativeBidPrice);
  //   // Advance user assets from change
  //   if (unsold != 0) {
  //     // The mtdAdvance coefficient is adjusted by the DAO to account for
  //     // volatility, and determines how much of an asset it will advance. For
  //     // example, if mtdAdvance is 75% this implies the DAO thinks the risk of
  //     // the asset rising in value by 25% in ~20s is sufficiently low that it
  //     // can advance 75% of the amount of assets it expects to eventually
  //     // provide.
  //     uint128 advanced = (unsold*daoIndicativeBidPrice*mtdAdvance)/100;
  //     // Add the unsold asset to the MTD Pool. The DAO periodically processes
  //     // the pool, redeeming the assets at a price determined by what the
  //     // liquidity providers will close their promissory coins. When the DAO
  //     // redeems the asset, it will send the currency minus the advance already
  //     // sent.
  //     askToMtdPool(account, id, unsold, advanced);
  //   }
  //   // Immediately transfer to the user the currency that has been received from
  //   // the exchange and advanced by the DAO.
  //   currencyLedger.transferFromDAO(account, received+advanced);

  //   return (received+advanced, unsold == 0);
  // }

  // private
  // Add new order to pile. Search for insertion point by moving downwards from
  // the top of the pile. We insert new orders below existing orders with the
  // same price, thus providing for FIFO matching where appropriate.
  function insertOrder(
    Pile storage pile,
    Order storage o,
    uint64 id
  )
    private
    returns (bool)
  {
    // Find insertion point...
    if (pile.size > 0) {
      // search...
      uint64 prev = NIL;
      uint64 curr = pile.top;
      while (true) {
        Order co = orders[curr];
        bool insert = false;
        // bid?
        if (pile.orderType == BID) {
          if (o.price > co.price)
            insert = true;
        }
        // ask
        else { // if (pile.type == ASK) {
          if (o.price < co.price)
            insert = true;
        }
        // insert here?
        if (insert) {
          o.above = prev;
          o.below = curr;
          break;
        } 
        // move down
        prev = curr;
        curr = co.below;
        // bottom?
        if (curr == NIL) {
          o.above = prev;
          o.below = NIL;
          break;
        }
      }
    }
    // Insert into place
    if (o.above == NIL)
      pile.top = id;
    else
      orders[o.above].below = id;
    if (o.below != NIL)
      orders[o.below].above = id;
    pile.size++;
    // Return whether is top of file
    return o.above == NIL;
  }

  function buyLimitOrdersToPrice(uint128 spend, uint128 price) returns (uint128, uint128) {
      // TODO
      return (0, spend);
  }

  function sellLimitOrdersToPrice(uint128 amount, uint128 price) returns (uint128, uint128) {
      // TODO
      return (0, amount);
  }

  function askToMtdPool(address sender, uint64 id, uint128 unsold, uint128 advanced) {
      return;
  }

  function bidToMtdPool(address sender, uint64 id, uint128 unbought, uint128 advanced) {
      // TODO
      return;
  }
}

