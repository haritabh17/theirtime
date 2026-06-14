<p align="center">
  <span style="font-size: 2.4em; font-weight: 650; letter-spacing: -0.03em;">their<span style="color: #2997ff;">time</span></span><br>
  <span style="font-size: 0.72em; letter-spacing: 0.14em; text-transform: uppercase; opacity: 0.55;">Teammate clocks · macOS menu bar</span>
</p>

<p align="center">
  <em>Stop opening Slack profiles to guess if it's a reasonable hour.</em><br>
  Teammate clocks and avatars in your menu bar — your Slack app, your Keychain, no server.
</p>

<div align="center">

![Teammate avatars and local times in the macOS menu bar](assets/theirtime-hero.png)

</div>

<p align="center">
  <a href="https://github.com/haritabh17/theirtime/releases"><img src="https://img.shields.io/github/v/release/haritabh17/theirtime?style=flat-square&label=release&color=111111&labelColor=111111" alt="Release"></a>
  &nbsp;
  <a href="#"><img src="https://img.shields.io/badge/macOS-only-111111?style=flat-square&logo=apple&logoColor=white" alt="macOS"></a>
  &nbsp;
  <a href="LICENSE"><img src="https://img.shields.io/badge/license-MIT-111111?style=flat-square" alt="MIT"></a>
</p>

## Quick start

```bash
curl -fsSL https://raw.githubusercontent.com/haritabh17/theirtime/main/scripts/install.sh | bash
theirtime onboard
theirtime team add bob U012ABCDEF
theirtime install-agents
```

| | |
|:--|:--|
| **Member ID** | Slack profile → **⋮** → *Copy member ID* (`U…`) |
| **Updates** | Menu bar every minute · Slack avatars every 15 minutes |

## Display

> Default: **`[avatar] 4.07 PM`** in the menu bar — turn names on with `show_name true`.

```bash
theirtime config set show_name true
theirtime config set format_24h true
theirtime config set time_precision hours
theirtime install-agents
```

| You want | Set |
|:--|:--|
| Names in the bar | `show_name true` |
| 24-hour clock | `format_24h true` |
| Hour only (`4 PM`) | `time_precision hours` |
| Avatar + name + time | `show_avatar true` · `show_name true` · `show_time true` |

`theirtime config get` · `theirtime config set show_avatar|show_name|show_time|format_24h|time_precision …`

## Commands

| Command | Purpose |
|:--|:--|
| `theirtime team list` · `team remove <label>` | Manage watched teammates |
| `theirtime status` | Config, Keychain, and agent state |
| `theirtime auth` | Re-authorize Slack |
| `theirtime offboard` | Uninstall everything |
| `theirtime menubar --demo` | Preview the menu bar without Slack |

Logs · `~/Library/Logs/theirtime/menubar.log`

## Privacy & develop

Keychain for secrets · `~/Library/Application Support/theirtime/config.yaml` for prefs · OAuth on `127.0.0.1:8765` only · no telemetry.

```bash
make build && ./bin/theirtime menubar --demo
```

[`manifest/theirtime.manifest.yaml`](manifest/theirtime.manifest.yaml) · [Releases](https://github.com/haritabh17/theirtime/releases) · [CLI plan](docs/PLAN-cli.md)
