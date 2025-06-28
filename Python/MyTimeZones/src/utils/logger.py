import logging
from pathlib import Path
from ..config import LOG_FILE, LOG_FORMAT, LOG_LEVEL

def setup_logger(name: str) -> logging.Logger:
    """Set up and return a logger instance.
    
    Args:
        name (str): Name of the logger
        
    Returns:
        logging.Logger: Configured logger instance
    """
    logger = logging.getLogger(name)
    logger.setLevel(getattr(logging, LOG_LEVEL))
    
    # Create handlers
    file_handler = logging.FileHandler(LOG_FILE)
    console_handler = logging.StreamHandler()
    
    # Create formatters and add it to handlers
    formatter = logging.Formatter(LOG_FORMAT)
    file_handler.setFormatter(formatter)
    console_handler.setFormatter(formatter)
    
    # Add handlers to the logger
    logger.addHandler(file_handler)
    logger.addHandler(console_handler)
    
    return logger 