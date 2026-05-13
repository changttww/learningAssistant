## ADDED Requirements

### Requirement: Provide a persistent team communication chat room
The system SHALL provide one persistent text chat room for each team so team members can coordinate outside a specific task.

#### Scenario: Open team chat room
- **WHEN** a team member opens the team communication chat entry for a team
- **THEN** the system opens the persistent chat room for that team

#### Scenario: Reuse existing team chat room
- **WHEN** a team member opens the team communication chat entry multiple times for the same team
- **THEN** the system reuses the same persistent team chat room and shows previous messages

#### Scenario: Create team chat room when missing
- **WHEN** a team member opens the team communication chat entry and no persistent room exists for the team
- **THEN** the system creates one persistent team chat room for that team and opens it

### Requirement: Restrict team chat to team members
The system SHALL allow only current team members to access a team's persistent communication chat room.

#### Scenario: Team member can access
- **WHEN** a current team member opens the persistent team chat room
- **THEN** the system allows access and returns chat history

#### Scenario: Non-member is rejected
- **WHEN** a user who is not a member of the team requests the persistent team chat room
- **THEN** the system rejects the request and does not expose messages

### Requirement: Replace quick meeting entry with team chat
The frontend SHALL replace the quick meeting experience with the persistent team communication chat experience.

#### Scenario: Team page shows team chat
- **WHEN** a user views team collaboration or study-room related team actions
- **THEN** the UI presents a team communication chat entry instead of a quick meeting entry

#### Scenario: Existing quick meeting route does not create temporary meeting state
- **WHEN** a user reaches the previous quick meeting route or entry point
- **THEN** the frontend routes the user to the persistent team chat experience for the selected team when a team context is available

### Requirement: Preserve team chat history
The system SHALL persist and reload team chat messages across page reloads and repeated entry.

#### Scenario: Message persists
- **WHEN** a team member sends a non-empty text message in the persistent team chat room
- **THEN** the system stores the message and returns it in later chat history

#### Scenario: History remains team-scoped
- **WHEN** a member opens one team's persistent chat room
- **THEN** the history does not include messages from another team's persistent chat room

### Requirement: Verify persistent team chat through browser testing
The implementation SHALL be verified with browser-based testing for the persistent team chat flow.

#### Scenario: Browser verifies team chat reuse
- **WHEN** the implemented team chat flow is tested in the browser
- **THEN** opening the team chat repeatedly for the same team reuses the same persistent room and displays previous messages
