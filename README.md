# VSModUpdater
Vintage Story Mod Updater - [ModDB](https://mods.vintagestory.at/modupdater)

## Installation
The program is distributed as a [zip file](https://github.com/rafalb8/VSModUpdater/releases) containing three executable files:
* VSModUpdater (for Linux)
* VSModUpdater.exe (for Windows)
* VSModUpdater_macOS (macOS Universal Binary)


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
* `-pre-release`
  * Allows updating to pre-release mod versions (e.g., alpha, beta). This functionality is also enabled automatically if a mod is already a pre-release version.
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
* `-webpage`
  * Generates a static HTML modlist webpage.


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

## Webpage Generation

The `-webpage` mode generates a static HTML file that can be hosted to share your server's modlist. The generated page includes live search, sortable columns, and automatic dark/light theme detection.

**Additional flags:**
* `-output <filename>` - Output filename (default: `modlist.html`)
* `-title <title>` - Page title (default: `"Server Modlist"`)
* `-deploy <project>` - Deploy to Cloudflare Pages project

**Examples:**

**Generate a basic modlist:**
```sh
./VSModUpdater -webpage
```
**Generate with custom title and filename:**
```sh
./VSModUpdater -webpage -title "Super Awesome Server Mods" -output public/mods.html
```
**Generate and deploy to Cloudflare Pages:**
```sh
./VSModUpdater -webpage -deploy my-server-mods
```

**Deploying to Cloudflare Pages:**

First-time setup:
```sh
npm install -g wrangler
wrangler login
```

Then generate and deploy:
```sh
./VSModUpdater -webpage -deploy your-project-name
```

The modlist will be available at `https://your-project-name.pages.dev`

The generated HTML file can also be uploaded to any web host.