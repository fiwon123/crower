<div align="center">

# 🍳 Recipes

**Copy, paste, adjust. Each block is a working sequence.**

</div>

## 🚀 Daily development loop

```bash
crower create pull "git pull --rebase --autostash"
crower create push -a p "git push origin HEAD"
crower create tidy "go mod tidy && go build ./..."
crower create ship -a s "crower execute tidy && crower execute push"

cr s
```

Nesting works: a stored command may call another stored command, and recursion is
refused rather than looping.

## 🐳 Docker, without the typing

```bash
crower create dps "docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'"
crower create dlog "docker logs -f {1}"
crower create dsh "docker exec -it {1} sh"

cr dps
cr dlog api
```

## 📦 Backup in one word

```bash
crower create backup "tar -czf ~/backups/$(date +%F).tar.gz ~/projects"
cr backup
crower list -i      # it is in the history, with the size and time
```

## 📁 File chores

```bash
crower create scratch --folder ~/scratch          # folder
crower create notes --file ~/notes.md             # empty file
crower copy ~/Downloads/*.zip ~/archive/          # several sources, one target
crower move ~/archive/old.zip ~/trash/            # move instead of copy
crower extract ~/Downloads/tool.tar.gz -o ~/opt/  # unpack somewhere specific
crower search -f "*.log" ~/projects               # find by name
```

Deleting is reversible:

```bash
crower delete ~/trash/old.zip
crower revert 1                 # undo that delete
```

## 🔍 Web search from the shell

```bash
crower search -b "golang context cancellation"
```

## 🧾 Notes that travel with the command

```bash
crower update dlog -e "docker logs -f {1}"
crower open                      # opens the configuration file
```

> [!TIP]
> A note is the cheapest documentation there is: the next person to run
> `crower list -i` sees what the command was for, next to when it ran.

## 👥 Onboard the team

```bash
crower export -o team.toml                # everything
crower export dps -o just-dps.toml        # one command
crower import team.toml --overwrite       # on the other machine
```
