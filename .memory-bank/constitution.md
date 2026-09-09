---
description: Project Constitution — governing principles for AI-first development.
status: active
version: 1
project_principles: skipped
ratified: null
last_updated: 2026-09-04
---
# Project Constitution

## Purpose

This Constitution defines the non-negotiable principles that guide AI agents when planning, implementing, verifying, and synchronizing project work.

## Core Principles

### 0. Project Principles Status

The operator explicitly skipped the contextual interview on 2026-09-04 and chose to continue with the existing framework principles and the project KISS policy from `AGENTS.md`. `ratified: null` records that project-specific principles were not separately ratified. The interview may be revisited later.

### I. AI-First Spec-Driven Development

Agents MUST derive implementation work from explicit product, requirement, feature, task, and workflow artifacts. Agents MUST NOT invent product scope without evidence or user instruction.

### II. Memory Bank Is Durable Project Knowledge

`.memory-bank/` is the durable source of project knowledge. Chat context is temporary. Agents MUST update Memory Bank after meaningful changes.

### III. Schema-Backed Task Execution

Tasks MUST use the current schema-backed JSON task record model. If the framework uses `tier: T0|T1|T2|T3`, agents MUST route execution and verification through that tier model.

### IV. Minimal Verifiable Change

Agents SHOULD prefer the smallest change that satisfies the task. Every completed task MUST have clear checks or evidence.

### V. Evidence Before Done

A task MUST NOT be marked done without verification evidence appropriate to its tier and scope.

### VI. No Legacy Fallback and No Speculation

Agents MUST NOT rely on deprecated task formats, old risk models, or undocumented assumptions. Unknowns MUST be recorded as blockers or explicit assumptions.

### VII. Context Discipline

Agents SHOULD read the smallest sufficient context for the task. Higher-tier or cross-cutting tasks MUST read relevant normative docs such as invariants, contracts, states, testing, and workflow policies.

### VIII. Synchronization

After meaningful changes, agents MUST synchronize affected Memory Bank docs, task state, changelog, and routing files.

## Governance

- Constitution has precedence over workflow habits and generated plans.
- MBB, spec-index, spec-backbone, invariants, contracts, states, testing, and workflow docs refine this Constitution; they must not contradict it.
- Amendments must include rationale and update affected docs if needed.
- Constitution should stay short. Put concrete project rules into `invariants.md`, `contracts/*`, `states/*`, or workflow policy docs.

## Governance Decisions

- 2026-09-04: keep the existing framework principles, do not add project-specific constitutional machinery, and apply the existing KISS gate from `AGENTS.md`.

**Version**: 1 | **Ratified**: skipped | **Last updated**: 2026-09-04
