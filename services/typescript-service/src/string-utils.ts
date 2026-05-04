/**
 * Capitalizes the first letter of a string
 * @param str - The input string
 * @returns The string with first letter capitalized
 */
export function capitalize(str: string): string {
  if (str.length === 0) {
    return str;
  }
  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Reverses a string
 * @param str - The input string
 * @returns The reversed string
 */
export function reverse(str: string): string {
  return str.split('').reverse().join('');
}

/**
 * Truncates a string to a maximum length, adding ellipsis if truncated
 * @param str - The input string
 * @param maxLength - The maximum allowed length
 * @returns The truncated string
 */
export function truncate(str: string, maxLength: number): string {
  if (str.length <= maxLength) {
    return str;
  }
  return str.slice(0, maxLength) + '...';
}
