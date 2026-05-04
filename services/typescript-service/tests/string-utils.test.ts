import { capitalize, reverse, truncate } from '../src/string-utils';

describe('String Utilities', () => {
  describe('capitalize', () => {
    it('should capitalize the first letter of a string', () => {
      expect(capitalize('hello')).toBe('Hello');
    });

    it('should handle empty string', () => {
      expect(capitalize('')).toBe('');
    });

    it('should not change already capitalized string', () => {
      expect(capitalize('Hello')).toBe('Hello');
    });

    it('should handle single character', () => {
      expect(capitalize('a')).toBe('A');
    });
  });

  describe('reverse', () => {
    it('should reverse a string', () => {
      expect(reverse('hello')).toBe('olleh');
    });

    it('should handle empty string', () => {
      expect(reverse('')).toBe('');
    });

    it('should handle palindrome', () => {
      expect(reverse('racecar')).toBe('racecar');
    });
  });

  describe('truncate', () => {
    it('should truncate string longer than maxLength', () => {
      expect(truncate('hello world', 5)).toBe('hello...');
    });

    it('should not truncate string shorter than maxLength', () => {
      expect(truncate('hello', 10)).toBe('hello');
    });

    it('should handle exact length', () => {
      expect(truncate('hello', 5)).toBe('hello');
    });

    it('should handle empty string', () => {
      expect(truncate('', 5)).toBe('');
    });
  });
});
