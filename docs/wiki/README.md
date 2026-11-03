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
is the navigation.

> **Note:** enabling the wiki is a manual GitHub setting (Settings → Features →
> Wikis), and CI cannot push to it. The sync above is a human step — run it in
> the same change that edits a page here, or the wiki keeps serving the old one.
