# BF4DB-Search-Tool 🎮
Simple tool for searching for Players and check their ban status, cheat score, and linked accounts (on IP) using the BF4DB API.

## Features ✨
- Quick search for players by name or IP address
- Supports searching by Discord User ID
- *Debug mode for tracking errors
## Prerequisites 📋
- A valid API key from BF4DB. You can get one [here](https://bf4db.com/patreon "here").
- [Go 1.19 or higher](https://go.dev/dl/) (for building)

## Installation 🛠️

1. Clone the repository:
    ```sh
    git clone https://github.com/pruu-networking/BF4DB-Search-Tool.git
    cd BF4DB-Search-Tool
    ```
2. Build the project:
    ```sh
    go build -o bf4db
    ```

## Usage 🚀

Run the executable with the player name or IP address you want to search for:
```sh
./bf4db <player_name_or_ip> OR <discord_user_id -dc>

# Prints:
#Searching for player_name

#player_name | Not reported | Cheat score = 0 | https://bf4db.com/player/0/
# Cheat Report: http://bf4cheatreport.com/?pid=0&uid=&cnt=200&startdate=202501010000
# BF Agency: https://battlefield.agency/player/by-persona_id/bf4/=0
```

Example
To search for a player named "Player123":
```sh
./bf4db Player123
```
To search for an IP:
```sh
./bf4db 192.168.0.1
```
To search with Discord User ID:
```sh
./bf4db 274247581801119745 -dc
# *Debug is not available for Discord Search
```


### Options
```
 -h:    Show the help message.
 -c:    Set a new BF4DB API key.
 -dc:   Search for a Discord User ID.
  dbg:  Enable debug mode.
```

# Contributions are welcome! 🤝
Contributions are welcome! Feel free to fork the repository and submit a pull request.

To get started, make sure you have Go 1.16 or later installed on your system. Then, follow these steps:
1. Fork this repository to your own GitHub account.
2. Clone your fork of the repository to your local machine.
3. Make your changes and commit them with clear commit messages.
4. Push your changes to your forked repository on GitHub.
5. Submit a pull request to our repository.

If you're not comfortable submitting code changes, you can still help by opening an issue or feature request [here](https://github.com/pruu-networking/BF4DB-Search-Tool/issues) :)

## License 📄
This project is licensed under the MIT License. See the [LICENSE](https://github.com/pruu-networking/BF4DB-Search-Tool/blob/main/LICENSE) file for details.