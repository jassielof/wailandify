# src/waylandify/config.py
from pathlib import Path
import importlib.resources

import tomlkit
from pydantic import BaseModel, ValidationError
from rich import print

CONFIG_DIR = Path.home() / ".config" / "waylandify"
CONFIG_FILE_PATH = CONFIG_DIR / "config.toml"
BACKUP_DIR = CONFIG_DIR / "backups"


class ProgramSettings(BaseModel):
    """Defines settings for a single program entry in the config."""

    names: list[str]
    flags: list[str]


class Config(BaseModel):
    """The root model for the entire config.toml file."""

    programs: dict[str, ProgramSettings]


def create_default_config():
    """Creates the config directory and default config file if they don't exist."""
    if CONFIG_FILE_PATH.exists():
        print(f"[yellow]Configuration file already exists at:[/] {CONFIG_FILE_PATH}")
        return

    print(f"Creating default config at {CONFIG_FILE_PATH}...")
    try:
        CONFIG_DIR.mkdir(parents=True, exist_ok=True)

        # Read the default config from the external file
        try:
            with importlib.resources.open_text(
                "wailandify", "default_config.toml"
            ) as f:
                default_config_content = f.read()
        except FileNotFoundError:
            # Fallback to reading from the same directory
            default_config_path = Path(__file__).parent / "default_config.toml"
            with open(default_config_path, "r") as f:
                default_config_content = f.read()

        with open(CONFIG_FILE_PATH, "w") as f:
            f.write(default_config_content)
        print("[green]✅ Successfully created configuration file.[/green]")
        print("Please edit it to match your needs before running 'apply'.")
    except Exception as e:
        print(f"[bold red]❌ Error creating config file: {e}[/bold red]")


def load_config() -> Config:
    """Loads and validates the configuration from the TOML file."""
    if not CONFIG_FILE_PATH.exists():
        print(
            f"[bold red]❌ Configuration file not found at {CONFIG_FILE_PATH}[/bold red]"
        )
        print("Please run 'waylandify init' to create a default config file.")
        raise FileNotFoundError

    try:
        with open(CONFIG_FILE_PATH, "r") as f:
            data = tomlkit.parse(f.read())

        # Validate the parsed data against our Pydantic model
        config_model = Config.model_validate(data)
        return config_model

    except ValidationError as e:
        print("[bold red]❌ Configuration file is invalid.[/bold red]")
        for error in e.errors():
            loc = " -> ".join(map(str, error["loc"]))
            print(f"  - {loc}: {error['msg']}")
        raise
    except Exception as e:
        print(f"[bold red]❌ Failed to load or parse config file: {e}[/bold red]")
        raise
