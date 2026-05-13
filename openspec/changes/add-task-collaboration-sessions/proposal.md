## Why

Team tasks and subtasks already have clear owners, but discussion is still disconnected from the task context. The existing study/meeting room experience is generic: members can chat, yet they cannot start a fresh task-scoped discussion, share personal knowledge into that discussion, or leave with an AI-generated record of who synchronized what and who should do what next.

This change adds one-off task collaboration sessions so each click from a team task creates a clean temporary discussion context with its own knowledge cards, AI minutes, action items, and silent dismissal.

## What Changes

- Add a `one-click collaboration` entry on team task cards and subtask rows.
- Create a new temporary collaboration session every time the entry is used, bound to the selected task or subtask.
- Bring the task creator, task owner, and visible assignees/owners into the intended participant set for the session.
- Reuse the existing real-time room/chat foundation while distinguishing task collaboration sessions from public study rooms.
- Allow room members to select entries from their personal knowledge base and send them into the session as structured knowledge-card messages.
- Add AI minutes for each session, covering discussion summary, synchronized knowledge, action items, owners, blockers, and next steps.
- Allow users to save AI minutes back to the task discussion history and optionally convert selected content into team knowledge-base entries.
- Add silent dismissal so a temporary session can be ended without broadcasting noisy notifications.

## Capabilities

### New Capabilities

- `task-collaboration-session`: One-off task-bound collaboration sessions with participant scoping, knowledge-card sharing, AI minutes, task follow-up capture, and silent dismissal.

### Modified Capabilities

None.

## Impact

- Backend models and routes for task-bound session creation, lookup, participant authorization, knowledge-card messages, AI minutes, task-history persistence, optional team knowledge creation, and room dismissal.
- Existing `Task`, `TaskAssignee`, `StudyRoom`, `StudyRoomMember`, `RoomSession`, `ChatMessage`, study-room WebSocket, and knowledge-base APIs are the primary integration points.
- Frontend team task/subtask UI gains a one-click collaboration action and session history affordances.
- Frontend collaboration room UI gains knowledge selection, knowledge-card rendering, AI minutes, save-to-task, optional save-to-team-knowledge, and silent dismissal controls.
- Verification must include backend tests plus browser-based flows for creating a session, sending chat/knowledge cards, generating AI minutes, saving results, and dismissing the session.
