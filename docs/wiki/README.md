# Wiki source

These pages are the source of the GitHub wiki
(`https://github.com/fiwon123/crower/wiki`). They live in the repository so a
change to them is reviewed like any other, and so the wiki cannot drift from the
release it describes.

## Publishing

The wiki is a second git repository. GitHub creates it the first time a page is
made in the web UI; after that:

```bash
git clone https://github.com/fiwon123/crower.wiki.git
cp docs/wiki/*.md crower.wiki/
cd crower.wiki && git add -A && git commit -m "wiki: sync from docs/wiki" && git push
```

The file names are the page names: `Home.md` is the landing page, `_Sidebar.md`
is the navigation, `_Footer.md` is the strip at the bottom, and everything else
is a page (a command page's file name is the command, e.g. `create.md`).

## Pages

| Page | Covers |
| --- | --- |
| `Home.md` | Hero, quick start, where to go, what is not proven |
| `How-to-Install.md` | Per-platform download, `PATH`, build from source |
| `Quick-Start.md` | First shortcut in four commands |
| `How-to-Use.md` | Command map, global flags, prefix shortcuts |
| `Recipes.md` | Copy-paste workflows |
| `Scripting.md` | stdout/stderr, exit codes, `--no-input`, `--config` |
| `Sharing.md` | Export, import, registry |
| `Configuration.md` | Data files, `--config`, upgrade |
| `Troubleshooting.md` | The failures people actually hit |
| `Development.md` | `make check`, sandbox, branches, releasing |
| `Verifying-claims.md` | The unproven claims and the steps that settle them |
| `How-to-Contribute.md` | Issues, discussions, the bar for a patch |
| `Documentation.md` | Editing the wiki, style guide |
| Command pages | `create`, `update`, `delete`, `list`, `execute`, `open`, `copy`, `move`, `search`, `extract`, `revert`, `restore`, `reset`, `share`, `gui` |

> **Note:** enabling the wiki is a manual GitHub setting (Settings → Features →
> Wikis), and CI cannot push to it. The sync above is a human step — run it in
> the same change that edits a page here, or the wiki keeps serving the old one.
