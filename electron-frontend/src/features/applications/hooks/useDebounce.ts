import { useEffect, useState } from 'react';

/**
 * A custom React hook that debounces a value.
 *
 * This hook is used to delay updating a value until a certain amount of time
 * has passed since the last time the source value changed. This is useful for
 * performance-critical operations like filtering a list based on user input,
 * where you don't want to re-filter on every keystroke.
 *
 * @param {T} value The value to be debounced (e.g., a search query from an input).
 * @param {number} delay The delay in milliseconds (e.g., 300).
 * @returns {T} The debounced value, which will only update after the delay.
 */
export const useDebounce = <T>(value: T, delay: number): T => {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
};