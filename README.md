# crower

  A dev tool that manages system commands by executing commands via custom aliases and
managing it with useful operations like add, edit, remove, list and more.

## OS Support

- Linux (sh as default to execute commands)

## version v0.0.1

- usage go run ./ --help to show availables parameters
- paramereters: add, update, delete, reset, list

## Example 

- go run ./ --add -n "command1" -a "com1,c" -e "echo hello"

now you can use name or alias:
- go run ./ command1 
- go run ./ com1
- go run ./ c
