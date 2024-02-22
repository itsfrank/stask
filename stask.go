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
    profile     list, show, load, save, delete profiles
    staskfile   print path to your staskfile

other topics:
    syntax      how to author stask tasks
    shell       how to configure what shel and shell flags stask will use`

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

const profileHelptext = `stask profile - list, load, save profiles

	usage: stask profile [subcommand] <args>

subcommands:
    list          - list available profiles
    show <name>   - print the state stored in the profile
    load <name>   - apply state stored in profile, overwiting only values in profile
    save <name>   - save current state as new profile
    delete <name> - save current state as new profile`

const staskfileHelptext = `stask staskfile - print the path to your staskfile

    you can override the default location with the "STASKFILE_PATH" environment variable

    usage: stask staskfile`

const syntaxHelptext = `stask task syntax:
    your staskfile has a "task" object, every field in that object is a runnable task

    task definitions look like this:
	    "my-task": "something {state} something else {other-state}"

    stored state can be used in task by wrapping the name in braces {}`

const shellHelptext = `stask shell config:
    stask requires that a shell be explicitely set via environment variables
    on unix systems this is typically already done by default via the "SHELL" environment variable
    stask will consider 2 different variable to determing the desired shell:
        STASK_SHELL   custom shell used only by stask
        SHELL         system default on unix systems (fallback is STASK_SHELL not set)

    stask passes flags to the shell to execute your tasks
    by default is passes "-ic" flags, you can customized the passed in flags with this variable:
        STASK_SHELL_FLAGS    custom flags passed to shell

	the complete command executed by "stask run" looks like this:
	    <$STASK_SHELL(or $SHELL)> <$STASK_SHELL_FLAGS> "<task>"`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(flag.CommandLine.Output(), mainHelptext)
		os.Exit(1)
	}

	switch os.Args[1] {

	case "-h", "--h", "-help", "--help":
		fmt.Fprintln(flag.CommandLine.Output(), mainHelptext)

	case "help":
		doHelp(os.Args)

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

	case "profile":
		doProfile(os.Args)

	case "staskfile":
		doStaskfile()

	default:
		fmt.Fprintf(flag.CommandLine.Output(), "error: unexpected command '%s'\n", os.Args[1])
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask --help\" for usage information")
		os.Exit(1)
	}
}

func doHelp(args []string) {
	if len(os.Args) < 3 {
		fmt.Fprintln(flag.CommandLine.Output(), helpHelptext)
		os.Exit(1)
	}

	topic := args[2]
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

	case "profile":
		fmt.Fprintln(flag.CommandLine.Output(), profileHelptext)

	case "staskfile":
		fmt.Fprintln(flag.CommandLine.Output(), staskfileHelptext)

	case "syntax":
		fmt.Fprintln(flag.CommandLine.Output(), syntaxHelptext)

	case "shell":
		fmt.Fprintln(flag.CommandLine.Output(), shellHelptext)

	default:
		fmt.Fprintf(flag.CommandLine.Output(), "error: unexpected help topic '%s'\n", topic)
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask --help\" for information")
		os.Exit(1)
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

	fmt.Fprintln(os.Stdout, "stask state:")
	for key, value := range sf.State {
		fmt.Fprintln(os.Stdout, "    ", key, ":", value)
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

	fmt.Fprintln(os.Stdout, "stask tasks:")
	for key := range sf.Tasks {
		fmt.Fprintln(os.Stdout, "    ", key)
	}
}

func doProfile(args []string) {
	exitProfileUsageError := func(message string) {
		fmt.Fprintf(flag.CommandLine.Output(), "error: %s\n", message)
		fmt.Fprintln(flag.CommandLine.Output(), "    use \"stask help profile\" for usage information")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		exitProfileUsageError("unexpected number of arguments")
	}

	subcommand := args[2]
	switch subcommand {

	case "list":
		doProfileList()

	case "show":
		if len(os.Args) < 4 {
			exitProfileUsageError("missing argument <name>")
		}
		doProfileShow(os.Args[3])

	case "load":
		if len(os.Args) < 4 {
			exitProfileUsageError("missing argument <name>")
		}
		doProfileLoad(os.Args[3])

	case "save":
		if len(os.Args) < 4 {
			exitProfileUsageError("missing argument <name>")
		}
		doProfileSave(os.Args[3])

	case "delete":
		if len(os.Args) < 4 {
			exitProfileUsageError("missing argument <name>")
		}
		doProfileDelete(os.Args[3])

	default:
		exitProfileUsageError(fmt.Sprintf("unexpected subcommand '%s'", subcommand))
	}
}

func doProfileList() {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	if len(sf.Profiles) == 0 {
		fmt.Fprintln(flag.CommandLine.Output(), "no profiles stored in staskfile")
		return
	}

	fmt.Fprintln(os.Stdout, "saved profiles:")
	for key := range sf.Profiles {
		fmt.Fprintln(os.Stdout, "    ", key)
	}
}

func doProfileShow(name string) {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	profile, found := sf.Profiles[name]
	if !found {
		fmt.Fprintf(flag.CommandLine.Output(), "no profile named '%s' in staskfile\n", name)
		return
	}

	fmt.Fprintf(os.Stdout, "%s - profile state:\n", name)
	for key, value := range profile {
		fmt.Fprintln(os.Stdout, "    ", key, ":", value)
	}
}

func doProfileLoad(name string) {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	profile, found := sf.Profiles[name]
	if !found {
		fmt.Fprintf(flag.CommandLine.Output(), "no profile named '%s' in staskfile\n", name)
		return
	}

	fmt.Fprintf(os.Stdout, "%s - applying profile...\n", name)
	for key, value := range profile {
		currentValue, found := sf.State[key]
		if !found {
			currentValue = "-"
		}

		sf.State[key] = value
		if value == currentValue {
			fmt.Fprintln(os.Stdout, "    ", key, ":", value, "[unchanged]")
		} else {
			fmt.Fprintln(os.Stdout, "    ", key, ":", currentValue, "->", value)
		}
	}

	err = staskfile.WriteStaskfile(getStaskfilePath(), sf)
	if err != nil {
		fmt.Println("error while writing staskfile, profile was not applied")
		panic(err)
	}

	fmt.Println("\nprofile applied sucessfully")
}

func doProfileSave(name string) {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	_, profileExists := sf.Profiles[name]

	newProfile := map[string]string{}
	for key, value := range sf.State {
		newProfile[key] = value
	}

	sf.Profiles[name] = newProfile

	err = staskfile.WriteStaskfile(getStaskfilePath(), sf)
	if err != nil {
		fmt.Println("error while writing staskfile, profile was not applied")
		panic(err)
	}

	if profileExists {
		fmt.Fprintf(os.Stdout, "profile '%s' overwritten sucessfully\n", name)
	} else {
		fmt.Fprintf(os.Stdout, "profile '%s' saved sucessfully\n", name)
	}
}

func doProfileDelete(name string) {
	sf, err := staskfile.ReadStaskfile(getStaskfilePath())
	if err != nil {
		panic(err)
	}

	_, profileExists := sf.Profiles[name]
	if !profileExists {
		fmt.Fprintf(flag.CommandLine.Output(), "no profile named '%s' in staskfile\n", name)
		return
	}

	delete(sf.Profiles, name)

	err = staskfile.WriteStaskfile(getStaskfilePath(), sf)
	if err != nil {
		fmt.Println("error while writing staskfile, profile was not applied")
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "profile '%s' deleted sucessfully\n", name)
}

func doStaskfile() {
	fmt.Fprintln(os.Stdout, getStaskfilePath())
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

func getShellConfig() (shellConfig struct {
	shell string
	flags string
}) {

	shell, _ := os.LookupEnv("STASK_SHELL")
	if len(shell) == 0 {
		envShell, prs := os.LookupEnv("SHELL")
		if len(envShell) == 0 || !prs {
			fmt.Fprintln(flag.CommandLine.Output(), "error: no shell set")
			fmt.Fprintln(flag.CommandLine.Output(), "    set either STASK_SHELL or SHELL env vars with path to shell STASK should use")
			os.Exit(1)
		}
		shell = envShell
	}

	shellFlags, _ := os.LookupEnv("STASK_SHELL_FLAGS")
	if len(shellFlags) == 0 {
		shellFlags = "-ic"
	}

	shellConfig.shell = shell
	shellConfig.flags = shellFlags
	return shellConfig
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
	shellConfig := getShellConfig()
	command = shellConfig.flags + " \"" + command + "\""
	args, err := shlex.Split(command)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(shellConfig.shell, args...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if exerr := err.(*exec.ExitError); exerr != nil {
			os.Exit(exerr.ExitCode())
		}
		fmt.Fprintln(flag.CommandLine.Output(), "stask error while running task:  ", err.Error())
		os.Exit(1)
	}
}
