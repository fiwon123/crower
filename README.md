# crower

&nbsp;&nbsp;&nbsp;&nbsp;A dev tool that manages system commands by executing commands via custom aliases and
managing it with useful operations like add, edit, remove, list and more.

Latest Releases: [Link](https://github.com/fiwon123/crower/releases)

More information on Wiki page: [Link](https://github.com/fiwon123/crower/wiki)

## OS Support

- Windows (cmd default to execute commands)
- Linux (sh default to execute commands)

## About

- usage `go run ./ --help` to show availables parameters
- paramereters: add, update, delete, reset, list...

## Example 

- `go run ./ --add -n "command1" -a "com1,c" -e "echo hello"`

now you can use name or alias:
- `go run ./ command1 `
- `go run ./ com1`
- `go run ./ c`
