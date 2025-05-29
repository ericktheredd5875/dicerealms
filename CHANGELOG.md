## 🧾 Changelog

### v0.1.0 – Core Engine Launch

#### ⚙️ System Architecture
- Built Go-powered TCP server with concurrency
- Structured MCP protocol parser (`#$#mcp-...`)
- Room-based player organization
- Graceful error handling and logging

#### 🧑‍🤝‍🧑 Player Interactions
- `mcp-say`: In-character speech
- `mcp-emote`: Descriptive actions
- `mcp-whisper`: Private IC messages
- `mcp-narrate`: DM-style storytelling
- Custom room-based prompts per player

#### 🎲 RPG Mechanics
- `mcp-roll`: Dice rolling with critical success/fail
- `mcp-stat`: Roll individual stats (4d6 drop lowest)
- `mcp-stat-gen`: Auto-generate all stats
- Stat locking after assignment
- Structured `Stats` model: STR, DEX, CON, INT, WIS, CHA

#### 🎒 Inventory System
- `mcp-take`: Add item to inventory
- `mcp-drop`: Remove item
- `mcp-inv`: List current inventory

#### 🎭 Scene System
- `mcp-scene-start`: Title, mood, startedBy
- `mcp-scene-end`: Ends current scene, saves log
- Per-room scene tracking
- Logs narration, says, emotes
- Scene summaries written to `/logs/scene_*.txt`

#### 🎨 Terminal UI / UX
- ANSI-styled output (color-coded prompts, says, emotes, narration)
- `[ANSI OFF]` mode via `mcp-client: supports_ansi=false`
- Re-rendered prompts after every broadcast
- Styled prompt: `🗡️ Name@Room >>>`

#### 🧪 Testing & Utilities
- Unit tests for `say`, `narrate`, `look`, and `roll`
- Color rendering logic isolated for portability
- Internal modules separated by feature domain:
  - `game/`, `mcp/`, `server/`