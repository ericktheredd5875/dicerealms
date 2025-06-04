package config

var (
	SupportsANSI = true
	SceneDir     = "logs/scenes"
)

const WelcomeBanner = `
   ______  _________ _______  _______    _______  _______  _______  _        _______  _______ 
  (  __  \ \__   __/(  ____ \(  ____ \  (  ____ )(  ____ \(  ___  )( \      (       )(  ____ \
  | (  \  )   ) (   | (    \/| (    \/  | (    )|| (    \/| (   ) || (      | () () || (    \/
  | |   ) |   | |   | |      | (__      | (____)|| (__    | (___) || |      | || || || (_____ 
  | |   | |   | |   | |      |  __)     |     __)|  __)   |  ___  || |      | |(_)| |(_____  )
  | |   ) |   | |   | |      | (        | (\ (   | (      | (   ) || |      | |   | |      ) |
  | (__/  )___) (___| (____/\| (____/\  | ) \ \__| (____/\| )   ( || (____/\| )   ( |/\____) |
  (______/ \_______/(_______/(_______/  |/   \__/(_______/|/     \|(_______/|/     \|\_______)
`

const TagLine = `
	A realm of shared imagination, structured storytelling, and dice-fueled destiny.`

const WelcomePrompt = `Type: #$#mcp-help to begin your journey.`

// WelcomeBanner += "\033[36mType: #$#mcp-help: to begin your journey.\033[0m"

const Menu = "\n|----------------------------------------------------------------------|\n" +
	"|                       \033[36mDICE REALMS COMMAND MENU\033[0m                       |\n" +
	"|----------------------------------------------------------------------|\n" +
	"| \033[32m#$#mcp-help:\033[0m               Show this menu                            |\n" +
	"| \033[32m#$#mcp-look:\033[0m               Describe your current surroundings        |\n" +
	"| \033[32m#$#mcp-go: dir=\"north\"\033[0m     Move to another room (e.g., north, south) |\n" +
	"| \033[32m#$#mcp-say: text=\"...\"\033[0m     Speak to others in the room               |\n" +
	"| \033[32m#$#mcp-emote: text=\"...\"\033[0m   Perform an action or gesture              |\n" +
	"| \033[32m#$#mcp-roll: dice=\"1d20+3\"\033[0m Roll a dice with optional reason          |\n" +
	"| \033[32m#$#mcp-stat: roll=\"str\"\033[0m    Roll a stat-based skill check             |\n" +
	"| \033[32m#$#mcp-inventory:\033[0m          View your inventory                       |\n" +
	"| \033[32m#$#mcp-narrate: text=\"...\"\033[0m Narrate a scene as the DM                 |\n" +
	"| \033[32m#$#mcp-exit:\033[0m               Save and leave the realm                  |\n" +
	"|----------------------------------------------------------------------|\033[0m\n"

	// const MenuFallBack = ""
