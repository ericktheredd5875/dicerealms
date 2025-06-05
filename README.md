
# 🎲 DiceRealms

**DiceRealms** is a modern, Go-powered multiplayer MUD-style roleplaying platform using a structured protocol inspired by classic MUDs and tabletop RPGs like Dungeons & Dragons.

Built from the ground up to support immersive text-based storytelling, structured MCP commands, and group-based roleplaying adventures — DiceRealms lets you emote, speak, roll, and act in shared virtual spaces.

---
[![codecov](https://codecov.io/gh/ericktheredd5875/dicerealms/graph/badge.svg?token=8Q1IB3P0UL)](https://codecov.io/gh/ericktheredd5875/dicerealms)
---

## 👩‍👨 Features

* 🧹 **Structured MCP Protocol**: Custom command parsing with tags like `mcp-emote`, `mcp-roll`, and `mcp-say`.
* 🎝️ **Room-Based Group Play**: Join others in shared scenes and interact in real time.
* 🎲 **Dice Rolling**: Support for expressions like `1d20+5`, with critical success/failure detection.
* 🗣️ **Emotes & In-Character Speech**: Express yourself with structured roleplay.
* 🔄 **Extensible Architecture**: Future-ready for AI integration, persistence, and DM tools.

---

## 🚀 Getting Started

### Prerequisites

* Go 1.20+
* Git
* (Optional) Telnet or netcat for testing

### Clone and Run

```bash
git clone https://github.com/YOURUSERNAME/dice-realms.git
cd dice-realms
go run ./cmd/server
```

### Connect to the Server

In another terminal:

```bash
telnet localhost 4000
# OR
nc localhost 4000
```

Then try:

```text
#$#mcp-emote: text="draws his sword"
#$#mcp-say: text="We must be ready!"
#$#mcp-roll: dice="1d20+3" reason="Perception"
#$#mcp-help
```

---

## 🥪 Running Tests

```bash
go test ./internal/...
```

---

## 📂 Project Structure

```
cmd/server/        → Main entrypoint
internal/server/   → TCP server, connection handling
internal/game/     → Player, room, dice logic
internal/mcp/      → MCP tag parsing
```

---

## 📜 License

MIT — feel free to fork and build your own realms.

---

## 🧠 Future Roadmap

* [ ] Player commands: ~~`look`~~, ~~`move`~~, ~~`inventory`~~
* [ ] DM tools: `mcp-narrate`, scene control
* [ ] WebSocket/Discord client
* [ ] Persistent storage with PostgreSQL
* [ ] AI-driven NPCs and dynamic storytelling
* [ ] Split command descriptions into a dedicated map for maintainability.
* [x] Add mcp-help: command="mcp-roll" to explain individual commands in detail.
* [ ] Let DMs define custom help menus for their sessions.
* [ ] More Look-like commands (examine, scene, etc.).
* [ ] Add mcp-ooc for out-of-character speech.
* [ ] Support speech tags like volume="shout" → shouts, mutters, etc.
* [ ] Auto-prompt players with mcp-say: text="" if they type untagged input.
* [ ] Restrict narrate to DM-role players.
* [ ] Allow styled moods (e.g., tense, calm).
* [x] Store narration logs by scene or timestamp.
* [ ] Stats: show a reminder of how to improve them later.
* [ ] mcp-stat-reset for DM use
* [ ] Password or public key auth
* [ ] Player authentication (tie SSH login to in-game identity)
* [ ] Session logging
* [ ] Multiple ports (e.g., 4000 for Telnet, 2222 for SSH)
* [ ] Color support (many SSH clients are ANSI-capable!)

---

## 💬 Join the Realm

This is an open project — PRs and ideas are welcome!