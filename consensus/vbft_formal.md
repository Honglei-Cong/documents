
BlockState = { Proposed, EndorserVoted, CommitterVoted, Candidate, Sealed }
MessageType = { Proposal, EndorseVote, CommitVote, CandidateVote }

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
  M_b = { ProposalMsgs_b, EndorseMsgs_b, CommitMsgs_b, VoteMsgs_b }
  ProposalMsgs_b = { proposal_{b,p} }
  EndorseMsgs_b = { endorse_{b,e} }
  CommitMsgs_b = { commit_{b,c} }
  VoteMsgs_b = { vote_{b,c} }
  ProposalMsg proposal_{b,p} = { b, i, p, sig_p }
  EndorseMsg endorose_{b,e} = { b, e, sig_e }
  CommitMsg commit_{b,c} = { b, c, sig_c }
  VoteMsg vote_{b,v} = { b, v, sig_v }












===
pos table
vrf

isProposer
isEndorser
isCommitter


isEquivocation


manyVoters

see(x, y)




