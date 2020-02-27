
# MBP

MBP as one blockchain-based unified account system.

## Components

![](./MBP.png)

## Stories

* User Register with KYC (KYC supports blockchain)
* User Register with KYC (KYC not supports blockchain)
* User Login with 3rd-party App (OAuth client)
* MOP User Login through MBP (OAuth server)

### User Register with KYC (KYC supports blockchain)

1. User create wallet in App
2. User create DID with wallet in App
3. User sign `DID-register` transaction (tx)
4. App send `register` request (tx included) to doorman
5. doorman send `register` request to MBP
6. MBP process `register` request
7. MBP send `DID-register` transaction to blockchain
8. User request KYC from KYC provider
9. KYC provider send claim to User
10. User request KYC provider as `DID recovery`
11. KYC provider send `take-DID-recovery` transaction to blockchain
12. User send `KYC-proof` to doorman
13. doorman send `KYC-proof` to MBP
14. MBP verifies `KYC-proof`, as user as `KYC-verified`

##### question:
* App send KYC request?  KYC service are usually not open service, and not free

```plantuml
actor App
entity KYC
entity Doorman
database Blockchain
App->App: create wallet
App->App: create DID, sign DID-register tx
App->Doorman: user register
Doorman->Doorman: process register request
Doorman->Blockchain: DID-register
Blockchain-->Doorman: transaction notification
Doorman-->App: register done
App->KYC: request KYC
KYC->App: claim and proof
App->App: store claim
App->KYC: request `DID recovery`
KYC->Blockchain: send `take-DID-recovery` tx
Blockchain-->KYC: transaction notification
App->Doorman: update user KYC
Doorman->Doorman: verify `KYC-proof`, update user status
Doorman-->App: verified
```

### User Register with KYC (KYC not supports blockchain)

1. User create wallet in App
2. User create DID with wallet in App
3. User Sign `DID-register` transaction (tx)
4. App send `register` request (tx included) to doorman
5. doorman send `register` request to MBP
6. MBP process `register` request
7. MBP send `DID-register` transaction to blockchain
8. User request KYC from KYC provider
9. KYC provider send claim to User
10. User send `KYC-proof` to doorman
11. doorman send `KYC-proof` to MBP
12. MBP verifies `KYC-proof`, as user as `KYC-verified`
13. MBP sets KYC provider as `DID-recovery`


##### Question: How DID smart contract implement recovery?
  * recovery can not base on signature from user wallet
  * KYC provider should provide signatures in claim/proof (some KYC not support)
  * blockchain need verify siganture of `DID-recovery`


```plantuml
actor App
entity KYC
entity Doorman
database Blockchain
App->App: create wallet
App->App: create DID, sign DID-register tx
App->Doorman: user register
Doorman->Doorman: process register request
Doorman->Blockchain: DID-register
Blockchain-->Doorman: transaction notification
Doorman-->App: register done
App->KYC: request KYC
KYC->App: claim and proof
App->App: store claim
App->Doorman: update user KYC (with KYC-proof)
Doorman->Doorman: verify `KYC-proof`, update user status
Doorman-->Blockchain: set KYC-provider as `DID-recovery`
Doorman-->App: verified
```


### User Login with 3rd-party App (OAuth client)

1. User login MBP with 3rd-party App (callback url redirected to MBP OAuth)
3. 3rd-party OAuth server callback Doorman
4. Doorman request `access-token` from 3rd-path OAuth server
5. Doorman make user mapping
6. Doorman creates DID on blockchain
7. Doorman notifies App user login OK

```plantuml
actor App
entity 3rd_OAuth
entity Doorman
database Blockchain

App->3rd_OAuth: Login (MBP_OAuth as callback)
3rd_OAuth->Doorman: callback with login code
Doorman->3rd_OAuth: request access token with (login-code/client-id/...)
3rd_OAuth->Doorman: return access token
Doorman->Doorman: create wallet, request user mapping
Doorman->Blockchain: DID-register
Blockchain-->Doorman: transaction notification
Doorman-->App: login OK
```

### MOP User Login through MBP (OAuth Server)

1. User login MOP through MBP
2. User request OAuth from Doorman(server)
3. Doorman verifies user request, callback MOP(OAuth client)
4. MOP get login code, request access token
5. Doorman grant access token to MOP
6. MOP request user information with access token
7. MOP notify User login OK

```plantuml
actor MOP_App
entity MBP_App
entity MOP
entity Doorman

MOP_App->MOP: Open
MOP-->MOP_App: Login through MBP supported
MOP_App->MBP_App: Login
MBP_App->Doorman: Login code request (MOP_OAuth as callback)
Doorman->Blockchain: verify login request
Blockchain-->Doorman: Login OK
Doorman->MOP: callback with login code
MOP->Doorman: request access token
Doorman-->MOP: return access token
MOP->MOP: set user access token
MOP->Doorman: request user info
Doorman-->MOP: return user info
MOP-->MOP_App: login OK
```

