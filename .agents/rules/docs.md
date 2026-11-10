# Docs

- **DOC3** `docs/` stays flat.
- **DOC4** A change to the exported API, to what `Load` reads, or to an error sentinel updates `docs/loading.md`, `docs/api.md`, `docs/errors.md` and the agent skill in the same commit, as needed.
- **DOC5** A change to the Makefile or `pipeline.yml` updates `docs/development.md`.
- **DOC6** `make docs` builds the site with the GitHub Pages image and must finish without errors before a docs change ships.
- **REL4** The agent skill in `.agents/skills/go-asyncapi/` changes together with any exported API it describes, and every example in it must compile against the code.
