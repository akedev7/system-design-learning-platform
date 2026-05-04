import pytest
from src.validators import is_email, is_phone, is_strong_password


class TestIsEmail:
    def test_valid_email(self):
        assert is_email("test@example.com") == True
        assert is_email("user.name+tag@domain.co.uk") == True

    def test_invalid_email_no_at(self):
        assert is_email("invalidemail") == False

    def test_invalid_email_no_domain(self):
        assert is_email("user@") == False

    def test_invalid_email_no_local_part(self):
        assert is_email("@domain.com") == False

    def test_empty_string(self):
        assert is_email("") == False


class TestIsPhone:
    def test_valid_phone_parentheses(self):
        assert is_phone("(123) 456-7890") == True

    def test_valid_phone_dashes(self):
        assert is_phone("123-456-7890") == True

    def test_valid_phone_dots(self):
        assert is_phone("123.456.7890") == True

    def test_invalid_phone_too_short(self):
        assert is_phone("123-456") == False

    def test_invalid_phone_letters(self):
        assert is_phone("123-ABC-7890") == False

    def test_empty_string(self):
        assert is_phone("") == False


class TestIsStrongPassword:
    def test_strong_password(self):
        assert is_strong_password("Abc123!@#") == True

    def test_strong_password_minimum_length(self):
        assert is_strong_password("Abc123!") == True

    def test_weak_password_no_uppercase(self):
        assert is_strong_password("abc123!@#") == False

    def test_weak_password_no_lowercase(self):
        assert is_strong_password("ABC123!@#") == False

    def test_weak_password_no_digits(self):
        assert is_strong_password("Abcdef!@#") == False

    def test_weak_password_too_short(self):
        assert is_strong_password("Ab1!") == False

    def test_empty_string(self):
        assert is_strong_password("") == False
