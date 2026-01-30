import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { type DocumentHead, useNavigate, useLocation } from "@builder.io/qwik-city";
import {
  login,
  getUser,
  getIsLoading,
  getIsAuthenticated,
  subscribe,
} from "~/lib/stores/auth";
import "./index.css";

export default component$(() => {
  const nav = useNavigate();
  const loc = useLocation();

  const email = useSignal("");
  const password = useSignal("");
  const error = useSignal("");
  const isSubmitting = useSignal(false);
  const isNavigating = useSignal(false);

  useVisibleTask$(({ track }) => {
    track(() => isNavigating.value);

    const unsub = subscribe(() => {
      const isAuth = getIsAuthenticated();
      const loading = getIsLoading();
      const user = getUser();
      if (!isNavigating.value && !loading && isAuth) {
        isNavigating.value = true;
        const redirect = loc.url.searchParams.get("redirect");
        if (redirect && redirect.startsWith("/")) {
          nav(redirect);
        } else if (user?.role === "admin") {
          nav("/admin");
        } else {
          nav("/dashboard");
        }
      }
    });
    return () => unsub();
  });

  return (
    <div class="auth-page">
      <div class="auth-card card">
        <h1>Sign In</h1>
        <p class="subtitle">Welcome back! Please sign in to continue.</p>

        {error.value && (
          <div class="alert alert-error">{error.value}</div>
        )}

        <form
          preventdefault:submit
          onSubmit$={async () => {
            error.value = "";
            isSubmitting.value = true;

            const result = await login(email.value, password.value);

            if (result.success) {
              isNavigating.value = true;
              const user = getUser();
              const redirect = loc.url.searchParams.get("redirect");
              if (redirect && redirect.startsWith("/")) {
                nav(redirect);
              } else if (user?.role === "admin") {
                nav("/admin");
              } else {
                nav("/dashboard");
              }
            } else {
              error.value = result.error || "Login failed";
            }

            isSubmitting.value = false;
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

          <div class="form-group">
            <label for="password">Password</label>
            <input
              type="password"
              id="password"
              value={password.value}
              onInput$={(e) => {
                password.value = (e.target as HTMLInputElement).value;
              }}
              placeholder="Your password"
              required
              disabled={isSubmitting.value}
            />
          </div>

          <button
            type="submit"
            class="btn-primary btn-full"
            disabled={isSubmitting.value}
          >
            {isSubmitting.value ? "Signing in..." : "Sign In"}
          </button>
        </form>

        <p class="forgot-password-link">
          <a href="/forgot-password">Forgot password?</a>
        </p>

        <p class="auth-footer">
          Don't have an account? <a href="/register">Sign up</a>
        </p>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Login | App",
};
