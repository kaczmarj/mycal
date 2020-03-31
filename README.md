# mycal

This is a command-line program that modifies mandatory events in an iCalendar file. It is meant to be used __ONLY__ with Stony Brook University School of Medicine calendars.

It is also my way of practicing using Go, so the code might not be the shiniest.

See the Releases tab for a program compiled for your distribution.

## Examples

The following example adds the prefix "[MANDATORY] " to the title of every mandatory event.

```bash
./mycal -in calendar.ics -out modified.ics -pre "[MANDATORY] "
```

The following example adds the prefix "[MANDATORY] " to the title of every mandatory event and adds a reminder 30 minutes before each mandatory event.

```bash
./mycal -in calendar.ics -out modified.ics -pre "[MANDATORY] " -remind -min 30
```
