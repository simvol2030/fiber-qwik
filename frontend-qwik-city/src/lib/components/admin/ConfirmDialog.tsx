import { component$ } from "@builder.io/qwik";

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  message: string;
  confirmLabel?: string;
  cancelLabel?: string;
  variant?: "danger" | "warning" | "info";
  onConfirm$: () => void;
  onCancel$: () => void;
}

export const ConfirmDialog = component$<ConfirmDialogProps>(
  ({
    open,
    title,
    message,
    confirmLabel = "Confirm",
    cancelLabel = "Cancel",
    variant = "danger",
    onConfirm$,
    onCancel$,
  }) => {
    if (!open) return null;

    const icon =
      variant === "danger"
        ? "\u26A0\uFE0F"
        : variant === "warning"
          ? "\u26A0\uFE0F"
          : "\u2139\uFE0F";

    return (
      <div
        class="dialog-backdrop"
        onClick$={onCancel$}
        role="presentation"
      >
        <div
          class="dialog"
          role="alertdialog"
          aria-modal="true"
          aria-labelledby="dialog-title"
          aria-describedby="dialog-message"
          onClick$={(e) => e.stopPropagation()}
        >
          <div class={`dialog-icon dialog-icon-${variant}`}>{icon}</div>

          <h2 id="dialog-title" class="dialog-title">
            {title}
          </h2>
          <p id="dialog-message" class="dialog-message">
            {message}
          </p>

          <div class="dialog-actions">
            <button class="btn-cancel" onClick$={onCancel$}>
              {cancelLabel}
            </button>
            <button
              class={`btn-confirm btn-${variant}`}
              onClick$={onConfirm$}
            >
              {confirmLabel}
            </button>
          </div>
        </div>
      </div>
    );
  }
);
