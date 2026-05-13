## ADDED Requirements

### Requirement: Start a one-off task collaboration session
The system SHALL allow an authorized team member to create a new temporary collaboration session from a team task or subtask every time the collaboration action is used.

#### Scenario: Create session from team task
- **WHEN** an authorized user clicks the collaboration action on a team task
- **THEN** the system creates a new collaboration session bound to that task and returns a room/session entry that can be opened immediately

#### Scenario: Create separate sessions for repeated clicks
- **WHEN** an authorized user clicks the collaboration action for the same task multiple times
- **THEN** the system creates distinct session records instead of reusing a previous active or dismissed session

#### Scenario: Reject unauthorized session creation
- **WHEN** a user without access to the selected team task requests a collaboration session
- **THEN** the system rejects the request and does not create a session

### Requirement: Scope session participants to task collaborators
The system SHALL associate each session with the task creator, task owner, child-task owner when applicable, and visible assignees/collaborators for the selected task context.

#### Scenario: Participant set is derived from task context
- **WHEN** a session is created for a team task with owners or assignees
- **THEN** the system records the intended participant set from the task creator, task owner, and assignees available to the requester

#### Scenario: Participant access follows task access
- **WHEN** a user opens a collaboration session
- **THEN** the system allows access only if the user has access to the bound task or is recorded as a session participant

### Requirement: Hide task sessions from ordinary study-room discovery
The system SHALL distinguish task collaboration sessions from ordinary public/private study rooms and SHALL NOT show them in the normal study-room discovery list.

#### Scenario: Study-room list excludes task sessions
- **WHEN** the user views the ordinary study-room list
- **THEN** task collaboration sessions are not included in the returned room list

#### Scenario: Task session is opened from task context
- **WHEN** the user opens a task collaboration session
- **THEN** the system loads it through task/session-specific navigation rather than public room discovery

### Requirement: Exchange text messages in a session
The system SHALL support real-time text chat and persisted chat history for each task collaboration session.

#### Scenario: Send text message
- **WHEN** a session participant sends a non-empty text message
- **THEN** the system persists the message under that session and broadcasts it to connected participants

#### Scenario: Load bounded chat history
- **WHEN** a participant opens a session
- **THEN** the system returns recent persisted messages for that session without including messages from other sessions on the same task

### Requirement: Share personal knowledge as cards
The system SHALL allow a session participant to select one of their own knowledge-base entries and send it into the session as a structured knowledge card.

#### Scenario: Share owned knowledge entry
- **WHEN** a participant selects an owned personal knowledge-base entry and sends it to the session
- **THEN** the system persists and broadcasts a knowledge-card message with the entry title, excerpt or summary, source owner, and source entry identifier

#### Scenario: Reject unowned knowledge sharing
- **WHEN** a participant attempts to send a knowledge-base entry they do not own and cannot access
- **THEN** the system rejects the knowledge-card message

#### Scenario: Render knowledge cards in history
- **WHEN** a participant loads session history containing knowledge-card messages
- **THEN** the frontend displays them as knowledge cards rather than raw JSON or plain text

### Requirement: Generate AI minutes for a session
The system SHALL generate AI minutes from a single collaboration session's text messages, knowledge-card metadata, and task context.

#### Scenario: Generate structured minutes
- **WHEN** a participant requests AI minutes for a session with discussion content
- **THEN** the system returns a structured result containing discussion summary, synchronized knowledge, action items with owners when inferable, blockers, and next steps

#### Scenario: Empty session minutes
- **WHEN** a participant requests AI minutes for a session without enough discussion content
- **THEN** the system returns a clear message that there is not enough content to summarize

#### Scenario: AI provider failure
- **WHEN** the AI provider fails or is unavailable during minutes generation
- **THEN** the system returns a non-destructive error and keeps the session usable

### Requirement: Save minutes to task follow-up
The system SHALL allow authorized participants to save AI minutes back to the bound task as discussion history or a task comment.

#### Scenario: Save minutes to task
- **WHEN** a participant saves generated minutes to the task
- **THEN** the system appends a readable minutes record to the bound task's discussion history without changing task status automatically

#### Scenario: Prevent duplicate accidental saves
- **WHEN** the same generated minutes are saved repeatedly without changes
- **THEN** the system prevents duplicate records or clearly marks subsequent saves as separate versions

### Requirement: Optionally create team knowledge from session output
The system SHALL allow a participant to explicitly convert selected session output into a team knowledge-base entry when the bound task belongs to a team.

#### Scenario: Save selected content as team knowledge
- **WHEN** a participant chooses to save selected minutes or shared knowledge content to team knowledge
- **THEN** the system creates a team-scoped knowledge-base entry linked to the task or session

#### Scenario: No automatic team knowledge copy
- **WHEN** a participant shares a personal knowledge card into a session
- **THEN** the system does not automatically copy that knowledge entry into the team knowledge base

### Requirement: Silently dismiss a temporary session
The system SHALL allow an authorized participant to silently dismiss a task collaboration session, marking it ended and preventing new discussion without broadcasting noisy notifications.

#### Scenario: Dismiss session
- **WHEN** an authorized participant dismisses an active collaboration session
- **THEN** the system marks the session ended, keeps history available, and prevents new text or knowledge-card messages

#### Scenario: Open dismissed session
- **WHEN** a participant opens a dismissed session
- **THEN** the system displays the session history and minutes in read-only mode

#### Scenario: No noisy broadcast on dismissal
- **WHEN** a session is silently dismissed
- **THEN** the system does not create broad team notifications or public system messages beyond the session state change

### Requirement: Verify the full flow through browser testing
The implementation SHALL be verified with browser-based testing in addition to backend or unit tests.

#### Scenario: Browser verifies collaboration flow
- **WHEN** the implemented feature is tested in the browser
- **THEN** the tester can create a session from a team task, send a text message, share a knowledge card, request AI minutes, save the result, dismiss the session, and observe read-only history after dismissal
