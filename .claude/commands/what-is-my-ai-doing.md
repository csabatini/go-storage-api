---
description: Live monitor of everything your AI is doing — tokens, cost, violations, tool usage
argument-hint: [--json]
allowed-tools: Bash, AskUserQuestion
---

# What Is My AI Doing?

Launch the RuleCatch AI-Pooler live monitor to see everything your AI is doing in real time.

**Arguments:** $ARGUMENTS

## Important

This monitor runs in a **separate terminal** — it does NOT run inside this Claude session. It's a persistent process that watches all AI activity.

## Step 1 — Launch the Monitor

Tell the user:

```
The AI-Pooler monitor needs to run in a separate terminal window.
Open a new terminal and run:

  npx @rulecatch/ai-pooler@latest monitor -v

This shows you a live view of:
  • Every tool call Claude makes (Read, Write, Edit, Bash, etc.)
  • Token usage per turn
  • Cost per session
  • Rule violations as they happen
  • Which files are being accessed

Press Ctrl+C to stop the monitor.
```

## Step 2 — Verify Installation

Check if the ai-pooler is already installed:

```bash
npx @rulecatch/ai-pooler@latest --version 2>/dev/null
```

If not installed or no API key configured, tell the user:

```
To set up the AI-Pooler for the first time:

  npx @rulecatch/ai-pooler init --api-key=dc_your_key --region=us

Get your API key from https://app.rybbit.io (RuleCatch dashboard).
After init, run the monitor in a separate terminal:

  npx @rulecatch/ai-pooler@latest monitor -v
```

## Step 3 — Remind

After providing the instructions, remind the user:

- The monitor runs **outside** Claude's context — zero token overhead
- It watches ALL Claude sessions, not just this one
- Violations are reported to the RuleCatch dashboard automatically
- You can query violations from within Claude using the RuleCatch MCP: "RuleCatch, what was violated today?"
