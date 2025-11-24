# Multiplayer & Lobby System Task Tracker

Comprehensive tracker for building the multiplayer, lobby, and social features for the Go Mage server. Modern gRPC-Web based system with real-time updates.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

---

## T001: Project Initialization
**Priority:** P0  
****Dependencies:** None

**Description:**  
Initialize a new SvelteKit project with TypeScript, configure Vite, and set up Tailwind CSS.

**Acceptance Criteria:**
- [x] SvelteKit project created with TypeScript template
- [x] Tailwind CSS configured and working
- [x] Vite config includes appropriate build settings
- [x] Basic `+page.svelte` renders with Tailwind styles
- [x] Dev server starts without errors (`npm run dev`)
- [x] Project builds successfully (`npm run build`)

**Files to Create:**
- `svelte.config.js`
- `tailwind.config.js`
- `vite.config.ts`
- `src/app.css` (Tailwind imports)
- `src/routes/+page.svelte` (basic test page)

---

## T002: Environment Configuration
**Priority:** P0  
****Dependencies:** T001

**Description:**  
Set up environment variable management for API endpoints and configuration.

**Acceptance Criteria:**
- [x] `.env.example` file created with all required variables
- [x] `.env` file in `.gitignore`
- [x] Environment variables accessible via `$env/static/public` or `$env/dynamic/public`
- [x] At minimum: `PUBLIC_API_URL` and `PUBLIC_WS_URL` defined
- [x] README includes instructions for environment setup

**Files to Create:**
- `.env.example`
- Update `.gitignore`
- `src/lib/config.ts` (exports typed config values)

---

## T003: ESLint and Prettier Setup
**Priority:** P0
****Dependencies:** T001

**Description:**
Configure ESLint and Prettier with appropriate rules for Svelte and TypeScript.

**Acceptance Criteria:**
- [x] ESLint configured with Svelte and TypeScript plugins
- [x] Prettier configured to work with Svelte
- [x] `eslint.config.js` created (ESLint v9 flat config)
- [x] `.prettierrc` created with sensible defaults
- [x] `npm run lint` command works
- [x] `npm run format` command works
- [x] VSCode settings recommended in `.vscode/settings.json`

**Files Created:**
- `eslint.config.js` (ESLint v9 flat config format)
- `.prettierrc`
- `.vscode/settings.json`
- Updated `package.json` with lint/format scripts (lint, lint:fix, format, format:check)

---

## T004: gRPC-Web Client Library Setup
**Priority:** P0
****Dependencies:** T001, T002

**Description:**
Install and configure gRPC-Web client library and set up TypeScript type generation from proto files.

**Acceptance Criteria:**
- [x] `@grpc/grpc-js` installed (with @grpc/proto-loader and google-protobuf)
- [x] `protoc` TypeScript plugin configured (ts-proto with grpc-tools)
- [x] Script to generate TypeScript types from `.proto` files (scripts/generate-proto.sh)
- [x] Example proto files compiled successfully (game.proto and lobby.proto)
- [x] Generated types available in `src/lib/generated/` (game.ts and lobby.ts)
- [x] README includes instructions for regenerating types

**Files Created:**
- `proto/game.proto` - Game service protocol definitions
- `proto/lobby.proto` - Lobby service protocol definitions
- `scripts/generate-proto.sh` - Proto generation script
- `src/lib/generated/game.ts` - Generated game service types and client
- `src/lib/generated/lobby.ts` - Generated lobby service types and client
- `src/lib/grpc/client.ts` - Basic gRPC client setup with factory functions
- Updated `package.json` with `proto:generate` script
- Updated `README.md` with gRPC/Protobuf development instructions

---

## T005: Basic Routing Structure
**Priority:** P0
****Dependencies:** T001

**Description:**
Set up the basic routing structure with placeholder pages for main sections.

**Acceptance Criteria:**
- [x] Route for `/` (home/landing)
- [x] Route for `/login`
- [x] Route for `/register`
- [x] Route for `/lobby`
- [x] Route for `/table/[id]`
- [x] Route for `/game/[id]`
- [x] Route for `/profile`
- [x] Route for `/decks`
- [x] Each route has a placeholder component with page title
- [x] Navigation between routes works

**Files Created:**
- `src/routes/+page.svelte` - Home/landing page with navigation grid
- `src/routes/login/+page.svelte` - Login page placeholder
- `src/routes/register/+page.svelte` - Registration page placeholder
- `src/routes/lobby/+page.svelte` - Lobby page placeholder with table list
- `src/routes/table/[id]/+page.svelte` - Table lobby page placeholder
- `src/routes/game/[id]/+page.svelte` - Game view page placeholder
- `src/routes/profile/+page.svelte` - Profile page placeholder with stats
- `src/routes/decks/+page.svelte` - Deck management page placeholder
- All routes tested and working via dev server

---

## T006: Authentication Store
**Priority:** P0
****Dependencies:** T001

**Description:**
Create a Svelte writable store to manage authentication state (JWT token, user info).

**Acceptance Criteria:**
- [x] Store tracks: `isAuthenticated`, `token`, `user` (username, email, id)
- [x] `login()` function stores token in localStorage and updates store
- [x] `logout()` function clears token and resets store
- [x] `loadAuthFromStorage()` function restores session on app load
- [x] Store is properly typed with TypeScript
- [x] Token expiry check implemented (basic JWT decode)
- [x] Test suite created for auth store functionality
- [x] All code passes typecheck, lint, and format checks

**Files Created:**
- `src/lib/types/auth.ts` - Authentication type definitions (User, AuthState, JwtPayload, LoginCredentials, RegisterData)
- `src/lib/utils/jwt.ts` - JWT utility functions (decodeJwt, isTokenExpired, getTokenTimeRemaining, getUserFromToken)
- `src/lib/stores/auth.ts` - Authentication store with login/logout/loadAuthFromStorage/checkTokenValidity/updateUser
- `src/lib/stores/__tests__/auth.test.ts` - Comprehensive test suite for auth store
- Updated `package.json` with vitest dependencies

---

## T007: Login Page Component
**Priority:** P0
****Dependencies:** T005, T006

**Description:**
Create a functional login page with form validation and authentication.

**Acceptance Criteria:**
- [x] Form has username and password fields
- [x] Client-side validation (required fields, min length)
- [x] "Remember me" checkbox
- [x] Submit button with loading state
- [x] Error message display for failed login
- [x] Link to registration page
- [x] "Guest login" button
- [x] Redirects to `/lobby` on successful login
- [x] Form is accessible (proper labels, ARIA attributes)
- [x] Simulated API calls with mock JWT tokens
- [x] Integration with auth store
- [x] Auto-redirect if already authenticated

**Files Modified:**
- `src/routes/login/+page.svelte` - Complete login page with form validation, loading states, guest login, and auth integration
- Includes simulated API calls that will be replaced with actual backend calls later

---

## T008: Registration Page Component
**Priority:** P0
****Dependencies:** T005, T006

**Description:**
Create a registration page with form validation for new user signup.

**Acceptance Criteria:**
- [x] Form has username, email, password, confirm password fields
- [x] Client-side validation:
  - Username: 3-20 characters, alphanumeric
  - Email: valid email format
  - Password: min 8 characters
  - Confirm password: matches password
- [x] Submit button with loading state
- [x] Error message display (username taken, etc.)
- [x] Success message on registration
- [x] Link back to login page
- [x] Auto-login after successful registration

**Files Modified:**
- `src/routes/register/+page.svelte` - Complete registration page with comprehensive validation, error handling, success messages, auth integration, and auto-login functionality

---

## T009: Auth Guard (Protected Routes)
**Priority:** P0
****Dependencies:** T006

**Description:**
Create a route guard that redirects unauthenticated users to login page.

**Acceptance Criteria:**
- [x] `+layout.ts` or `+page.ts` load function checks auth state
- [x] Redirects to `/login` if not authenticated
- [x] Preserves original URL for redirect after login
- [x] Works with SSR and client-side navigation
- [x] Does not protect `/login` and `/register` routes

**Files Created/Modified:**
- `src/routes/(protected)/+layout.ts` - Server-side auth guard with redirect logic
- `src/routes/(protected)/+layout.svelte` - Client-side auth guard with loading state and periodic token validation
- `src/lib/utils/auth-guard.ts` - Reusable auth check utilities (isAuthenticated, isTokenValid, clearInvalidToken)
- `src/routes/login/+page.svelte` - Updated to handle returnUrl query parameter
- `src/routes/register/+page.svelte` - Updated to handle returnUrl query parameter
- Moved protected routes into `(protected)` group: lobby, decks, profile, game, table
---

## T010: Main App Layout Component
**Priority:** P0  
****Dependencies:** T001, T006

**Description:**  
Create the main application layout with navigation bar and content area.

**Acceptance Criteria:**
- [ ] Top navigation bar with logo/title
- [ ] User menu dropdown (username, logout button)
- [ ] Navigation links: Lobby, My Decks, Profile
- [ ] Connection status indicator (online/offline/reconnecting)
- [ ] Responsive design (hamburger menu on mobile)
- [ ] Layout wraps all protected routes
- [ ] Smooth page transitions

---

## T011: Toast Notification System
**Priority:** P0  
****Dependencies:** T001

**Description:**  
Create a global toast notification system for user feedback.

**Acceptance Criteria:**
- [ ] Notification types: success, error, info, warning
- [ ] Auto-dismiss after configurable duration (default 3s)
- [ ] Manual dismiss button (X icon)
- [ ] Stack multiple notifications
- [ ] Slide-in animation
- [ ] Positioned top-right or bottom-right
- [ ] Accessible (ARIA live region)
- [ ] Global store for managing notifications


---

## T012: Modal Dialog Component
**Priority:** P0  
****Dependencies:** T001

**Description:**  
Create a reusable modal dialog component for various use cases.

**Acceptance Criteria:**
- [ ] Backdrop overlay (semi-transparent dark background)
- [ ] Modal content area (centered)
- [ ] Close button (X icon in top-right)
- [ ] Close on backdrop click (optional prop)
- [ ] Close on ESC key press
- [ ] Prevent body scroll when modal open
- [ ] Fade-in/scale animation
- [ ] Accessible (focus trap, ARIA dialog role)
- [ ] Configurable size (small, medium, large)


---

## T013: Confirmation Dialog Component
**Priority:** P0  
****Dependencies:** T012

**Description:**  
Create a confirmation dialog component that wraps the Modal with Yes/No actions.

**Acceptance Criteria:**
- [ ] Uses Modal component as base
- [ ] Title and message props
- [ ] Confirm and cancel button text customizable
- [ ] Returns promise that resolves to boolean
- [ ] Destructive action styling (red confirm button) optional
- [ ] ESC key maps to cancel
- [ ] Enter key maps to confirm

---

## T014: Loading Spinner Component
**Priority:** P0  
****Dependencies:** T001

**Description:**  
Create a reusable loading spinner component with different sizes.

**Acceptance Criteria:**
- [ ] CSS-only spinner animation (no images)
- [ ] Size variants: small, medium, large
- [ ] Optional centered overlay mode (fullscreen)
- [ ] Optional label text below spinner
- [ ] Accessible (ARIA live region)
- [ ] Works with light and dark backgrounds

---

## T015: gRPC Client Service Factory
**Priority:** P0  
****Dependencies:** T004, T006

**Description:**  
Create a factory for gRPC service clients with authentication and error handling.

**Acceptance Criteria:**
- [ ] Factory function creates gRPC service clients
- [ ] Automatically injects JWT token in metadata
- [ ] Wraps calls with error handling
- [ ] Converts gRPC errors to user-friendly messages
- [ ] Logs errors to console (dev mode)
- [ ] Handles connection errors gracefully
- [ ] Implements request timeout (30s default)
- [ ] Typed with TypeScript

---

## T016: Connection Status Store
**Priority:** P0  
****Dependencies:** T015

**Description:**  
Create a store to track WebSocket/gRPC connection status with auto-reconnect.

**Acceptance Criteria:**
- [ ] States: `connected`, `connecting`, `disconnected`, `reconnecting`
- [ ] Automatic reconnection with exponential backoff
- [ ] Max reconnection attempts (10)
- [ ] Manual reconnect function
- [ ] Connection health check (ping/pong)
- [ ] Emit events on status change
- [ ] Display connection status in UI

---

## T017: Lobby Page - Table List Display
**Priority:** P0  
****Dependencies:** T009, T010, T015

**Description:**  
Create the lobby page with a list of active tables.

**Acceptance Criteria:**
- [ ] Fetch table list from API on page load
- [ ] Display tables in a grid or list layout
- [ ] Each table shows: format, host username, player count (2/4), status
- [ ] Empty state when no tables available
- [ ] Loading state while fetching
- [ ] Refresh button to manually reload tables
- [ ] Tables clickable to view details or join
- [ ] Responsive design (mobile-friendly)

**Files to Create/Modify:**
- `src/routes/(protected)/lobby/+page.svelte`
- `src/lib/components/TableCard.svelte`
- `src/lib/api/lobby.ts` (API calls)

---

## T018: Lobby Page - Real-Time Updates
**Priority:** P0  
****Dependencies:** T017, T016

**Description:**  
Implement gRPC streaming to receive real-time lobby updates.

**Acceptance Criteria:**
- [ ] Establish gRPC stream on lobby page mount
- [ ] Listen for table created/updated/deleted events
- [ ] Update table list in real-time without full refresh
- [ ] Add new tables to list with animation
- [ ] Remove closed tables from list
- [ ] Update existing table info (player count, status)
- [ ] Close stream on page unmount
- [ ] Handle stream errors and reconnection


---

## T019: Lobby Page - Filters and Search
**Priority:** P0  
****Dependencies:** T017

**Description:**  
Add filtering and search controls to the lobby page.

**Acceptance Criteria:**
- [ ] Format dropdown filter (All, Standard, Commander, Modern, etc.)
- [ ] "Open tables only" toggle checkbox
- [ ] Search by host username (text input)
- [ ] Filters applied client-side (no API call)
- [ ] Filter state persists across table updates
- [ ] Clear filters button
- [ ] Show filtered count vs total count

---

## T020: Lobby Page - Online Players Display
**Priority:** P0  
****Dependencies:** T017

**Description:**  
Display count and list of online players in the lobby.

**Acceptance Criteria:**
- [ ] Show online player count in header
- [ ] Collapsible sidebar or section with player list
- [ ] Player list shows usernames
- [ ] Online status indicator (green dot)
- [ ] Real-time updates when players join/leave
- [ ] Scroll if player list exceeds height
- [ ] Show "You" indicator for current user

---

## T021: Create Table Modal - Basic Structure
**Priority:** P0  
****Dependencies:** T012, T017

**Description:**  
Create a modal dialog for creating a new table with basic options.

**Acceptance Criteria:**
- [ ] Opens when "Create Table" button clicked
- [ ] Format selector dropdown (Standard, Commander, Modern, etc.)
- [ ] Player count selector (2, 3, 4, etc.)
- [ ] Optional password field with show/hide toggle
- [ ] Form validation (format required)
- [ ] "Create & Join" submit button
- [ ] Cancel button to close modal
- [ ] Loading state during submission
- [ ] Error handling for failed creation

---

## T022: Create Table Modal - Deck Selection
**Priority:** P0  
****Dependencies:** T021

**Description:**  
Add deck selection to the create table modal.

**Acceptance Criteria:**
- [ ] Dropdown showing user's saved decks for selected format
- [ ] "Upload new deck" option if no decks available
- [ ] Opens deck upload modal (separate component)
- [ ] Validates deck is for selected format
- [ ] Shows deck card count (60 cards, etc.)
- [ ] Deck selection required before creating table
- [ ] Fetches user decks on modal open

---

## T023: Lobby Chat Component
**Priority:** P0  
****Dependencies:** T017

**Description:**  
Create a chat panel for the lobby with message display and input.

**Acceptance Criteria:**
- [ ] Chat panel (side panel or bottom)
- [ ] Message list with auto-scroll to bottom
- [ ] Message input field with send button
- [ ] Send on Enter key press
- [ ] Display username and timestamp for each message
- [ ] System messages styled differently (gray, italic)
- [ ] Load last 50 messages on mount
- [ ] Real-time message updates via gRPC stream
- [ ] Scroll to bottom button when scrolled up
- [ ] Empty state when no messages

---

## T024: Chat - Whisper Command
**Priority:** P0  
****Dependencies:** T023

**Description:**  
Implement whisper functionality for private messages in chat.

**Acceptance Criteria:**
- [ ] Detect `/w username message` format
- [ ] Parse username and message
- [ ] Send whisper via API
- [ ] Display whispers in italic with "(whisper)" prefix
- [ ] Show sent whispers as "To [username]: message"
- [ ] Show received whispers as "From [username]: message"
- [ ] Whispers in different color (muted purple/blue)
- [ ] Error if username not found
- [ ] Cannot whisper to self

---

## T025: Chat - Rate Limiting Feedback
**Priority:** P0  
****Dependencies:** T023

**Description:**  
Implement client-side rate limiting feedback for chat messages.

**Acceptance Criteria:**
- [ ] Track message count per time window (10 messages per 60 seconds)
- [ ] Disable send button when limit reached
- [ ] Show countdown timer when rate limited
- [ ] Display warning message "Sending too fast, wait X seconds"
- [ ] Reset counter after time window expires
- [ ] Visual feedback (red text, disabled button)

---

## T026: Table View - Pre-Game Lobby Component
**Priority:** P0  
****Dependencies:** T009, T015

**Description:**  
Create the table lobby view where players wait before game starts.

**Acceptance Criteria:**
- [ ] Display table info header (format, host, table ID)
- [ ] Show player list with slots (occupied and empty)
- [ ] Each player shows: username, ready status indicator
- [ ] Local player has "Ready" toggle button
- [ ] Empty slots show "Waiting for player..."
- [ ] Host sees "Start Game" button (enabled when all ready)
- [ ] Non-host players cannot start game
- [ ] Real-time updates when players join/leave/ready
- [ ] "Leave Table" button with confirmation
- [ ] Password indicator if table is password-protected
---

## T027: Table View - Host Controls
**Priority:** P0  
****Dependencies:** T026

**Description:**  
Add host-specific controls to the table lobby.

**Acceptance Criteria:**
- [ ] Only visible to table host
- [ ] "Kick Player" button next to each non-host player
- [ ] Confirmation dialog before kicking
- [ ] API call to kick player
- [ ] Player removed from table immediately
- [ ] Toast notification on kick success/failure
- [ ] Cannot kick self
---

## T028: Table View - Table Chat
**Priority:** P0  
****Dependencies:** T026, T023

**Description:**  
Add table-specific chat panel to the table lobby.

**Acceptance Criteria:**
- [ ] Reuse Chat component from lobby
- [ ] Table chat scope (only players at table)
- [ ] Separate chat stream per table
- [ ] Chat persists when players join/leave
- [ ] System messages for player join/leave events
- [ ] Clear chat when table closes

---

## T029: Table View - Game Start Countdown
**Priority:** P0  
****Dependencies:** T026

**Description:**  
Add countdown timer before game starts when all players ready.

**Acceptance Criteria:**
- [ ] When host clicks "Start Game", show 5 second countdown
- [ ] Display countdown overlay (modal or banner)
- [ ] Count down: 5... 4... 3... 2... 1... Starting!
- [ ] Cancel countdown if player unreadies
- [ ] Navigate to game view after countdown completes
- [ ] Host can cancel countdown

---

## T030: Table View - Leave Table Confirmation
**Priority:** P0  
****Dependencies:** T026, T013

**Description:**  
Add confirmation dialog when player tries to leave table.

**Acceptance Criteria:**
- [ ] "Leave Table" button triggers confirmation
- [ ] Dialog shows warning: "Are you sure you want to leave?"
- [ ] Confirm button sends leave request to API
- [ ] On success, navigate back to lobby
- [ ] Show toast notification on error
- [ ] If host leaves, table closes for all players

---

## T031: Deck Management Page - Deck List Display
**Priority:** P0  
****Dependencies:** T009, T015

**Description:**  
Create the deck management page showing user's saved decks.

**Acceptance Criteria:**
- [ ] Fetch user's decks on page load
- [ ] Display decks in a grid layout
- [ ] Each deck card shows: format, card count, last modified date
- [ ] Empty state when no decks saved
- [ ] Loading state while fetching
- [ ] "Upload New Deck" button
- [ ] Decks grouped by format (optional)
- [ ] Click deck to view details

---

## T032: Deck Upload Modal - Text Import
**Priority:** P0  
****Dependencies:** T031, T012

**Description:**  
Create a modal for uploading/importing decks via text.

**Acceptance Criteria:**
- [ ] Large text area for deck list input
- [ ] Format selector dropdown
- [ ] Parse deck list format: `4 Lightning Bolt`
- [ ] Show card count as user types (real-time)
- [ ] Display validation errors:
  - Invalid card names
  - Wrong deck size (not 60 cards)
  - Illegal cards for format
- [ ] "Clear" button to reset text area
- [ ] "Save Deck" button (disabled if invalid)
- [ ] Loading state during save
- [ ] Success toast on save
- [ ] Close modal and refresh deck list on success

---

## T033: Deck Upload Modal - Validation Display
**Priority:** P0  
****Dependencies:** T032

**Description:**  
Add inline validation feedback for deck upload.

**Acceptance Criteria:**
- [ ] Real-time validation as user types
- [ ] Show validation errors in a list below text area
- [ ] Error types:
  - "Invalid card: [Card Name]"
  - "Deck must be 60 cards (currently: X)"
  - "[Card Name] is not legal in [Format]"
  - "Too many copies of [Card Name] (max 4)"
- [ ] Errors highlighted in red
- [ ] Green checkmark when deck is valid
- [ ] Disable save button while errors present

---

## T034: Deck Viewer Component
**Priority:** P0  
****Dependencies:** T031

**Description:**  
Create a component to view deck details and card list.

**Acceptance Criteria:**
- [ ] Display deck name and format
- [ ] Show total card count and breakdown (creatures, instants, etc.)
- [ ] Group cards by type (Creatures, Instants, Sorceries, etc.)
- [ ] Display card quantity and name
- [ ] Mana curve visualization (bar chart)
- [ ] Color distribution (pie chart or bar)
- [ ] "Export" button to download as text
- [ ] "Delete" button with confirmation
- [ ] "Edit" button to modify deck (future)

---

## T035: Deck Deletion with Confirmation
**Priority:** P0  
****Dependencies:** T034, T013

**Description:**  
Add deck deletion functionality with confirmation dialog.

**Acceptance Criteria:**
- [ ] "Delete" button in deck viewer
- [ ] Confirmation dialog: "Delete [Deck Name]?"
- [ ] Warning: "This action cannot be undone"
- [ ] API call to delete deck
- [ ] Remove deck from list on success
- [ ] Show toast notification
- [ ] Handle errors gracefully

---

## T036: User Profile Page - Basic Info Display
**Priority:** P0  
****Dependencies:** T009, T015

**Description:**  
Create user profile page displaying basic account information and stats.

**Acceptance Criteria:**
- [ ] Display username and email
- [ ] Show join date / account created date
- [ ] Display stats:
  - Total games played
  - Wins / Losses
  - Win rate percentage
  - Quit ratio (prominently displayed)
- [ ] Stats cards with icons
- [ ] Loading state while fetching profile
- [ ] Error state if fetch fails


---

## T037: User Profile - Change Password Form
**Priority:** P0  
****Dependencies:** T036

**Description:**  
Add change password functionality to user profile.

**Acceptance Criteria:**
- [ ] Form with three fields: current password, new password, confirm new password
- [ ] Client-side validation:
  - Current password required
  - New password min 8 characters
  - Confirm password matches new password
- [ ] "Change Password" submit button
- [ ] Loading state during submission
- [ ] Success toast on password changed
- [ ] Error messages for:
  - Wrong current password
  - Server errors
- [ ] Clear form on success
---

## T038: User Profile - Recent Match History
**Priority:** P0  
****Dependencies:** T036

**Description:**  
Display list of recent matches on profile page.

**Acceptance Criteria:**
- [ ] Show last 10 matches
- [ ] Each match displays:
  - Opponent username
  - Format
  - Result (Win/Loss/Draw)
  - Date/time
- [ ] Result badge colored (green win, red loss, gray draw)
- [ ] Sorted by most recent first
- [ ] Empty state if no matches played
- [ ] Link to match details (future)

---

## T039: Game View - Basic Layout Structure
**Priority:** P0  
****Dependencies:** T009

**Description:**  
Create the basic game view layout with zones for opponent, battlefield, and player.

**Acceptance Criteria:**
- [ ] Three main sections:
  - Top: Opponent area (hand placeholder, life, library count)
  - Middle: Battlefield (shared zone)
  - Bottom: Player area (hand, life, library count)
- [ ] Side panel for game chat
- [ ] Game info header (format, turn count, current phase)
- [ ] "Concede" button with prominent placement
- [ ] Responsive layout (stack vertically on mobile)
- [ ] Fixed positions for zones (no scrolling inside zones)
---

## T040: Game View - Game Info Header
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create the game info header showing game state and current phase.

**Acceptance Criteria:**
- [ ] Display format name
- [ ] Show current turn number
- [ ] Display active player indicator
- [ ] Show current phase (Upkeep, Main, Combat, etc.)
- [ ] Phase highlighted/animated during transition
- [ ] Priority indicator (whose turn to act)
- [ ] Timer display (if game has timer)
- [ ] Compact design (doesn't take too much space)

---

## T041: Game View - Player Hand Component
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create component to display player's hand with draggable cards.

**Acceptance Criteria:**
- [ ] Display cards in a horizontal row
- [ ] Cards overlap slightly (fan layout)
- [ ] Hover to preview card (enlarge and lift)
- [ ] Click to select card (highlight border)
- [ ] Multi-select with Shift+Click
- [ ] Drag card to battlefield or other zone
- [ ] Show card count badge
- [ ] Responsive (stack on mobile)
- [ ] Empty state when hand is empty

---

## T042: Game View - Card Component with Hover Preview
**Priority:** P0  
****Dependencies:** T041

**Description:**  
Create the card component with hover preview and tooltips.

**Acceptance Criteria:**
- [ ] Display card image (placeholder if not loaded)
- [ ] Hover shows enlarged card preview
- [ ] Preview positioned to not go off-screen
- [ ] Show card name as tooltip
- [ ] Display mana cost in corner
- [ ] Show tapped state (90° rotation animation)
- [ ] Display counters (+1/+1, etc.) as badges
- [ ] Selection highlight (border glow)
- [ ] Loading state for card images
- [ ] Fallback for missing images

---

## T043: Game View - Battlefield Component
**Priority:** P0  
****Dependencies:** T039, T042

**Description:**  
Create the battlefield zone where permanents are displayed.

**Acceptance Criteria:**
- [ ] Grid layout for permanents (automatic positioning)
- [ ] Separate sections for lands and nonlands
- [ ] Cards draggable within battlefield (reorder)
- [ ] Support for tapped cards (rotated)
- [ ] Grouping by type (optional)
- [ ] Zoom in/out controls (if many cards)
- [ ] Empty state when battlefield is empty
- [ ] Hover to preview any card
- [ ] Click to select card for actions

---

## T044: Game View - Opponent Hand Placeholder
**Priority:** P0  
****Dependencies:** T039, T042

**Description:**  
Create component to display opponent's hand as card backs.

**Acceptance Criteria:**
- [ ] Show card backs (not visible)
- [ ] Display count of cards in hand
- [ ] Horizontal row layout
- [ ] No hover preview (cards hidden)
- [ ] Tooltip shows "Opponent's hand (X cards)"
- [ ] Empty state when opponent hand is empty

**Files to Create:**
- `src/lib/components/game/OpponentHand.svelte`

---

## T045: Game View - Life Total Display
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create life total display component for both players.

**Acceptance Criteria:**
- [ ] Large, prominent life number
- [ ] Color-coded (green high, yellow medium, red low)
- [ ] Life change animation (flash on change)
- [ ] Show life change delta (e.g., "-3" or "+5")
- [ ] Positioned in player/opponent zones
- [ ] Accessible (ARIA labels)
- [ ] Optional: life history graph (last 10 changes)

**Files to Create:**
- `src/lib/components/game/LifeTotal.svelte`

---

## T046: Game View - Graveyard Display
**Priority:** P0  
****Dependencies:** T039, T042

**Description:**  
Create graveyard zone component for both players.

**Acceptance Criteria:**
- [ ] Shows top card of graveyard (if any)
- [ ] Card count badge
- [ ] Click to expand and view all cards
- [ ] Modal or side panel shows full graveyard
- [ ] Cards in graveyard are hoverable (preview)
- [ ] Close button to collapse graveyard view
- [ ] Empty state when graveyard is empty
- [ ] Separate graveyards for each player

**Files to Create:**
- `src/lib/components/game/Graveyard.svelte`
- `src/lib/components/game/GraveyardModal.svelte`

---

## T047: Game View - Exile Zone Display
**Priority:** P0  
****Dependencies:** T046

**Description:**  
Create exile zone component (similar to graveyard).

**Acceptance Criteria:**
- [ ] Shows exiled cards count
- [ ] Click to expand and view all cards
- [ ] Modal/panel shows exiled cards
- [ ] Cards are hoverable (preview)
- [ ] Empty state when no exiled cards
- [ ] Different visual style from graveyard

**Files to Create:**
- `src/lib/components/game/ExileZone.svelte`

---

## T048: Game View - Library Counter
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create library (deck) counter display showing remaining cards.

**Acceptance Criteria:**
- [ ] Shows card back image
- [ ] Badge with card count
- [ ] Positioned in player/opponent zones
- [ ] Updates in real-time as cards drawn
- [ ] Warning state when low cards (< 5)
- [ ] Empty state when library is empty

**Files to Create:**
- `src/lib/components/game/LibraryCounter.svelte`

---

## T049: Game View - Mana Pool Display
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create mana pool display showing available mana.

**Acceptance Criteria:**
- [ ] Shows mana symbols with counts (W, U, B, R, G, C)
- [ ] Animated when mana added/spent
- [ ] Positioned near player hand
- [ ] Compact display (mana icons + numbers)
- [ ] Updates in real-time
- [ ] Empty state when no mana available
- [ ] Accessible (ARIA labels)

**Files to Create:**
- `src/lib/components/game/ManaPool.svelte`
- `src/lib/components/game/ManaSymbol.svelte`

---

## T050: Game View - Game Chat Panel
**Priority:** P0  
****Dependencies:** T039, T023

**Description:**  
Add in-game chat panel (reuse Chat component).

**Acceptance Criteria:**
- [ ] Reuse existing Chat component
- [ ] Positioned on right side or bottom
- [ ] Game-specific chat scope
- [ ] Show game events in chat (card played, damage dealt)
- [ ] Collapsible panel (hide to maximize game view)
- [ ] System messages for game events

**Files to Modify:**
- `src/routes/(protected)/game/[id]/+page.svelte`
- Reuse `src/lib/components/Chat.svelte`

---

## T051: Game View - Action Log
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create scrollable action log showing game events.

**Acceptance Criteria:**
- [ ] List of game actions (played card, attacked, etc.)
- [ ] Timestamped entries
- [ ] Color-coded by player
- [ ] Icons for action types (sword for attack, etc.)
- [ ] Auto-scroll to latest action
- [ ] Scrollable to view history
- [ ] Positioned in side panel or bottom
- [ ] Collapsible to save space

**Files to Create:**
- `src/lib/components/game/ActionLog.svelte`
- `src/lib/components/game/ActionLogItem.svelte`

---

## T052: Game View - Concede Button with Confirmation
**Priority:** P0  
****Dependencies:** T039, T013

**Description:**  
Add concede button with confirmation dialog.

**Acceptance Criteria:**
- [ ] Prominent "Concede" button (red, top of UI)
- [ ] Confirmation dialog: "Are you sure you want to concede?"
- [ ] Warning: "You will lose the match"
- [ ] API call to concede game
- [ ] On success, show game end overlay
- [ ] Navigate back to lobby after conceding
- [ ] Toast notification on error

**Files to Modify:**
- `src/routes/(protected)/game/[id]/+page.svelte`
- `src/lib/api/game.ts` (concedeGame API call)

---

## T053: Game View - Priority Indicator
**Priority:** P0  
****Dependencies:** T039

**Description:**  
Create visual indicator showing which player has priority.

**Acceptance Criteria:**
- [ ] Shows whose turn it is to act
- [ ] Glowing border around active player's zone
- [ ] Text indicator: "Your turn" or "Opponent's turn"
- [ ] Animated pulse effect
- [ ] Updates in real-time with game state
- [ ] Different visual for priority vs. just active player

**Files to Create:**
- `src/lib/components/game/PriorityIndicator.svelte`

---

## T054: Game View - Pass Priority Button
**Priority:** P0  
****Dependencies:** T053

**Description:**  
Add "Pass Priority" button for player actions.

**Acceptance Criteria:**
- [ ] Button enabled only when player has priority
- [ ] Prominent placement (center bottom, next to hand)
- [ ] Keyboard shortcut (Space bar)
- [ ] Loading state while waiting for response
- [ ] Disabled state when waiting for opponent
- [ ] Tooltip explains action
- [ ] API call to pass priority

**Files to Modify:**
- `src/routes/(protected)/game/[id]/+page.svelte`
- `src/lib/api/game.ts` (passPriority API call)

---

## T055: Game State Store
**Priority:** P0  
****Dependencies:** T039, T015

**Description:**  
Create Svelte store to manage game state synchronized via gRPC.

**Acceptance Criteria:**
- [ ] Store tracks complete game state:
  - Players (life, hand size, library size, etc.)
  - Battlefield permanents
  - Graveyards, exile zones
  - Stack (spells being cast)
  - Current phase and turn
  - Priority holder
- [ ] Subscribe to gRPC stream for game updates
- [ ] Update store on each message received
- [ ] Handle partial updates (delta, not full state)
- [ ] Typed with TypeScript interfaces
- [ ] Reactive updates to UI components

**Files to Create:**
- `src/lib/stores/game.ts`
- `src/lib/types/game-state.ts`

---

## T056: Game State Synchronization via gRPC
**Priority:** P0  
****Dependencies:** T055, T016

**Description:**  
Implement gRPC streaming to receive real-time game state updates.

**Acceptance Criteria:**
- [ ] Establish bidirectional gRPC stream on game mount
- [ ] Send client actions (play card, pass priority, etc.)
- [ ] Receive server state updates
- [ ] Parse game state messages into store
- [ ] Handle connection errors (show reconnect overlay)
- [ ] Reconnect on disconnect with state restore
- [ ] Close stream on game end or page unmount

**Files to Modify:**
- `src/routes/(protected)/game/[id]/+page.svelte`
- `src/lib/stores/game.ts`
- `src/lib/grpc/game-stream.ts`

---

## T057: Card Interaction - Click to Select
**Priority:** P0  
****Dependencies:** T042, T055

**Description:**  
Implement card selection via click with visual feedback.

**Acceptance Criteria:**
- [ ] Click card to select (toggle selection)
- [ ] Visual feedback (border glow, highlight)
- [ ] Multi-select with Shift+Click
- [ ] Clear selection when clicking elsewhere
- [ ] Update selection state in store
- [ ] Show available actions for selected card
- [ ] Deselect when action taken

**Files to Modify:**
- `src/lib/components/game/Card.svelte`
- `src/lib/stores/game-selection.ts`

---

## T058: Card Interaction - Drag and Drop
**Priority:** P0  
****Dependencies:** T057

**Description:**  
Implement drag-and-drop for playing cards from hand to battlefield.

**Acceptance Criteria:**
- [ ] Drag card from hand
- [ ] Show drag preview (ghost card follows cursor)
- [ ] Valid drop zones highlighted (battlefield, graveyard)
- [ ] Drop to play card on battlefield
- [ ] Cancel drag with ESC key
- [ ] Snap back to origin if invalid drop
- [ ] API call to play card on valid drop
- [ ] Optimistic UI update (show immediately)
- [ ] Rollback on server rejection

**Files to Modify:**
- `src/lib/components/game/Card.svelte`
- `src/lib/components/game/PlayerHand.svelte`
- `src/lib/components/game/Battlefield.svelte`
- `src/lib/utils/drag-drop.ts`

---

## T059: Card Interaction - Target Selection Mode
**Priority:** P0  
****Dependencies:** T057

**Description:**  
Implement target selection mode for spells and abilities.

**Acceptance Criteria:**
- [ ] Enter targeting mode when spell/ability requires targets
- [ ] Valid targets highlighted (glow effect)
- [ ] Invalid targets grayed out
- [ ] Cursor changes to crosshair
- [ ] Click to select target
- [ ] Show selected target with indicator
- [ ] Cancel targeting with ESC or right-click
- [ ] Confirm target selection (button or Enter key)
- [ ] Support multi-target (select multiple)

**Files to Create:**
- `src/lib/components/game/TargetingMode.svelte`
- `src/lib/stores/game-targeting.ts`

**Files to Modify:**
- `src/lib/components/game/Card.svelte`

---

## T060: Card Interaction - Tap/Untap Animation
**Priority:** P0  
****Dependencies:** T042

**Description:**  
Add smooth rotation animation for tapping/untapping cards.

**Acceptance Criteria:**
- [ ] Smooth 90° rotation for tap
- [ ] Smooth -90° rotation for untap
- [ ] Animation duration: 200ms
- [ ] CSS transform for performance
- [ ] Visual state persists after animation
- [ ] Works on mobile (touch)

**Files to Modify:**
- `src/lib/components/game/Card.svelte`

---

## T061: Game Actions - Play Card from Hand
**Priority:** P0  
****Dependencies:** T058, T055

**Description:**  
Implement action to play a card from hand to battlefield.

**Acceptance Criteria:**
- [ ] Drag card from hand to battlefield or click "Play" button
- [ ] Show mana cost payment UI (if applicable)
- [ ] API call to play card
- [ ] Optimistic UI update (card moves immediately)
- [ ] Rollback if server rejects (invalid play)
- [ ] Update game state store
- [ ] Animation for card movement
- [ ] Toast notification on error

**Files to Modify:**
- `src/lib/components/game/PlayerHand.svelte`
- `src/lib/api/game.ts` (playCard API call)
- `src/lib/stores/game.ts`

---

## T062: Game Actions - Activate Ability
**Priority:** P0  
****Dependencies:** T057, T059

**Description:**  
Implement UI for activating card abilities.

**Acceptance Criteria:**
- [ ] Right-click card to show context menu (or long-press mobile)
- [ ] List available abilities
- [ ] Click ability to activate
- [ ] Show targeting mode if ability requires targets
- [ ] Show mana/cost payment UI
- [ ] API call to activate ability
- [ ] Update game state on success
- [ ] Error handling and rollback

**Files to Create:**
- `src/lib/components/game/CardContextMenu.svelte`

**Files to Modify:**
- `src/lib/components/game/Card.svelte`
- `src/lib/api/game.ts` (activateAbility API call)

---

## T063: Game Actions - Declare Attackers
**Priority:** P0  
****Dependencies:** T057, T055

**Description:**  
Implement UI for declaring attackers during combat.

**Acceptance Criteria:**
- [ ] During combat phase, show "Declare Attackers" button
- [ ] Click creatures to toggle attacker status
- [ ] Visual indicator on attacking creatures (red border)
- [ ] Show which creatures can attack (valid attackers highlighted)
- [ ] "Confirm Attackers" button to submit
- [ ] API call to declare attackers
- [ ] Lock creatures until combat ends
- [ ] Cancel button to reset selection

**Files to Create:**
- `src/lib/components/game/CombatPhase.svelte`
- `src/lib/components/game/DeclareAttackers.svelte`

**Files to Modify:**
- `src/routes/(protected)/game/[id]/+page.svelte`
- `src/lib/api/game.ts` (declareAttackers API call)

---

## T064: Game Actions - Declare Blockers
**Priority:** P0  
****Dependencies:** T063

**Description:**  
Implement UI for declaring blockers during combat.

**Acceptance Criteria:**
- [ ] After attackers declared, show "Declare Blockers" button
- [ ] Click creature to select as blocker
- [ ] Click attacking creature to assign blocker
- [ ] Visual indicators (arrows from blocker to attacker)
- [ ] Support multiple blockers per attacker
- [ ] Show which creatures can block
- [ ] "Confirm Blockers" button to submit
- [ ] API call to declare blockers

**Files to Create:**
- `src/lib/components/game/DeclareBlockers.svelte`

**Files to Modify:**
- `src/lib/api/game.ts` (declareBlockers API call)

---

## T065: Game Actions - Assign Combat Damage
**Priority:** P0  
****Dependencies:** T064

**Description:**  
Implement UI for assigning combat damage in complex blocking scenarios.

**Acceptance Criteria:**
- [ ] Show damage assignment UI when needed (trample, multiple blockers)
- [ ] Number inputs for damage to each blocker
- [ ] Validate total damage equals creature power
- [ ] Visual feedback on valid/invalid assignment
- [ ] "Confirm Damage" button
- [ ] API call to assign damage
- [ ] Show damage calculation preview

**Files to Create:**
- `src/lib/components/game/AssignDamage.svelte`

**Files to Modify:**
- `src/lib/api/game.ts` (assignDamage API call)

---

## T066: Game Actions - Mana Payment Interface
**Priority:** P0  
****Dependencies:** T055

**Description:**  
Create UI for paying mana costs when casting spells or activating abilities.

**Acceptance Criteria:**
- [ ] Modal shows when mana payment required
- [ ] Display cost to pay (e.g., "2RR")
- [ ] List available mana sources (lands, mana abilities)
- [ ] Click to tap source and add mana to pool
- [ ] Show current mana in pool vs. cost required
- [ ] Disable invalid sources (wrong color)
- [ ] "Pay Cost" button (enabled when enough mana)
- [ ] "Cancel" button to cancel action
- [ ] API call to pay mana
- [ ] Auto-pay option (let server choose sources)

**Files to Create:**
- `src/lib/components/game/ManaPayment.svelte`

**Files to Modify:**
- `src/lib/api/game.ts` (payMana API call)

---

## T067: Game Actions - Choice Dialog
**Priority:** P0  
****Dependencies:** T012

**Description:**  
Create modal for making game choices (e.g., choose one mode).

**Acceptance Criteria:**
- [ ] Modal displays when choice required
- [ ] Show choice prompt (question text)
- [ ] List of options (radio buttons or buttons)
- [ ] Single-select or multi-select (configurable)
- [ ] "Confirm" button to submit choice
- [ ] API call to submit choice
- [ ] Cannot close modal without choosing (forced choice)
- [ ] Timeout indicator if timed choice

**Files to Create:**
- `src/lib/components/game/ChoiceDialog.svelte`

---

## T068: Game Actions - Number Input Dialog
**Priority:** P0  
****Dependencies:** T012

**Description:**  
Create modal for inputting numbers (e.g., X costs, damage distribution).

**Acceptance Criteria:**
- [ ] Modal displays when number input required
- [ ] Number input field with min/max validation
- [ ] Increment/decrement buttons
- [ ] Show context (what the number is for)
- [ ] "Confirm" button
- [ ] API call to submit number
- [ ] Validation (cannot exceed max value)

**Files to Create:**
- `src/lib/components/game/NumberInputDialog.svelte`

---

## T069: Game View - Stack Visualization
**Priority:** P0  
****Dependencies:** T055

**Description:**  
Create visualization for the stack showing spells and abilities.

**Acceptance Criteria:**
- [ ] Vertical stack display (last on top)
- [ ] Each item shows card/ability and controller
- [ ] Arrows indicate order of resolution (bottom to top)
- [ ] Highlight current item resolving
- [ ] Empty state when stack is empty
- [ ] Positioned prominently (center or side)
- [ ] Updates in real-time as spells cast/resolve
- [ ] Click item to view details

**Files to Create:**
- `src/lib/components/game/Stack.svelte`
- `src/lib/components/game/StackItem.svelte`

---

## T070: Reconnection - Detect Disconnect
**Priority:** P0  
****Dependencies:** T016, T056

**Description:**  
Detect when connection to game server is lost.

**Acceptance Criteria:**
- [ ] Detect WebSocket/gRPC stream disconnect
- [ ] Show "Disconnected" overlay immediately
- [ ] Disable all game interactions
- [ ] Start reconnection timer
- [ ] Show reconnection attempt count
- [ ] Cancel ongoing actions on disconnect

**Files to Modify:**
- `src/lib/stores/connection.ts`
- `src/routes/(protected)/game/[id]/+page.svelte`

---

## T071: Reconnection - Auto-Reconnect Logic
**Priority:** P0  
****Dependencies:** T070

**Description:**  
Implement automatic reconnection with exponential backoff.

**Acceptance Criteria:**
- [ ] Attempt reconnection automatically
- [ ] Exponential backoff (1s, 2s, 4s, 8s, max 30s)
- [ ] Max 10 reconnection attempts
- [ ] Show "Reconnecting..." overlay with countdown
- [ ] Request game state snapshot on reconnect
- [ ] Restore game view from snapshot
- [ ] Success: hide overlay and resume game
- [ ] Failure after max attempts: show error and return to lobby

**Files to Modify:**
- `src/lib/stores/connection.ts`
- `src/lib/grpc/reconnect.ts`

---

## T072: Reconnection - Restore Game State
**Priority:** P0  
****Dependencies:** T071, T055

**Description:**  
Restore full game state after reconnecting.

**Acceptance Criteria:**
- [ ] Request game state snapshot from server on reconnect
- [ ] Update game store with received state
- [ ] Redraw all zones (hand, battlefield, graveyard, etc.)
- [ ] Restore selections and UI state
- [ ] Resume pending actions if any
- [ ] Show toast: "Reconnected successfully"
- [ ] Log events that occurred during disconnect

**Files to Modify:**
- `src/lib/stores/game.ts`
- `src/lib/api/game.ts` (getGameState API call)

---

## T073: Reconnection - Manual Reconnect Button
**Priority:** P0  
****Dependencies:** T070

**Description:**  
Add manual reconnect button to reconnection overlay.

**Acceptance Criteria:**
- [ ] "Reconnect Now" button in overlay
- [ ] Click to immediately attempt reconnection
- [ ] Reset backoff timer on manual click
- [ ] Disable button during reconnection attempt
- [ ] Show loading spinner on button

**Files to Create:**
- `src/lib/components/ReconnectOverlay.svelte`

---

## T074: AFK Detection and Warning
**Priority:** P0  
****Dependencies:** T055

**Description:**  
Detect AFK (away from keyboard) and warn player before auto-forfeit.

**Acceptance Criteria:**
- [ ] Detect no activity for 2 minutes when player has priority
- [ ] Show warning modal: "Are you still there?"
- [ ] Countdown timer (30 seconds)
- [ ] "I'm here" button to dismiss warning
- [ ] Auto-forfeit if no response
- [ ] Reset timer on any game action
- [ ] Show toast before forfeit: "You've been AFK too long"

**Files to Create:**
- `src/lib/components/game/AfkWarning.svelte`
- `src/lib/utils/afk-detector.ts`

---

## T075: Game End Overlay
**Priority:** P0  
****Dependencies:** T055

**Description:**  
Create overlay shown when game ends with results.

**Acceptance Criteria:**
- [ ] Modal shows game result (You Won! / You Lost / Draw)
- [ ] Display winner username
- [ ] Show game stats (turns, duration, cards played)
- [ ] "Return to Lobby" button
- [ ] "View Match Details" button (future)
- [ ] Cannot close without clicking button
- [ ] Confetti animation on win (optional)

**Files to Create:**
- `src/lib/components/game/GameEndOverlay.svelte`

---

## T076: Optimistic UI Updates
**Priority:** P0  
****Dependencies:** T055

**Description:**  
Implement optimistic UI updates for smoother gameplay.

**Acceptance Criteria:**
- [ ] Update UI immediately on user action (before server confirms)
- [ ] Track pending actions in store
- [ ] Show visual indicator for pending actions (spinner, gray overlay)
- [ ] Rollback state if server rejects action
- [ ] Show error message on rollback
- [ ] Merge server state with optimistic updates
- [ ] Handle race conditions (multiple actions in flight)

**Files to Modify:**
- `src/lib/stores/game.ts`
- `src/lib/utils/optimistic-updates.ts`

---

## T077: Error Handling - Global Error Boundary
**Priority:** P1  
****Dependencies:** T001

**Description:**  
Create error boundary component to catch unhandled errors.

**Acceptance Criteria:**
- [ ] Wrap app in error boundary
- [ ] Catch unhandled exceptions
- [ ] Display friendly error page
- [ ] Show error message and stack trace (dev mode only)
- [ ] "Report Error" button to send to server
- [ ] "Reload Page" button
- [ ] Log errors to console
- [ ] Does not crash entire app

**Files to Create:**
- `src/lib/components/ErrorBoundary.svelte`
- `src/routes/+error.svelte`

---

## T078: Error Handling - API Error Interceptor
**Priority:** P1  
****Dependencies:** T015

**Description:**  
Create interceptor to handle API errors globally.

**Acceptance Criteria:**
- [ ] Intercept all gRPC/API errors
- [ ] Convert error codes to user-friendly messages
- [ ] Show toast notification for errors
- [ ] Log errors for debugging
- [ ] Handle specific error codes:
  - 401: Redirect to login
  - 403: Show "Permission denied"
  - 404: Show "Not found"
  - 500: Show "Server error, try again"
- [ ] Retry logic for transient errors

**Files to Modify:**
- `src/lib/grpc/error-handler.ts`

---

## T079: Loading States - Skeleton Screens
**Priority:** P1  
****Dependencies:** T001

**Description:**  
Create skeleton loading components for better perceived performance.

**Acceptance Criteria:**
- [ ] Skeleton for table list (lobby)
- [ ] Skeleton for deck list
- [ ] Skeleton for profile stats
- [ ] Skeleton for match history
- [ ] Animated shimmer effect
- [ ] Matches actual component layout
- [ ] Accessible (ARIA busy state)

**Files to Create:**
- `src/lib/components/skeletons/TableListSkeleton.svelte`
- `src/lib/components/skeletons/DeckListSkeleton.svelte`
- `src/lib/components/skeletons/ProfileSkeleton.svelte`

---

## T080: Form Validation - Reusable Utility
**Priority:** P1  
****Dependencies:** T001

**Description:**  
Create reusable form validation utilities and components.

**Acceptance Criteria:**
- [ ] Validation rules (required, minLength, maxLength, email, etc.)
- [ ] Composable validators
- [ ] Error message generation
- [ ] Form validation helper function
- [ ] Input component with validation state
- [ ] Show errors on blur or submit
- [ ] Clear errors on input change

**Files to Create:**
- `src/lib/utils/validation.ts`
- `src/lib/components/forms/ValidatedInput.svelte`

**Example Usage:**
```typescript
const rules = {
  username: [required(), minLength(3), maxLength(20)],
  email: [required(), email()],
};
```

---

## T081: Settings Page - Basic Structure
**Priority:** P1  
****Dependencies:** T009

**Description:**  
Create settings page with placeholder sections.

**Acceptance Criteria:**
- [ ] Settings page route (`/settings`)
- [ ] Sections: Audio, Graphics, Game Preferences, Chat
- [ ] Collapsible sections (accordion)
- [ ] "Save Settings" button at bottom
- [ ] "Reset to Defaults" button
- [ ] Toast notification on save success

**Files to Create:**
- `src/routes/(protected)/settings/+page.svelte`

---

## T082: Settings - Audio Preferences
**Priority:** P1  
****Dependencies:** T081

**Description:**  
Add audio settings section.

**Acceptance Criteria:**
- [ ] "Enable Sound Effects" toggle
- [ ] Volume slider (0-100%)
- [ ] "Test Sound" button
- [ ] Settings saved to localStorage
- [ ] Applied immediately on change

**Files to Modify:**
- `src/routes/(protected)/settings/+page.svelte`
- `src/lib/stores/settings.ts`

---

## T083: Settings - Graphics Preferences
**Priority:** P1  
****Dependencies:** T081

**Description:**  
Add graphics/visual settings section.

**Acceptance Criteria:**
- [ ] "Enable Animations" toggle
- [ ] Card quality selector (low, medium, high)
- [ ] "Reduce Motion" toggle (accessibility)
- [ ] Theme selector (light, dark, auto)
- [ ] Settings saved to localStorage

**Files to Modify:**
- `src/routes/(protected)/settings/+page.svelte`
- `src/lib/stores/settings.ts`

---

## T084: Settings - Game Preferences
**Priority:** P1  
****Dependencies:** T081

**Description:**  
Add game-specific settings section.

**Acceptance Criteria:**
- [ ] "Auto-pass priority" toggle
- [ ] "Auto-order triggers" toggle
- [ ] "Confirm concede" toggle
- [ ] Default time per turn preference
- [ ] Settings saved to localStorage
- [ ] Applied in game view

**Files to Modify:**
- `src/routes/(protected)/settings/+page.svelte`
- `src/lib/stores/settings.ts`

---

## T085: Responsive Design - Mobile Layout Testing
**Priority:** P1  
****Dependencies:** Multiple (most UI components)

**Description:**  
Test and fix responsive layouts for mobile devices.

**Acceptance Criteria:**
- [ ] Test all pages on mobile viewport (375px width)
- [ ] Fix layout issues (overlapping, overflow)
- [ ] Ensure buttons are touch-friendly (min 44px height)
- [ ] Hamburger menu works on mobile
- [ ] Tables/decks displayed in stacked layout
- [ ] Chat collapsible on mobile
- [ ] Game view optimized for vertical layout
- [ ] No horizontal scrolling

**Files to Modify:**
- Various component files
- `src/app.css` (responsive utilities)

---

## T086: Accessibility - Keyboard Navigation
**Priority:** P2  
****Dependencies:** Multiple

**Description:**  
Ensure all interactive elements are keyboard accessible.

**Acceptance Criteria:**
- [ ] Tab through all interactive elements
- [ ] Visible focus indicators (outline)
- [ ] ESC key closes modals
- [ ] Enter key submits forms
- [ ] Arrow keys navigate lists
- [ ] Keyboard shortcuts documented
- [ ] No keyboard traps

**Files to Modify:**
- Various component files (add tabindex, key handlers)

---

## T087: Accessibility - ARIA Labels and Roles
**Priority:** P2  
****Dependencies:** Multiple

**Description:**  
Add proper ARIA attributes for screen reader accessibility.

**Acceptance Criteria:**
- [ ] All buttons have aria-label or accessible text
- [ ] Forms have proper labels (not just placeholders)
- [ ] Modals have role="dialog" and aria-labelledby
- [ ] Toasts have role="alert" or aria-live="polite"
- [ ] Loading states have aria-busy
- [ ] Lists have role="list" and role="listitem"
- [ ] Test with screen reader (NVDA or VoiceOver)

**Files to Modify:**
- Various component files

---

## T088: Performance - Lazy Loading Routes
**Priority:** P2  
****Dependencies:** T005

**Description:**  
Implement lazy loading for route components to reduce initial bundle size.

**Acceptance Criteria:**
- [ ] Routes loaded on-demand (not in initial bundle)
- [ ] Loading indicator shown while route loads
- [ ] Preload next route on hover (link prefetch)
- [ ] Analyze bundle size before/after
- [ ] Reduce initial bundle by at least 30%

**Files to Modify:**
- `svelte.config.js` (configure route splitting)
- `src/routes/` (add dynamic imports if needed)

---

## T089: Performance - Image Optimization
**Priority:** P2  
****Dependencies:** T042

**Description:**  
Optimize card image loading and caching.

**Acceptance Criteria:**
- [ ] Lazy load card images (only when visible)
- [ ] Use responsive image sizes (srcset)
- [ ] Cache images in browser (service worker)
- [ ] Show low-res placeholder while loading
- [ ] Progressive image loading (blur-up)
- [ ] Use CDN for card images
- [ ] Limit concurrent image requests

**Files to Modify:**
- `src/lib/components/game/Card.svelte`
- `src/lib/utils/image-loader.ts`

---

## T090: Performance - Virtual Scrolling for Long Lists
**Priority:** P2  
****Dependencies:** T017, T038

**Description:**  
Implement virtual scrolling for table list and match history.

**Acceptance Criteria:**
- [ ] Only render visible items + buffer
- [ ] Smooth scrolling performance
- [ ] Works with dynamic item heights
- [ ] Scroll position preserved on updates
- [ ] Library: svelte-virtual or custom implementation
- [ ] Handles 1000+ items smoothly

**Files to Modify:**
- `src/routes/(protected)/lobby/+page.svelte`
- `src/lib/components/MatchHistory.svelte`

---

## T091: Testing - Unit Test Setup
**Priority:** P2  
****Dependencies:** T001

**Description:**  
Set up Vitest for unit testing Svelte components and utilities.

**Acceptance Criteria:**
- [ ] Vitest configured and working
- [ ] Svelte testing library installed
- [ ] Example test for utility function
- [ ] Example test for component
- [ ] `npm test` command runs tests
- [ ] Coverage reporting configured
- [ ] CI integration ready (GitHub Actions)

**Files to Create:**
- `vitest.config.ts`
- `src/lib/utils/__tests__/example.test.ts`
- `src/lib/components/__tests__/Button.test.ts`

---

## T092: Documentation - README
**Priority:** P2  
****Dependencies:** T001

**Description:**  
Write comprehensive README with setup and development instructions.

**Acceptance Criteria:**
- [ ] Project description and purpose
- [ ] Prerequisites (Node.js version, etc.)
- [ ] Installation instructions (`npm install`)
- [ ] Environment setup (copy .env.example)
- [ ] Development commands (`npm run dev`)
- [ ] Build and deployment instructions
- [ ] Project structure overview
- [ ] Contributing guidelines
- [ ] License information

**Files to Create/Modify:**
- `README.md`

---

## T093: Documentation - Component Documentation
**Priority:** P2  
****Dependencies:** Multiple

**Description:**  
Document reusable components with usage examples.

**Acceptance Criteria:**
- [ ] JSDoc comments for all props
- [ ] Usage examples in component files
- [ ] Storybook setup (optional)
- [ ] Component props documented with types
- [ ] Events and slots documented
- [ ] Accessibility notes included

**Files to Modify:**
- Various component files (add JSDoc comments)

---

## T094: Sound Effects - Basic Implementation
**Priority:** P3  
****Dependencies:** T082

**Description:**  
Add sound effects for game events.

**Acceptance Criteria:**
- [ ] Sound for card played
- [ ] Sound for damage dealt
- [ ] Sound for game won/lost
- [ ] Sound for notification
- [ ] Volume controlled by settings
- [ ] Sounds can be toggled on/off
- [ ] Use Web Audio API or Howler.js
- [ ] Sounds are subtle and not annoying

**Files to Create:**
- `src/lib/utils/sound.ts`
- `static/sounds/card-play.mp3`
- `static/sounds/damage.mp3`
- etc.

---

## T095: Animations - Card Movement
**Priority:** P3  
****Dependencies:** T061

**Description:**  
Add smooth animations for card movement between zones.

**Acceptance Criteria:**
- [ ] Card slides from hand to battlefield
- [ ] Card flies to graveyard when destroyed
- [ ] Card fades when exiled
- [ ] Animation duration: 300-500ms
- [ ] Uses CSS transitions or Web Animations API
- [ ] Can be disabled in settings
- [ ] Doesn't block game state updates

**Files to Modify:**
- `src/lib/components/game/Card.svelte`
- `src/lib/utils/animations.ts`

---

## T096: 404 and Error Pages
**Priority:** P1  
****Dependencies:** T001

**Description:**  
Create custom 404 and error pages.

**Acceptance Criteria:**
- [ ] Custom 404 page with friendly message
- [ ] "Return to Lobby" button
- [ ] Generic error page for 500 errors
- [ ] Matches app theme and style
- [ ] Accessible (proper heading structure)

**Files to Create:**
- `src/routes/+error.svelte` (already in T077)
- Custom styling and messaging

---

## T097: Production Build and Deployment
**Priority:** P0  
****Dependencies:** T001

**Description:**  
Configure production build and deployment process.

**Acceptance Criteria:**
- [ ] `npm run build` creates optimized production build
- [ ] Environment variables work in production
- [ ] Static assets cached with long TTL
- [ ] Gzip/Brotli compression enabled
- [ ] Source maps available for debugging
- [ ] Docker container for deployment (optional)
- [ ] Deployment guide in README

**Files to Create:**
- `Dockerfile` (optional)
- `nginx.conf` (if using nginx)
- Deployment scripts

---

## Phase 2: Post-MVP Tasks

### Package A: Enhanced Deck Management
- [ ] T098: Multiple Saved Decks per Format
- [ ] T099: Deck Naming and Tagging System
- [ ] T100: Import from MTGO/Arena Formats
- [ ] T101: Deck Sharing and Export Codes
- [ ] T102: Deck Statistics and Metadata

### Package B: Spectating
- [ ] T103: Spectator Mode - Join as Observer
- [ ] T104: Spectator List Display
- [ ] T105: Hide Hidden Information in Spectator View
- [ ] T106: Spectator Chat (Separate from Players)
- [ ] T107: Leave Spectate Functionality

### Package C: Matchmaking Queue
- [ ] T108: Matchmaking Queue Page
- [ ] T109: Queue Join by Format
- [ ] T110: Queue Position and Wait Time Display
- [ ] T111: Match Found Modal (Accept/Decline)
- [ ] T112: ELO Rating Display

### Package D: Friend System
- [ ] T113: Friends List Page
- [ ] T114: Friend Request System
- [ ] T115: Online Status Indicators
- [ ] T116: Friend Invite to Table
- [ ] T117: Private Messaging Interface

### Package E: Match History & Replays
- [ ] T118: Detailed Match History Page with Pagination
- [ ] T119: Filter Match History (Opponent, Format, Date)
- [ ] T120: Replay Viewer Component
- [ ] T121: Replay Controls (Play, Pause, Step, Speed)
- [ ] T122: Share Replay Link

### Package F: Tournaments
- [ ] T123: Tournament List Page
- [ ] T124: Tournament Details View (Brackets, Standings)
- [ ] T125: Tournament Registration Flow
- [ ] T126: Round Pairings Display
- [ ] T127: Tournament Chat
