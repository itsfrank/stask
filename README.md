# stask

stateful task running utility

I built this because my company's build system has a number of commands with a
bunch of flags. The flags stay the same between commands depending on which
flavor of the build you are currently using. I was tired of having to edit
multiple commands in my shell history, so I made `stask`.

## Install

stask is not on any package manager for the moment. The best way to install it
is to use `go install`

Requirement: make sure go is installed (use your favourite package manager, or
the [official instructions](https://go.dev/doc/install))

Make sure `$GOPATH/bin` is on your `PATH` (note, use the
[default GOPATH](https://go.dev/doc/gopath_code#GOPATH) if it is not set:
`$HOME/go`)

Install stask: `go install github.com/itsFrank/stask`

## Usage

Everything you need should be in the helptext: `stask help`. If you rather
markdown, or aren't convinced yet, the basic workflow is outlined below.

First, make a your staskfile: `stask init`

Add some tasks to your staskfile with your favorite text editor:

```json
// staskfile.json (get it's path with `stask staskfile`)
{
  "Tasks": {
    "echo": "echo {message}",
    "cp": "cp {from} {to}"
  }
}
```

Set some state

```shell
> stask set message "hello from stask!"
> stask set from ~/myfile
> stask set to ~/myfile_copy
```

Run tasks!

```shell
> stask run echo
hello from stask

> stask dryrun cp #dryrun prints the command instead of executing it
cp ~/myfile ~/myfile_copy
```

## Commands

stask has a bunch of commands, here is the list from the help text

```text
commands:
    help        show help for a command or topic
    init        create a default staskfile
    state       print current stored state
    set         set a value to state
    clear       remove a value from state
    run         run a command with state
    dryrun      print command with state inserted
    tasks       show list of available tasks
    staskfile   print path to your staskfile
```

use `stask help <command>` for more information about any command
