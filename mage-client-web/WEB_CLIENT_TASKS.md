# XMage Web Client TODO List

Svelte-based web client for the Go XMage port. Organized by priority and dependencies.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

## Project Setup & Infrastructure (P0)
- [ ] Initialize SvelteKit project with TypeScript
- [ ] Configure Vite with appropriate build settings
- [ ] Set up Tailwind CSS (or CSS framework of choice)
- [ ] Configure gRPC-Web client library (@improbable-eng/grpc-web or @grpc/grpc-js)
- [ ] Generate TypeScript types from proto files (protoc-gen-ts)
- [ ] Set up environment configuration (.env for API endpoints)
- [ ] Configure ESLint and Prettier
- [ ] Set up basic routing structure
- [ ] Add favicon and app metadata
- [ ] Configure build for production deployment

## Authentication & Session Management (P0)
- [ ] Create login page (username/password)
- [ ] Implement login form with validation
- [ ] Create registration page (username, email, password)
- [ ] Add registration form with client-side validation
- [ ] Implement JWT token storage (localStorage or sessionStorage)
- [ ] Create auth store (Svelte store for user session)
- [ ] Add "Remember me" checkbox for persistent login
- [ ] Implement logout functionality
- [ ] Create guest login option (anonymous play)
- [ ] Add auth interceptor for gRPC calls (inject JWT)
- [ ] Implement token refresh logic
- [ ] Add session expiry handling (redirect to login)
- [ ] Create auth guard for protected routes
- [ ] Add loading states during authentication

## gRPC Client Setup (P0)
- [ ] Initialize gRPC-Web transport layer
- [ ] Create gRPC service client instances
- [ ] Implement error handling wrapper for gRPC calls
- [ ] Add reconnection logic with exponential backoff
- [ ] Create gRPC streaming handler for bidirectional updates
- [ ] Implement connection state tracking (connected, disconnected, reconnecting)
- [ ] Add timeout handling for requests
- [ ] Create debug logging for gRPC messages (dev mode only)
- [ ] Implement request/response interceptors
- [ ] Add network status detection (online/offline)

## Layout & Navigation (P0)
- [ ] Create main app layout component
- [ ] Implement top navigation bar (logo, username, logout)
- [ ] Add connection status indicator (online/reconnecting/offline)
- [ ] Create sidebar navigation (Lobby, My Profile, Settings)
- [ ] Implement responsive layout (mobile, tablet, desktop)
- [ ] Add loading spinner component (global)
- [ ] Create toast notification system (success, error, info)
- [ ] Add modal dialog component (reusable)
- [ ] Implement confirmation dialog component
- [ ] Add page transition animations

## Lobby View (P0)
- [ ] Create lobby page component
- [ ] Implement table list display (scrollable grid/list)
- [ ] Add table card component (format, players, status, host)
- [ ] Show online player count
- [ ] Implement real-time table updates (gRPC streaming)
- [ ] Add table filtering controls (format dropdown, "open only" toggle)
- [ ] Create "Create Table" button with modal
- [ ] Implement table sorting (by created time, player count)
- [ ] Add refresh button (manual table list reload)
- [ ] Show "No tables available" empty state
- [ ] Add table join animation/transition
- [ ] Implement lobby chat panel (side panel or bottom)
- [ ] Add lobby connection error handling

## Create Table Modal (P0)
- [ ] Create table creation modal component
- [ ] Add format selection dropdown (Standard, Commander, Modern, etc.)
- [ ] Implement player count selector (2, 3, 4, etc.)
- [ ] Add optional password field (toggle visibility)
- [ ] Create deck selection dropdown (from saved decks)
- [ ] Add deck upload option (if no saved deck)
- [ ] Implement form validation (format required, deck required)
- [ ] Add "Create & Join" submit button
- [ ] Show loading state during table creation
- [ ] Handle creation errors (show error message)
- [ ] Close modal on success and navigate to table view

## Table View (Pre-Game) (P0)
- [ ] Create table lobby component (waiting room)
- [ ] Display table info header (format, host, player count)
- [ ] Show player list with ready status indicators
- [ ] Implement player avatars/placeholders
- [ ] Add local player "Ready/Unready" toggle button
- [ ] Show "Waiting for players..." status
- [ ] Display "Start Game" button for host (when all ready)
- [ ] Implement table chat panel (right side or bottom)
- [ ] Add "Leave Table" button with confirmation
- [ ] Show host controls (kick player button)
- [ ] Display password indicator if table is private
- [ ] Implement real-time player join/leave updates
- [ ] Add countdown timer before game starts (5 seconds)
- [ ] Handle table transition to game (navigate to game view)

## Chat System (P0)
- [ ] Create reusable chat component
- [ ] Implement message list (scrollable, auto-scroll to bottom)
- [ ] Add message input field with send button
- [ ] Support Enter key to send message
- [ ] Display message timestamp
- [ ] Show username with each message
- [ ] Implement system messages styling (different color/format)
- [ ] Add message history (load last 50 messages)
- [ ] Implement whisper command (`/w username message`)
- [ ] Show whisper messages differently (italic, muted color)
- [ ] Add rate limiting feedback (disable send if limit hit)
- [ ] Implement chat scroll to bottom button (when scrolled up)
- [ ] Add "user is typing..." indicator (optional, low priority)
- [ ] Handle long messages (word wrap, max length)

## Deck Management (P0)
- [ ] Create "My Decks" page
- [ ] Display saved deck list (one per format)
- [ ] Add deck card component (format, card count, last modified)
- [ ] Implement deck upload modal (text area for deck list)
- [ ] Add deck format validation (before upload)
- [ ] Show validation errors (card legality, deck size, etc.)
- [ ] Create deck viewer/editor component
- [ ] Display deck list grouped by card type (Creatures, Instants, etc.)
- [ ] Show mana curve visualization (bar chart)
- [ ] Add card count display (60 cards, 15 sideboard, etc.)
- [ ] Implement deck deletion with confirmation
- [ ] Add deck download/export (text format)
- [ ] Show deck last modified date
- [ ] Handle empty state (no decks yet)

## Deck Upload & Validation (P0)
- [ ] Create deck import text area (support plain text lists)
- [ ] Add format selector for validation
- [ ] Implement client-side deck parsing (quantity, card name)
- [ ] Show real-time card count as user types
- [ ] Display validation errors inline (specific cards, total count)
- [ ] Add example deck link (show format example)
- [ ] Implement "Clear" button for text area
- [ ] Add "Save Deck" button (sends to server)
- [ ] Show loading state during upload
- [ ] Handle server validation errors (display message)
- [ ] Add success notification on save
- [ ] Clear form after successful save

## Game View - Basic Structure (P0)
- [ ] Create game page component (main game container)
- [ ] Implement game board layout (opponent area, battlefield, player area)
- [ ] Add game info header (format, turn count, timer)
- [ ] Create opponent hand placeholder (card back count)
- [ ] Display battlefield (shared zone for permanents)
- [ ] Show player hand (draggable cards)
- [ ] Implement graveyard display (both players)
- [ ] Add exile zone display (both players)
- [ ] Create library counter (cards remaining)
- [ ] Display life total for both players
- [ ] Add mana pool display (current available mana)
- [ ] Show phase indicator (Upkeep, Main, Combat, etc.)
- [ ] Implement game chat panel (side panel)
- [ ] Add "Concede" button with confirmation
- [ ] Create game action log (scrollable event list)

## Game View - Card Rendering (P0)
- [ ] Create card component (display card image)
- [ ] Implement card hover preview (enlarged view)
- [ ] Add card tooltip (show oracle text)
- [ ] Display tapped/untapped state (rotation)
- [ ] Show +1/+1 counters on creatures
- [ ] Implement card glow/highlight for selections
- [ ] Add card dragging (for playing to battlefield)
- [ ] Show card casting cost in corner
- [ ] Display card types (Creature, Instant, etc.)
- [ ] Implement double-faced card flip (transform display)
- [ ] Add card search/filter (find card in hand)
- [ ] Show card legality indicator (grayed out if unplayable)

## Game View - Player Interactions (P0)
- [ ] Implement priority system (wait for priority, show indicator)
- [ ] Add "Pass Priority" button
- [ ] Create target selection mode (click to target)
- [ ] Implement card selection (click to select, multi-select with shift)
- [ ] Add drag-and-drop for playing cards
- [ ] Create choice dialogs (modal for game choices)
- [ ] Implement number input for X costs
- [ ] Add mana payment interface (tap lands, mana abilities)
- [ ] Show available actions for selected card (Play, Activate, etc.)
- [ ] Implement combat damage assignment
- [ ] Add declare attackers interface
- [ ] Create declare blockers interface
- [ ] Show combat phase visualizations
- [ ] Implement stack visualization (cards on stack)

## Game State Synchronization (P0)
- [ ] Set up gRPC streaming for game updates
- [ ] Parse incoming game state messages
- [ ] Update Svelte stores with game state
- [ ] Implement optimistic UI updates (predict actions)
- [ ] Handle rollback on invalid actions
- [ ] Add animation for card movement (hand → battlefield)
- [ ] Implement smooth counter updates (life, +1/+1)
- [ ] Show opponent actions in real-time
- [ ] Add sound effects for game events (card played, damage dealt)
- [ ] Handle simultaneous updates (queue actions)

## Reconnection & Error Handling (P0)
- [ ] Detect disconnection (WebSocket/gRPC connection lost)
- [ ] Show reconnection overlay (modal with spinner)
- [ ] Implement automatic reconnection attempts
- [ ] Request game state snapshot on reconnect
- [ ] Restore game view from snapshot
- [ ] Show "You were disconnected" message
- [ ] Add manual reconnect button
- [ ] Handle timeout during reconnection (redirect to lobby)
- [ ] Implement AFK warning (show timer before auto-forfeit)
- [ ] Add "connection unstable" indicator
- [ ] Handle game ended while disconnected

## User Profile (P0)
- [ ] Create user profile page
- [ ] Display username and email
- [ ] Show basic stats (games played, wins, losses, win rate)
- [ ] Display quit ratio prominently
- [ ] Add "Change Password" form
- [ ] Implement password change validation
- [ ] Show recent match history (last 10 games)
- [ ] Display match results (opponent, format, result, date)
- [ ] Add logout button on profile
- [ ] Show account created date
- [ ] Implement simple settings section (future expansion)

## Settings & Preferences (P1)
- [ ] Create settings page
- [ ] Add audio settings (enable/disable sound effects)
- [ ] Implement volume slider for sounds
- [ ] Add graphics settings (card quality, animations on/off)
- [ ] Create game preferences (auto-pass priority, etc.)
- [ ] Add chat preferences (show timestamps, message size)
- [ ] Implement theme selection (light/dark mode)
- [ ] Add language selector (if supporting i18n)
- [ ] Show current version number
- [ ] Add "Save Settings" button
- [ ] Implement settings persistence (localStorage)

## Responsive Design (P1)
- [ ] Test and fix mobile layout (lobby, table, game)
- [ ] Add touch-friendly controls (larger buttons)
- [ ] Implement swipe gestures for mobile (optional)
- [ ] Create mobile-optimized card layout
- [ ] Add hamburger menu for mobile navigation
- [ ] Test tablet layout (landscape/portrait)
- [ ] Implement adaptive card sizes (scale to screen)
- [ ] Add mobile chat (collapsible panel)
- [ ] Test on iOS Safari and Chrome
- [ ] Test on Android Chrome

## Error Handling & User Feedback (P1)
- [ ] Create error boundary component
- [ ] Add global error handler (catch unhandled errors)
- [ ] Implement retry logic for failed requests
- [ ] Show friendly error messages (not technical jargon)
- [ ] Add error report button (send logs to server)
- [ ] Create 404 page (page not found)
- [ ] Add 500 error page (server error)
- [ ] Implement form validation feedback (inline errors)
- [ ] Add loading skeletons (during data fetch)
- [ ] Show progress indicators for long operations

## Performance Optimization (P2)
- [ ] Implement lazy loading for routes
- [ ] Add code splitting (separate chunks per route)
- [ ] Optimize card image loading (lazy load, CDN)
- [ ] Implement virtual scrolling for long lists (table list, chat)
- [ ] Add service worker for offline support
- [ ] Cache static assets (images, fonts)
- [ ] Optimize bundle size (tree shaking, minification)
- [ ] Add preloading for critical resources
- [ ] Implement debouncing for search/filter inputs
- [ ] Profile and optimize render performance

## Testing (P2)
- [ ] Set up Vitest for unit tests
- [ ] Add component tests for key components
- [ ] Implement integration tests for user flows
- [ ] Add E2E tests with Playwright
- [ ] Test authentication flow
- [ ] Test lobby and table creation
- [ ] Test game interactions
- [ ] Add accessibility tests (screen reader, keyboard nav)
- [ ] Test cross-browser compatibility
- [ ] Implement visual regression tests

## Accessibility (P2)
- [ ] Add ARIA labels to interactive elements
- [ ] Implement keyboard navigation (tab through UI)
- [ ] Add keyboard shortcuts (space to pass priority, etc.)
- [ ] Test with screen readers (NVDA, JAWS)
- [ ] Ensure color contrast meets WCAG standards
- [ ] Add focus indicators (visible focus rings)
- [ ] Implement skip to content link
- [ ] Add alt text for images and icons
- [ ] Test with keyboard-only navigation
- [ ] Add reduced motion option (disable animations)

## Documentation (P2)
- [ ] Write README with setup instructions
- [ ] Document project structure and conventions
- [ ] Add component documentation (Storybook or similar)
- [ ] Create developer onboarding guide
- [ ] Document gRPC integration patterns
- [ ] Write deployment guide
- [ ] Add troubleshooting section
- [ ] Document environment variables
- [ ] Create user guide (how to play)
- [ ] Add contribution guidelines

---

## Phase 2: Post-MVP Feature Packages

### Package A: Enhanced Deck Management
- [ ] Multiple saved decks per format
- [ ] Deck naming and tagging
- [ ] Import from MTGO/Arena formats
- [ ] Deck sharing (export code)
- [ ] Deck statistics display
- [ ] Deck archetypes/colors display

### Package B: Spectating
- [ ] Spectator mode UI (join as observer)
- [ ] Spectator list display
- [ ] Hide hidden information (hands) in spectator view
- [ ] Spectator chat (separate from players)
- [ ] Leave spectate button

### Package C: Matchmaking Queue
- [ ] Matchmaking queue page
- [ ] Queue join button (by format)
- [ ] Show queue position and wait time
- [ ] Match found modal (accept/decline)
- [ ] Queue leave button
- [ ] ELO display

### Package D: Friend System
- [ ] Friends list page
- [ ] Friend request system (send, accept, decline)
- [ ] Online status indicators
- [ ] Friend invite to table button
- [ ] Private messaging interface

### Package E: Match History & Replays
- [ ] Match history page (paginated)
- [ ] Filter by opponent, format, date
- [ ] Replay viewer (step through game)
- [ ] Replay controls (play, pause, step, speed)
- [ ] Share replay link

### Package F: Tournaments
- [ ] Tournament list page
- [ ] Tournament details view (brackets, standings)
- [ ] Tournament registration flow
- [ ] Round pairings display
- [ ] Tournament chat
