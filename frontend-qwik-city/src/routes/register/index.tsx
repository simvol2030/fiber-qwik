import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { type DocumentHead, useNavigate } from "@builder.io/qwik-city";
import {
  register,
  getIsLoading,
  getIsAuthenticated,
  subscribe,
} from "~/lib/stores/auth";
import "./index.css";

export default component$(() => {
  const nav = useNavigate();

  const name = useSignal("");
  const email = useSignal("");
  const password = useSignal("");
  const confirmPassword = useSignal("");
  const error = useSignal("");
  const isSubmitting = useSignal(false);

  useVisibleTask$(() => {
    const unsub = subscribe(() => {
      if (!getIsLoading() && getIsAuthenticated()) {
        nav("/dashboard");
      }
    });
    return () => unsub();
  });

  return (
    <div class="auth-page">
      <div class="auth-card card">
        <h1>Create Account</h1>
        <p class="subtitle">Sign up to get started with App.</p>

        {error.value && (
          <div class="alert alert-error">{error.value}</div>
        )}

        <form
          preventdefault:submit
          onSubmit$={async () => {
            error.value = "";

            if (password.value !== confirmPassword.value) {
              error.value = "Passwords do not match";
              return;
            }

            if (password.value.length < 8) {
              error.value = "Password must be at least 8 characters";
              return;
            }

            isSubmitting.value = true;

            const result = await register(
              email.value,
              password.value,
              name.value || undefined
            );

            if (result.success) {
              nav("/dashboard");
            } else {
              error.value = result.error || "Registration failed";
            }

            isSubmitting.value = false;
          }}
        >
          <div class="form-group">
            <label for="name">Name (optional)</label>
            <input
              type="text"
              id="name"
              value={name.value}
              onInput$={(e) => {
                name.value = (e.target as HTMLInputElement).value;
              }}
              placeholder="Your name"
              disabled={isSubmitting.value}
            />
          </div>

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

          <div class="form-group">
            <label for="password">Password</label>
            <input
              type="password"
              id="password"
              value={password.value}
              onInput$={(e) => {
                password.value = (e.target as HTMLInputElement).value;
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
                confirmPassword.value = (e.target as HTMLInputElement).value;
              }}
              placeholder="Confirm your password"
              required
              disabled={isSubmitting.value}
            />
          </div>

          <button
            type="submit"
            class="btn-primary btn-full"
            disabled={isSubmitting.value}
          >
            {isSubmitting.value ? "Creating account..." : "Create Account"}
          </button>
        </form>

        <p class="auth-footer">
          Already have an account? <a href="/login">Sign in</a>
        </p>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Register | App",
};
