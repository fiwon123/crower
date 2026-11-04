# 🩺 `doctor`

**Show the paths, permissions and environment a support conversation asks for.**

```bash
crower doctor
crower doctor --json
```

`doctor` reports the active setup, the configuration file and folder, whether
the folder is writable, how many commands and history entries exist, the
update-check mode and last check, the detected shell and editor, whether the
diagnostics stream is a terminal, and the platform.

```bash
crower doctor --json | jq '{setup, config_file, config_writable}'
```

It reads only: nothing is created beyond the config folder a normal run needs.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · [📖 Commands](How-to-Use)
