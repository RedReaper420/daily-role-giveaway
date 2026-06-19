{{/* Trigger type: Command */}}
{{/* Trigger: `gwstats` */}}

{{$user := .User}}
{{$found := true}}

{{if .CmdArgs}}
	{{$target := userArg (index .CmdArgs 0)}}

	{{if $target}}
		{{$user = $target}}
	{{else}}
		{{$query := joinStr " " .CmdArgs}}
		{{$targetMulti := userArg $query}}
		{{if $targetMulti}}
			{{$user = $targetMulti}}
		{{else}}
			{{$found = false}}
		{{end}}
	{{end}}
{{end}}

{{if not $found}}
	{{sendMessage nil "❌ Failed to find user with such name. Try to set their **ID** or **mention** them via `@`."}}
{{else}}
	{{$stats := dbGet $user.ID "user_stats"}}
	{{if not $stats}}
		{{sendMessage nil (print "**" $user.Mention "** didn't participate in giveaways yet.")}}
	{{else}}
		{{$data := sdict $stats.Value}}
		
		{{/* Winstreak reset */}}
		{{$currentWinner := dbGet 0 "daily_lottery_winner"}}
		{{$isCurrentWinner := false}}
		{{if $currentWinner}}
			{{if eq (print $currentWinner.Value) (print $user.ID)}}{{$isCurrentWinner = true}}{{end}}
		{{end}}
		{{if and (not $isCurrentWinner) (ne (index $data "last_win_date" | print) "")}}
			{{$data.Set "current_streak" 0}}
			{{$data.Set "last_win_date" ""}}
			{{dbSet $user.ID "user_stats" $data}}
		{{end}}
		
		{{$joined := index $data "joined"}}
		{{$wins := index $data "wins"}}
		
		{{/* Calculating winrate percentage */}}
		{{$winRate := 0}}
		{{if gt $joined 0}}
			{{$winRate = roundFloor (mult (div (toFloat $wins) (toFloat $joined)) 100)}}
		{{end}}
		
		{{sendMessage nil (complexMessage
			"embed" (cembed
				"title" (print "📊 **Giveaway statistics of `" $user.Username "`**:")
				"description" (print 
					"👤 **User**: " $user.Mention "\n"
					"* **Total participations:** `" $joined "`" "\n"
					"* **Total wins:** `" $wins "` *(Winrate: " $winRate "%)*" "\n"
					"* **Current winstreak:** `" (index $data "current_streak") "`" "\n"
					"* **Max winstreak:** `" (index $data "max_streak") "` *(Achieved: " (index $data "max_streak_date") ")*"
				)
			)
		)}}
	{{end}}
{{end}}
