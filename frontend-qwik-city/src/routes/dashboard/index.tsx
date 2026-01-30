import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { type DocumentHead, useNavigate } from "@builder.io/qwik-city";
import type { User } from "~/lib/api/client";
import {
  getUser,
  getIsLoading,
  getIsAuthenticated,
  subscribe,
} from "~/lib/stores/auth";
import "./index.css";

export default component$(() => {
  const nav = useNavigate();
  const user = useSignal<User | null>(null);

  useVisibleTask$(() => {
    const unsub = subscribe(() => {
      const loading = getIsLoading();
      const isAuth = getIsAuthenticated();
      if (!loading && !isAuth) {
        nav("/login?redirect=/dashboard");
        return;
      }
      user.value = getUser();
    });
    // Check initial state
    if (!getIsLoading() && !getIsAuthenticated()) {
      nav("/login?redirect=/dashboard");
    } else {
      user.value = getUser();
    }
    return () => unsub();
  });

  return (
    <div class="dashboard">
      <header class="dashboard-header">
        <h1>Dashboard</h1>
        <p class="welcome">
          Welcome back, {user.value?.name || user.value?.email}!
        </p>
      </header>

      <div class="dashboard-content">
        <div class="stats-grid">
          <div class="stat-card card">
            <h3>Account Status</h3>
            <p class="stat-value text-success">Active</p>
          </div>

          <div class="stat-card card">
            <h3>Email</h3>
            <p class="stat-value">{user.value?.email}</p>
          </div>

          <div class="stat-card card">
            <h3>Member Since</h3>
            <p class="stat-value">
              {user.value?.createdAt
                ? new Date(user.value.createdAt).toLocaleDateString()
                : "N/A"}
            </p>
          </div>

          {user.value?.role && (
            <div class="stat-card card">
              <h3>Role</h3>
              <p
                class={`stat-value ${user.value.role === "admin" ? "text-success" : ""}`}
              >
                {user.value.role === "admin" ? "Administrator" : "User"}
              </p>
            </div>
          )}
        </div>

        {user.value?.role === "admin" && (
          <section class="dashboard-section">
            <a href="/admin" class="admin-panel-link">
              <span class="admin-panel-icon">&#9881;</span>
              <span>
                <strong>Admin Panel</strong>
                <small>Manage users, files and settings</small>
              </span>
              <span class="admin-panel-arrow">&rarr;</span>
            </a>
          </section>
        )}

        <section class="dashboard-section">
          <h2>Quick Actions</h2>
          <div class="actions">
            <a href="/dashboard/profile" class="btn-secondary">
              Edit Profile
            </a>
            <a href="/dashboard/profile" class="btn-secondary">
              Change Password
            </a>
          </div>
        </section>

        <section class="dashboard-section">
          <h2>Recent Activity</h2>
          <div class="activity-list card">
            <p class="no-activity">No recent activity to show.</p>
          </div>
        </section>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Dashboard | App",
};
