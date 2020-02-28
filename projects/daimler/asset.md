
title Asset creation and attestations for user asset


```plantuml
UserMobilityApp->UserMobilityApp: Create private and public key
UserMobilityApp->KYCProvider: Request verification of on-chain claim for user (pubkey, driver license + real-world proof of claims) note right of KYC Provider: Checks claim and proof
KYCProvider->MBP: Create user asset / register user DID for pubkey
MBP-->KYCProvider: user DID
KYCProvider->MBP: Set claim for user DID, topic driver license, signed by KYC DID note right of MBP: Stores claim on ledger
MBP-->KYCProvider:ok
KYCProvider-->UserMobilityApp:user DID
UserMobilityApp->KYCProvider: Request verification of off-chain claim for user DID (name + real-world proof of claims) note right of KYC Provider: KVC Provider checks claim and proof
KYCProvider-->UserMobilityApp:Verifiable credential for user DID containing name, signed by KYC DID note right of User Mobility App:Store verifiable credential in wallet
```