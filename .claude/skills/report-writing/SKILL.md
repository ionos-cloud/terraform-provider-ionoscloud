---
name: report-writing
description: Write a structured report file when an autonomous task cannot be completed. Use when an agent has exhausted reasonable approaches on a task and needs to hand the situation back to the user with full context — what the task was, what was tried, why it stopped, and what input would unblock it. Takes a single argument: the filename the report should be saved as.
arguments: [filename]
---

# Report Writing

Use this skill to produce a single Markdown report that explains why a task could not be finished. Reports are the only artifact an agent leaves for a `BLOCKED` task; for `DONE` tasks no report is written.

## Argument

This skill receives one positional argument, `$filename` — the name the report file should be saved as (e.g., `blocked-k8s-nodepools-bug.md`). Use `$filename` exactly as given. Do not invent a different name.

## Where to write

Write the report to `claude_generated/reports/$filename` relative to the repository root. If the `claude_generated/reports/` directory does not exist, create it. If `claude_generated/reports/$filename` already exists, **do not overwrite it** — derive a new name by appending a numeric suffix before the extension (e.g., if `$filename` is `blocked-k8s-nodepools-bug.md`, write to `blocked-k8s-nodepools-bug-2.md`, then `-3.md`, etc., until you find an unused name) and use that name instead.

## Required sections

Every report must contain these sections, in this order, as Markdown headings:

```markdown
# <Task title — terse, one line>

## Status
BLOCKED

## Task
<Verbatim or near-verbatim restatement of the task as the user gave it. Include any clarifying answers the user provided up-front that shaped the work.>

## What I tried
<Concrete, ordered list of approaches attempted. Each entry: what you did, what you observed. No platitudes ("I investigated thoroughly"), no vague claims. Cite file paths, function names, commands, and command output. If you ran tests, name the tests and their failure modes. If you read SDK code, name the files.>

## What I learned
<Specific findings that future work should build on: API quirks, undocumented behavior, mismatches between the request and reality, patterns in existing code that confused or misled. One bullet per finding.>

## Why I stopped
<The actual blocker, in one or two paragraphs. Be precise — distinguish "I don't have enough information" from "the API does not support this" from "this requires a decision only the user can make".>

## What would unblock me
<Concrete asks of the user. Questions to answer, decisions to make, credentials/access to provide, scope changes to approve. Phrase each as something the user can act on directly.>

## Branch & local state
<Name of the task branch. Whether changes are committed, uncommitted, or stashed. If uncommitted, summarize what's on disk so the user can resume or discard.>
```

## Style rules

- **Be specific.** "Investigated the bug" is useless; "ran `TF_ACC=1 go test ./internal/framework/services/k8s -run TestAccK8sNodePool -v`, observed the data source returns `node_pools = []` even when `len(pools) > 0` in the API response (see `nodepool_data_source.go:147`)" is useful.
- **No padding.** Skip introductions, conclusions, and meta-commentary about the report itself.
- **No emojis.** No decorative formatting. Plain Markdown only.
- **Quote outputs sparingly.** Include error messages or command output only when they are load-bearing for the user's understanding. Trim unrelated lines.
- **One report per blocked task.** Do not bundle multiple tasks into one file.

## What this skill does NOT do

- Does not write reports for completed (`DONE`) tasks. Successful tasks produce a commit on a branch and nothing else.
- Does not write progress updates, status check-ins, or summaries between tasks.
- Does not modify the original task list — that's the agent's responsibility (marking `DONE` / `BLOCKED`).
- Does not push to a remote, open a PR, or notify anyone. The report is a local file.
