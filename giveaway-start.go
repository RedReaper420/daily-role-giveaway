{{/* Trigger type: Command */}}
{{/* Trigger: `gwstart` */}}
{{/* Role restrictions: only mods & admins */}}

{{$targetHour := 17}} {{/* Giveaway finish time (UTC) */}}

{{deleteTrigger 0}}

{{/* Making an embed message with a button */}}
{{/* Similar one is being made in the `button-handler` */}}

{{/* Participants counter */}}
{{$count := 0}}

{{/* Calculating Unix time to the nearest 17:00 UTC */}}
{{$now := currentTime}}
{{$targetUnix := 0}}
{{if lt $now.Hour $targetHour}}
    {{$targetUnix = add $now.Unix (sub 61200 (add (mult $now.Hour 3600) (add (mult $now.Minute 60) $now.Second)))}}
{{else}}
    {{$targetUnix = add $now.Unix (sub 147600 (add (mult $now.Hour 3600) (add (mult $now.Minute 60) $now.Second)))}}
{{end}}
{{$timestamp := (print "<t:" $targetUnix ":t>")}}

{{sendMessage nil (complexMessage
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
