# CLAUDE.web.md -- Developer Workflow

> Workflow for Claude Code Web (Developer role)

---

## Phases

### Phase 2: Research
1. Study the specification in `project-doc/session-N/`
2. Explore the codebase
3. Document findings

### Phase 3: Implementation
1. Create a branch: `claude/session-N-implementation`
2. Implement changes
3. Commit with clear messages
4. Push and notify via hooks

---

## Git Workflow

```bash
# Create branch
git checkout -b claude/session-N-implementation

# Commit
git add .
git commit -m "feat: description of changes"

# Push
git push origin claude/session-N-implementation
```

---

## Commit Rules

- `feat:` -- new functionality
- `fix:` -- bug fix
- `docs:` -- documentation
- `refactor:` -- refactoring
- `test:` -- tests

---

## Pre-push Checklist

- [ ] Code compiles without errors
- [ ] Tests pass
- [ ] Linter is clean
- [ ] Documentation updated
- [ ] Commit message is clear

---

*Workflow v6.0*
