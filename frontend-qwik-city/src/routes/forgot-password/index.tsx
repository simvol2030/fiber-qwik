import { component$, useSignal } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { api } from "~/lib/api/client";

export default component$(() => {
  const email = useSignal("");
  const error = useSignal("");
  const success = useSignal(false);
  const isSubmitting = useSignal(false);

  return (
    <div class="auth-page">
      <div class="auth-card card">
        {success.value ? (
          <div class="success-state">
            <h1>Check Your Email</h1>
            <p class="success-message">
              If an account with that email exists, a password reset link has
              been sent. Please check your inbox and follow the instructions.
            </p>
            <a href="/login" class="btn-primary btn-full">
              Back to Login
            </a>
          </div>
        ) : (
          <>
            <h1>Forgot Password</h1>
            <p class="subtitle">
              Enter your email address and we'll send you a reset link.
            </p>

            {error.value && (
              <div class="alert alert-error">{error.value}</div>
            )}

            <form
              preventdefault:submit
              onSubmit$={async () => {
                error.value = "";
                isSubmitting.value = true;

                try {
                  const response = await api.post("/auth/forgot-password", {
                    email: email.value,
                  });
                  if (response.success) {
                    success.value = true;
                  } else {
                    error.value =
                      response.error?.message ||
                      "Failed to send reset link";
                  }
                } catch {
                  error.value = "Network error. Please try again.";
                } finally {
                  isSubmitting.value = false;
                }
              }}
            >
              <div class="form-group">
                <label for="email">Email</label>
                <input
                  type="email"
                  id="email"
                  value={email.value}
                  onInput$={(e) => {
                    email.value = (e.target as HTMLInputElement).value;
                  }}
                  placeholder="you@example.com"
                  required
                  disabled={isSubmitting.value}
                />
              </div>

              <button
                type="submit"
                class="btn-primary btn-full"
                disabled={isSubmitting.value}
              >
                {isSubmitting.value ? "Sending..." : "Send Reset Link"}
              </button>
            </form>

            <p class="auth-footer">
              <a href="/login">Back to Login</a>
            </p>
          </>
        )}
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Forgot Password | App",
};
