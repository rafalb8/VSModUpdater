# VSModUpdater
Vintage Story Mod Updater - [ModDB](https://mods.vintagestory.at/modupdater)

## Installation
The program is distributed as a [zip file](https://github.com/rafalb8/VSModUpdater/releases) containing two executable files:
* VSModUpdater (for Linux)
* VSModUpdater.exe (for Windows)

Simply download the zip file, extract the contents, and run the appropriate executable for your operating system.

## Usage

### Flag Reference
* `-mod-path <path>`
  * Specifies the path to your Vintage Story mods directory.
  * **Default:** `~/.config/VintagestoryData/Mods` (on Linux) or `%APPDATA%/VintagestoryData/Mods` (on Windows).
* `-backup-path <path>`
  * Specifies where to store mod backups.
  * **Default:** `~/.config/VintagestoryData/ModBackups` (on Linux) or `%APPDATA%/VintagestoryData/ModBackups` (on Windows).
* `-dry-run`
  * If this flag is set, the program will just print the updates.
* `-backup`
  * If this flag is set, the program will move old mods to the backup directory instead of deleting them.
* `-interactive`
  * Starts the program in an interactive mode, allowing you to select which mods to update. This is the default behavior on Windows.
* `-ignore <modID1,modID2,...>`
  * Provides a comma-separated list of mod IDs to skip during updates.

### Modes
The program can run in several modes. You should only use one mode at a time.

* `-version`
  * Prints the program's version and exits.
* `-self`
  * Updates the Vintage Story Mod Updater program itself.
* `-list`
  * Lists all installed mods and their versions.


### Examples
**Update all mods:**
```sh
./VSModUpdater
```
**Update all mods, but back up old ones instead of deleting them:**
```sh
./VSModUpdater -backup
```
**Update all mods except for two specific ones:**
```sh
./VSModUpdater -ignore "some-mod-id,another-mod-id"
```
**List all installed mods:**
```sh
./VSModUpdater -list
```
**Check the program's version:**
```sh
./VSModUpdater -version
```