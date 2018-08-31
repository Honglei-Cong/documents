
datatype validator = Validator int

datatype estimate  = Estimate bool

datatype bet       = Bet "estimate * valdiator * bet list"



weight = "validator -> int"


definition tie_breaking  :: "validator set -> weight -> bool"
                             where
                                "tie_breaking V w = 
                                (\<forall> S0 S1,
                                        S0 \<subseteq> V -->
                                        S1 \<subseteq> V -->
                                        S0 != S1 -->
                                        sum w S0 != sum w S1
                                )"




fun sender    :: 'bet -> validator' where "sender (Bet (_, v, _)) = v"

fun est       :: 'bet -> estimate' where "est (Bet (e, _, _v)) = e"

fun justifies :: 'bet -> (bet -> bool)' where "justified b (Bet (e, v, js)) = (b \<in> set js)"




inductive tran :: '(bet -> bet -> bool) -> (bet -> bet -> bool)' 
                  where
                        tran_simple: "R x y --> tran R x y"
                      | tran_composit: "tran R x y --> R y z --> tran R x z"



definition is_dependency :: "bet -> bet -> bool" where "is_dependency = tran justifies"




definition equivocation :: "bet -> bet -> bool" 
                    where
                    "equivocation b0 b1 = 
                       (sender b0 = sender b1 
                        \<and> \<not> is_dependency b0 b1 
                        \<and> \<not> is_dependency b1 b0 
                        \<and> b0 \<noteq> b1 )"



definition is_view :: "bet set -> bool"
                        where
                        "is_view bs = (\<forall> b0 b1.  b0 \<in> bs --> b1 \<in> bs --> 
                                        not equivocation b0 b1)"



definition latest_bets :: "bet set -> validator -> bet set"
                          where
                          "latest_bets bs v = { l . l \<in> bs \<and> sender l = v \<and> (\<not> (\<exists> b' . b' \<in> bs \<and> sender b' = v \<and> is_dependency l b')) } "





definition is_latest_in :: "bet -> bet set -> bool"
                           where
                           "is_latest_in b bs = (b \<in> bs \<and> (\<not> (\<exists> b' . b' \<in> bs \<and> is_dependency b b'))) "


definition is_non_empty :: "a set -> bool"
                           where
                           "is_non_empty bs = (\<exists> b. b <in> bs)"





definition at_most_one :: "a set -> bool"
                          where
                          "at_most_one s = (\<forall> x y. x \<in> s --> y \<in> s --> x = y)"






definition observed_validators :: "bet set -> validator set"
                                  where
                                  "oberserved_validators bs = ({v :: validator . \<exists> b. b \<in> bs \<and> v = sender b})





definition has_a_latest_bet_on :: "bet set -> validator -> estimate -> bool"
                                  where
                                  "has_a_latest_bet_on bs v e = (\<exists> b. b <in> latest_bets bs v \<and> est b = e)"





definition weight_of_estimate :: "bet set -> weight -> estimate -> int"
                                 where
                                 "weight_of_estimate bs w e = sum w {v. has_a_latest_bet_on bs v e}"





definition is_max_weight_estimate :: "bet set -> weight -> estimate -> bool"
                                     where
                                     "is_max_weight_estimate bs w e = (\<forall> e'.
                                                                       weight_of_estimate bs w e \<ge> weight_of_estimate bs w e')



definition positive_weights :: "validator set -> weight -> bool"
                               where
                               "positive_weights vs w = (\<forall> v. v \<in> vs --> w v > 0)"



fun is_valid :: "weight -> bet -> bool"
                where
                "is_valid w (Bet (e, v, js)) = is_max_weight_estimate (set js) w e"





definition is_valid_view :: "weight -> bet set -> bool"
                            where
                            "is_valid_view w bs = (is_view bs \<and> (\<forall> b \<in> bs. is_valid w b))"








definition is_future_view :: "weight -> bet set -> bet set -> bool"
                             where
                             "is_future_view w b0, b1 = 
                             (b0 \<supseteq> b1 \<and> is_valid_view w b0 \<and> is_valid_view w b1)"





definition is_estimate_safe :: "weight -> bet set -> estimate -> bool"
                               where
                               "is_estimate_safe w bs e = 
                                        (\<forall> bf. is_future_view w bf bs --> is_max_weight_estimate bf w e)"






definition consistent_views :: "weight -> bet set -> bet set -> bool"
                                where
                                "consistent_views w b0, b1 = 
                                        (is_valid_view w b0 \<and> is_valid_view w b1 \<and> is_valid_view w (b0 \<union> b1))"
















