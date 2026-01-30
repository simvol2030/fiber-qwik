import { component$, useSignal, $, useOnDocument } from "@builder.io/qwik";

interface SearchInputProps {
  value?: string;
  placeholder?: string;
  onSearch$: (value: string) => void;
  debounceMs?: number;
}

export const SearchInput = component$<SearchInputProps>(
  ({ value = "", placeholder = "Search...", onSearch$, debounceMs = 300 }) => {
    const inputValue = useSignal(value);
    const timerRef = useSignal<number | null>(null);

    const handleInput = $((e: Event) => {
      const target = e.target as HTMLInputElement;
      inputValue.value = target.value;

      if (timerRef.value !== null) {
        clearTimeout(timerRef.value);
      }
      timerRef.value = window.setTimeout(() => {
        onSearch$(inputValue.value);
      }, debounceMs) as unknown as number;
    });

    const handleClear = $(() => {
      inputValue.value = "";
      onSearch$("");
    });

    const handleKeydown = $((e: KeyboardEvent) => {
      if (e.key === "Enter") {
        if (timerRef.value !== null) {
          clearTimeout(timerRef.value);
        }
        onSearch$(inputValue.value);
      }
    });

    return (
      <div class="search-input">
        <span class="search-icon-label">&#128269;</span>
        <input
          type="text"
          placeholder={placeholder}
          value={inputValue.value}
          onInput$={handleInput}
          onKeyDown$={handleKeydown}
          aria-label="Search"
        />
        {inputValue.value && (
          <button
            class="clear-btn"
            onClick$={handleClear}
            aria-label="Clear search"
          >
            &#10005;
          </button>
        )}
      </div>
    );
  }
);
