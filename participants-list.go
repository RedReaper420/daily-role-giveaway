{{/* Trigger type: Command */}}
{{/* Trigger: `gwlist` */}}

{{$participants := dbGet 0 "daily_lottery_participants"}}

{{$hasParticipants := false}}
{{if $participants}}
	{{if gt (len $participants.Value) 0}}
		{{$hasParticipants = true}}
	{{end}}
{{end}}

{{if not $hasParticipants}}
	{{sendMessage nil "📋 **Noone has joined the today's giveaway yet!"}}
{{else}}
	{{$list := $participants.Value}}
	{{$count := len $list}}

	{{$displayList := ""}}

	{{range $index, $id := $list}}
		{{$user := userArg $id}}
		
		{{$userText := print "<@" $id ">"}}
		{{if $user}}
			{{$userText = print "**" $user.Mention "** (`" $user.Username "`)"}}
		{{end}}
		
		{{$num := add $index 1}}
		
		{{$displayList = print $displayList "\n" $num ". " $userText}}
	{{end}}

	{{sendMessage nil (complexMessage
		"embed" (cembed
			"title" "📋 Giveaway participants list 📋"
			"description" (print "List of users who are applying for the role's giveaway today:\n" $displayList)
			"color" 2866010
			"footer" (sdict "text" (print "Total participants: " $count))
			"timestamp" currentTime
		)
	)}}
{{end}}