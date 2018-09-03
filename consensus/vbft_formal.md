
BlockState = { Proposed, Endorserd, Committerd, Candidate, Sealed }
MessageType = { Proposal, Endorse, Commit, Candidate }

n = the number of members in the population

B = the set of all blocks


block b = {h, t, p, []tx, sig_p }

  h = hash(b) : a list of hashes of the block's parents, self-parent first
  t = time(b) : proposer's claimed date and time of the block's creation
  p = proposer(b) : proposer's ID number
  []tx = 
  sig_p = sig(b) : proposer's digital signature of { h, t, p, []tx }



node state s = { B, M }
  B = { b }
  M = { M_b }
  M_b = { ProposalMsgs_b, EndorseMsgs_b, CommitMsgs_b, SealMsgs_b }
  ProposalMsgs_b = { proposal_{b,p} }
  EndorseMsgs_b = { endorse_{b,e} }
  CommitMsgs_b = { commit_{b,c} }
  SealMsgs_b = { Seal_{b,c} }
  ProposalMsg proposal_{b,p} = { b, i, p, sig_p }
  EndorseMsg endorose_{b,e} = { b, e, sig_e }
  CommitMsg commit_{b,c} = { b, c, sig_c }
  SealMsg Seal_{b,v} = { b, v, sig_v }


### Global Variables

P = the set of proposers

E = the set of endorsers

C = the set of committers

V = the set of validators


definition proposer_priority :: "int -> int -> int"        // (vrf, proposer_id) to priority



### Procedure

1. proposer selects proposal number n and send a proposal request to V

2. endorser receives proposal, 
   **determine if the proposal should be endorsed,**
   if yes: send endorse message

3. committer receives proposal/endorsement,
   **determine if the proposal is ready to commit,**
   if yes: send commit message

4. validator receives proposal/endorsement/commitment,
   **determine if the proposal is ready to finalize,**
   if yes, send seal message



### Specification

datatype 

fun blockhash :: "block -> hash"
                 where
                 "

fun prevhash :: "block -> hash"
                where
                "

fun prevblock :: "block -> block"
                  where
                  "


fun block_proposer :: "block -> int"
                      where


fun block_validators :: "block -> set


















