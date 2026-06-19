# Daily Role Giveaway

This is a system of Custom Commands for [YAGPDB.xyz Discord Bot](https://yagpdb.xyz/), implementing an automatic daily role giveaway system. 

## Features

* Giveaway sign up in one button click.
* Automatic assign and removal of a temporary giveaway role.
* Daily automatic giveaway reset at the set time (default 17:00 UTC).
* Stats tracking:
  * Participants list (`gwlist` command).
  * User statistics: participations, wins and winrate, current and maximal winstreaks (`gwstats` command). The command may take user mention or ID as an argument.

## Discord server setup

* Add a giveaway channel to your server.
* Add YAGPDB.xyz bot to your server.
* Import Custom Commands from this repository.
* Set up some variables in `giveaway-start`, `button-handler`, and `giveaway-handler`. Tweak the scripts further if needed.
* Enter `gwstart` command in the giveaway channel to start automatic giveaway.
