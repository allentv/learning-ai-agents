Fix all failing checks in this project:

1. Run `mise run v` to identify all failures
2. Fix formatting issues: `mise run f` then `mise run fc` to confirm
3. Fix lint issues: `mise run lf` then `mise run l` to confirm
4. Fix type errors: follow `.agents/skills/fix-type-errors.md`, then `mise run tc` to confirm
5. Fix failing tests: read error output, fix code (not test assertions unless asked), then `mise run t` to confirm
6. Final verification: `mise run v` — all checks must pass
