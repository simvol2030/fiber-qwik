import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { StatCard } from "~/lib/components/admin/StatCard";
import "~/lib/components/admin/StatCard.css";
import { adminApi, type DashboardStats } from "~/lib/api/admin";
import { toast } from "~/lib/stores/admin";
import "./index.css";

export default component$(() => {
  const stats = useSignal<DashboardStats | null>(null);
  const loading = useSignal(true);
  const error = useSignal<string | null>(null);

  const loadStats = $(async () => {
    loading.value = true;
    error.value = null;

    try {
      const response = await adminApi.getDashboard();
      if (response.success && response.data) {
        stats.value = response.data;
      } else {
        error.value = response.error?.message || "Failed to load dashboard";
        toast.error(error.value);
      }
    } catch (e) {
      error.value = "Failed to load dashboard";
      toast.error(error.value);
    } finally {
      loading.value = false;
    }
  });

  useVisibleTask$(() => {
    loadStats();
  });

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  };

  return (
    <div class="dashboard">
      {loading.value ? (
        <div class="loading-state">
          <div class="spinner" />
          <span>Loading dashboard...</span>
        </div>
      ) : error.value ? (
        <div class="error-state">
          <span class="error-icon">{"\u26A0\uFE0F"}</span>
          <span>{error.value}</span>
          <button class="btn-retry" onClick$={loadStats}>
            Retry
          </button>
        </div>
      ) : stats.value ? (
        <>
          {/* Stats Grid */}
          <div class="stats-grid">
            <StatCard
              title="Total Users"
              value={stats.value.totalUsers}
              icon={"\uD83D\uDC65"}
              variant="primary"
            />
            <StatCard
              title="Active Users"
              value={stats.value.activeUsers}
              icon={"\u2705"}
              variant="success"
            />
            <StatCard
              title="Admin Users"
              value={stats.value.adminUsers}
              icon={"\uD83D\uDEE1\uFE0F"}
              variant="warning"
            />
            <StatCard
              title="New Today"
              value={stats.value.newUsersToday}
              icon={"\uD83D\uDCC8"}
              change={
                stats.value.newUsersThisWeek > 0
                  ? Math.round(
                      (stats.value.newUsersToday /
                        stats.value.newUsersThisWeek) *
                        100
                    )
                  : 0
              }
              changeLabel="of weekly"
            />
          </div>

          {/* Quick Stats */}
          <div class="quick-stats">
            <div class="admin-card">
              <h3 class="card-title">Registration Trends</h3>
              <div class="trend-stats">
                <div class="trend-item">
                  <span class="trend-value">{stats.value.newUsersToday}</span>
                  <span class="trend-label">Today</span>
                </div>
                <div class="trend-item">
                  <span class="trend-value">
                    {stats.value.newUsersThisWeek}
                  </span>
                  <span class="trend-label">This Week</span>
                </div>
                <div class="trend-item">
                  <span class="trend-value">
                    {stats.value.newUsersThisMonth}
                  </span>
                  <span class="trend-label">This Month</span>
                </div>
              </div>
            </div>
          </div>

          {/* Recent Users & Activity */}
          <div class="dashboard-grid">
            {/* Recent Users */}
            <div class="admin-card">
              <div class="card-header">
                <h3 class="card-title">Recent Users</h3>
                <a href="/admin/users" class="view-all">
                  View All {"\u2192"}
                </a>
              </div>
              {stats.value.recentUsers.length === 0 ? (
                <p class="empty-message">No users yet</p>
              ) : (
                <div class="recent-list">
                  {stats.value.recentUsers.map((user) => (
                    <div key={user.id} class="recent-item">
                      <div class="user-avatar">
                        {user.name?.[0]?.toUpperCase() ||
                          user.email[0].toUpperCase()}
                      </div>
                      <div class="user-info">
                        <span class="user-name">
                          {user.name || "No name"}
                        </span>
                        <span class="user-email">{user.email}</span>
                      </div>
                      <span class="user-date">
                        {formatDate(user.createdAt)}
                      </span>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Quick Actions */}
            <div class="admin-card">
              <h3 class="card-title">Quick Actions</h3>
              <div class="quick-actions">
                <a href="/admin/users/new" class="action-link">
                  <span class="action-icon">{"\u2795"}</span>
                  <span>Add New User</span>
                </a>
                <a href="/admin/users" class="action-link">
                  <span class="action-icon">{"\uD83D\uDC65"}</span>
                  <span>Manage Users</span>
                </a>
                <a href="/admin/files" class="action-link">
                  <span class="action-icon">{"\uD83D\uDCC1"}</span>
                  <span>View Files</span>
                </a>
                <a href="/admin/settings" class="action-link">
                  <span class="action-icon">{"\u2699\uFE0F"}</span>
                  <span>App Settings</span>
                </a>
              </div>
            </div>
          </div>
        </>
      ) : null}
    </div>
  );
});

export const head: DocumentHead = {
  title: "Admin Dashboard | Box App",
};
