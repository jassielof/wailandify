import configparser
import io
from pathlib import Path


def add_flags_to_exec_command(exec_cmd: str, flags: list[str]) -> str:
    """
    Intelligently adds flags to an Exec command string, avoiding duplicates.
    """
    # Don't modify empty commands
    if not exec_cmd.strip():
        return ""

    parts = exec_cmd.split()

    executable = parts[0]
    original_args = parts[1:]

    new_flags = []
    # Only add flags that are not already present in the command
    for flag in flags:
        if flag not in exec_cmd:
            new_flags.append(flag)

    # Reconstruct the command: executable + new_flags + original_args
    final_parts = [executable] + new_flags + original_args
    return " ".join(final_parts)


def apply_flags_to_desktop_file(path: Path, flags: list[str]) -> str:
    """
    Parses a .desktop file, applies flags to all Exec keys, and returns the new content.
    """
    parser = configparser.ConfigParser(
        delimiters=("="),
        interpolation=None,
        allow_no_value=True,  # Makes parsing more robust for .desktop files
    )
    # Preserve the case of keys
    parser.optionxform = str

    content = path.read_text()

    # Some .desktop files (especially PWAs) start with a shebang.
    # configparser sees this as an error. We must filter it out before parsing.
    if content.startswith("#!"):
        content = "\n".join(content.splitlines()[1:])

    parser.read_string(content)

    # Apply flags to all 'Exec' entries in every section
    for section in parser.sections():
        if parser.has_option(section, "Exec"):
            original_exec = parser.get(section, "Exec")
            modified_exec = add_flags_to_exec_command(original_exec, flags)
            parser.set(section, "Exec", modified_exec)

    # Write the modified configuration to a string
    string_io = io.StringIO()
    parser.write(string_io, space_around_delimiters=False)
    return string_io.getvalue().strip()
