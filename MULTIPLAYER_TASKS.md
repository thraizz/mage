# Multiplayer & Lobby System Task Tracker

Comprehensive tracker for building the multiplayer, lobby, and social features for the Go Mage server. Modern gRPC-Web based system with real-time updates.

Status legend:
- `[x]` Completed
- `[ ]` Pending / not yet started
- `[~]` In progress or partially implemented

## Phase 1: MVP (XMage Feature Parity)

**Goal:** Get players into games with the core XMage experience. No fancy features, just functional multiplayer.

### Infrastructure Setup
- [ ] Create Dockerfile for Backend
- [ ] Create Dockerfile for Client Webserver
- [ ] Create docker-compose for Database

### Core Authentication & Sessions
- [ ] JWT-based authentication (username/password)
- [ ] Guest mode for anonymous play
- [ ] Session management with reconnection support
- [ ] Basic user state tracking (online, in-game)

### Basic User Management
- [ ] User profiles (username, email)
- [ ] Simple stats tracking (games played, wins/losses)
- [ ] Quit ratio tracking and display (no enforcement yet)

### gRPC-Web Foundation
- [ ] Proto messages for: CreateTable, JoinTable, LeaveTable, StartGame
- [ ] Proto messages for chat (SendMessage, ReceiveMessage)
- [ ] Bidirectional streaming for real-time updates
- [ ] Basic error handling and status messages

### Simple Lobby
- [ ] List of active tables (format, open slots, host name)
- [ ] Real-time updates when tables created/removed
- [ ] Show online player count
- [ ] Basic filtering (format, open tables only)

### Table Management
- [ ] Create table (format, player count, optional password)
- [ ] Join table with deck submission
- [ ] Simple table states (WAITING → ACTIVE → FINISHED)
- [ ] Ready/unready system
- [ ] Auto-start when all players ready
- [ ] Table host can kick players
- [ ] Table chat

### Basic Chat
- [ ] Lobby chat (global)
- [ ] Table chat (pre-game)
- [ ] In-game chat (during match)
- [ ] Private whispers (/w username message)
- [ ] Message history (last 50 per channel)
- [ ] Simple rate limiting (10 messages/minute)

### Deck Management
- [ ] Submit deck for table
- [ ] Deck format validation (card legality, deck size)
- [ ] Text-based deck import (simple list format)
- [ ] Save one deck per format per user

### Game Lifecycle
- [ ] Table → game transition
- [ ] Record game results (winner, duration)
- [ ] Update win/loss stats
- [ ] Track quits (increment quit ratio)
- [ ] Concede option

### Reconnection
- [ ] Game state snapshot on disconnect
- [ ] Rejoin active game after disconnect (5 min window)
- [ ] Timeout/AFK detection (2 min → auto-forfeit)

### Minimal Persistence
- [ ] User accounts (PostgreSQL)
- [ ] Deck storage
- [ ] Game results (basic stats only)
- [ ] Database migrations

### Security Basics
- [ ] Password hashing (bcrypt)
- [ ] Rate limiting on auth endpoints
- [ ] Basic input validation

---

## Phase 2: Feature Packages (Post-MVP)

Pick and choose based on user demand after launch:

### Package A: Matchmaking & Rankings
- Auto-matchmaking queue (format-specific)
- ELO rating system
- Global leaderboards
- Match acceptance flow

### Package B: Enhanced Social
- Friend system (requests, list, online status)
- Friend invites to tables
- Friend-only tables
- Activity notifications

### Package C: Spectating & Replays
- Spectator mode (join as observer)
- Game replay system
- Match history with details
- Shareable replay links

### Package D: Tournaments
- Swiss tournament support
- Tournament brackets
- Round pairings
- Standings/leaderboards

### Package E: Moderation Tools
- Admin dashboard
- User ban system
- Chat moderation (mute, filter)
- Report system

### Package F: Enhanced Deck Management
- Multiple saved decks per format
- Deck sharing/export
- Import from MTGO/Arena formats
- Deck statistics and metadata

### Package G: User Profiles & Achievements
- Public profile pages
- Avatar uploads
- Achievement system
- Match history display

### Package H: Advanced Features
- Tournament draft support
- Scheduled events
- Seasonal rankings
- Advanced analytics

---

**Key Changes:**

1. **Removed from MVP:** ELO, tournaments, spectating, friend system, admin tools, achievements, advanced moderation
2. **Simplified:** One deck per format (not unlimited saves), basic chat (no emotes/markdown initially), minimal stats
3. **Deferred:** All the "nice to have" social features until you validate core gameplay works

**MVP Timeline Estimate:** 6-8 weeks for a solo developer focusing on the essentials.

Does this feel like the right scope? We can adjust the MVP boundaries if needed.
