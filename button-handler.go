{{/* Trigger type: Message Component */}}
{{/* Component ID: `lottery_join` */}}

{{$targetHour := 17}} {{/* Giveaway finish time (UTC) */}}

{{/* Safe interaction response that avoids crash after 3 seconds of timeout */}}
{{$responded := true}}
{{try}}
	{{/* Respond to Discord immediately to save the token from expiration */}}
	{{sendResponse nil (complexMessage "content" "Processing your entry... ⚙️" "ephemeral" true)}}
{{catch}}
	{{$responded = false}}
{{end}}

{{/* Now we can leisurely perform all necessary database actions */}}
{{$stats := dbGet .User.ID "user_stats"}}
{{$data := dict "joined" 0 "wins" 0 "current_streak" 0 "max_streak" 0 "max_streak_date" "—"}}
{{if $stats}}{{$data = sdict $stats.Value}}{{end}}

{{/* If a user clicks the button but is not the current winner, and their profile still has 
the date of their last win listed, then their winning streak has been broken. */}}
{{$currentWinner := dbGet 0 "daily_lottery_winner"}}
{{$isCurrentWinner := false}}
{{if $currentWinner}}
	{{if eq (print $currentWinner.Value) (print .User.ID)}}{{$isCurrentWinner = true}}{{end}}
{{end}}

{{if and (not $isCurrentWinner) (ne (index $data "last_win_date" | print) "")}}
	{{$data.Set "current_streak" 0}}
	{{$data.Set "last_win_date" ""}}
{{end}}

{{$participants := dbGet 0 "daily_lottery_participants"}}
{{$list := cslice}}
{{if $participants}}
	{{$list = $participants.Value}}
{{end}}

{{/* Checking user's participation */}}
{{if in $list .User.ID}}
	{{if $responded}}
		{{editResponse nil nil "You're already entered in today's giveaway! Wait for the results. ⏳"}}
	{{end}}
{{else}}
	{{/* Updating the list */}}
	{{$list = $list.Append .User.ID}}
	{{dbSet 0 "daily_lottery_participants" $list}}
	
	{{/* Updating user's profile */}}
	{{$data.Set "joined" (add (index $data "joined") 1)}}
	{{dbSet .User.ID "user_stats" $data}}
	
	{{/* Calculating participants number */}}
	{{$count := len $list}}

	{{/* Calculating Unix time to the nearest target hour (17:00 UTC by default) */}}
	{{$now := currentTime}}
	{{$targetUnix := 0}}
	{{if lt $now.Hour $targetHour}}
		{{$targetUnix = add $now.Unix (sub 61200 (add (mult $now.Hour 3600) (add (mult $now.Minute 60) $now.Second)))}}
	{{else}}
		{{$targetUnix = add $now.Unix (sub 147600 (add (mult $now.Hour 3600) (add (mult $now.Minute 60) $now.Second)))}}
	{{end}}
	{{$timestamp := (print "<t:" $targetUnix ":t>")}}
	
	{{/* Updating the button message */}}
	{{editMessage nil .Message.ID (complexMessage
		"embed" (cembed
			"title" "🎉 Daily Role Giveaway 🎉"
			"description" (print 
				"Press on the button below to join today's role giveaway." "\n"
				"\n"
				"Every day (" $timestamp "), the bot picks one lucky winner at random!" "\n"
				"\n"
				"🔹 **Today's participants:** `" $count "`"
			)
			"color" 2866010
		)
		"buttons" (cslice
			(cbutton "custom_id" "lottery_join" "label" "Participate! 🎲" "style" "Success")
		)
	)}}
	
	{{if $responded}}
		{{editResponse nil nil "You've successfully entered today's giveaway! ✅"}}
	{{end}}
{{end}}
