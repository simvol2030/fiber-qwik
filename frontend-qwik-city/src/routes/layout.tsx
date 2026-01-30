import { component$, Slot, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { useNavigate } from "@builder.io/qwik-city";
import "./layout.css";
import {
  initAuth,
  getUser,
  getIsLoading,
  getIsAuthenticated,
  logout,
  subscribe,
} from "~/lib/stores/auth";

export default component$(() => {
  const nav = useNavigate();
  const userEmail = useSignal<string | null>(null);
  const isLoading = useSignal(true);
  const isAuthenticated = useSignal(false);

  useVisibleTask$(() => {
    initAuth();
    const unsub = subscribe(() => {
      const user = getUser();
      userEmail.value = user?.email ?? null;
      isLoading.value = getIsLoading();
      isAuthenticated.value = getIsAuthenticated();
    });
    return () => unsub();
  });

  return (
    <div class="app">
      <header class="app-header">
        <nav class="container app-nav">
          <a href="/" class="logo">
            App
          </a>
          <div class="nav-links">
            {isLoading.value ? (
              <span class="nav-loading">Loading...</span>
            ) : isAuthenticated.value ? (
              <>
                <a href="/dashboard">Dashboard</a>
                <span class="user-email">{userEmail.value}</span>
                <button
                  class="btn-secondary"
                  onClick$={async () => {
                    await logout();
                    nav("/login");
                  }}
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <a href="/login">Login</a>
                <a href="/register">Register</a>
              </>
            )}
          </div>
        </nav>
      </header>

      <main class="container app-main">
        <Slot />
      </main>

      <footer class="app-footer">
        <div class="container">
          <p>&copy; {new Date().getFullYear()} App. All rights reserved.</p>
        </div>
      </footer>
    </div>
  );
});
