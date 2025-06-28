import os
from pathlib import Path

# Application settings
APP_NAME = "KeepMyTime"
APP_VERSION = "1.0.0"
DEFAULT_WINDOW_SIZE = "650x250"
UPDATE_INTERVAL = 60000  # milliseconds

# File paths
BASE_DIR = Path(__file__).parent
DATA_DIR = BASE_DIR / "data"
CONFIG_DIR = BASE_DIR / "config"

# Ensure directories exist
DATA_DIR.mkdir(exist_ok=True)
CONFIG_DIR.mkdir(exist_ok=True)

# Theme settings
DEFAULT_THEME = "yaru"
THEME_FILE_PATTERN = "*timezones*.json"

# Timezone settings
DEFAULT_TIMEZONE = "UTC"
TIME_FORMAT = "%H:%M"
DATE_FORMAT = "%Y-%m-%d"

# GUI settings
TREEVIEW_COLUMNS = {
    'Name': {'width': 100, 'anchor': 'center'},
    'Description': {'width': 200, 'anchor': 'center'},
    'Date': {'width': 100, 'anchor': 'center'},
    'Time': {'width': 50, 'anchor': 'w'},
    'HoursDiff': {'width': 100, 'anchor': 'center'}
}

# Logging settings
LOG_FILE = DATA_DIR / "app.log"
LOG_FORMAT = "%(asctime)s - %(name)s - %(levelname)s - %(message)s"
LOG_LEVEL = "INFO" 