package main

/*
  rlimiter - helper program to limit program resource consumption
  Copyright (C) 2017  Maximilian Schott (github.com/ms-xy)

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

/*
#include <sys/resource.h>
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var license_notification = `
Copyright (C) 2017  Maximilian Schott (github.com/ms-xy)
This program comes with ABSOLUTELY NO WARRANTY.
This is free software, and you are welcome to redistribute it under certain conditions.
License: GNU GPLv3
`

func main() {

	arguments := []string{}
	for i, arg := range os.Args {
		if arg == "--" {
			if len(os.Args) > i+1 {
				arguments = os.Args[i+1:]
			}
			os.Args = os.Args[:i]
			break
		}
	}

	flags := []T{
		T{"as", int(C.RLIMIT_AS), Flag{}, Flag{}},
		T{"core", int(C.RLIMIT_CORE), Flag{}, Flag{}},
		T{"cpu", int(C.RLIMIT_CPU), Flag{}, Flag{}},
		T{"data", int(C.RLIMIT_DATA), Flag{}, Flag{}},
		T{"fsize", int(C.RLIMIT_FSIZE), Flag{}, Flag{}},
		T{"locks", int(C.RLIMIT_LOCKS), Flag{}, Flag{}},
		T{"memlock", int(C.RLIMIT_MEMLOCK), Flag{}, Flag{}},
		T{"msgqueue", int(C.RLIMIT_MSGQUEUE), Flag{}, Flag{}},
		T{"nice", int(C.RLIMIT_NICE), Flag{}, Flag{}},
		T{"nofile", int(C.RLIMIT_NOFILE), Flag{}, Flag{}},
		T{"nproc", int(C.RLIMIT_NPROC), Flag{}, Flag{}},
		T{"rss", int(C.RLIMIT_RSS), Flag{}, Flag{}},
		T{"rtprio", int(C.RLIMIT_RTPRIO), Flag{}, Flag{}},
		T{"rttime", int(C.RLIMIT_RTTIME), Flag{}, Flag{}},
		T{"sigpending", int(C.RLIMIT_SIGPENDING), Flag{}, Flag{}},
		T{"stack", int(C.RLIMIT_STACK), Flag{}, Flag{}},
	}

	for _, t := range flags {
		flag.Var(&t.Soft, "S"+t.Name,
			"Set soft limit for RLIMIT_"+strings.ToUpper(t.Name))
		flag.Var(&t.Hard, "H"+t.Name,
			"Set hard limit for RLIMIT_"+strings.ToUpper(t.Name))
	}

	executable := flag.String("executable", "",
		"The executable to be invoked after setting the resource limits")
	workingDir := flag.String("wdir", ".",
		"The working directory")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", license_notification)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()

	for _, t := range flags {
		if t.Hard.Isset || t.Soft.Isset {
			resource := t.Resource
			rlimit := syscall.Rlimit{}

			if err := syscall.Getrlimit(resource, &rlimit); err != nil {
				panic(err)
			}

			if t.Hard.Isset && t.Hard.Value < rlimit.Max {
				rlimit.Max = t.Hard.Value
				if !t.Soft.Isset {
					t.Soft.Isset = true
					t.Soft.Value = t.Hard.Value
				}
			}
			if t.Soft.Isset {
				if t.Soft.Value <= rlimit.Max {
					rlimit.Cur = t.Soft.Value
				} else {
					rlimit.Cur = rlimit.Max
				}
			}

			if err := syscall.Setrlimit(resource, &rlimit); err != nil {
				panic(err)
			}
		}
	}

	if err := os.Chdir(*workingDir); err != nil {
		panic(err)
	}

	err := syscall.Exec(
		*executable,
		append([]string{*executable}, arguments...),
		os.Environ(),
	)
	panic(err)
}

type T struct {
	Name     string
	Resource int
	Soft     Flag
	Hard     Flag
}

type Flag struct {
	Isset bool
	Value uint64
}

func (f *Flag) Set(x string) error {
	if val, err := strconv.ParseUint(x, 10, 64); err == nil {
		f.Value = val
		f.Isset = true
		return nil
	} else {
		return err
	}
}

func (f *Flag) String() string {
	return strconv.FormatUint(f.Value, 10)
}
