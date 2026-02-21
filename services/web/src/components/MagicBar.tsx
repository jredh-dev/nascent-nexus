import { useState, useEffect, useRef, useCallback } from 'preact/hooks';

interface Action {
  title: string;
  description: string;
  type: 'navigation' | 'function';
  target: string;
}

const DEBOUNCE_MS = 150;

export default function MagicBar() {
  const [query, setQuery] = useState('');
  const [items, setItems] = useState<Action[]>([]);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [isOpen, setIsOpen] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const resultsRef = useRef<HTMLDivElement>(null);
  const debounceRef = useRef<ReturnType<typeof setTimeout>>();

  const close = useCallback(() => {
    setIsOpen(false);
    setItems([]);
    setSelectedIndex(-1);
  }, []);

  const search = useCallback(async (q: string) => {
    try {
      const res = await fetch(`/api/actions?q=${encodeURIComponent(q)}`);
      const data: Action[] = await res.json();
      setItems(data || []);
      setSelectedIndex(data && data.length > 0 ? 0 : -1);
      setIsOpen(true);
    } catch {
      setItems([]);
      setSelectedIndex(-1);
      setIsOpen(true);
    }
  }, []);

  const executeAction = useCallback((action: Action) => {
    close();
    setQuery('');
    inputRef.current?.blur();

    if (action.type === 'navigation') {
      window.location.href = action.target;
    }
  }, [close]);

  const onInput = useCallback((e: Event) => {
    const value = (e.target as HTMLInputElement).value;
    setQuery(value);

    clearTimeout(debounceRef.current);
    const trimmed = value.trim();
    if (trimmed.length === 0) {
      close();
      return;
    }
    debounceRef.current = setTimeout(() => search(trimmed), DEBOUNCE_MS);
  }, [close, search]);

  const onKeydown = useCallback((e: KeyboardEvent) => {
    switch (e.key) {
      case 'ArrowDown':
        e.preventDefault();
        setSelectedIndex(prev => Math.min(items.length - 1, prev + 1));
        break;
      case 'ArrowUp':
        e.preventDefault();
        setSelectedIndex(prev => Math.max(-1, prev - 1));
        break;
      case 'Enter':
        e.preventDefault();
        if (selectedIndex >= 0 && selectedIndex < items.length) {
          executeAction(items[selectedIndex]);
        }
        break;
      case 'Escape':
        e.preventDefault();
        close();
        inputRef.current?.blur();
        break;
    }
  }, [items, selectedIndex, executeAction, close]);

  // Ctrl+K / Cmd+K global shortcut
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        inputRef.current?.focus();
        inputRef.current?.select();
      }
    };
    document.addEventListener('keydown', handler);
    return () => document.removeEventListener('keydown', handler);
  }, []);

  // Close on outside click
  useEffect(() => {
    const handler = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      if (!target.closest('.magic-bar')) {
        close();
      }
    };
    document.addEventListener('click', handler);
    return () => document.removeEventListener('click', handler);
  }, [close]);

  // Scroll selected item into view
  useEffect(() => {
    if (selectedIndex >= 0 && resultsRef.current) {
      const children = resultsRef.current.querySelectorAll('.magic-bar-item');
      children[selectedIndex]?.scrollIntoView({ block: 'nearest' });
    }
  }, [selectedIndex]);

  return (
    <>
      <i class="fas fa-search magic-bar-icon"></i>
      <input
        ref={inputRef}
        class="magic-bar-input"
        type="text"
        placeholder="Search actions..."
        aria-label="Search actions"
        value={query}
        onInput={onInput}
        onKeyDown={onKeydown}
        onFocus={() => {
          if (query.trim().length > 0) search(query.trim());
        }}
      />
      <span class="magic-bar-hint">Ctrl+K</span>
      <div ref={resultsRef} class={`magic-bar-results${isOpen ? ' visible' : ''}`}>
        {isOpen && items.length === 0 && (
          <div class="magic-bar-empty">No results</div>
        )}
        {items.map((action, i) => (
          <div
            class={`magic-bar-item${i === selectedIndex ? ' selected' : ''}`}
            onMouseEnter={() => setSelectedIndex(i)}
            onClick={() => executeAction(action)}
          >
            <span class="magic-bar-item-icon">
              {action.type === 'navigation'
                ? <i class="fas fa-arrow-right"></i>
                : <i class="fas fa-bolt"></i>
              }
            </span>
            <span class="magic-bar-item-title">{action.title}</span>
            <span class="magic-bar-item-desc">{action.description}</span>
          </div>
        ))}
      </div>
    </>
  );
}
