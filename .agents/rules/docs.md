# Docs

Refines `~/.claude/rule-details/writing-docs.md`.

- **DOC1** `README.md` is the front door: what it is, install, a short usage example, the Documentation section and the Sponsor and License sections. Reference material lives in `docs/`, and each fact has one home.
- **DOC2** The README links to the docs site with full `https://psyb0t.github.io/go-asyncapi/...` URLs and to repository files with full `https://github.com/psyb0t/go-asyncapi/blob/main/...` URLs, because the README is also shown where relative links break.
- **DOC3** `docs/` stays flat. Every content page has front matter with a unique `title`, a `description` of 50 to 170 characters, a `nav_order`, and one `#` heading. Pages link to each other with relative `.md` paths and to files outside `docs/` with full GitHub URLs.
- **DOC4** A change to the exported API, to what `Load` reads, or to an error sentinel updates `docs/loading.md`, `docs/api.md`, `docs/errors.md` and the agent skill in the same commit, as needed.
- **DOC5** A change to the Makefile or `pipeline.yml` updates `docs/development.md`.
- **DOC6** `make docs` builds the site with the GitHub Pages image and must finish without errors before a docs change ships.
