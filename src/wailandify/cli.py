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
    Applies Wayland flags to configured applications based on the config file.
    """
    if dry_run:
        print(
            "[bold yellow]Running in dry-run mode. No files will be changed.[/bold yellow]"
        )

    try:
        cfg = config.load_config()
    except Exception:
        raise typer.Exit(code=1)

    # --- DEBUGGING LINE ADDED HERE ---
    print("[bold blue]DEBUG: Loaded config:[/bold blue]", cfg)

    all_desktop_files = discovery.get_all_desktop_files()
    user_desktop_dir = Path.home() / ".local/share/applications"

    if not user_desktop_dir.exists() and not dry_run:
        print(
            f"🔧 User application directory not found. Creating {user_desktop_dir}..."
        )
        user_desktop_dir.mkdir(parents=True, exist_ok=True)

    print("-" * 30)

    for key, program_settings in cfg.programs.items():
        print(f"[bold magenta]Processing '{key}'...[/bold magenta]")

        exec_path = discovery.find_executable_path(program_settings.names)
        if not exec_path:
            print(
                f"  [yellow]⚠️  Could not find executable for any of: {program_settings.names}. Skipping.[/yellow]"
            )
            continue

        print(f"  [dim]Found executable: {exec_path}[/dim]")

        related_files = discovery.find_related_desktop_files(
            exec_path, program_settings, all_desktop_files
        )

        if not related_files:
            print("  [yellow]No related .desktop files found.[/yellow]")
            continue

        for source_path in related_files:
            target_path = user_desktop_dir / source_path.name
            print(f"  -> Found desktop file: [cyan]{source_path}[/cyan]")

            # Since you are running in dry-run, we disable the 'already applied' check to always see what it *would* do.
            # original_content = source_path.read_text().strip()
            modified_content = desktop.apply_flags_to_desktop_file(
                source_path, program_settings.flags
            )

            # if original_content == modified_content.strip():
            #      print("     [green]✅ Flags already applied. Skipping.[/green]")
            #      continue

            print(f"     [bold]Target path:[/bold] {target_path}")
            print(f"     [bold]Flags to add:[/bold] {' '.join(program_settings.flags)}")

            if not dry_run:
                try:
                    # Backup the original file if it exists in the target location
                    if target_path.exists():
                        backup.create_backup(target_path)

                    # If the source is not in the user dir, copy it first to ensure we aren't modifying system files directly
                    if source_path != target_path:
                        shutil.copy2(source_path, target_path)

                    # Now, write the modified content to the user's local directory
                    with open(target_path, "w") as f:
                        f.write(modified_content)
                    print("     [green]✅ Applied flags successfully.[/green]")

                except Exception as e:
                    print(f"     [bold red]❌ Error applying flags: {e}[/bold red]")
                    # Fail fast as requested
                    raise typer.Exit(code=1)
        print("-" * 30)

    if not dry_run:
        print("\n[bold green]✨ All operations completed successfully! ✨[/bold green]")
