<div align="center">

# 📚 Editing this wiki

**These pages have a source of truth in the repository, so they can be reviewed.**

</div>

## 📍 Where the pages live

Every page here is a file under
[`docs/wiki/`](https://github.com/fiwon123/crower/tree/dev/docs/wiki) in the
repository. A change to a page is therefore a normal pull request: reviewed,
attributed, and kept in step with the code it describes.

## 🔁 How the sync works

The wiki is a **second git repository** (`crower.wiki.git`). GitHub creates it
when the first page is made in the web UI; after that:

```bash
git clone https://github.com/fiwon123/crower.wiki.git
cp docs/wiki/*.md crower.wiki/
cd crower.wiki && git add -A && git commit -m "wiki: sync from docs/wiki" && git push
```

The file name is the page name: `Home.md` is the landing page, `_Sidebar.md` is
the navigation, `_Footer.md` is the strip at the bottom.

## ✏️ Propose a change

1. Edit the page in `docs/wiki/` and open a pull request.
2. A maintainer runs the sync above after the merge.

Editing directly in the wiki UI works, but the next sync overwrites it — put the
change in the repository and it survives.

## 🎨 Style

- One emoji per heading, one blank line around every table and code block.
- Coloured badges come from [shields.io](https://shields.io); the palette is
  blue `4C8BF5`, green `22C55E`, amber `F59E0B`, purple `8250DF`.
- Callouts use GitHub alerts: `> [!NOTE]`, `> [!TIP]`, `> [!WARNING]`,
  `> [!IMPORTANT]`.
- Command pages follow one shape: what it does → synopsis → flags → examples →
  related.

---

🤝 [Contributing](How-to-Contribute) ·
🐛 [Issues](https://github.com/fiwon123/crower/issues)
