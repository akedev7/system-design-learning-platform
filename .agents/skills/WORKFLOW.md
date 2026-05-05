# Skill Workflow Guide

How to use these skills together, in order, for a feature from idea to shipped code.

---

## Step 0 — Setup (once per repo)

```
/setup-matt-pocock-skills
```

Run this once when you first bring these skills into a repo. It configures:
- Which issue tracker to use (GitHub, Linear, or local files)
- The label vocabulary for triage
- Where domain docs (`CONTEXT.md`, ADRs) live

Every other engineering skill reads from this config. Nothing else works correctly without it.

---

## Step 1 — Align before you build

```
/grill-with-docs        ← for code changes
/grill-me               ← for non-code plans
```

Run one of these at the start of every feature. The agent interviews you relentlessly — one question at a time — until every decision branch is resolved.

`/grill-with-docs` also:
- Updates **`CONTEXT.md`** as domain terms are agreed on (a shared glossary the other skills use)
- Writes **ADRs** for hard-to-reverse architectural decisions

This is the most important step. Skipping it is the primary cause of "it built the wrong thing."

---

## Step 2 — Capture the plan as a PRD

```
/to-prd
```

After the grilling session the agent synthesises everything discussed into a PRD and publishes it as an issue on your tracker, labeled `needs-triage`.

No further interview — it works from what's already in the conversation.

---

## Step 3 — Break the PRD into issues

```
/to-issues #<prd-issue-number>
```

The agent reads the PRD issue and proposes vertical-slice sub-issues — each a thin, end-to-end slice through every layer (schema → API → UI → tests). It quizzes you on granularity and dependencies before publishing anything.

Each sub-issue is published labeled `needs-triage` and references the parent PRD.

---

## Step 4 — Triage the sub-issues

```
/triage
```

Move each sub-issue through the state machine:

| State | Meaning |
|---|---|
| `needs-triage` | Not yet evaluated |
| `needs-info` | Waiting on reporter for clarification |
| `ready-for-agent` | Fully specified, agent can implement autonomously |
| `ready-for-human` | Needs human judgment — design call, manual test, external access |
| `wontfix` | Will not be actioned |

When you move an issue to `ready-for-agent`, triage writes an **agent brief** — a self-contained comment with everything an AFK agent needs to implement it without asking questions.

---

## Step 5 — Build with TDD

```
/tdd
```

Pick up a `ready-for-agent` issue and implement it with a red-green-refactor loop:

1. Write **one** failing test for the first behavior
2. Write the minimal code to make it pass
3. Repeat for the next behavior
4. Refactor once all tests are green

Tests verify **behavior through public interfaces**, not implementation details. A good test survives a full internal refactor.

---

## Step 6 — Debug when something breaks

```
/diagnose
```

When a bug resists a quick fix, run `/diagnose`. It follows a disciplined six-phase loop:

1. **Build a feedback loop** — a fast, repeatable pass/fail signal (test, curl, CLI script, Playwright)
2. **Reproduce** — confirm the loop shows the exact failure
3. **Hypothesise** — generate 3–5 ranked, falsifiable hypotheses
4. **Instrument** — test one hypothesis at a time
5. **Fix + regression test** — write the test before applying the fix
6. **Cleanup** — remove debug logs, write a post-mortem note in the commit

If the fix reveals an architectural problem (no good test seam, tangled callers), hand off to `/improve-codebase-architecture`.

---

## Ongoing — Keep the codebase healthy

```
/improve-codebase-architecture
```

Run this every few days. It reads `CONTEXT.md` and the ADRs to find deepening opportunities — places where shallow modules can be extracted into deep ones with simple, stable interfaces.

---

## Any time — Get your bearings in unfamiliar code

```
/zoom-out
```

When you land in a section of code you don't know, run `/zoom-out`. The agent maps the relevant modules, callers, and relationships using the project's domain glossary.

---

## Any time — Reduce token usage

```
/caveman
```

Activates ultra-compressed communication mode (~75% fewer tokens). Toggle it on when you're doing repetitive work and don't need verbose responses.

---

## Full sequence at a glance

```
/setup-matt-pocock-skills       ← once per repo

/grill-with-docs                ← start of every feature
  └─ updates CONTEXT.md + ADRs

/to-prd                         ← after grilling; publishes PRD issue
  └─ issue labeled needs-triage

/to-issues #<prd>               ← breaks PRD into sub-issues
  └─ each labeled needs-triage

/triage                         ← evaluate each sub-issue
  └─ move to ready-for-agent or ready-for-human

/tdd                            ← implement, one vertical slice at a time
  └─ red → green → refactor

/diagnose                       ← when something breaks
  └─ hands off to /improve-codebase-architecture if architectural problem found

/improve-codebase-architecture  ← every few days
/zoom-out                       ← any time you're lost in unfamiliar code
/caveman                        ← any time you want terse output
```

---

## The connective tissue

`CONTEXT.md` and the ADRs (in `docs/adr/`) are written by `/grill-with-docs` and read by every other skill. They are what make the skills compose rather than just stack — a shared language that keeps naming consistent, narrows the thinking space for the agent, and documents decisions so you don't re-litigate them.

If you skip the grilling step, the downstream skills still work, but the output will be noisier and less accurate to your actual intent.
