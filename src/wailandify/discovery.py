# src/waylandify/discovery.py
import shutil
from pathlib import Path


# Standard directories where .desktop files are stored
DESKTOP_FILE_DIRS = [
    Path("/usr/share/applications"),
    Path("/usr/local/share/applications"),
    Path.home() / ".local/share/applications",
]


def find_executable_path(names: list[str]) -> str | None:
    """
    Finds the full path of an executable from a list of possible names.
    Returns the first one found, or None.
    """
    for name in names:
        path = shutil.which(name)
        if path:
            return path
    return None


def get_all_desktop_files() -> list[Path]:
    """Scans standard directories and returns a list of all .desktop files."""
    all_files = []
    for directory in DESKTOP_FILE_DIRS:
        if directory.is_dir():
            all_files.extend(directory.glob("*.desktop"))
    return all_files


# src/waylandify/discovery.py
# ... (imports and first two functions are fine)


def find_related_desktop_files(
    exec_path: str,
    executables: list[str],  # Changed this parameter
    all_desktop_files: list[Path],
) -> list[Path]:
    """
    Finds all .desktop files that reference the given program.
    This simulates a `grep` for the executable name/aliases in all .desktop files.
    """
    found_files: set[Path] = set()
    search_terms = set(executables)  # Use the list directly

    for desktop_file in all_desktop_files:
        try:
            content = desktop_file.read_text()
            for line in content.splitlines():
                line = line.strip()
                if line.startswith("Exec="):
                    command_str = line.split("=", 1)[1].strip()
                    if not command_str:
                        continue  # Handle empty Exec=
                    executable_in_file = command_str.split()[0]

                    if Path(executable_in_file).name in search_terms:
                        found_files.add(desktop_file)
                        break
        except (IOError, UnicodeDecodeError):
            continue

    return list(found_files)
