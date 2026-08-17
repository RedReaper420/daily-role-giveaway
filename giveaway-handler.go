{{/* Trigger type: Crontab */}}
{{/* Crontab: `0 17 * * *` for 17:00 UTC */}}

{{/* Settings: set IDs for the role, the channel, and the start command */}}
{{$roleID := 799169211901149225}}
{{$channelID := 799164744229453824}}
{{$buttonCommandID := 13}}

{{/* Taking away the role of the previous winner */}}
{{$oldWinner := dbGet 0 "daily_lottery_winner"}}
{{if and $oldWinner (ne (print $oldWinner.Value) "0")}}
	{{$oldWinnerID := toInt64 $oldWinner.Value}}
	{{$void := targetHasRoleID $oldWinnerID $roleID}}
	{{takeRoleID $oldWinnerID $roleID}}
{{end}}

{{/* Requesting participants list */}}
{{$participants := dbGet 0 "daily_lottery_participants"}}

{{$hasParticipants := false}}
{{if $participants}}
	{{if gt (len $participants.Value) 0}}
		{{$hasParticipants = true}}
	{{end}}
{{end}}

{{if not $hasParticipants}}
	{{/* !!! No participants logic !!! */}}
	{{sendMessage $channelID (complexMessage
		"embed" (cembed
			"title" "It appears that noone participated in giveaway."
			"description" "Noone gets the role today! 😢"
			"color" 16368640
		)
	)}}
	{{dbSet 0 "daily_lottery_winner" "0"}}
{{else}}
	{{/* !!! Present participants logic !!! */}}
	{{$list := $participants.Value}}
	{{$winnerID := index $list (randInt (len $list))}}
	
	{{/* Giving the role to the new winner */}}
	{{giveRoleID (toInt64 $winnerID) $roleID}}
	{{dbSet 0 "daily_lottery_winner" (print $winnerID)}}
	
	{{$todayStr := currentTime.Format "02.01.2006"}}
	{{$newStreak := 1}}

	{{/* Updating the winner's stats */}}
	{{$winnerStats := dbGet (toInt64 $winnerID) "user_stats"}}
	{{$wData := dict "joined" 1 "wins" 0 "current_streak" 0 "max_streak" 0 "max_streak_date" "—" "last_win_date" ""}}
	{{if $winnerStats}}{{$wData = sdict $winnerStats.Value}}{{end}}
	
	{{$currentStreak := index $wData "current_streak" | toInt}}
	{{if eq (index $wData "last_win_date" | print) ""}}
		{{$newStreak = 1}}
	{{else}}
		{{$newStreak = add $currentStreak 1}}
	{{end}}

	{{$newWins := add (index $wData "wins") 1}}
	{{$maxStreak := index $wData "max_streak" | toInt}}
	{{$maxStreakDate := index $wData "max_streak_date"}}
	
	{{if gt $newStreak $maxStreak}}
		{{$maxStreak = $newStreak}}
		{{$maxStreakDate = $todayStr}}
	{{end}}
	
	{{$wData.Set "wins" $newWins}}
	{{$wData.Set "current_streak" $newStreak}}
	{{$wData.Set "max_streak" $maxStreak}}
	{{$wData.Set "max_streak_date" $maxStreakDate}}
	{{$wData.Set "last_win_date" $todayStr}}
	
	{{dbSet (toInt64 $winnerID) "user_stats" $wData}}

	{{/* Reseting the previous other winner's streak */}}
	{{if and $oldWinner (ne (print $oldWinner.Value) "0") (ne (print $oldWinner.Value) (print $winnerID))}}
		{{$oldStats := dbGet (toInt64 $oldWinner.Value) "user_stats"}}
		{{if $oldStats}}
			{{$oData := sdict $oldStats.Value}}
			{{$oData.Set "current_streak" 0}}
			{{$oData.Set "last_win_date" ""}}
			{{dbSet (toInt64 $oldWinner.Value) "user_stats" $oData}}
		{{end}}
	{{end}}
	
	{{/* Announcing results  */}}
	{{$winnerUser := userArg $winnerID}}
	{{sendMessage $channelID (print 
		"· · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·" "\n"
		"🎉 **The daily giveaway is over!** 🎉" "\n"
		"· · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·" "\n"
		"Today, luck was on __" $winnerUser.Mention "__! The winner gets a role for the next 24 hours!" "\n"
		"\n"
		"Winner's winstreak: **" $newStreak "**. The winner was selected from **" (len $list) "** participant(s)." "\n"
		"· · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · · ·"
	)}}
	
	{{/* Reseting the participants list for the new giveaway */}}
	{{dbDel 0 "daily_lottery_participants"}}
{{end}}

{{/* Calling the participation button command (starting the new giveaway) */}}
{{execCC $buttonCommandID nil 3 .ExecData}}
