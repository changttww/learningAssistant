## Context

Team task collaboration currently happens through task comments, a generic quick meeting entry, and the public/private study-room chat infrastructure. The task model already supports team tasks, child tasks, owners, and assignee records. The study-room model already supports rooms, room members, chat messages, and WebSocket chat delivery. Knowledge-base entries can belong to a user or team and can be created from tasks or notes.

The gap is product semantics: a task discussion should start from a task or subtask, create a fresh temporary context every time, allow participants to bring personal knowledge into the discussion, and end with AI-generated minutes that can be saved back to the task.

## Goals / Non-Goals

**Goals:**

- Create a new task-bound collaboration session every time a user clicks the collaboration action on a team task or subtask.
- Reuse existing real-time room and chat infrastructure where possible.
- Keep each session's chat history, knowledge cards, AI minutes, and dismissal state independent from other sessions on the same task.
- Allow personal knowledge-base entries to be shared as structured cards without changing ownership of the original entries.
- Generate AI minutes that identify summary, synchronized knowledge, action items, owners, blockers, and next steps.
- Support saving AI minutes to the task discussion history and optionally saving selected content to the team knowledge base.
- Require browser-based verification for the core user flow after implementation.

**Non-Goals:**

- Real audio/video conferencing is not part of this change.
- Automatic notifications to all team members are not required for the first version.
- Automatic task creation or task assignment from AI action items is not required for the first version.
- Automatic saving of every shared personal knowledge entry into the team knowledge base is not required.
- Cross-session semantic search over all prior collaboration sessions is not required.

## Decisions

### Model task collaboration as one-off sessions

Each collaboration action creates a new session instead of reopening a long-lived task room. This keeps AI minutes scoped to one discussion and avoids mixing unrelated conversations from different days.

Alternative considered: one active room per task. That is simpler for navigation, but it makes summaries drift over time and makes "silent dismissal" ambiguous because closing the room would close the task's primary collaboration surface.

### Reuse study-room infrastructure with task-session metadata

The implementation should reuse `StudyRoom`, `StudyRoomMember`, `RoomSession`, `ChatMessage`, and the study-room WebSocket path for real-time delivery. A companion task-collaboration model or additional room metadata should record `task_id`, `team_id`, session status, creator, dismissed time, and saved minutes.

Alternative considered: a fully separate task-chat subsystem. That has cleaner naming, but it duplicates WebSocket, member, history, and message delivery logic already present in the app.

### Store structured messages while preserving text chat compatibility

Knowledge cards should be represented as structured messages that include at least `knowledge_entry_id`, title, summary/content excerpt, source user, category/tags, and original ownership metadata. Text messages remain first-class chat messages. If `ChatMessage.Content` remains the storage field, structured payloads must be encoded with an explicit message type so old text rendering is not confused with JSON.

Alternative considered: sending knowledge cards as plain text links. That is easier, but AI minutes and UI rendering would lose reliable attribution and source metadata.

### Treat personal knowledge sharing as view-only by default

When a member sends a personal knowledge entry into a session, other participants can view the shared card in that session. The original entry remains owned by the sender. Team knowledge is created only when a user explicitly saves selected content or minutes to the team knowledge base.

Alternative considered: automatically copying every shared card to the team knowledge base. That risks leaking unfinished personal notes and polluting team knowledge.

### Generate minutes from bounded session context

AI minutes should be generated from the selected session's messages and knowledge-card metadata only. The prompt should request structured JSON or a predictable schema with summary, synchronized knowledge, action items, blockers, and next steps. The backend should validate and normalize the result, falling back to a readable summary when the model response cannot be parsed.

Alternative considered: letting the frontend call the generic `/ai/chat` endpoint with a free-form prompt. That exists today, but it hides business rules in the UI and makes saving/verifying minutes harder.

### Silent dismissal changes state without noisy broadcasts

Dismissal should mark the session ended and prevent new messages. Participants can still view history and saved minutes if they have task access. The action should not create broad notifications or public room-list entries.

Alternative considered: deleting the room and messages. That would make the feature feel temporary, but it destroys the minutes/history that make the task discussion useful.

## Risks / Trade-offs

- **Risk: Reusing `StudyRoom` leaks task sessions into public study-room lists.** → Filter task-collaboration sessions from normal study-room listing, or add metadata/status that excludes them by default.
- **Risk: Chat message JSON payloads break existing text UI.** → Add explicit message types and normalize history responses so text and knowledge-card messages render predictably.
- **Risk: Personal knowledge sharing exposes private content too broadly.** → Only allow sharing by the owner, only expose the selected card to participants with task/session access, and avoid automatic team-knowledge persistence.
- **Risk: AI minutes hallucinate owners or action items.** → Include task/member context in the prompt, ask for evidence-linked action items when possible, and present generated minutes for user review before saving.
- **Risk: Silent dismissal leaves stale WebSocket clients connected.** → Backend should reject new messages after dismissal and frontend should switch the room into read-only ended state.
- **Risk: Browser testing is brittle if local AI provider tokens are missing.** → Verify non-AI flows independently and test AI minutes with configured provider when available; if unavailable, assert graceful error messaging.
