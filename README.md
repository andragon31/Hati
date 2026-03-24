# Hati

**Task Planning & Human-in-the-Loop Layer for AI Development**

<p align="center">
<em>Structured Planning, Checkpoints, and Quality Gates</em>
</p>

Hati provides structured planning with phases, checkpoints, and developer approval workflows.

```
OpenCode / Claude Code / Cursor / ...
    ↓ MCP stdio
Hati (single Go binary)
    ↓
Plans & Phases + Quality Snapshot
```

## Features

- **Plan Management** - Create, revise, and track multi-phase plans
- **Checkpoints** - Human approval gates (pre, post, plan, final)
- **Spec Delta** - Track spec changes per plan
- **Module Hints** - Get context from Skoll
- **Quality Snapshot** - Integration with Tyr for quality gates
- **Learning Mode** - Capture and reuse Q&A

## Quick Start

### Install (One-liner)

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/andragon31/Hati/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/andragon31/Hati/main/install.ps1 | iex
```

### Setup Your Agent

| Agent | Command |
|-------|---------|
| OpenCode | `hati setup opencode` |
| Claude Code | `hati setup claude-code` |
| Cursor | `hati setup cursor` |
| Windsurf | `hati setup windsurf` |

## MCP Tools (20 total)

### Plan Tools
| Tool | Description |
|------|-------------|
| `plan_create` | Create new plan with phases |
| `plan_get` | Get plan details |
| `plan_revise` | Revise existing plan |
| `plan_abandon` | Abandon a plan |
| `plan_completeness` | Check plan completeness score |
| `plan_quality` | Get plan quality score |

### Checkpoint Tools
| Tool | Description |
|------|-------------|
| `checkpoint_open` | Open checkpoint for approval |
| `checkpoint_decide` | Developer decision on checkpoint |
| `checkpoint_status` | Get checkpoint status |

### Phase Tools
| Tool | Description |
|------|-------------|
| `phase_start` | Start executing a phase |
| `phase_report` | Report phase completion |

### Feedback Tools
| Tool | Description |
|------|-------------|
| `feedback_request` | Request feedback from developer |
| `feedback_receive` | Receive developer feedback |
| `feedback_escalate` | Escalate to human review |

### Record Tools
| Tool | Description |
|------|-------------|
| `record_list` | List approval records |
| `record_get` | Get approval record details |
| `record_export` | Export approval record (markdown/json) |

### System & Integration
| Tool | Description |
|------|-------------|
| `hati_status` | Get Hati system status |
| `hati_stats` | Get Hati statistics |
| `hati_commit_info` | Get commit info for traceability |
| `hati_register_commit` | Register commit with plan ID |
| `module_hints` | Get active module hints from Skoll |
| `spec_impact` | Get specs affected by plan |
| `quality_snapshot` | Get Quality Snapshot from Tyr |
| `learning_answer` | Submit learning mode answer |

## CLI Reference

```bash
hati setup [agent]   # Setup for an AI agent
hati init           # Initialize in project
hati mcp            # Start MCP server
hati version        # Show version
```

## Workflow Example

```bash
# 1. Create a plan
hati plan_create --request "Add user authentication"

# 2. Start first phase
hati phase_start --phase_id "phase-1"

# 3. Open checkpoint for approval
hati checkpoint_open --type "pre"

# 4. Developer approves
hati checkpoint_decide --decision "approved"

# 5. Complete phase
hati phase_report --phase_id "phase-1" --summary "Implemented auth"

# 6. Get quality snapshot
hati quality_snapshot
```

## Architecture

```
┌─────────────────────────────────────────────┐
│                 OpenCode                     │
│              Claude Code                     │
│                Cursor                        │
└─────────────────┬───────────────────────────┘
                  │ MCP stdio
                  ▼
┌─────────────────────────────────────────────┐
│                   Hati                       │
├─────────────────────────────────────────────┤
│  Plans    │  Phases   │ Checkpoints │Feedback│
├───────────┼───────────┼─────────────┼────────┤
│           │           │             │        │
│    ←──────┴───────────┴─────────────┴─────→  │
│                 │                            │
│        ┌────────┴────────┐                   │
│        ▼                 ▼                   │
│   ┌──────────┐    ┌──────────┐              │
│   │  Skoll   │    │   Tyr    │              │
│   │(hints)   │    │(quality) │              │
│   └──────────┘    └──────────┘              │
└─────────────────────────────────────────────┘
```

## Integration

Hati integrates with:
- **Fenrir** - Memory and knowledge tracking
- **Skoll** - Module hints and agent context
- **Tyr** - Quality snapshots and security checks

## License

MIT
