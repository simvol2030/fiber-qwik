import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { getIsAuthenticated, subscribe } from "~/lib/stores/auth";
import "./index.css";

export default component$(() => {
  const isAuthenticated = useSignal(false);

  useVisibleTask$(() => {
    isAuthenticated.value = getIsAuthenticated();
    const unsub = subscribe(() => {
      isAuthenticated.value = getIsAuthenticated();
    });
    return () => unsub();
  });

  return (
    <div class="hero">
      <h1>Welcome to App</h1>
      <p class="subtitle">
        Production-ready Qwik City template with authentication
      </p>

      <div class="features">
        <div class="feature">
          <h3>Qwik City</h3>
          <p>
            Built with Qwik for resumable, instant-loading applications
          </p>
        </div>
        <div class="feature">
          <h3>Authentication</h3>
          <p>JWT-based auth with refresh tokens and secure cookies</p>
        </div>
        <div class="feature">
          <h3>TypeScript</h3>
          <p>Full type safety throughout the application</p>
        </div>
      </div>

      <div class="cta">
        {isAuthenticated.value ? (
          <a href="/dashboard" class="btn-primary cta-btn">
            Go to Dashboard
          </a>
        ) : (
          <>
            <a href="/register" class="btn-primary cta-btn">
              Get Started
            </a>
            <a href="/login" class="btn-secondary cta-btn">
              Sign In
            </a>
          </>
        )}
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Home | App",
};
