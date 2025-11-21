# Multiplayer & Lobby System Task Tracker

Comprehensive tracker for building the multiplayer, lobby, and social features for the Go Mage server. Modern gRPC-Web based system with real-time updates.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

## Session & Authentication (P0)
- [ ] Implement JWT-based authentication system (replace RMI session IDs)
- [ ] Add user registration with email verification
- [ ] Implement password hashing (bcrypt/argon2)
- [ ] Add session management with reconnection support
- [ ] Implement "guest mode" for anonymous play
- [ ] Add rate limiting for auth endpoints
- [ ] Implement user state tracking (online, in-game, idle)
- [ ] Add multi-device session handling

## User Management (P0)
- [ ] Create user profile system (username, email, avatar URL)
- [ ] Implement user statistics tracking (matches played, win rate, quit ratio)
- [ ] Add ELO rating system (general, constructed, limited)
- [ ] Implement quit ratio calculation and enforcement
- [ ] Add user preferences storage (default deck, notification settings)
- [ ] Track connection history (last login, session duration)
- [ ] Implement user account deletion/deactivation

## gRPC-Web API Design (P0)
- [ ] Define proto messages for lobby operations (CreateTable, JoinTable, LeaveTable)
- [ ] Define proto messages for game state synchronization
- [ ] Define proto messages for chat system (SendMessage, ChatHistory)
- [ ] Define proto messages for matchmaking (JoinQueue, LeaveQueue, MatchFound)
- [ ] Implement bidirectional streaming for real-time updates
- [ ] Add health check and version compatibility endpoints
- [ ] Implement proper error codes and status messages
- [ ] Add request/response logging and metrics
- [ ] Design proto messages for tournaments (Create, Join, Submit, Standings)
- [ ] Add proto messages for spectating and replays

## Lobby System (P0)
- [ ] Implement main lobby room with active tables list
- [ ] Add server-sent events for lobby updates (new tables, games starting)
- [ ] Implement lobby filtering (format, game type, open slots)
- [ ] Add lobby sorting (by created time, players, format)
- [ ] Display recent finished matches (last 25)
- [ ] Show online player count and list
- [ ] Add "find game" quick-match button
- [ ] Implement lobby view updates via gRPC streaming
- [ ] Add lobby announcement system (server messages)

## Table Management (P0)
- [ ] Implement table creation with configuration (format, player count, password)
- [ ] Add table browser with real-time updates
- [ ] Implement table joining with deck submission
- [ ] Add table lifecycle states (WAITING → STARTING → ACTIVE → FINISHED)
- [ ] Implement player kick/invite system for table hosts
- [ ] Add table chat for pre-game coordination
- [ ] Implement table password protection for private games
- [ ] Add table limits per user (prevent spam: max 2 open tables)
- [ ] Implement table requirements (minimum ELO, max quit ratio)
- [ ] Add ready/unready system for players at table
- [ ] Implement table spectator slots (observers)
- [ ] Add table options (timer settings, game variants)
- [ ] Implement table auto-start when all players ready
- [ ] Handle table cleanup when host disconnects

## Auto-Matchmaking (P1)
- [ ] Implement ELO-based matchmaking queue
- [ ] Add format-specific queues (Standard, Commander, Draft, Sealed)
- [ ] Implement matchmaking algorithm (closest ELO ±200 range)
- [ ] Add queue timeout escalation (widen search after 2 minutes)
- [ ] Implement match acceptance/decline flow
- [ ] Add matchmaking cancel option
- [ ] Track matchmaking wait times for analytics
- [ ] Add queue position display for players
- [ ] Implement queue backfill (replace declined players)
- [ ] Add priority queue for reconnecting players

## Chat System (P1)
- [ ] Implement multi-channel chat system (lobby, table, game, whisper)
- [ ] Add gRPC streaming for real-time chat delivery
- [ ] Implement chat message types (user, system, game event, error)
- [ ] Add message history (last 100 messages per channel)
- [ ] Implement chat commands (`/help`, `/whisper`, `/mute`, `/friend`)
- [ ] Add card reference support (`[[Card Name]]` → card preview)
- [ ] Support Markdown formatting in chat messages
- [ ] Implement emote system (predefined reactions)
- [ ] Add typing indicators for active chat
- [ ] Implement chat search/filter functionality

## Chat Moderation (P1)
- [ ] Implement chat mute system (admin can mute users)
- [ ] Add profanity filter (configurable word list)
- [ ] Implement rate limiting (max 5 messages per 10 seconds)
- [ ] Add chat report system for moderation queue
- [ ] Implement user blocking (client-side mute)
- [ ] Add chat logs for admin review
- [ ] Implement automatic spam detection
- [ ] Add warning system before mute/ban
- [ ] Track moderation actions per user

## Deck Management (P1)
- [ ] Implement deck submission API endpoint
- [ ] Add deck format validation (card legality, deck size, banned list)
- [ ] Implement deck import formats (MTGO, Arena, text list)
- [ ] Add deck storage per user (save multiple decks)
- [ ] Implement deck sharing (export deck code)
- [ ] Add deck validation error reporting with specific rule violations
- [ ] Implement sideboard validation for constructed formats
- [ ] Add deck legality checking for multiple formats
- [ ] Implement deck list endpoint (user's saved decks)
- [ ] Add deck metadata (name, format, colors, archetype)
- [ ] Implement deck thumbnail generation (mana curve visualization)
- [ ] Add deck version history (track changes)
- [ ] Implement deck cloning/copying
- [ ] Add deck archiving (hide old decks)

## Game Lifecycle (P0)
- [ ] Wire table → game transition (move players from lobby to game)
- [ ] Implement game result recording (winner, turn count, duration)
- [ ] Update ELO ratings after game completion
- [ ] Track quit events and update quit ratio
- [ ] Implement rematch system (best-of-3, sideboarding)
- [ ] Add game concession handling
- [ ] Implement timeout/AFK detection and auto-forfeit
- [ ] Add draw agreement system
- [ ] Track game statistics (turns, spells cast, damage dealt)
- [ ] Implement postgame summary view

## Reconnection & Spectating (P1)
- [ ] Implement game state snapshot for reconnection
- [ ] Add "rejoin game" after disconnect (restore game view)
- [ ] Implement spectator system (join as observer)
- [ ] Add spectator permissions (view hands with consent)
- [ ] Implement spectator chat (separate from player chat)
- [ ] Add spectator delay option (30 second delay for competitive)
- [ ] Implement game replay system (step through previous game)
- [ ] Add replay speed controls (pause, step, fast-forward)
- [ ] Track spectator count per game
- [ ] Implement spectator kick/ban for hosts

## Tournament System (P2)
- [ ] Implement tournament types (Swiss, Single Elimination, Double Elimination)
- [ ] Add draft tournament support (pods, draft rounds)
- [ ] Implement tournament registration with deck submission
- [ ] Add tournament bracket generation
- [ ] Implement round timer system (chess clock style)
- [ ] Add bye assignments for odd player counts
- [ ] Implement tournament standings view (live leaderboard)
- [ ] Add tournament chat channel
- [ ] Implement tournament spectating (view all matches)
- [ ] Track tournament statistics (participation, completion rate)

## Tournament Management (P2)
- [ ] Implement round pairing algorithm (Swiss or bracket)
- [ ] Add match result reporting from game engine
- [ ] Implement tournament progression (advance winners)
- [ ] Add tournament admin controls (pause, cancel, restart round)
- [ ] Implement tiebreaker calculation (match win %, game win %, opponent win %)
- [ ] Add tournament quit penalties (higher than casual)
- [ ] Implement tournament prizes/rewards tracking
- [ ] Add tournament export (results, pairings, standings)
- [ ] Implement tournament seeding options (random, by rating)
- [ ] Add top 8 playoff system for large tournaments

## Admin Tools (P1)
- [ ] Implement admin authentication (separate admin token)
- [ ] Add user management dashboard (view all users, stats)
- [ ] Implement user ban/unban system with duration
- [ ] Add chat moderation tools (view reports, mute users)
- [ ] Implement server broadcast messages
- [ ] Add table force-close for stuck games
- [ ] Implement user impersonation for debugging
- [ ] Add server metrics dashboard (active games, players online)
- [ ] Implement admin action logging/audit trail
- [ ] Add bulk moderation actions (ban multiple users)

## Monitoring & Analytics (P1)
- [ ] Track concurrent player count over time
- [ ] Monitor average matchmaking wait times
- [ ] Track game completion rates (vs quits/timeouts)
- [ ] Measure average game duration by format
- [ ] Track chat message volume and moderation actions
- [ ] Monitor server resource usage (CPU, memory, goroutines)
- [ ] Add error rate tracking and alerting
- [ ] Implement performance metrics (latency percentiles)
- [ ] Track user retention metrics (daily/weekly active users)
- [ ] Add conversion funnel tracking (lobby → table → game)

## Friend System (P2)
- [ ] Implement friend requests and acceptance
- [ ] Add friends list with online status
- [ ] Implement friend notifications (came online, started game)
- [ ] Add friend invite to table/game
- [ ] Implement private messaging between friends
- [ ] Add friend removal/blocking
- [ ] Implement friend activity feed
- [ ] Add friend search by username
- [ ] Track friend game history (games played together)
- [ ] Implement friend list sorting/grouping

## User Profiles (P2)
- [ ] Add public profile page (stats, recent matches, favorite decks)
- [ ] Implement avatar upload system with image validation
- [ ] Add achievement/badge system (100 games played, 10-win streak, etc.)
- [ ] Display match history with deck archetypes
- [ ] Add privacy settings (hide stats, friends-only profile)
- [ ] Implement profile customization (bio, favorite colors)
- [ ] Add profile sharing (shareable URL)
- [ ] Track and display win rate by format
- [ ] Add "most played cards" statistics
- [ ] Implement profile themes/backgrounds

## Database Persistence (P2)
- [ ] Design schema for users, decks, games, tournaments
- [ ] Implement user account persistence (PostgreSQL)
- [ ] Add deck storage in database with versioning
- [ ] Implement game result history storage
- [ ] Add tournament result archival
- [ ] Implement chat log persistence for moderation
- [ ] Add database migration system (goose or golang-migrate)
- [ ] Implement database backup and recovery
- [ ] Add database connection pooling and optimization
- [ ] Implement soft deletion for user accounts

## Replay & History (P2)
- [ ] Store game action logs for replay
- [ ] Implement replay playback API with step controls
- [ ] Add replay sharing (shareable game ID)
- [ ] Implement match history pagination
- [ ] Add advanced search (filter by opponent, deck, format)
- [ ] Implement replay bookmarking (save key moments)
- [ ] Add replay export (download action log)
- [ ] Implement replay trimming (cut to interesting moments)
- [ ] Add replay commentary/annotation system
- [ ] Track most-watched replays

## Notifications (P2)
- [ ] Implement notification system (friend requests, game invites, match found)
- [ ] Add notification preferences per user
- [ ] Implement push notification support (web push API)
- [ ] Add notification history/inbox
- [ ] Implement notification grouping (multiple friend requests)
- [ ] Add notification read/unread status
- [ ] Implement notification muting (do not disturb mode)
- [ ] Add notification sound preferences
- [ ] Track notification delivery and read rates

## Leaderboards & Rankings (P2)
- [ ] Implement global leaderboard (top 100 by ELO)
- [ ] Add format-specific leaderboards
- [ ] Implement seasonal rankings (monthly/quarterly resets)
- [ ] Add leaderboard filtering (region, format, time period)
- [ ] Track leaderboard position changes
- [ ] Implement leaderboard rewards/titles
- [ ] Add friend leaderboard (compare with friends)
- [ ] Track historical leaderboard positions
- [ ] Implement leaderboard decay (inactive players drop)

## Security & Anti-Cheat (P1)
- [ ] Implement rate limiting on all endpoints
- [ ] Add CAPTCHA for registration
- [ ] Implement IP-based abuse detection
- [ ] Add suspicious activity flagging (impossible play speed)
- [ ] Implement game integrity validation (legal moves only)
- [ ] Add account age restrictions for competitive play
- [ ] Implement report system for suspected cheating
- [ ] Add automated pattern detection (bot behavior)
- [ ] Track multiple accounts from same IP
- [ ] Implement hardware fingerprinting for ban evasion

## Performance & Scaling (P2)
- [ ] Implement horizontal scaling for game servers
- [ ] Add load balancing for lobby servers
- [ ] Implement Redis for distributed session storage
- [ ] Add database read replicas for scaling
- [ ] Implement caching layer (user profiles, deck lists)
- [ ] Add CDN support for static assets
- [ ] Implement graceful shutdown and player migration
- [ ] Add health checks and auto-recovery
- [ ] Implement circuit breakers for external services
- [ ] Track and optimize hot paths (profiling)

## Client Integration (P1)
- [ ] Define client-server protocol documentation
- [ ] Implement error handling guidelines for clients
- [ ] Add versioning and backward compatibility
- [ ] Implement client capability negotiation
- [ ] Add client SDK/library for common operations
- [ ] Implement offline mode support (cached data)
- [ ] Add reconnection logic with exponential backoff
- [ ] Implement optimistic UI updates
- [ ] Add client-side validation before server calls
- [ ] Track client versions and force update if needed
