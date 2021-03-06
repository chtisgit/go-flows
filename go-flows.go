package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"sort"
	"strings"
	"text/tabwriter"
)

var (
	cpuprofile     = flag.String("cpuprofile", "", "Write cpu profile")
	heapprofile    = flag.String("heapprofile", "", "Write heap profile")
	memprofile     = flag.String("memprofile", "", "Write memory profile")
	memprofilerate = flag.Int("memprofilerate", 0, "Set MemProfileRate")
	blockprofile   = flag.String("blockprofile", "", "Write goroutine blocking profile")
	mutexprofile   = flag.String("mutexprofile", "", "Write mutex blocking profile")
	tracefile      = flag.String("trace", "", "Turn on tracing")
)

func flags() {
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
}

func cmdString(append string) {
	fmt.Fprintf(os.Stderr, "Usage:\n  %s [flags] %s\n", os.Args[0], append)
}

func usage() {
	cmdString("command [args]")
	fmt.Fprint(os.Stderr, `
A generic and fully customizable flow exporter written in go. Use one of
the provided subcommands.
`)
	fmt.Fprintf(os.Stderr, "\nAvailable Commands:\n")
	t := tabwriter.NewWriter(os.Stderr, 8, 4, 4, ' ', 0)
	for _, command := range commands {
		line := new(bytes.Buffer)
		fmt.Fprintf(line, "  %s\t%s\n", command.cmd, command.desc)
		t.Write(line.Bytes())
	}
	t.Flush()
	flags()

	os.Exit(-1)
}

type command struct {
	cmd  string
	desc string
	run  func(string, []string)
}

var commands []*command

func addCommand(cmd, desc string, run func(string, []string)) {
	commands = append(commands, &command{cmd, desc, run})
}

type plugins []string

func (p *plugins) String() string {
	return "plugins"
}

func (p *plugins) Set(a string) error {
	(*p) = append(*p, a)
	return nil
}

func main() {
	sort.Slice(commands, func(i, j int) bool {
		return strings.Compare(commands[i].cmd, commands[j].cmd) < 0
	})
	flag.CommandLine.Usage = usage
	var defs plugins
	flag.Var(&defs, "defs", "Load definitions from directory. Definition files must follow the name scheme go-flows*.defs.so.")
	flag.Parse()
	if *memprofilerate != 0 {
		runtime.MemProfileRate = *memprofilerate
	}
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			runtime.GC()
			err := pprof.Lookup("heap").WriteTo(f, 0)
			if err != nil {
				fmt.Println(err)
			}
			f.Close()
		}()
	}
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *blockprofile != "" {
		f, err := os.Create(*blockprofile)
		if err != nil {
			log.Fatal(err)
		}
		runtime.SetBlockProfileRate(1)
		defer pprof.Lookup("block").WriteTo(f, 0)
	}
	if *mutexprofile != "" {
		f, err := os.Create(*mutexprofile)
		if err != nil {
			log.Fatal(err)
		}
		runtime.SetMutexProfileFraction(1)
		defer pprof.Lookup("mutex").WriteTo(f, 0)
	}
	if *tracefile != "" {
		f, err := os.Create(*tracefile)
		if err != nil {
			log.Fatal(err)
		}
		trace.Start(f)
		defer trace.Stop()
	}

	if len(defs) != 0 {
		for _, def := range defs {
			if defs, err := filepath.Glob(filepath.Join(def, "go-flows*.defs.so")); err != nil {
				log.Fatalf("Couldn't find definitions: %s\n", err)
			} else {
				for _, def := range defs {
					if _, err := plugin.Open(def); err != nil {
						log.Fatalf("Couldn't load '%s' because of: %s", def, err)
					}
				}
			}
		}
	}

	for _, command := range commands {
		if flag.Arg(0) == command.cmd {
			command.run(command.cmd, flag.Args()[1:])
			return
		}
	}
	if flag.Arg(0) != "" {
		log.Fatalf("Unknown command '%s'", flag.Arg(0))
	}

	usage()
}
