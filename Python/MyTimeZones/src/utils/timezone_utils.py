from datetime import datetime
from dateutil.relativedelta import relativedelta
import pytz
from typing import Dict, List, Tuple
from ..config import DEFAULT_TIMEZONE, TIME_FORMAT, DATE_FORMAT

def get_timezone_info(timezone_name: str) -> Dict:
    """Get timezone information including abbreviation and offset.
    
    Args:
        timezone_name (str): Name of the timezone
        
    Returns:
        Dict: Timezone information including abbreviation and offset
    """
    try:
        tz = pytz.timezone(timezone_name)
        now = datetime.now(tz)
        return {
            'abbreviation': now.strftime('%Z'),
            'offset': now.strftime('%z'),
            'name': timezone_name
        }
    except pytz.exceptions.UnknownTimeZoneError:
        return {
            'abbreviation': 'UTC',
            'offset': '+0000',
            'name': DEFAULT_TIMEZONE
        }

def calculate_time_difference(local_tz: str, target_tz: str) -> Tuple[str, str]:
    """Calculate time difference between two timezones.
    
    Args:
        local_tz (str): Local timezone name
        target_tz (str): Target timezone name
        
    Returns:
        Tuple[str, str]: Time difference in hours and minutes
    """
    try:
        utc_now = pytz.UTC.localize(datetime.utcnow())
        local_time = utc_now.astimezone(pytz.timezone(local_tz))
        target_time = utc_now.astimezone(pytz.timezone(target_tz))
        
        time_diff = relativedelta(local_time, target_time)
        hours = str(abs(time_diff.hours)).zfill(2)
        minutes = str(abs(time_diff.minutes)).zfill(2)
        
        return hours, minutes
    except pytz.exceptions.UnknownTimeZoneError:
        return "00", "00"

def format_datetime(dt: datetime, tz: str) -> Tuple[str, str]:
    """Format datetime in the specified timezone.
    
    Args:
        dt (datetime): Datetime object
        tz (str): Timezone name
        
    Returns:
        Tuple[str, str]: Formatted date and time
    """
    try:
        tz_obj = pytz.timezone(tz)
        localized_dt = dt.astimezone(tz_obj)
        return (
            localized_dt.strftime(DATE_FORMAT),
            localized_dt.strftime(TIME_FORMAT)
        )
    except pytz.exceptions.UnknownTimeZoneError:
        return dt.strftime(DATE_FORMAT), dt.strftime(TIME_FORMAT)

def validate_timezone(timezone_name: str) -> bool:
    """Validate if the timezone name is valid.
    
    Args:
        timezone_name (str): Timezone name to validate
        
    Returns:
        bool: True if valid, False otherwise
    """
    return timezone_name in pytz.all_timezones 