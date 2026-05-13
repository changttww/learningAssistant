## 1. Data Model And Backend Foundation

- [x] 1.1 Add task-collaboration session metadata that links each session to `task_id`, optional `team_id`, creator, room/session identifiers, status, dismissal time, and saved minutes state.
- [x] 1.2 Add support for structured chat message payloads or message metadata so text messages and knowledge-card messages can be stored and returned distinctly.
- [x] 1.3 Add database migration/model registration updates for the new session metadata and any message fields or companion tables.
- [x] 1.4 Add authorization helpers that reuse existing task access rules for session creation, opening, message posting, minutes generation, saving, and dismissal.

## 2. Session Lifecycle APIs

- [x] 2.1 Add an endpoint to create a new collaboration session from a team task or subtask every time it is called.
- [x] 2.2 Derive and persist intended participants from task creator, task owner, child-task owner when applicable, and assignees/collaborators visible to the requester.
- [x] 2.3 Add endpoints to fetch a collaboration session, list sessions for a task, and open a session through task-specific navigation.
- [x] 2.4 Exclude task collaboration sessions from ordinary study-room discovery responses.
- [x] 2.5 Add silent dismissal endpoint that marks a session ended, keeps history available, and prevents new messages.

## 3. Chat And Knowledge Sharing

- [x] 3.1 Update room chat history responses to normalize text and knowledge-card messages for frontend rendering.
- [x] 3.2 Update WebSocket and HTTP message posting paths to reject new messages after dismissal.
- [x] 3.3 Add an endpoint or WebSocket event for sending an owned personal knowledge-base entry into the session as a knowledge-card message.
- [x] 3.4 Validate knowledge-card sharing so users can only share entries they own or can access.
- [x] 3.5 Ensure messages from one session never appear in another session on the same task.

## 4. AI Minutes And Persistence

- [x] 4.1 Add backend AI minutes generation for a single collaboration session using task context, text messages, and knowledge-card metadata.
- [x] 4.2 Normalize AI minutes into discussion summary, synchronized knowledge, action items, owners, blockers, and next steps.
- [x] 4.3 Return clear empty-content and provider-failure responses without damaging the session.
- [x] 4.4 Add save-to-task behavior that appends readable minutes to task discussion history or comments without changing task status.
- [x] 4.5 Add optional save-to-team-knowledge behavior for selected minutes or shared knowledge content, with no automatic copy of personal knowledge.
- [x] 4.6 Prevent accidental duplicate saves of unchanged minutes or mark repeated saves as distinct versions.

## 5. Frontend Task Entry Points

- [x] 5.1 Add `一键协作` action to team task cards.
- [x] 5.2 Add `一键协作` action to visible subtask rows.
- [x] 5.3 Create frontend API methods for creating sessions, listing task sessions, opening session detail, dismissing sessions, generating minutes, saving minutes, and sharing knowledge cards.
- [x] 5.4 Route users from task/subtask collaboration actions into the newly created session every time.
- [x] 5.5 Add a task-level session history affordance so users can revisit prior session minutes/history.

## 6. Frontend Session Experience

- [x] 6.1 Build or adapt the collaboration room view so it presents task context, participant context, chat history, and ended/read-only state.
- [x] 6.2 Render text messages and knowledge-card messages with distinct layouts.
- [x] 6.3 Add a personal knowledge picker that lets the current user select and share an owned knowledge-base entry.
- [x] 6.4 Add AI minutes panel with generate, loading, error, review, save-to-task, and optional save-to-team-knowledge states.
- [x] 6.5 Add silent dismissal control for authorized users and switch dismissed sessions to read-only mode.
- [x] 6.6 Keep ordinary study-room UI behavior unchanged for non-task rooms.

## 7. Automated Verification

- [x] 7.1 Add backend tests for session creation, repeated-session behavior, authorization, participant derivation, and study-room list exclusion.
- [x] 7.2 Add backend tests for text message persistence, knowledge-card validation, dismissed-session write rejection, and session history isolation.
- [x] 7.3 Add backend tests for AI minutes empty-content handling, provider-failure handling, save-to-task behavior, optional team-knowledge save, and duplicate-save behavior.
- [x] 7.4 Add frontend tests or component-level checks for task entry actions, knowledge-card rendering, minutes states, and read-only dismissed sessions where the existing frontend test setup supports it.

## 8. Browser Verification

- [x] 8.1 Start the backend and frontend locally with valid development configuration.
- [x] 8.2 Use the browser to log in and open the team tasks page.
- [x] 8.3 Create a collaboration session from a team task and verify a new session opens.
- [x] 8.4 Create another collaboration session from the same task and verify it is distinct from the first session.
- [x] 8.5 Send a text message and verify it appears in real time and after reload.
- [x] 8.6 Share a personal knowledge-base entry and verify it renders as a knowledge card.
- [x] 8.7 Generate AI minutes and verify the structured sections are shown, or verify graceful provider error handling if no AI token is available.
- [x] 8.8 Save minutes to the task and verify the task discussion/history shows the saved record.
- [x] 8.9 Silently dismiss the session and verify the room becomes read-only without broad notifications.
- [x] 8.10 Confirm ordinary study-room discovery does not show task collaboration sessions.
