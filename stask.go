package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/shlex"
	"github.com/itsFrank/stask/internal/staskfile"
	"github.com/itsFrank/stask/internal/template"
)

const mainHelptext = `stask - little stateful task runner

usage: stask [command] <args> options

state and tasks are stored in your staskfile, use the staskfile command to get its path

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

other topics:
    syntax      how to author stask tasks`

const helpHelptext = `stask help - print help for stask command or topic

    usage: stask help <topic>`

const helpInit = `stask init - create a default staskfile at the staskfile path

    note: will fail if file already exists

    usage: stask init <topic>`

const stateHelptext = `stask state - print current stored state

    usage: stask state`

const setHelptext = `stask set -  set a stored state value

    usage: stask set <name> <value>`

const clearHelptext = `stask clear - remove a stored state value

    usage: stask clear <name>`

const runHelptext = `stask run - run a task using stored state

    usage: stask run <task>`

const dryrunHelptext = `stask dryrun - print the command that would be executed with "stask run"

    usage: stask dryrun <task>`

const tasksHelptext = `stask tasks - print list of available tasks

    usage: stask tasks`

const staskfileHelptext = `stask staskfile - print the path to your staskfile

    you can override the default location with the "STASKFILE_PATH" environment variable

    usage: stask staskfile`

const syntaxHelptext = `stast task syntax:
    your staskfile has a "task" object, every field in that object is a runnable task

    task definitions look like this:
	    "my-task": "something {state} something else {other-state}"

    stored state can be used in task by wrapping the name in braces {}`

func main() {
	// var setCmd = flag.NewFlagSet("set", flag.ExitOnError)
	//
	// var clearCmd = flag.NewFlagSet("clear", flag.ExitOnError)
	//
	// var stateCmd = flag.NewFlagSet("state", flag.ExitOnError)
	//
	// var runCmd = flag.NewFlagSet("run", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Fprintln(flag.CommandLine.Output(), mainHelptext)
		os.Exit(1)
	}

	switch os.Args[1] {

	case "-h", "--h", "-help", "--help":
		fmt.Fprintln(flag.CommandLine.Output(), mainHelptext)

	case "help":
		if len(os.Args) < 3 {
			fmt.Fprintln(flag.CommandLine.Output(), helpHelptext)
			os.Exit(1)
		}
		doHelp(os.Args[2])

	case "init":
		doInit()

	case "state":
		doState()

	case "set":
		doSet(os.Args)

	case "clear":
		doClear(os.Args)

	case "run":
		doRun(os.Args)

	case "dryrun":
		doDryrun(os.Args)

	case "tasks":
		doTasks()
	case "staskfile":
		doStaskfile()

	default:
		fmt.Fprintf(flag.CommandLine.Output(), "error: unexpected command '%s'\n", os.Args[1])
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask --help\" for usage information")
		os.Exit(1)
	}
}

func doHelp(topic string) {
	switch topic {

	case "help":
		fmt.Fprintln(flag.CommandLine.Output(), helpHelptext)

	case "state":
		fmt.Fprintln(flag.CommandLine.Output(), stateHelptext)

	case "set":
		fmt.Fprintln(flag.CommandLine.Output(), setHelptext)

	case "clear":
		fmt.Fprintln(flag.CommandLine.Output(), clearHelptext)

	case "run":
		fmt.Fprintln(flag.CommandLine.Output(), runHelptext)

	case "dryrun":
		fmt.Fprintln(flag.CommandLine.Output(), dryrunHelptext)

	case "tasks":
		fmt.Fprintln(flag.CommandLine.Output(), tasksHelptext)

	case "staskfile":
		fmt.Fprintln(flag.CommandLine.Output(), staskfileHelptext)

	case "syntax":
		fmt.Fprintln(flag.CommandLine.Output(), syntaxHelptext)
	}
}

func doInit() {
	var staskfilePath = getStaskfilePath()
	var staskfileDir = filepath.Dir(staskfilePath)
	var err = os.MkdirAll(staskfileDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(staskfilePath); errors.Is(err, os.ErrNotExist) {
		staskfile.WriteStaskfile(staskfilePath, staskfile.Staskfile{Tasks: make(map[string]string), State: make(map[string]string)})
		fmt.Fprintln(flag.CommandLine.Output(), "success - wrote default staskfile at path:")
		fmt.Fprintln(flag.CommandLine.Output(), "    ", staskfilePath)
	} else {
		fmt.Fprintln(flag.CommandLine.Output(), "error: staskfile already exists at path:")
		fmt.Fprintln(flag.CommandLine.Output(), "    ", staskfilePath)
		os.Exit(1)
	}
}

func doState() {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	if len(sf.State) == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "no state stored in staskfile")
		return
	}

	fmt.Fprintln(flag.CommandLine.Output(), "stask state:")
	for key, value := range sf.State {
		fmt.Fprintln(flag.CommandLine.Output(), "    ", key, ":", value)
	}
}

func doSet(args []string) {
	if len(args) != 4 {
		fmt.Fprintln(flag.CommandLine.Output(), "error: unexpected number of arguments")
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask help set\" for usage information")
		os.Exit(1)
	}
	var key = args[2]
	var value = args[3]

	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	sf.State[key] = value
	err = staskfile.WriteStaskfile(getStaskfilePath(), sf)
	if err != nil {
		panic(err)
	}
}

func doClear(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(flag.CommandLine.Output(), "error: unexpected number of arguments")
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask help clear\" for usage information")
		os.Exit(1)
	}
	var key = args[2]

	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	delete(sf.State, key)
	err = staskfile.WriteStaskfile(getStaskfilePath(), sf)
	if err != nil {
		panic(err)
	}
}

func doRun(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(flag.CommandLine.Output(), "error: unexpected number of arguments")
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask help run\" for usage information")
		os.Exit(1)
	}
	var key = args[2]

	execCommand(getFormattedTask(key))
}

func doDryrun(args []string) {
	if len(args) != 3 {
		fmt.Fprintln(flag.CommandLine.Output(), "error: unexpected number of arguments")
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask help run\" for usage information")
		os.Exit(1)
	}
	var key = args[2]

	fmt.Println(getFormattedTask(key))
}

func doTasks() {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	if len(sf.Tasks) == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "no tasks found in staskfile")
		return
	}

	fmt.Fprintln(flag.CommandLine.Output(), "stask tasks:")
	for key := range sf.Tasks {
		fmt.Fprintln(flag.CommandLine.Output(), "    ", key)
	}
}

func doStaskfile() {
	fmt.Fprintln(flag.CommandLine.Output(), getStaskfilePath())
}

func getStaskfilePath() string {
	path, _ := os.LookupEnv("STASKFILE_PATH")
	if len(path) > 0 {
		return path
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	return filepath.Join(homedir, ".config", "stask", "staskfile.json")
}

func readStaskfile() staskfile.Staskfile {
	return staskfile.Staskfile{}
}

func getFormattedTask(key string) string {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	task, prs := sf.Tasks[key]
	if !prs {
		fmt.Fprintln(flag.CommandLine.Output(), "error: task '", key, "' was not found in staskfile")
		os.Exit(1)
	}

	tmpl, err := template.ParseTemplate(task)
	if err != nil {
		panic(err)
	}

	str, missing := template.ApplyTemplate(tmpl, sf.State)
	if len(missing) > 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "error: task keys not found in state: ", missing)
		os.Exit(1)
	}

	return str
}

func execCommand(command string) {
	args, err := shlex.Split(command)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
