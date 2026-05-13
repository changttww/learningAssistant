## Why

The first task-collaboration iteration introduced task-bound rooms, note sharing, AI minutes, and silent dismissal, but its "new temporary session on every click" behavior does not match the desired team workflow. Task collaboration should feel like a persistent group conversation tied to the task: users can re-enter the same active room, see who is present, share "My Notes", and only start a new room after the active one is silently dismissed.

The existing quick meeting room also has weak product value for this app. Teams need a persistent team chat space for everyday text coordination, separate from task-specific collaboration rooms.

## What Changes

- Change the product semantics of task one-click collaboration from "always create a fresh session" to "ensure and enter the active task collaboration room".
- Show task collaboration actions based on lifecycle state: `One-click collaboration` when no active room exists, `Enter collaboration` when an active room exists, and a new active room only after silent dismissal.
- Preserve the user's team/task context when returning from a collaboration room to the team task page.
- Add visible participant context in task collaboration rooms and create a system chat message when members are pulled into the collaboration.
- Treat the shared knowledge source as `My Notes` in the user experience, while continuing to reuse the existing note/knowledge data where appropriate.
- Allow participants to open shared note cards from the room and view the shared content without leaving the collaboration flow.
- Save AI collaboration minutes to tasks and team knowledge as readable text or Markdown, never as raw JSON.
- Replace the quick meeting-room experience with a persistent team communication chat room.

## Capabilities

### New Capabilities

- `task-persistent-collaboration-room`: Persistent task-bound collaboration rooms, active-room reuse, participant visibility, My Notes sharing, readable AI minutes, and state-aware navigation.
- `team-persistent-chat-room`: A persistent team-level text chat room that replaces the weak quick meeting-room experience.

### Modified Capabilities

None. The previous task-collaboration change has not been archived into the baseline `openspec/specs` tree, so this follow-up is modeled as new incremental capabilities that override the previous one-off semantics during implementation.

## Impact

- Backend task collaboration APIs must find or create the active room instead of creating a new room for every click.
- Backend session models, participant creation, chat-message handling, minutes formatting, and team-knowledge persistence are affected.
- Frontend team task cards and subtask rows must render lifecycle-aware collaboration buttons and return to the prior team/task context.
- Frontend collaboration-room UI must show participants, system messages, My Notes picker labels, clickable shared note cards, and readable saved minutes.
- Study-room or meeting-room frontend surfaces must be adjusted so quick meeting is replaced by a persistent team chat room.
- Browser verification must cover re-entering the same active task room, silent dismissal and recreation, My Notes sharing, readable saves, and team persistent chat.
