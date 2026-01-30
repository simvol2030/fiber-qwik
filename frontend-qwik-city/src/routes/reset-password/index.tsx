import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { type DocumentHead, useLocation } from "@builder.io/qwik-city";
import { api } from "~/lib/api/client";
import "./index.css";

export default component$(() => {
  const loc = useLocation();

  const newPassword = useSignal("");
  const confirmPassword = useSignal("");
  const error = useSignal("");
  const success = useSignal(false);
  const isSubmitting = useSignal(false);
  const tokenValid = useSignal<boolean | null>(null);
  const tokenEmail = useSignal("");
  const token = useSignal("");

  useVisibleTask$(() => {
    const t = loc.url.searchParams.get("token") || "";
    token.value = t;

    if (!t) {
      tokenValid.value = false;
      return;
    }

    api
      .post<{ valid: boolean; email: string }>("/auth/validate-reset-token", {
        token: t,
      })
      .then((response) => {
        if (response.success && response.data?.valid) {
          tokenValid.value = true;
          tokenEmail.value = response.data.email;
        } else {
          tokenValid.value = false;
          error.value =
            response.error?.message || "Invalid or expired reset token";
        }
      })
      .catch(() => {
        tokenValid.value = false;
        error.value = "Failed to validate token";
      });
  });

  return (
    <div class="auth-page">
      <div class="auth-card card">
        {success.value ? (
          <div class="success-state">
            <h1>Password Reset</h1>
            <p class="success-message">
              Your password has been reset successfully. You can now sign in
              with your new password.
            </p>
            <a href="/login" class="btn-primary btn-full">
              Go to Login
            </a>
          </div>
        ) : tokenValid.value === null ? (
          <div class="loading-state">
            <p>Validating reset token...</p>
          </div>
        ) : tokenValid.value === false ? (
          <div class="error-state">
            <h1>Invalid Token</h1>
            <p class="error-description">
              {error.value ||
                "This password reset link is invalid or has expired."}
            </p>
            <a href="/forgot-password" class="btn-primary btn-full">
              Request New Link
            </a>
            <p class="auth-footer">
              <a href="/login">Back to Login</a>
            </p>
          </div>
        ) : (
          <>
            <h1>Reset Password</h1>
            <p class="subtitle">
              Enter a new password for {tokenEmail.value}.
            </p>

            {error.value && (
              <div class="alert alert-error">{error.value}</div>
            )}

            <form
              preventdefault:submit
              onSubmit$={async () => {
                error.value = "";

                if (newPassword.value.length < 8) {
                  error.value = "Password must be at least 8 characters";
                  return;
                }

                if (newPassword.value !== confirmPassword.value) {
                  error.value = "Passwords do not match";
                  return;
                }

                isSubmitting.value = true;

                try {
                  const response = await api.post("/auth/reset-password", {
                    token: token.value,
                    newPassword: newPassword.value,
                  });
                  if (response.success) {
                    success.value = true;
                  } else {
                    error.value =
                      response.error?.message || "Failed to reset password";
                  }
                } catch {
                  error.value = "Network error. Please try again.";
                } finally {
                  isSubmitting.value = false;
                }
              }}
            >
              <div class="form-group">
                <label for="newPassword">New Password</label>
                <input
                  type="password"
                  id="newPassword"
                  value={newPassword.value}
                  onInput$={(e) => {
                    newPassword.value = (e.target as HTMLInputElement).value;
                  }}
                  placeholder="At least 8 characters"
                  required
                  minLength={8}
                  disabled={isSubmitting.value}
                />
              </div>

              <div class="form-group">
                <label for="confirmPassword">Confirm Password</label>
                <input
                  type="password"
                  id="confirmPassword"
                  value={confirmPassword.value}
                  onInput$={(e) => {
                    confirmPassword.value = (
                      e.target as HTMLInputElement
                    ).value;
                  }}
                  placeholder="Confirm new password"
                  required
                  disabled={isSubmitting.value}
                />
              </div>

              <button
                type="submit"
                class="btn-primary btn-full"
                disabled={isSubmitting.value}
              >
                {isSubmitting.value ? "Resetting..." : "Reset Password"}
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
  title: "Reset Password | App",
};
