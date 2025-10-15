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


DEFAULT_CONFIG = (Path(__file__).parent / "data" / "config.toml").read_text()


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
