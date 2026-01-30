import { component$, Slot } from "@builder.io/qwik";

interface ModalProps {
  open: boolean;
  title: string;
  size?: "sm" | "md" | "lg";
  onClose$: () => void;
}

export const Modal = component$<ModalProps>(
  ({ open, title, size = "md", onClose$ }) => {
    if (!open) return null;

    return (
      <div
        class="modal-backdrop"
        onClick$={(e) => {
          if (e.target === e.currentTarget) onClose$();
        }}
        role="presentation"
      >
        <div
          class={`modal modal-${size}`}
          role="dialog"
          aria-modal="true"
          aria-labelledby="modal-title"
        >
          <div class="modal-header">
            <h2 id="modal-title" class="modal-title">
              {title}
            </h2>
            <button
              class="modal-close"
              onClick$={onClose$}
              aria-label="Close"
            >
              &#10005;
            </button>
          </div>

          <div class="modal-body">
            <Slot />
          </div>
        </div>
      </div>
    );
  }
);
