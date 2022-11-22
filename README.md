To run, type `go run main.go`. You will need Go installed.

The goal here is to simplify round tracking, effects like afflictions or spells,
and enemy HP/afflictions.

The critical command is `step`, which steps forward 1 round by default. But you
can also specify a duration with `step <duration>`. For duration, `10r` is 10
rounds, but other units include `s` for seconds (e.g., `100s`), `m` for minutes, `d` for days,
`y` for years. It will automatically convert the duration to rounds and step
forward that number of rounds.

Critically, it will automatically notify you when effects end. To add an effect,
type `in <duration> <effect>`, where `<effect>` is something like "Slick is no longer
poisoned". If you step past that time, it will notify you that effect has ended
by printing the effect to the screen.

`normalize` subtracts the current round number from all effects, essentially
bringing you back to round "1", with the relative effect round numbers remaining
the same.

`add <name> <hp>` adds an entity with the given name and HP. You can then
afflict entities with effects with `afflict <name> <duration> <effect>`.
Stepping may end entity effects too.

You can delete entities with `del <name>`, and damage entities with `dmg <name>
<amount>` where amount is a number.
