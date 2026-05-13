## Context

The current implementation builds on `StudyRoom`, `StudyRoomMember`, `RoomSession`, `ChatMessage`, task-collaboration metadata, and knowledge-base APIs. The first design intentionally created a fresh task collaboration session on every click, but the product direction has shifted: a task should have one active collaboration room at a time, and members should re-enter it until someone silently dismisses it.

The app also has a "quick meeting" concept that is too generic for a pure text-collaboration product. A long-lived team chat room better matches how teams coordinate outside a specific task.

## Goals / Non-Goals

**Goals:**

- Ensure one active task collaboration room per task or subtask until silent dismissal.
- Make one-click collaboration idempotent for active rooms: repeated clicks return the existing active room.
- Preserve return navigation to the user's previous team/task context.
- Pull the appropriate team members into task collaboration and show them in the room UI.
- Emit a readable system message in the chat when members are invited into a task collaboration room.
- Rename collaboration-room knowledge actions to `My Notes` in the UI.
- Let room participants open shared note cards and inspect shared content from inside the room.
- Save generated minutes to tasks and team knowledge as readable text or Markdown.
- Replace quick meeting with one persistent team chat room per team.

**Non-Goals:**

- Audio/video conferencing is not part of this change.
- Real-time typing indicators, read receipts, message reactions, and threaded replies are not required.
- Automatic task creation from AI action items is not required.
- Automatic copying of every shared personal note into team knowledge is not required.
- Cross-room semantic search across all team chats and task rooms is not required.

## Decisions

### Use an active-room lifecycle for task collaboration

Task collaboration SHALL use a lifecycle of `none -> active -> dismissed -> active`. The create action should behave as an "ensure active room" operation: if an active room already exists for the task, return it; otherwise create one. Dismissed rooms remain readable history and do not receive new messages.

Alternative considered: keep every click as a fresh session and expose a history chooser. That preserves the first implementation but conflicts with the user's expectation that one-click collaboration becomes the living room for a task until it is dismissed.

### Pull all team members into task collaboration by default

For team tasks, the room should include all current team members so the collaboration does not silently exclude reviewers or helpers who are not directly assigned to the subtask. The room still remains task-bound and hidden from normal study-room discovery.

Alternative considered: only include task owner and assignees. That is more restrictive but makes the user wonder why some team members cannot see or enter a task collaboration room.

### Store system invite messages as chat messages

When a task collaboration room is created, the system should add a system message such as `You invited A, B, C into collaboration`. This makes the room explain itself without requiring a separate notification system.

Alternative considered: show the invite only as transient UI toast. That disappears and does not help later readers understand who was pulled into the room.

### Keep backend names stable but change user-facing copy to My Notes

The existing knowledge-base data can still be used as the underlying source because notes and knowledge entries are already connected in the app. In the collaboration UX, buttons, empty states, modals, cards, and generated minutes copy should use `My Notes` when referring to the user's personal material.

Alternative considered: rename backend models and API paths immediately. That would create unnecessary churn without improving the user-facing behavior.

### Render shared notes as openable cards

A shared note card should carry enough metadata to open a detail modal or side panel inside the room. Participants should be able to inspect the shared title, owner, summary, tags, and shared content/excerpt without leaving the collaboration page.

Alternative considered: navigate to the global knowledge-base page. That breaks the chat context and makes return navigation harder.

### Format AI minutes at the backend persistence boundary

AI minutes may remain structured internally for display and duplicate detection, but save-to-task and save-to-team-knowledge operations must convert the structure into readable text or Markdown. The persistence layer should reject or normalize raw JSON strings before storing user-visible records.

Alternative considered: let the frontend stringify or format minutes before save. That spreads formatting rules across clients and risks storing raw JSON again.

### Model team chat as a persistent room per team

Each team should have one long-lived text chat room. The previous quick meeting entry should route to this persistent team chat instead of creating a short-lived meeting-like experience.

Alternative considered: keep quick meeting and add a new team chat entry. That leaves the weak existing surface in place and splits team communication into two similar places.

## Risks / Trade-offs

- **Risk: Active-room lookup can accidentally reopen dismissed rooms.** -> Query only active status for the default entry path; expose dismissed rooms only as history.
- **Risk: Pulling all team members exposes shared personal notes too broadly.** -> Make sharing explicit per note, show the source owner, and do not automatically copy shared notes to team knowledge.
- **Risk: Reusing study-room infrastructure leaks rooms into discovery.** -> Keep room-kind filtering for task collaboration and team chat where normal study-room discovery should not show them.
- **Risk: Saved minutes can still become raw JSON through old clients.** -> Centralize readable formatting in backend save endpoints and add tests for JSON-shaped input.
- **Risk: Return navigation can become stale if team/task data changes.** -> Prefer query or route state for `teamId`/`taskId`, with a safe fallback to team tasks.
- **Risk: Team chat and task rooms may feel visually identical.** -> Use clear labels, task context headers for task rooms, and team context headers for team chat.

## Migration Plan

- Existing active task collaboration sessions may be reused if they match the new active-room rule.
- If multiple active sessions already exist for the same task from earlier testing, choose the newest active session as the canonical room or dismiss/mark older sessions inactive during migration.
- Existing raw JSON minutes saved by earlier local testing do not need automatic migration for this change, but new saves must be readable.
- Existing quick meeting routes should redirect to or be replaced by the new team persistent chat route.

## Open Questions

- Should task-collaboration history show older dismissed sessions directly on the task card, or only from inside the current/active room?
- Should team chat support sharing My Notes in the same way as task collaboration, or remain pure text for this iteration?
