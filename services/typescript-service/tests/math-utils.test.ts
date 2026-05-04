import { add, multiply, factorial } from '../src/math-utils';

describe('Math Utilities', () => {
  describe('add', () => {
    it('should add two positive numbers', () => {
      expect(add(2, 3)).toBe(5);
    });

    it('should handle negative numbers', () => {
      expect(add(-1, -2)).toBe(-3);
    });

    it('should handle zero', () => {
      expect(add(0, 5)).toBe(5);
      expect(add(5, 0)).toBe(5);
    });

    it('should handle decimal numbers', () => {
      expect(add(1.5, 2.5)).toBe(4);
    });
  });

  describe('multiply', () => {
    it('should multiply two positive numbers', () => {
      expect(multiply(2, 3)).toBe(6);
    });

    it('should handle negative numbers', () => {
      expect(multiply(-2, 3)).toBe(-6);
      expect(multiply(-2, -3)).toBe(6);
    });

    it('should handle zero', () => {
      expect(multiply(0, 5)).toBe(0);
    });
  });

  describe('factorial', () => {
    it('should return 1 for 0', () => {
      expect(factorial(0)).toBe(1);
    });

    it('should return 1 for 1', () => {
      expect(factorial(1)).toBe(1);
    });

    it('should calculate factorial for positive numbers', () => {
      expect(factorial(5)).toBe(120);
    });

    it('should throw error for negative numbers', () => {
      expect(() => factorial(-1)).toThrow('Factorial is not defined for negative numbers');
    });
  });
});
