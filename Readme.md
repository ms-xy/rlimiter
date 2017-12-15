## rlimiter

**rlimiter** is a helper program to limit the resource consumption of a program
on Linux based OS. It uses the `setrlimit` syscall to enforce the requested
resource limits.

1. Set hard (and soft) limit for each given option
2. Load the program into the current process, replacing rlimiter (using `exec`)

## Usage

```shell
go get github.com/ms-xy/rlimiter
rlimiter [<options>]* -executable <executable> -- <parameters>*
```

## Options

```shell
rlimiter -h
```

For detailed information on options please refer to
```shell
man setrlimit
```

## License

GNU GPLv3. Please see the attached License.txt file for details.
Different license terms can be arranged on request.
