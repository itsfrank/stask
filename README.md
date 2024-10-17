# stask

stateful task running utility

I built this because my company's build system has a number of commands with a
bunch of flags. The flags stay the same between commands depending on which
flavor of the build you are currently using. I was tired of having to edit
multiple commands in my shell history, so I made `stask`.

**Update**: added `stask profile` subcommand! See below for more information

## Install

stask is not on any package manager for the moment. The best way to install it
is to use `go install`

Requirement: make sure go is installed (use your favourite package manager, or
the [official instructions](https://go.dev/doc/install))

Make sure `$GOPATH/bin` is on your `PATH` (note, use the
[default GOPATH](https://go.dev/doc/gopath_code#GOPATH) if it is not set:
`$HOME/go`)

Install stask: `go install github.com/itsfrank/stask`

## Usage

Everything you need should be in the helptext: `stask help`. If you rather
markdown, or aren't convinced yet, the basic workflow is outlined below.

First, make a your staskfile: `stask init`

Add some tasks to your staskfile with your favorite text editor (get its path
with `stask staskfile`):

```json
{
  "Tasks": {
    "echo": "echo {message}",
    "cp": "cp {from} {to}"
  },
  "State": {}
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

**new!** Forward args with `run` and `dryrun`!

```shell
> stask run echo -- more args here
hello from stask more args here
```

**new!** Save and load profiles!

```shell
> stask set name Frank
> stask profile save frank
profile 'frank' saved sucessfully

> stask set name Joe
> stask run greet
hello Joe!

> stask profile load frank
frank - applying profile...
    name : Joe -> Frank

> stask run greet
hello Frank!
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
    profile     list, show, load, save, delete profiles
    staskfile   print path to your staskfile
```

use `stask help <command>` for more information about any command

## Shell

stask requires that a shell be set via environment variables, on unix systems
this is typically set by default so no extra config should be necessary. However
if you wish to customize the shell used by stask (or troubleshoot errors) here
is how stask determines the shell

Shell env vars:

- `STASK_SHELL`: if set will use this shell
- `SHELL`: if `STASK_SHELL` is not set, will look for this var

stask also passes flags to the shell, by default it passes `-ic`. If those flags
don't work or you want to customize them, you can override them by setting
`STASK_SHELL_FLAGS`
