# src/waylandify/cli.py
import shutil
from pathlib import Path
from typing_extensions import Annotated

import typer
from rich import print

from . import config, discovery, desktop, backup

app = typer.Typer(
    help="A CLI tool to apply Wayland flags to Chromium-based applications."
)


@app.command()
def init():
    """
    Creates a default configuration file at ~/.config/waylandify/config.toml
    """
    config.create_default_config()


@app.command()
def apply(
    dry_run: Annotated[
        bool,
        typer.Option(
            "--dry-run", help="Show what would be changed without applying anything."
        ),
    ] = False,
):
    """
    Applies Wayland flags to the applications defined in the config file.
    """

    if dry_run:
        print(
            "[bold yellow]Running in dry-run mode. No files will be changed.[/bold yellow]"
        )

    try:
        cfg = config.load_config()
    except Exception:
        raise typer.Exit(code=1)

    print("[bold blue]DEBUG: Loaded config:[/bold blue]", cfg)

    all_desktop_files = discovery.get_all_desktop_files()
    user_desktop_dir = Path.home() / ".local/share/applications"

    if not user_desktop_dir.exists() and not dry_run:
        user_desktop_dir.mkdir(parents=True, exist_ok=True)

    print("-" * 30)

    # --- MODIFIED LOOP HERE ---
    for program_settings in cfg.programs:
        print(f"[bold magenta]Processing '{program_settings.name}'...[/bold magenta]")

        exec_path = discovery.find_executable_path(program_settings.executables)
        if not exec_path:
            print(
                f"  [yellow]⚠️  Could not find executable for any of: {program_settings.executables}. Skipping.[/yellow]"
            )
            continue

        print(f"  [dim]Found executable: {exec_path}[/dim]")

        # We need to pass ProgramSettings to discovery now
        related_files = discovery.find_related_desktop_files(
            exec_path, program_settings.executables, all_desktop_files
        )

        if not related_files:
            print("  [yellow]No related .desktop files found.[/yellow]")
            continue

        for source_path in related_files:
            # ... (the rest of the file from here is mostly the same)
            target_path = user_desktop_dir / source_path.name
            print(f"  -> Found desktop file: [cyan]{source_path}[/cyan]")

            modified_content = desktop.apply_flags_to_desktop_file(
                source_path, program_settings.flags
            )

            print(f"     [bold]Target path:[/bold] {target_path}")
            print(f"     [bold]Flags to add:[/bold] {' '.join(program_settings.flags)}")

            if not dry_run:
                try:
                    if target_path.exists():
                        backup.create_backup(target_path)

                    if source_path != target_path:
                        shutil.copy2(source_path, target_path)

                    target_path.write_text(modified_content)
                    print("     [green]✅ Applied flags successfully.[/green]")

                except Exception as e:
                    print(f"     [bold red]❌ Error applying flags: {e}[/bold red]")
                    raise typer.Exit(code=1)
        print("-" * 30)

    if not dry_run:
        print("\n[bold green]✨ All operations completed successfully! ✨[/bold green]")
