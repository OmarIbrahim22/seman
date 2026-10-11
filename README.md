# seman

Terminal-based student planner. Manage exams, todos, projects and subjects from the keyboard.

## Install

```bash
go install github.com/OmarIbrahim22/seman/cmd/seman@latest
```

Make sure `$GOPATH/bin` is in your PATH (one-time setup):

```bash
export PATH=$PATH:$(go env GOPATH)/bin
go install github.com/OmarIbrahim22/seman/cmd/seman@release

```

Or build from source:

```bash
git clone git@github.com:OmarIbrahim22/seman.git
cd seman
go install ./cmd/seman
```

Data is stored at `~/.local/share/seman/semester.json`.

## Run

```bash
seman
```

For development (no install):

```bash
go run ./cmd/seman
```

## Dependencies

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) — lofi player (optional)
- [mvp](https://github.com/nicowillis/mvp) — lofi player (optional)

## Tabs

| Key | Tab                                                        |
| --- | ---------------------------------------------------------- |
| `1` | Dashboard — overview of upcoming exams, todos and projects |
| `2` | Exams — manage subjects and their exams                    |
| `3` | Todos — weekly checklist                                   |
| `4` | Projects — track assignments                               |
| `5` | Settings                                                   |
| `6` | Subjects — dedicated subject management                    |
| `7` | Lofi — music player (requires lofi enabled in Settings)    |

## Global keys

| Key     | Action                                       |
| ------- | -------------------------------------------- |
| `S`     | Add subject                                  |
| `A`     | Add exam                                     |
| `P`     | Add project                                  |
| `E`     | Edit selected item                           |
| `D`     | Delete selected item                         |
| `G`     | Toggle global view (all weeks) / weekly view |
| `←` `→` | Previous / next week                         |
| `1`–`7` | Switch tabs                                  |
| `Q`     | Quit                                         |

## Per-tab keys

### Subjects (`6`)

| Key       | Action                       |
| --------- | ---------------------------- |
| `j` / `k` | Navigate                     |
| `Enter`   | Edit selected / add if empty |

### Exams (`2`)

| Key       | Action                                  |
| --------- | --------------------------------------- |
| `j` / `k` | Navigate                                |
| `Tab`     | Switch focus between subjects and exams |

### Todos (`3`)

| Key                     | Action      |
| ----------------------- | ----------- |
| `j` / `k`               | Navigate    |
| `Space` / `Enter` / `X` | Toggle done |
| `N`                     | Add todo    |

### Projects (`4`)

| Key       | Action   |
| --------- | -------- |
| `j` / `k` | Navigate |

### Settings (`5`)

| Key | Action                                |
| --- | ------------------------------------- |
| `O` | Toggle delete confirmation            |
| `W` | Cycle week span (1 → 2 → 3 → 4 → all) |
| `L` | Toggle lofi player                    |
| `U` | Set lofi playlist URL                 |

### Lofi (`7`)

| Key       | Action            |
| --------- | ----------------- |
| `j` / `k` | Navigate playlist |
| `Enter`   | Play selected     |
| `Space`   | Pause / resume    |
| `N`       | Next track        |
| `B`       | Previous track    |
| `X`       | Stop              |

## Subject autocomplete

When adding an exam or project, the **Subject** field supports autocomplete:

- Type to fuzzy-filter subjects (searches code and name)
- `↑` / `↓` or `Tab` / `Shift+Tab` to cycle matches
- `Enter` to confirm and advance to the next field

The matching panel appears to the right of the modal.
