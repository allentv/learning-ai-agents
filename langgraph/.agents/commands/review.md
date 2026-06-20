Review the code in the current working directory or specified files for:

1. **Correctness**: Logic errors, edge cases, race conditions
2. **Type safety**: Ensure all types are annotated (mypy strict)
3. **Style**: Black formatting, isort, ruff lint rules
4. **Security**: Hardcoded secrets, injection risks, unsafe patterns
5. **Performance**: Unnecessary async/await, blocking calls, N+1 queries
6. **Testing**: Are new functions covered by tests?

For each issue found:
- State the file and line
- Explain the problem
- Suggest a concrete fix

Run `mise run v` at the end to verify all checks pass.
