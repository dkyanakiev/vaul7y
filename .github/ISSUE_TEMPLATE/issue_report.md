---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug
assignees: dkyanakiev

---

**Describe the bug**
A clear and concise description of what the bug is.

**To reproduce**
Steps to reproduce the behavior:
1. Go to '...' view
2. Enter '...' key strokes
3. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots or a screen recording to help explain your problem.

**Please complete the following information:**
 - OS: [e.g. macOS 14, Ubuntu 22.04, Windows 11]
 - Terminal: [e.g. iTerm2, Alacritty, Windows Terminal]
 - Terminal dimensions: [run `echo "${COLUMNS}x${LINES}"` and paste the result]
 - Vaul7y version: [e.g. v0.2.0 or "built from source @ commit abc1234"]
 - Vault server version: [full version string, e.g. `Vault v1.15.4+ent` — indicates OSS vs Enterprise]
 - KV engine version: [v1 or v2, if the issue is secret-related]
 - Auth method in use: [e.g. token, AppRole, LDAP, OIDC]
 - Namespaces involved: [yes / no — if yes, describe the namespace path]
 - Configuration: [any relevant cli arguments, env vars, or `.vaul7y.yaml` options]
 - Vault permission policy: [paste the relevant policy if the issue may be permission-related]

**Log output**
Enable debug logging before reproducing the issue:
```shell
export VAULTY_LOG_LEVEL=debug
export VAULTY_LOG_FILE=/tmp/vaul7y-debug.log
vaul7y
```
Then paste the relevant lines from `/tmp/vaul7y-debug.log` here:
```
<paste log lines>
```

**Additional context**
Add any other context about the problem here.
