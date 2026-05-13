## 1. Backend Task Collaboration Lifecycle

- [x] 1.1 Update task collaboration creation to behave as an ensure-active-room operation that returns the existing active room when present.
- [x] 1.2 Enforce at most one active collaboration room per task or subtask in lookup and creation paths.
- [x] 1.3 Keep dismissed rooms readable while excluding them from active-room lookup.
- [x] 1.4 Add or update response fields so the frontend can distinguish active, dismissed, and missing collaboration state.

## 2. Participants And System Messages

- [x] 2.1 Include all current team members as participants when creating a team task collaboration room.
- [x] 2.2 Return participant display data in collaboration session detail responses.
- [x] 2.3 Create a single system invitation chat message when a new collaboration room is created.
- [x] 2.4 Ensure re-entering an existing active room does not duplicate the system invitation message.

## 3. My Notes Sharing Experience

- [x] 3.1 Rename collaboration-room knowledge picker, buttons, empty states, and cards to use My Notes copy.
- [x] 3.2 Add shared-note card detail support so participants can open a shared card inside the collaboration room.
- [x] 3.3 Validate that shared-note detail access follows task collaboration room access.
- [x] 3.4 Keep shared-note cards distinct from plain text and avoid rendering raw structured payloads.

## 4. Readable AI Minutes Persistence

- [x] 4.1 Add backend formatting that converts structured AI minutes into readable text or Markdown.
- [x] 4.2 Use readable formatting when saving AI minutes back to task history or comments.
- [x] 4.3 Use readable formatting when saving AI minutes as team knowledge.
- [x] 4.4 Add tests that JSON-shaped minutes input is normalized before user-visible storage.

## 5. Frontend Task Collaboration Flow

- [x] 5.1 Render task and subtask collaboration actions as one-click collaboration or enter collaboration based on active state.
- [x] 5.2 Preserve selected team and originating task context when navigating into and back from a collaboration room.
- [x] 5.3 Show participant list in the collaboration room header or side panel.
- [x] 5.4 Keep dismissed collaboration rooms read-only and allow new collaboration creation from the task after dismissal.

## 6. Persistent Team Chat

- [x] 6.1 Add backend support for one persistent team chat room per team.
- [x] 6.2 Restrict team chat access to current team members.
- [x] 6.3 Replace quick meeting UI entry points with the persistent team communication chat entry.
- [x] 6.4 Route previous quick meeting entry points to the persistent team chat when team context is available.
- [x] 6.5 Persist and reload team chat messages independently per team.

## 7. Tests And Browser Verification

- [x] 7.1 Add backend tests for active-room reuse, dismissal and recreation, participant creation, and system invitation messages.
- [x] 7.2 Add backend tests for shared-note detail access and readable minutes persistence.
- [x] 7.3 Add backend tests for persistent team chat creation, reuse, team-member authorization, and team-scoped history.
- [x] 7.4 Build the frontend and run existing automated checks available in the project.
- [x] 7.5 Use browser testing to verify repeated task collaboration entry reuses the same active room.
- [x] 7.6 Use browser testing to verify silent dismissal makes the room read-only and allows a new collaboration room afterwards.
- [x] 7.7 Use browser testing to verify My Notes cards are openable and saved minutes are readable.
- [x] 7.8 Use browser testing to verify the quick meeting surface has been replaced by persistent team chat and that chat history persists.
