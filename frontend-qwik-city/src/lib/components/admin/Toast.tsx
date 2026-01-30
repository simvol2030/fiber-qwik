import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import {
  getToasts,
  toast,
  subscribe as adminSubscribe,
} from "~/lib/stores/admin";

export const Toast = component$(() => {
  const toasts = useSignal(getToasts());

  useVisibleTask$(() => {
    const unsub = adminSubscribe(() => {
      toasts.value = [...getToasts()];
    });
    return unsub;
  });

  const icons: Record<string, string> = {
    success: "\u2713",
    error: "\u2715",
    warning: "\u26A0",
    info: "\u2139",
  };

  if (toasts.value.length === 0) return null;

  return (
    <div class="toast-container">
      {toasts.value.map((t) => (
        <div key={t.id} class={`toast toast-${t.type}`}>
          <span class="toast-icon">{icons[t.type]}</span>
          <span class="toast-message">{t.message}</span>
          <button
            class="toast-close"
            onClick$={() => toast.remove(t.id)}
            aria-label="Close"
          >
            &#10005;
          </button>
        </div>
      ))}
    </div>
  );
});
