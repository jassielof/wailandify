# src/waylandify/config.py
from pathlib import Path

import tomlkit
from pydantic import BaseModel, ValidationError
from rich import print

CONFIG_DIR = Path.home() / ".config" / "waylandify"
CONFIG_FILE_PATH = CONFIG_DIR / "config.toml"
BACKUP_DIR = CONFIG_DIR / "backups"


class ProgramSettings(BaseModel):
    name: str
    executables: list[str]
    flags: list[str]


class Config(BaseModel):
    programs: list[ProgramSettings]


# Updated default config to match the new structure
DEFAULT_CONFIG = """
#
# Waylandify Configuration File
#
# Use [[programs]] to define a list of applications to modify.
#
# - 'name': A unique name for this entry (e.g., "vscode").
# - 'executables': A list of binary names to search for (e.g., ["code", "code-insiders"]).
# - 'flags': A list of command-line flags to apply.
#

[[programs]]
name = "vscode"
executables = ["code", "code-insiders"]
flags = [
    "--enable-features=UseOzonePlatform",
    "--ozone-platform=wayland",
]

[[programs]]
name = "brave"
executables = ["brave-browser", "brave-browser-stable"]
flags = [
    "--enable-features=TouchpadOverscrollHistoryNavigation",
]
"""


# ... (the rest of the file remains the same)
def create_default_config():
    if CONFIG_FILE_PATH.exists():
        print(f"[yellow]Configuration file already exists at:[/] {CONFIG_FILE_PATH}")
        return
    print(f"Creating default config at {CONFIG_FILE_PATH}...")
    try:
        CONFIG_DIR.mkdir(parents=True, exist_ok=True)
        CONFIG_FILE_PATH.write_text(DEFAULT_CONFIG)
        print("[green]✅ Successfully created configuration file.[/green]")
    except Exception as e:
        print(f"[bold red]❌ Error creating config file: {e}[/bold red]")


def load_config() -> Config:
    if not CONFIG_FILE_PATH.is_file():
        print(
            f"[bold red]❌ Configuration file not found at {CONFIG_FILE_PATH}[/bold red]"
        )
        raise FileNotFoundError
    try:
        data = tomlkit.parse(CONFIG_FILE_PATH.read_text())
        return Config.model_validate(data)
    except ValidationError as e:
        print("[bold red]❌ Configuration file is invalid.[/bold red]")
        for error in e.errors():
            loc = " -> ".join(map(str, error["loc"]))
            print(f"  - {loc}: {error['msg']}")
        raise
    except Exception as e:
        print(f"[bold red]❌ Failed to load or parse config file: {e}[/bold red]")
        raise
