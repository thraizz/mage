# Multiplayer & Lobby System Task Tracker

Comprehensive tracker for building the multiplayer, lobby, and social features for the Go Mage server. Modern gRPC-Web based system with real-time updates.

Status legend:

- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

---

## T001: Project Initialization

**Priority:** P0  
\***\*Dependencies:** None

**Description:**  
Initialize a new SvelteKit project with TypeScript, configure Vite, and set up Tailwind CSS.

**Acceptance Criteria:**

- [x] SvelteKit project created with TypeScript template
- [x] Vite config includes appropriate build settings
- [x] Basic `+page.svelte` renders with Tailwind styles
- [x] Dev server starts without errors (`npm run dev`)
- [x] Project builds successfully (`npm run build`)
- [x] CSS Styling in place, ready to extract and style components

**Files to Create:**

- `svelte.config.js`
- `tailwind.config.js`
- `vite.config.ts`
- `src/app.css` (Tailwind imports)
- `src/routes/+page.svelte` (basic test page)

---

## T002: Environment Configuration

**Priority:** P0  
\***\*Dependencies:** T001

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

**Priority:** P0 \***\*Dependencies:** T001

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

**Priority:** P0 \***\*Dependencies:** T001, T002

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

**Priority:** P0 \***\*Dependencies:** T001

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

**Priority:** P0 \***\*Dependencies:** T001

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

**Priority:** P0 \***\*Dependencies:** T005, T006

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

**Priority:** P0 \***\*Dependencies:** T005, T006

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

**Priority:** P0 \***\*Dependencies:** T006

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

**Priority:** P0 \***\*Dependencies:** T001, T006

**Description:**
Create the main application layout with navigation bar and content area.

**Acceptance Criteria:**

- [x] Top navigation bar with logo/title
- [x] User menu dropdown (username, logout button)
- [x] Navigation links: Lobby, My Decks, Profile
- [x] Connection status indicator (online/offline/reconnecting)
- [x] Responsive design (hamburger menu on mobile)
- [x] Layout wraps all protected routes
- [x] Smooth page transitions

**Files Created:**

- `src/lib/components/Navbar.svelte` - Main navigation bar with brand, links, connection status, and user menu
- `src/lib/components/UserMenu.svelte` - Dropdown user menu with profile links and logout
- `src/lib/components/ConnectionStatus.svelte` - Connection status indicator with tooltip

**Files Modified:**

- `src/routes/(protected)/+layout.svelte` - Updated to include Navbar and smooth page transitions

---

## T011: Toast Notification System

**Priority:** P0 \***\*Dependencies:** T001

**Description:**
Create a global toast notification system for user feedback.

**Acceptance Criteria:**

- [x] Notification types: success, error, info, warning
- [x] Auto-dismiss after configurable duration (default 3s)
- [x] Manual dismiss button (X icon)
- [x] Stack multiple notifications
- [x] Slide-in animation
- [x] Positioned top-right or bottom-right
- [x] Accessible (ARIA live region)
- [x] Global store for managing notifications

**Files Created:**

- `src/lib/types/toast.ts` - Toast type definitions (Toast, ToastType, ToastOptions)
- `src/lib/stores/toast.ts` - Global toast store with add/dismiss/success/error/warning/info methods
- `src/lib/components/Toast.svelte` - Individual toast component with icon, message, and dismiss button
- `src/lib/components/ToastContainer.svelte` - Container component for stacking toasts

**Files Modified:**

- `src/routes/+layout.svelte` - Added ToastContainer to root layout for global availability
- `src/routes/+page.svelte` - Added test buttons to demonstrate toast functionality

---

## T012: Modal Dialog Component

**Priority:** P0 \***\*Dependencies:** T001

**Description:**
Create a reusable modal dialog component for various use cases.

**Acceptance Criteria:**

- [x] Backdrop overlay (semi-transparent dark background)
- [x] Modal content area (centered)
- [x] Close button (X icon in top-right)
- [x] Close on backdrop click (optional prop)
- [x] Close on ESC key press
- [x] Prevent body scroll when modal open
- [x] Fade-in/scale animation
- [x] Accessible (focus trap, ARIA dialog role)
- [x] Configurable size (small, medium, large)

**Files Created:**

- `src/lib/types/modal.ts` - Modal type definitions (ModalSize, ModalProps)
- `src/lib/components/Modal.svelte` - Reusable modal component with backdrop, animations, focus trap, and accessibility

**Files Modified:**

- `src/routes/+page.svelte` - Added modal test section with 5 different modal examples demonstrating all features

---

## T013: Confirmation Dialog Component

**Priority:** P0 \***\*Dependencies:** T012

**Description:**
Create a confirmation dialog component that wraps the Modal with Yes/No actions.

**Acceptance Criteria:**

- [x] Uses Modal component as base
- [x] Title and message props
- [x] Confirm and cancel button text customizable
- [x] Returns promise that resolves to boolean
- [x] Destructive action styling (red confirm button) optional
- [x] ESC key maps to cancel
- [x] Enter key maps to confirm

**Files Created:**

- `src/lib/components/ConfirmDialog.svelte` - Reusable confirmation dialog component (callback-based)
- `src/lib/stores/confirm.ts` - Global confirmation store with promise-based API
- `src/lib/components/GlobalConfirmDialog.svelte` - Global confirmation dialog instance

**Files Modified:**

- `src/routes/+layout.svelte` - Added GlobalConfirmDialog to root layout
- `src/routes/+page.svelte` - Added confirmation dialog test section with 4 different examples (basic, destructive, custom text, component-based)

---

## T014: Loading Spinner Component

**Priority:** P0 \***\*Dependencies:** T001

**Description:**
Create a reusable loading spinner component with different sizes.

**Acceptance Criteria:**

- [x] CSS-only spinner animation (no images)
- [x] Size variants: small, medium, large
- [x] Optional centered overlay mode (fullscreen)
- [x] Optional label text below spinner
- [x] Accessible (ARIA live region)
- [x] Works with light and dark backgrounds

**Files Created:**

- `src/lib/types/loading.ts` - Loading spinner type definitions (LoadingSize, LoadingSpinnerProps)
- `src/lib/components/LoadingSpinner.svelte` - Reusable loading spinner with CSS animation, size variants, overlay mode, and accessibility

**Files Modified:**

- `src/routes/+page.svelte` - Added loading spinner test section with overlay and inline examples

---

## T015: gRPC Client Service Factory

**Priority:** P0 \***\*Dependencies:** T004, T006

**Description:**
Create a factory for gRPC service clients with authentication and error handling.

**Acceptance Criteria:**

- [x] Factory function creates gRPC service clients
- [x] Automatically injects JWT token in metadata
- [x] Wraps calls with error handling
- [x] Converts gRPC errors to user-friendly messages
- [x] Logs errors to console (dev mode)
- [x] Handles connection errors gracefully
- [x] Implements request timeout (30s default)
- [x] Typed with TypeScript

**Files Created:**

- `src/lib/types/grpc.ts` - gRPC type definitions (GrpcStatusCode, GrpcError, UserError, GrpcClientOptions, GrpcMetadata)
- `src/lib/utils/grpc-errors.ts` - Error handling utilities (toGrpcError, toUserError, isRetryableError, isAuthError, logGrpcError)
- `src/lib/grpc/service-factory.ts` - Service factory with auth injection, timeout, retry logic, and error handling (grpcCall, grpcCallWithToast, createServiceClient, createGrpcMetadata)
- `src/lib/grpc/__tests__/service-factory.test.ts` - Comprehensive test suite (18 tests covering metadata, timeouts, auth errors, retries)
- `vitest.config.ts` - Vitest configuration for testing

**Files Modified:**

- `package.json` - Added test scripts (test, test:watch, test:ui)

---

## T016: Connection Status Store

**Priority:** P0 \***\*Dependencies:** T015

**Description:**
Create a store to track WebSocket/gRPC connection status with auto-reconnect.

**Acceptance Criteria:**

- [x] States: `connected`, `connecting`, `disconnected`, `reconnecting`
- [x] Automatic reconnection with exponential backoff
- [x] Max reconnection attempts (10)
- [x] Manual reconnect function
- [x] Connection health check (ping/pong)
- [x] Emit events on status change
- [x] Display connection status in UI

**Files Created:**

- `src/lib/types/connection.ts` - Connection type definitions (ConnectionStatus, ConnectionState, ConnectionOptions, ConnectionEvent, ConnectionEventCallback)
- `src/lib/stores/connection.ts` - Connection store with exponential backoff, health check, event system (connection, connectionStatus, isConnected, connectionLatency)
- `src/lib/stores/__tests__/connection.test.ts` - Comprehensive test suite (17 tests covering states, reconnection, health check, events)

**Files Modified:**

- `src/lib/components/ConnectionStatus.svelte` - Updated to use connection store with latency display and manual reconnect button

---

## T017: Lobby Page - Table List Display

**Priority:** P0  
\***\*Dependencies:** T009, T010, T015

**Description:**  
Create the lobby page with a list of active tables.

**Acceptance Criteria:**

- [x] Fetch table list from API on page load
- [x] Display tables in a grid or list layout
- [x] Each table shows: format, host username, player count (2/4), status
- [x] Empty state when no tables available
- [x] Loading state while fetching
- [x] Refresh button to manually reload tables
- [x] Tables clickable to view details or join
- [x] Responsive design (mobile-friendly)

**Files Created:**

- `src/routes/(protected)/lobby/+page.svelte` - Full lobby page with filters, real-time updates, and responsive design
- `src/lib/components/TableCard.svelte` - Table card with format icons, status colors, and join overlay
- `src/lib/components/JoinTableModal.svelte` - Modal for joining tables with deck selection
- `src/lib/api/lobby.ts` - API functions (fetchTables, createTable, leaveTable, etc.)

---

## T018: Lobby Page - Real-Time Updates

**Priority:** P0  
\***\*Dependencies:** T017, T016

**Description:**  
Implement gRPC streaming to receive real-time lobby updates.

**Acceptance Criteria:**

- [x] Establish gRPC stream on lobby page mount
- [x] Listen for table created/updated/deleted events
- [x] Update table list in real-time without full refresh
- [x] Add new tables to list with animation
- [x] Remove closed tables from list
- [x] Update existing table info (player count, status)
- [x] Close stream on page unmount
- [x] Handle stream errors and reconnection

**Files Created:**

- `src/lib/services/lobby-updates.ts` - Lobby updates service with WebSocket subscriptions for TABLE_WAITING, JOINED_TABLE events
- `src/lib/stores/websocket.ts` - WebSocket store with exponential backoff reconnection (max 10 attempts)

**Files Modified:**

- `src/routes/(protected)/lobby/+page.svelte` - Integrated real-time updates with `connectWebSocket()` on mount and cleanup on destroy

---

## T019: Lobby Page - Filters and Search

**Priority:** P0  
\***\*Dependencies:** T017

**Description:**  
Add filtering and search controls to the lobby page.

**Acceptance Criteria:**

- [x] Format dropdown filter (All, Standard, Commander, Modern, etc.)
- [x] "Open tables only" toggle checkbox
- [x] Search by host username (text input)
- [x] Filters applied client-side (no API call)
- [x] Filter state persists across table updates
- [x] Clear filters button
- [x] Show filtered count vs total count

**Implementation:** Filters integrated into `src/routes/(protected)/lobby/+page.svelte` with `$derived` computed `filteredTables`

---

## T020: Lobby Page - Online Players Display

**Priority:** P0  
\***\*Dependencies:** T017

**Description:**  
Display count and list of online players in the lobby.

**Acceptance Criteria:**

- [x] Show online player count in header
- [x] Collapsible sidebar or section with player list
- [x] Player list shows usernames
- [x] Online status indicator (green dot)
- [x] Real-time updates when players join/leave
- [x] Scroll if player list exceeds height
- [x] Show "You" indicator for current user

**Files Created:**

- `src/lib/components/OnlinePlayersList.svelte` - Collapsible player list with green status dots and "You" badge
- `src/lib/types/player.ts` - OnlinePlayer type definitions

**Files Modified:**

- `src/routes/(protected)/lobby/+page.svelte` - Integrated OnlinePlayersList in left sidebar with polling updates

---

## T021: Create Table Modal - Basic Structure

**Priority:** P0  
\***\*Dependencies:** T012, T017

**Description:**  
Create a modal dialog for creating a new table with basic options.

**Acceptance Criteria:**

- [x] Opens when "Create Table" button clicked
- [x] Format selector dropdown (Standard, Commander, Modern, etc.)
- [x] Player count selector (2, 3, 4, etc.)
- [x] Optional password field with show/hide toggle
- [x] Form validation (format required)
- [x] "Create & Join" submit button
- [x] Cancel button to close modal
- [x] Loading state during submission
- [x] Error handling for failed creation

**Files Created:**

- `src/lib/components/CreateTableModal.svelte` - Complete create table modal with format, player count, password, deck selection, and validation

---

## T022: Create Table Modal - Deck Selection

**Priority:** P0 \***\*Dependencies:** T021

**Description:**
Add deck selection to the create table modal.

**Acceptance Criteria:**

- [x] Dropdown showing user's saved decks for selected format
- [x] "Upload new deck" option if no decks available
- [x] Opens deck upload modal (separate component)
- [x] Validates deck is for selected format
- [x] Shows deck card count (60 cards, etc.)
- [x] Deck selection required before creating table
- [x] Fetches user decks on modal open

---

## T023: Lobby Chat Component

**Priority:** P0 \***\*Dependencies:** T017

**Description:**
Create a chat panel for the lobby with message display and input.

**Acceptance Criteria:**

- [x] Chat panel (side panel or bottom)
- [x] Message list with auto-scroll to bottom
- [x] Message input field with send button
- [x] Send on Enter key press
- [x] Display username and timestamp for each message
- [x] System messages styled differently (gray, italic)
- [x] Load last 50 messages on mount
- [x] Real-time message updates via gRPC stream
- [x] Scroll to bottom button when scrolled up
- [x] Empty state when no messages

---

## T024: Chat - Whisper Command

**Priority:** P0 \***\*Dependencies:** T023

**Description:**
Implement whisper functionality for private messages in chat.

**Acceptance Criteria:**

- [x] Detect `/w username message` format
- [x] Parse username and message
- [x] Send whisper via API
- [x] Display whispers in italic with "(whisper)" prefix
- [x] Show sent whispers as "To [username]: message"
- [x] Show received whispers as "From [username]: message"
- [x] Whispers in different color (muted purple/blue)
- [x] Error if username not found
- [x] Cannot whisper to self

---

## T025: Chat - Rate Limiting Feedback

**Priority:** P0 \***\*Dependencies:** T023

**Description:**
Implement client-side rate limiting feedback for chat messages.

**Acceptance Criteria:**

- [x] Track message count per time window (10 messages per 60 seconds)
- [x] Disable send button when limit reached
- [x] Show countdown timer when rate limited
- [x] Display warning message "Sending too fast, wait X seconds"
- [x] Reset counter after time window expires
- [x] Visual feedback (red text, disabled button)

---

## T026: Table View - Pre-Game Lobby Component

**Priority:** P0 \***\*Dependencies:** T009, T015

**Description:**
Create the table lobby view where players wait before game starts.

**Acceptance Criteria:**

- [x] Display table info header (format, host, table ID)
- [x] Show player list with slots (occupied and empty)
- [x] Each player shows: username, ready status indicator
- [x] Local player has "Ready" toggle button
- [x] Empty slots show "Waiting for player..."
- [x] Host sees "Start Game" button (enabled when all ready)
- [x] Non-host players cannot start game
- [x] Real-time updates when players join/leave/ready
- [x] "Leave Table" button with confirmation
- [x] Password indicator if table is password-protected

---

## T027: Table View - Host Controls

**Priority:** P0 \***\*Dependencies:** T026

**Description:**
Add host-specific controls to the table lobby.

**Acceptance Criteria:**

- [x] Only visible to table host
- [x] "Kick Player" button next to each non-host player
- [x] Confirmation dialog before kicking
- [x] API call to kick player
- [x] Player removed from table immediately
- [x] Toast notification on kick success/failure
- [x] Cannot kick self

---

## T028: Table View - Table Chat

**Priority:** P0 \***\*Dependencies:** T026, T023

**Description:**
Add table-specific chat panel to the table lobby.

**Acceptance Criteria:**

- [x] Reuse Chat component from lobby
- [x] Table chat scope (only players at table)
- [x] Separate chat stream per table
- [x] Chat persists when players join/leave
- [x] System messages for player join/leave events
- [x] Clear chat when table closes

**Implementation:** See `TABLE_VIEW_FEATURES.md` for complete details

---

## T029: Table View - Game Start Countdown

**Priority:** P0 \***\*Dependencies:** T026

**Description:**
Add countdown timer before game starts when all players ready.

**Acceptance Criteria:**

- [x] When host clicks "Start Game", show 5 second countdown
- [x] Display countdown overlay (modal or banner)
- [x] Count down: 5... 4... 3... 2... 1... Starting!
- [x] Cancel countdown if player unreadies
- [x] Navigate to game view after countdown completes
- [x] Host can cancel countdown

**Implementation:** See `TABLE_VIEW_FEATURES.md` for complete details

---

## T030: Table View - Leave Table Confirmation

**Priority:** P0 \***\*Dependencies:** T026, T013

**Description:**
Add confirmation dialog when player tries to leave table.

**Acceptance Criteria:**

- [x] "Leave Table" button triggers confirmation
- [x] Dialog shows warning: "Are you sure you want to leave?"
- [x] Confirm button sends leave request to API
- [x] On success, navigate back to lobby
- [x] Show toast notification on error
- [x] If host leaves, table closes for all players

**Implementation:** Already implemented in T026-T027

---

## T031: Deck Management Page - Deck List Display

**Priority:** P0  
\***\*Dependencies:** T009, T015

**Description:**  
Create the deck management page showing user's saved decks.

**Acceptance Criteria:**

- [x] Fetch user's decks on page load
- [x] Display decks in a grid layout
- [x] Each deck card shows: format, card count, last modified date
- [x] Empty state when no decks saved
- [x] Loading state while fetching
- [x] "Upload New Deck" button
- [x] Decks grouped by format (optional)
- [x] Click deck to view details

---

## T032: Deck Upload Modal - Text Import

**Priority:** P0 \***\*Dependencies:** T031, T012

**Description:**
Create a modal for uploading/importing decks via text.

**Acceptance Criteria:**

- [x] Large text area for deck list input
- [x] Format selector dropdown
- [x] Parse deck list format: `4 Lightning Bolt`
- [x] Show card count as user types (real-time)
- [x] Display validation errors:
  - Invalid card names
  - Wrong deck size (not 60 cards)
  - Illegal cards for format
- [x] "Clear" button to reset text area
- [x] "Save Deck" button (disabled if invalid)
- [x] Loading state during save
- [x] Success toast on save
- [x] Close modal and refresh deck list on success

**Files Created:**

- `mage-client-web/src/lib/components/DeckUploadModal.svelte` - Complete upload modal with text parsing and server integration

---

## T033: Deck Upload Modal - Validation Display

**Priority:** P0 \***\*Dependencies:** T032

**Description:**
Add inline validation feedback for deck upload.

**Acceptance Criteria:**

- [x] Real-time validation as user types
- [x] Show validation errors in a list below text area
- [x] Error types:
  - "Invalid card: [Card Name]"
  - "Deck must be 60 cards (currently: X)"
  - "[Card Name] is not legal in [Format]"
  - "Too many copies of [Card Name] (max 4)"
- [x] Errors highlighted in red
- [x] Green checkmark when deck is valid
- [x] Disable save button while errors present
- [x] Commander format validation (1 commander + 99 deck = 100 total)
- [x] Format-specific examples that append instead of replace
- [x] 4-of rule enforcement (except Commander and basic lands)
- [x] Comprehensive test suite (8 passing tests)

**Files Created/Modified:**

- `mage-client-web/src/lib/components/DeckUploadModal.svelte` - Enhanced validation with Commander rules, 4-of enforcement, categorized errors
- `mage-client-web/src/lib/components/__tests__/DeckUploadModal.test.ts` - Complete test suite for validation logic

---

## T034: Deck Viewer Component

**Priority:** P0 \***\*Dependencies:** T031

**Description:**
Create a component to view deck details and card list.

**Acceptance Criteria:**

- [x] Display deck name and format
- [x] Show total card count and breakdown (creatures, instants, etc.)
- [x] Group cards by type (Creatures, Instants, Sorceries, etc.)
- [x] Display card quantity and name
- [x] Mana curve visualization (bar chart)
- [x] Color distribution (pie chart or bar)
- [x] "Export" button to download as text
- [x] "Delete" button with confirmation
- [ ] "Edit" button to modify deck (future)

**Files Created/Modified:**

- `mage-client-web/src/lib/components/DeckViewer.svelte` - Complete deck viewer component with stats, visualizations, and card list
- `mage-client-web/src/routes/(protected)/decks/[id]/+page.svelte` - Dynamic route for viewing individual decks
- `mage-client-web/src/routes/(protected)/decks/+page.svelte` - Updated to navigate to deck viewer on click

---

## T035: Deck Deletion with Confirmation

**Priority:** P0 \***\*Dependencies:** T034, T013

**Description:**
Add deck deletion functionality with confirmation dialog.

**Acceptance Criteria:**

- [x] "Delete" button in deck viewer
- [x] Confirmation dialog: "Delete [Deck Name]?"
- [x] Warning: "This action cannot be undone"
- [x] API call to delete deck
- [x] Remove deck from list on success
- [x] Show toast notification
- [x] Handle errors gracefully

**Implementation:** Already implemented in `src/routes/(protected)/decks/[id]/+page.svelte`

---

## T036: User Profile Page - Basic Info Display

**Priority:** P0 \***\*Dependencies:** T009, T015

**Description:**
Create user profile page displaying basic account information and stats.

**Acceptance Criteria:**

- [x] Display username and email
- [x] Show join date / account created date
- [x] Display stats:
  - Total games played
  - Wins / Losses
  - Win rate percentage
  - Quit ratio (prominently displayed)
- [x] Stats cards with icons
- [x] Loading state while fetching profile
- [x] Error state if fetch fails

**Files Created:**

- `src/lib/api/profile.ts` - Profile API functions
- `src/lib/types/profile.ts` - Profile type definitions

**Files Modified:**

- `src/routes/(protected)/profile/+page.svelte` - Complete profile page implementation

---

## T037: User Profile - Change Password Form

**Priority:** P0 \***\*Dependencies:** T036

**Description:**
Add change password functionality to user profile.

**Acceptance Criteria:**

- [x] Form with three fields: current password, new password, confirm new password
- [x] Client-side validation:
  - Current password required
  - New password min 8 characters
  - Confirm password matches new password
- [x] "Change Password" submit button
- [x] Loading state during submission
- [x] Success toast on password changed
- [x] Error messages for:
  - Wrong current password
  - Server errors
- [x] Clear form on success

## **Implementation:** Integrated into profile page (`src/routes/(protected)/profile/+page.svelte`)

## T038: User Profile - Recent Match History

**Priority:** P0 \***\*Dependencies:** T036

**Description:**
Display list of recent matches on profile page.

**Acceptance Criteria:**

- [x] Show last 10 matches
- [x] Each match displays:
  - Opponent username
  - Format
  - Result (Win/Loss/Draw)
  - Date/time
- [x] Result badge colored (green win, red loss, gray draw)
- [x] Sorted by most recent first
- [x] Empty state if no matches played
- [ ] Link to match details (future)

**Implementation:** Integrated into profile page (`src/routes/(protected)/profile/+page.svelte`)

---

## T039: Game View - Basic Layout Structure

**Priority:** P0 \***\*Dependencies:** T009

**Description:**
Create the basic game view layout with zones for opponent, battlefield, and player.

**Acceptance Criteria:**

- [x] Three main sections:
  - Top: Opponent area (hand placeholder, life, library count)
  - Middle: Battlefield (shared zone)
  - Bottom: Player area (hand, life, library count)
- [x] Side panel for game chat
- [x] Game info header (format, turn count, current phase)
- [x] "Concede" button with prominent placement
- [x] Responsive layout (stack vertically on mobile)
- [x] Fixed positions for zones (no scrolling inside zones)

**Files Modified:**

- `src/routes/(protected)/game/[id]/+page.svelte` - Complete game view layout with dark theme, responsive design, and placeholder game state

---

## T040: Game View - Game Info Header

**Priority:** P0  
\***\*Dependencies:** T039

**Description:**  
Create the game info header showing game state and current phase.

**Acceptance Criteria:**

- [ ] Display format name (available in state, not shown in header)
- [x] Show current turn number
- [x] Display active player indicator
- [x] Show current phase (Upkeep, Main, Combat, etc.)
- [x] Phase highlighted/animated during transition
- [x] Priority indicator (whose turn to act)
- [ ] Timer display (if game has timer)
- [x] Compact design (doesn't take too much space)

**Files Created:**

- `src/lib/components/game/GameHeader.svelte` - Game header with turn info, phase track, priority indicator, and concede button

**Notes:** Format name available in game state but not displayed in header. Timer not yet implemented.

---

## T041: Game View - Player Hand Component

**Priority:** P0 \***\*Dependencies:** T039

**Description:**
Create component to display player's hand with draggable cards.

**Acceptance Criteria:**

- [x] Display cards in a horizontal row
- [x] Cards overlap slightly (fan layout)
- [x] Hover to preview card (enlarge and lift)
- [x] Click to select card (highlight border)
- [x] Multi-select with Shift+Click
- [ ] Drag card to battlefield or other zone (future enhancement)
- [x] Show card count badge
- [x] Responsive (stack on mobile)
- [x] Empty state when hand is empty

**Files Created:**

- `src/lib/components/game/PlayerHand.svelte` - Player hand component with card selection and multi-select

---

## T042: Game View - Card Component with Hover Preview

**Priority:** P0 \***\*Dependencies:** T041

**Description:**
Create the card component with hover preview and tooltips.

**Acceptance Criteria:**

- [x] Display card image (placeholder if not loaded)
- [x] Hover shows enlarged card preview
- [x] Preview positioned to not go off-screen
- [x] Show card name as tooltip
- [x] Display mana cost in corner
- [x] Show tapped state (90° rotation animation)
- [x] Display counters (+1/+1, etc.) as badges
- [x] Selection highlight (border glow)
- [x] Loading state for card images
- [x] Fallback for missing images

**Files Created:**

- `src/lib/components/game/Card.svelte` - Reusable card component with hover preview, tapped state, counters, and mana symbols
- `src/lib/types/game.ts` - Game state type definitions (GameCard, GameState, GamePlayer, etc.)

---

## T043: Game View - Battlefield Component

**Priority:** P0  
\***\*Dependencies:** T039, T042

**Description:**  
Create the battlefield zone where permanents are displayed.

**Acceptance Criteria:**

- [x] Grid layout for permanents (automatic positioning)
- [ ] Separate sections for lands and nonlands
- [ ] Cards draggable within battlefield (reorder)
- [x] Support for tapped cards (rotated)
- [ ] Grouping by type (optional)
- [ ] Zoom in/out controls (if many cards)
- [x] Empty state when battlefield is empty
- [x] Hover to preview any card
- [x] Click to select card for actions

**Implementation:** Battlefield is inline in game page (`src/routes/(protected)/game/[id]/+page.svelte`) with flex-wrap layout. Land/nonland separation and drag-drop are future enhancements.

---

## T044: Game View - Opponent Hand Placeholder

**Priority:** P0  
\***\*Dependencies:** T039, T042

**Description:**  
Create component to display opponent's hand as card backs.

**Acceptance Criteria:**

- [ ] Show card backs (not visible)
- [x] Display count of cards in hand
- [ ] Horizontal row layout
- [x] No hover preview (cards hidden)
- [x] Tooltip shows "Opponent's hand (X cards)"
- [x] Empty state when opponent hand is empty

**Implementation:** Opponent hand count shown in `OpponentPanel.svelte` header stats. Card back visualization is a future enhancement.

---

## T045: Game View - Life Total Display

**Priority:** P0  
\***\*Dependencies:** T039

**Description:**  
Create life total display component for both players.

**Acceptance Criteria:**

- [ ] Large, prominent life number
- [ ] Color-coded (green high, yellow medium, red low)
- [ ] Life change animation (flash on change)
- [ ] Show life change delta (e.g., "-3" or "+5")
- [x] Positioned in player/opponent zones
- [x] Accessible (ARIA labels)
- [ ] Optional: life history graph (last 10 changes)

**Implementation:** Life is displayed inline in player stats bar and OpponentPanel. Animations and color-coding are future enhancements.

---

## T046: Game View - Graveyard Display

**Priority:** P0 \***\*Dependencies:** T039, T042

**Description:**
Create graveyard zone component for both players.

**Acceptance Criteria:**

- [x] Shows top card of graveyard (if any)
- [x] Card count badge
- [x] Click to expand and view all cards
- [x] Modal or side panel shows full graveyard
- [x] Cards in graveyard are hoverable (preview)
- [x] Close button to collapse graveyard view
- [x] Empty state when graveyard is empty
- [x] Separate graveyards for each player

**Files Created:**

- `src/lib/components/game/Graveyard.svelte` - Complete graveyard component with modal viewer and card selection

---

## T047: Game View - Exile Zone Display

**Priority:** P0 \***\*Dependencies:** T046

**Description:**
Create exile zone component (similar to graveyard).

**Acceptance Criteria:**

- [x] Shows exiled cards count
- [x] Click to expand and view all cards
- [x] Modal/panel shows exiled cards
- [x] Cards are hoverable (preview)
- [x] Empty state when no exiled cards
- [x] Different visual style from graveyard (purple theme with sparkle animation)

**Files Created:**

- `src/lib/components/game/ExileZone.svelte` - Exile zone component with purple theme, sparkle effects, and modal viewer

---

## T048: Game View - Library Counter

**Priority:** P0  
\***\*Dependencies:** T039

**Description:**  
Create library (deck) counter display showing remaining cards.

**Acceptance Criteria:**

- [ ] Shows card back image
- [x] Badge with card count
- [x] Positioned in player/opponent zones
- [x] Updates in real-time as cards drawn
- [ ] Warning state when low cards (< 5)
- [x] Empty state when library is empty

**Implementation:** Library count shown inline in player stats and OpponentPanel. Card back visual and warning state are future enhancements.

---

## T049: Game View - Mana Pool Display

**Priority:** P0 \***\*Dependencies:** T039

**Description:**
Create mana pool display showing available mana.

**Acceptance Criteria:**

- [x] Shows mana symbols with counts (W, U, B, R, G, C)
- [ ] Animated when mana added/spent (future enhancement)
- [x] Positioned near player hand
- [x] Compact display (mana icons + numbers)
- [x] Updates in real-time
- [x] Empty state when no mana available
- [x] Accessible (ARIA labels)

**Files Created:**

- `src/lib/components/game/ManaPool.svelte` - Complete mana pool with clickable colored orbs, three size variants, and empty state

---

## T050: Game View - Game Chat Panel

**Priority:** P0 \***\*Dependencies:** T039, T023

**Description:**
Add in-game chat panel for player communication and game events.

**Acceptance Criteria:**

- [x] Created GameChat component based on existing chat patterns
- [x] Positioned on right side sidebar
- [x] Game-specific chat scope (uses gameId)
- [x] Show game events in chat as system messages
- [x] Collapsible panel (hide/expand to maximize game view)
- [x] System messages for game events with distinct styling
- [x] Rate limiting (10 messages per 60 seconds)
- [x] Scroll to bottom functionality
- [x] Dark theme matching game view
- [x] Export addGameEvent function for programmatic event adding

**Files Created:**

- `src/lib/components/game/GameChat.svelte` - Complete game chat with collapse, rate limiting, and event support

**Files Modified:**

- `src/routes/(protected)/game/[id]/+page.svelte` - Integrated GameChat component, replaced game log sidebar

---

## T051: Game View - Action Log

**Priority:** P0 \***\*Dependencies:** T039

**Description:**
Create scrollable action log showing game events.

**Acceptance Criteria:**

- [x] List of game actions (played card, attacked, etc.)
- [x] Timestamped entries
- [x] Color-coded by player
- [x] Icons for action types (sword for attack, etc.)
- [x] Auto-scroll to latest action
- [x] Scrollable to view history
- [x] Positioned in side panel or bottom
- [x] Collapsible to save space

**Files Created:**

- `src/lib/components/game/ActionLog.svelte` - Main action log component with collapsible left sidebar
- `src/lib/components/game/ActionLogItem.svelte` - Individual action log entry with timestamp, icon, and player color
- Updated `src/lib/types/game.ts` - Added ActionLogEntry and ActionType types
- Updated `src/routes/(protected)/game/[id]/+page.svelte` - Integrated ActionLog component

---

## T052: Game View - Priority Indicator

**Priority:** P0 \***\*Dependencies:** T039

**Description:**
Create visual indicator showing which player has priority.

**Acceptance Criteria:**

- [x] Shows whose turn it is to act
- [x] Text indicator: "Your Priority", "Opponent's Priority", or "Waiting..."
- [x] Animated pulse effect when player has priority
- [x] Updates in real-time with game state
- [x] Different visual for priority vs. just active player
- [x] Icon indicators (⚡ for priority, ⏳ for active, ⏸️ for waiting)
- [x] Priority hint text when player has priority

**Files Created:**

- `src/lib/components/game/PriorityIndicator.svelte` - Complete priority indicator with animations and states

**Files Modified:**

- `src/routes/(protected)/game/[id]/+page.svelte` - Integrated priority indicator with phase display

---

## T053: Game View - Game Actions Panel

**Priority:** P0 \***\*Dependencies:** T052

**Description:**
Add game actions panel with Pass Priority, Cast Spell, and Activate Ability buttons.

**Acceptance Criteria:**

- [x] Pass Priority button enabled only when player has priority
- [x] Keyboard shortcuts (Space for pass, C for cast, A for activate)
- [x] Loading state while waiting for response
- [x] Disabled state when waiting for opponent
- [x] Visual feedback for available actions
- [x] Priority badge showing current priority status
- [x] Waiting message when opponent has priority
- [x] Cast Spell and Activate Ability buttons (placeholders)
- [x] Responsive design (mobile-friendly button layout)
- [x] Shortcut key indicators on buttons

**Files Created:**

- `src/lib/components/game/GameActionsPanel.svelte` - Complete actions panel with keyboard shortcuts

**Files Modified:**

- `src/routes/(protected)/game/[id]/+page.svelte` - Replaced old button layout with GameActionsPanel

---

## T055: Game State Store

**Priority:** P0  
**Dependencies:** T039, T015  
**Status:** ✅ Complete

**Description:**  
Svelte store managing game state synchronized via WebSocket. Server sends full `GameView` on each update (source of truth per `ARCHITECTURE.md`).

**Implementation:**

- `src/lib/stores/game.ts` - Main store with WebSocket subscriptions for 16+ event types (GAME_INIT, GAME_UPDATE, GAME_TARGET, GAME_OVER, etc.)
- `src/lib/generated/mage/v1/models.ts` - Server-generated types (`GameView`, `PlayerView`, `CardView`, `ManaPoolView`)
- 25+ derived stores: `battlefield`, `stack`, `exile`, `myHand`, `hasPriority`, `isMyTurn`, `currentPhase`, `pendingPrompt`, etc.
- Prompt handling for user interactions (target selection, ability choice, mana payment)
- Game-over state, error handling, card selection UI state

---

## T056: Game State Synchronization

**Priority:** P0  
**Dependencies:** T055, T016  
**Status:** ✅ Complete

**Description:**  
Real-time game state sync. Architecture uses gRPC-Web (HTTP POST) for actions, WebSocket for server push events.

**Implementation:**

- `src/lib/api/game.ts` - 15+ game actions: `passPriority`, `playLand`, `sendPlayerUUID`, `mulligan`, `concedeGame`, `sendSpecialAction`, etc.
- `src/routes/(protected)/game/[id]/+page.svelte` - WebSocket connect on mount, `gameStore.initGame()`, cleanup via `gameStore.reset()`
- Actions sent via gRPC-Web → server processes → WebSocket pushes GAME_UPDATE → store updates reactively
- Reconnection overlay and state restore tracked separately in T070-T073

---

## T057: Card Interaction - Click to Select

**Priority:** P0  
**Dependencies:** T042, T055  
**Status:** ✅ Complete

**Description:**  
Card selection via click with visual feedback and multi-select support.

**Implementation:**

- `src/lib/stores/game.ts` - `toggleCardSelection()`, `clearSelection()`, `selectedCardIds` state
- `src/lib/components/game/Card.svelte` - `isSelected` prop → CSS `'selected'` class with border glow
- `src/lib/components/game/PlayerHand.svelte` - Shift+Click multi-select via `event?.shiftKey`
- Game page clears selection after actions, uses selection for action buttons

---

## T058: Card Interaction - Drag and Drop (optional)

**Priority:** P3  
\***\*Dependencies:** T057

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

## T059: Card Interaction - Target Selection Mode ✅ COMPLETED

**Priority:** P0  
\***\*Dependencies:** T057

**Description:**  
Implement target selection mode for spells and abilities.

**Acceptance Criteria:**

- [x] Enter targeting mode when spell/ability requires targets
- [x] Valid targets highlighted (glow effect)
- [x] Invalid targets grayed out
- [x] Cursor changes to crosshair
- [x] Click to select target
- [x] Show selected target with indicator
- [x] Cancel targeting with ESC or right-click
- [x] Confirm target selection (button or Enter key)
- [x] Support multi-target (select multiple)

**Files Created:**

- `src/lib/components/game/TargetingMode.svelte` - Target selection overlay UI
- `src/lib/stores/game-targeting.ts` - Target selection state management

**Files Modified:**

- `src/lib/components/game/Card.svelte` - Added targeting visual states
- `src/lib/components/game/PlayerHand.svelte` - Added targeting integration
- `src/lib/components/game/OpponentPanel.svelte` - Added targeting integration
- `src/routes/(protected)/game/[id]/+page.svelte` - Added targeting mode handling
- `mage-server-go/internal/game/mage_engine.go` - Added target selection server logic
- `mage-server-go/internal/server/grpc.go` - Added GAME_TARGET event sending

**Implementation Notes:**

- Full client-server target selection flow implemented
- Server sends GAME_TARGET event with valid targets when spell/ability requires targets
- Client enters targeting mode with visual feedback (glow, grayscale, cursor)
- Player can select/deselect targets with click, cancel with ESC/right-click
- Confirmation via Enter key or button when min targets selected
- Multi-target support with configurable min/max targets
- Target validation happens both client-side and server-side

---

## T060: Card Interaction - Tap/Untap Animation

**Priority:** P0  
\***\*Dependencies:** T042

**Description:**  
Add smooth rotation animation for tapping/untapping cards.

**Acceptance Criteria:**

- [x] Smooth 90° rotation for tap
- [x] Smooth -90° rotation for untap
- [x] Animation duration: 200ms
- [x] CSS transform for performance
- [x] Visual state persists after animation
- [x] Works on mobile (touch)

**Files Modified:**

- `src/lib/components/game/Card.svelte` - Added `tap-rotate-in` and `tap-rotate-out` keyframe animations with overshoot effect, glow pulse, and touch device optimization

---

## T061: Game Actions - Play Card from Hand

**Priority:** P0  
\***\*Dependencies:** T058, T055

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
\***\*Dependencies:** T057, T059

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
\***\*Dependencies:** T057, T055

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
\***\*Dependencies:** T063

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
\***\*Dependencies:** T064

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
\***\*Dependencies:** T055

**Description:**  
Create UI for paying mana costs when casting spells or activating abilities.

**Acceptance Criteria:**

- [x] Modal shows when mana payment required
- [x] Display cost to pay (e.g., "2RR")
- [x] List available mana sources (server provides manaOptions)
- [x] Click mana type to pay
- [x] Show available mana counts per color
- [x] Disable unavailable mana colors
- [x] "Cancel" button to cancel action
- [x] API call to pay mana (sendPlayerManaType)
- [x] X mana selection UI (XManaSelector component)

**Files Created:**

- `src/lib/components/game/ManaPayment.svelte` - Mana payment modal with color selection grid
- `src/lib/components/game/XManaSelector.svelte` - X mana value selector with number picker and slider

**Files Modified:**

- `src/routes/(protected)/game/[id]/+page.svelte` - Integrated ManaPayment and XManaSelector for 'mana' and 'xmana' prompts

---

## T067: Game Actions - Choice Dialog

**Priority:** P0  
\***\*Dependencies:** T012

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
\***\*Dependencies:** T012

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
\***\*Dependencies:** T055

**Description:**  
Create visualization for the stack showing spells and abilities.

**Acceptance Criteria:**

- [x] Vertical stack display (last on top via flex-direction: column-reverse)
- [x] Each item shows card/ability name and controller
- [x] Arrows indicate order of resolution ("Bottom ↓ resolves first ↓ Top")
- [x] Highlight current item resolving (top item with yellow border)
- [x] Empty state when stack is empty
- [x] Positioned in modal overlay (accessible via floating button)
- [x] Updates in real-time via game store subscription
- [x] Click item to select/view details

**Files Created:**

- `src/lib/components/game/Stack.svelte` - Complete stack visualization with spell/ability display, position numbers, target info, and resolving indicator

**Files Modified:**

- `src/routes/(protected)/game/[id]/+page.svelte` - Integrated Stack component in overlay modal

---

## T070: Reconnection - Detect Disconnect

**Priority:** P0  
\***\*Dependencies:** T016, T056

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
\***\*Dependencies:** T070

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
\***\*Dependencies:** T071, T055

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
\***\*Dependencies:** T070

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
\***\*Dependencies:** T055

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
\***\*Dependencies:** T055

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
\***\*Dependencies:** T055

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
\***\*Dependencies:** T001

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
\***\*Dependencies:** T015

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
\***\*Dependencies:** T001

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
\***\*Dependencies:** T001

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
\***\*Dependencies:** T009

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
\***\*Dependencies:** T081

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
\***\*Dependencies:** T081

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
\***\*Dependencies:** T081

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
\***\*Dependencies:** Multiple (most UI components)

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
\***\*Dependencies:** Multiple

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
\***\*Dependencies:** Multiple

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
\***\*Dependencies:** T005

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
\***\*Dependencies:** T042

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
\***\*Dependencies:** T017, T038

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
\***\*Dependencies:** T001

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
\***\*Dependencies:** T001

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
\***\*Dependencies:** Multiple

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
\***\*Dependencies:** T082

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
- [ ] TTS for "Match Started", "Your Game Can Start" etc

**Files to Create:**

- `src/lib/utils/sound.ts`
- `static/sounds/card-play.mp3`
- `static/sounds/damage.mp3`
- etc.

---

## T095: Animations - Card Movement

**Priority:** P3  
\***\*Dependencies:** T061

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
\***\*Dependencies:** T001

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

## T097: UX and UI Rework

With everything set up, we should revisist our styling and make this a professional, modern application.
Validate major UX flows. Optimize our UI work to be minimal, but delightful. It should be a breeze to play.

## T098: Production Build and Deployment

**Priority:** P0  
\***\*Dependencies:** T001

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
