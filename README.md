# VSModUpdater
Vintage Story Mod Updater - [ModDB](https://mods.vintagestory.at/modupdater)

## Installation
The program is distributed as a [zip file](https://github.com/rafalb8/VSModUpdater/releases) containing three executable files:
* VSModUpdater (for Linux)
* VSModUpdater.exe (for Windows)
* VSModUpdater_macOS (macOS Universal Binary)


Simply download the zip file, extract the contents, and run the appropriate executable for your operating system.

If you have `go` installed, you can install it using:
```sh
go install -trimpath github.com/rafalb8/VSModUpdater/v2@latest
```

## Usage

### Flag Reference
* `-m, --mod-path <path>`
  * Specifies the path to your Vintage Story mods directory.
  * **Default:** `~/.config/VintagestoryData/Mods` (on Linux), `%APPDATA%\VintagestoryData\Mods` (on Windows), or the equivalent OS user config directory.
* `--backup-path <path>`
  * Specifies where to store mod backups. If not set, defaults to a sibling directory of your `mod-path` named `ModBackups`.
  * **Default:** `~/.config/VintagestoryData/ModBackups` (on Linux) or `%APPDATA%\VintagestoryData\ModBackups` (on Windows).
* `-p, --dry-run`
  * Runs the updater without actually making any changes (print only).
* `-b, --backup`
  * Backs up old mods to the backup directory instead of deleting them.
* `--pre-release`
  * Allows updating to pre-release mod versions (e.g., alpha, beta). This functionality is also enabled automatically if an installed mod is already a pre-release version.
* `-y, --no-confirm`
  * Automatically confirms all update actions, skipping exclusion prompts.
* `-x, --ignore <modID1,modID2,...>`
  * Disables updates for a comma-separated list of specific mod IDs.

### Modes
The program can run in several modes. You should only use one mode at a time.

* `-v, --version`
  * Prints the program's version and exits.
* `--self`
  * Updates the `VSModUpdater` program itself.
* `-l, --list`
  * Lists all installed mods and their versions.
* `-s, --simple`
  * Runs the updater in a simple update mode.
* `-i, --import <file>`
  * Imports and downloads a mod list from the specified file to your `-mod-path`.
* `-e, --export <file>`
  * Exports your current mod list from `-mod-path` to the specified file.


### Examples
**Update all mods (Standard run):**
```sh
./VSModUpdater
```

**Run a preview of updates without changing any files:**
```sh
./VSModUpdater -p
```

**Update all mods without prompting for exclusions:**
```sh
./VSModUpdater -y
```

**Update all mods, but back up old ones instead of deleting them:**
```sh
./VSModUpdater -b
```

**Update all mods except for two specific ones:**
```sh
./VSModUpdater -x some-mod-id -x another-mod-id
```

**List all installed mods:**
```sh
./VSModUpdater -l
```

**Check the program's version:**
```sh
./VSModUpdater -v
```

**Export modlist to a file:**
```sh
./VSModUpdater -e modlist.txt
```

**Download modlist from a file to `mods` directory :**
```sh
./VSModUpdater -i modlist.txt -m mods
```
