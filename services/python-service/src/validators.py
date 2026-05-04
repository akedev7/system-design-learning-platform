import re


def is_email(string: str) -> bool:
    """
    Validates email format.
    
    Args:
        string: String to validate as email
        
    Returns:
        True if valid email format, False otherwise
    """
    if not string:
        return False
    
    # Basic email regex pattern
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return bool(re.match(pattern, string))


def is_phone(string: str) -> bool:
    """
    Validates US phone number format.
    Accepts formats: (123) 456-7890, 123-456-7890, 123.456.7890
    
    Args:
        string: String to validate as phone number
        
    Returns:
        True if valid phone format, False otherwise
    """
    if not string:
        return False
    
    # Remove spaces and check various phone formats
    cleaned = string.strip()
    
    # Patterns for different phone formats
    patterns = [
        r'^\(\d{3}\)\s\d{3}-\d{4}$',  # (123) 456-7890
        r'^\d{3}-\d{3}-\d{4}$',        # 123-456-7890
        r'^\d{3}\.\d{3}\.\d{4}$',     # 123.456.7890
    ]
    
    return any(re.match(pattern, cleaned) for pattern in patterns)


def is_strong_password(string: str) -> bool:
    """
    Validates password strength.
    Requirements:
    - At least 7 characters long
    - Contains at least one uppercase letter
    - Contains at least one lowercase letter
    - Contains at least one digit
    - Contains at least one special character
    
    Args:
        string: String to validate as strong password
        
    Returns:
        True if strong password, False otherwise
    """
    if not string or len(string) < 7:
        return False
    
    has_upper = bool(re.search(r'[A-Z]', string))
    has_lower = bool(re.search(r'[a-z]', string))
    has_digit = bool(re.search(r'\d', string))
    has_special = bool(re.search(r'[!@#$%^&*(),.?":{}|<>]', string))
    
    return has_upper and has_lower and has_digit and has_special
