## ADDED Requirements

### Requirement: Ensure a single active task collaboration room
The system SHALL maintain at most one active collaboration room for a given team task or subtask, and the one-click collaboration action SHALL return the active room when one already exists.

#### Scenario: Create active room when none exists
- **WHEN** an authorized team member clicks one-click collaboration for a task without an active collaboration room
- **THEN** the system creates an active task collaboration room bound to that task and opens it immediately

#### Scenario: Re-enter existing active room
- **WHEN** an authorized team member clicks one-click collaboration for a task that already has an active collaboration room
- **THEN** the system opens the existing active room instead of creating a second active room

#### Scenario: Create new room after dismissal
- **WHEN** the previous collaboration room for a task has been silently dismissed and an authorized member clicks one-click collaboration again
- **THEN** the system creates a new active collaboration room while preserving the dismissed room as read-only history

### Requirement: Render lifecycle-aware task collaboration entry points
The frontend SHALL render task and subtask collaboration actions according to the collaboration lifecycle state.

#### Scenario: No active room action
- **WHEN** a task has no active collaboration room
- **THEN** the task UI shows a one-click collaboration action

#### Scenario: Active room action
- **WHEN** a task has an active collaboration room
- **THEN** the task UI shows an enter-collaboration action that opens the existing room

#### Scenario: Dismissed room action
- **WHEN** a task has only dismissed collaboration rooms
- **THEN** the task UI allows starting a new collaboration while keeping dismissed history accessible

### Requirement: Preserve team task return context
The system SHALL return users from a task collaboration room to the team task page with the previously selected team and task context restored when possible.

#### Scenario: Return to selected team
- **WHEN** a user enters a collaboration room from a team task page and then clicks return to team tasks
- **THEN** the frontend returns to the team task page with the original team selected

#### Scenario: Return with task context
- **WHEN** a user returns from a task collaboration room that was opened from an expanded or focused task
- **THEN** the frontend restores or highlights the originating task when the task still exists

### Requirement: Show task collaboration participants
The task collaboration room SHALL show the members who are part of the room.

#### Scenario: Room displays participants
- **WHEN** a participant opens an active or dismissed task collaboration room
- **THEN** the frontend displays the participant list with member names or display names

#### Scenario: Team members are included by default
- **WHEN** a collaboration room is created for a team task
- **THEN** the system includes all current team members as room participants unless access rules explicitly reject a member

### Requirement: Record system invitation message
The system SHALL add a system chat message when members are invited into a newly created task collaboration room.

#### Scenario: Invitation message created
- **WHEN** the system creates a new task collaboration room and adds participants
- **THEN** the chat history includes a system message describing which members were invited into the collaboration

#### Scenario: Invitation message is not repeated on re-entry
- **WHEN** a user re-enters an existing active task collaboration room
- **THEN** the system does not create another invitation system message

### Requirement: Present personal shared material as My Notes
The task collaboration UI SHALL label personal shared material as My Notes rather than Knowledge Base.

#### Scenario: Picker uses My Notes copy
- **WHEN** a participant opens the personal material picker in a task collaboration room
- **THEN** the frontend labels the picker and action as My Notes

#### Scenario: Shared card uses My Notes copy
- **WHEN** a participant shares a note into the room
- **THEN** the message is displayed as a My Notes card with source owner information

### Requirement: Open shared note cards in the collaboration room
Participants SHALL be able to open shared My Notes cards from chat history and inspect the shared content in the collaboration room.

#### Scenario: Open shared note card
- **WHEN** a participant clicks a shared My Notes card in chat history
- **THEN** the frontend opens a detail view showing the shared title, owner, summary or excerpt, tags when available, and shared content

#### Scenario: Respect room access for shared notes
- **WHEN** a user without access to the task collaboration room requests shared note details
- **THEN** the system rejects the request and does not expose note content

### Requirement: Save AI minutes as readable records
The system SHALL save AI collaboration minutes to tasks and team knowledge as readable text or Markdown rather than raw JSON.

#### Scenario: Save readable minutes to task
- **WHEN** a participant saves generated AI minutes back to the task
- **THEN** the task record stores a readable minutes entry with sections for summary, synchronized notes, action items, blockers, and next steps where available

#### Scenario: Save readable minutes to team knowledge
- **WHEN** a participant saves generated AI minutes as team knowledge
- **THEN** the created team knowledge entry contains readable text or Markdown rather than a raw JSON payload

#### Scenario: Normalize JSON-shaped minutes input
- **WHEN** a client sends structured minutes or JSON-shaped minutes to a save endpoint
- **THEN** the backend normalizes the value into readable text before storing user-visible content

### Requirement: Keep dismissed task rooms readable and inactive
The system SHALL keep silently dismissed task collaboration rooms readable while preventing new discussion.

#### Scenario: Dismissed room is read-only
- **WHEN** a participant opens a dismissed task collaboration room
- **THEN** the frontend shows history and saved minutes in read-only mode and disables new text or My Notes messages

#### Scenario: Dismissed room is not active
- **WHEN** the task collaboration action checks for an active room
- **THEN** dismissed rooms are not treated as active rooms

### Requirement: Verify persistent task collaboration through browser testing
The implementation SHALL be verified with browser-based testing for the persistent task collaboration lifecycle.

#### Scenario: Browser verifies active-room reuse
- **WHEN** the implemented flow is tested in the browser
- **THEN** repeated one-click collaboration on the same task opens the same active room until the room is silently dismissed

#### Scenario: Browser verifies readable collaboration outputs
- **WHEN** the implemented flow is tested in the browser
- **THEN** shared My Notes cards are openable and saved AI minutes appear as readable text rather than raw JSON
