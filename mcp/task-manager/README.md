# Task Manager MCP Server

Copyright (c) 2026 Michael Lechner. All rights reserved.

A sophisticated state-management server that allows agents to maintain a persistent checklist of objectives and architectural decisions.

## Core Concepts

### Plan Mode
A specialized state where the agent is restricted to "read-only" operations. During this phase, the agent explores the codebase, designs an implementation strategy, and populates the task list. The agent must present this plan to the user and receive approval before exiting Plan Mode to start implementation.

## Tools

### 1. `mlc_task_create`
Creates a new structured task in the checklist.
- **Use Case**: Track progress during complex multi-step refactors.

### 2. `mlc_task_update`
Updates an existing task's status (`pending`, `in_progress`, `completed`, `deleted`) or manages dependencies (`blocks`, `blocked_by`).
- **Use Case**: Mark a step as finished or add new tasks as they are discovered.

### 3. `mlc_task_list`
Displays all tasks in the current session.
- **Use Case**: Overview of what is done and what remains.

### 4. `mlc_task_get`
Retrieves detailed information for a specific task.

### 5. `mlc_enter_plan_mode`
Transitions the agent into Plan Mode. Use this proactively for any non-trivial task.

### 6. `mlc_exit_plan_mode`
Transitions the agent back to Implementation Mode after the user has approved the plan.

## Installation

Built as part of the main project:
```bash
task build
```
