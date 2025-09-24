# src/waylandify/backup.py
import datetime
import shutil
from pathlib import Path

from .config import BACKUP_DIR


def create_backup(file_path: Path) -> Path | None:
    """
    Creates a timestamped backup of a file.
    Returns the path to the backup, or None on failure.
    """
    try:
        timestamp = datetime.datetime.now().strftime("%Y%m%d_%H%M%S_%f")
        backup_subdir = BACKUP_DIR / f"backup_{timestamp}"
        backup_subdir.mkdir(parents=True, exist_ok=True)

        backup_path = backup_subdir / file_path.name
        shutil.copy2(file_path, backup_path)

        return backup_path
    except Exception as e:
        print(f"⚠️  Could not create backup for {file_path}: {e}")
        return None
