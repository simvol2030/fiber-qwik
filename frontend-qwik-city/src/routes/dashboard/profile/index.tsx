import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { type DocumentHead, useNavigate } from "@builder.io/qwik-city";
import { api, type User } from "~/lib/api/client";
import {
  getUser,
  getIsLoading,
  getIsAuthenticated,
  subscribe,
} from "~/lib/stores/auth";

export default component$(() => {
  const nav = useNavigate();
  const user = useSignal<User | null>(null);

  const editName = useSignal("");
  const nameError = useSignal("");
  const nameSuccess = useSignal("");
  const nameSaving = useSignal(false);

  const currentPassword = useSignal("");
  const newPassword = useSignal("");
  const confirmPassword = useSignal("");
  const passwordError = useSignal("");
  const passwordSuccess = useSignal("");
  const passwordSaving = useSignal(false);

  useVisibleTask$(() => {
    const unsub = subscribe(() => {
      const loading = getIsLoading();
      const isAuth = getIsAuthenticated();
      if (!loading && !isAuth) {
        nav("/login?redirect=/dashboard/profile");
        return;
      }
      const u = getUser();
      user.value = u;
      if (u && !editName.value) {
        editName.value = u.name || "";
      }
    });
    if (!getIsLoading()) {
      if (!getIsAuthenticated()) {
        nav("/login?redirect=/dashboard/profile");
      } else {
        const u = getUser();
        user.value = u;
        editName.value = u?.name || "";
      }
    }
    return () => unsub();
  });

  return (
    <div class="profile-page">
      <header class="profile-header">
        <h1>Profile</h1>
        <p class="profile-subtitle">Manage your account settings</p>
      </header>

      <section class="profile-card card">
        <h2>Account Information</h2>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">Email</span>
            <span class="info-value">{user.value?.email}</span>
          </div>
          <div class="info-item">
            <span class="info-label">Role</span>
            <span class="info-value">
              {user.value?.role === "admin" ? "Administrator" : "User"}
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">Member Since</span>
            <span class="info-value">
              {user.value?.createdAt
                ? new Date(user.value.createdAt).toLocaleDateString()
                : "N/A"}
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">Last Login</span>
            <span class="info-value">
              {user.value?.lastLoginAt
                ? new Date(user.value.lastLoginAt).toLocaleString()
                : "N/A"}
            </span>
          </div>
        </div>
      </section>

      <section class="profile-card card">
        <h2>Edit Name</h2>

        {nameError.value && (
          <div class="alert alert-error">{nameError.value}</div>
        )}
        {nameSuccess.value && (
          <div class="alert alert-success">{nameSuccess.value}</div>
        )}

        <form
          preventdefault:submit
          onSubmit$={async () => {
            nameError.value = "";
            nameSuccess.value = "";
            nameSaving.value = true;

            try {
              const response = await api.put<User>("/auth/profile", {
                name: editName.value || null,
              });

              if (response.success && response.data) {
                nameSuccess.value = "Name updated successfully";
                user.value = response.data;
              } else {
                nameError.value =
                  response.error?.message || "Failed to update name";
              }
            } catch {
              nameError.value = "Network error. Please try again.";
            } finally {
              nameSaving.value = false;
            }
          }}
        >
          <div class="form-group">
            <label for="name">Display Name</label>
            <input
              type="text"
              id="name"
              value={editName.value}
              onInput$={(e) => {
                editName.value = (e.target as HTMLInputElement).value;
              }}
              placeholder="Your name"
              maxLength={100}
              disabled={nameSaving.value}
            />
          </div>
          <button
            type="submit"
            class="btn-primary"
            disabled={nameSaving.value}
          >
            {nameSaving.value ? "Saving..." : "Save Name"}
          </button>
        </form>
      </section>

      <section class="profile-card card">
        <h2>Change Password</h2>

        {passwordError.value && (
          <div class="alert alert-error">{passwordError.value}</div>
        )}
        {passwordSuccess.value && (
          <div class="alert alert-success">{passwordSuccess.value}</div>
        )}

        <form
          preventdefault:submit
          onSubmit$={async () => {
            passwordError.value = "";
            passwordSuccess.value = "";

            if (newPassword.value.length < 8) {
              passwordError.value =
                "New password must be at least 8 characters";
              return;
            }

            if (newPassword.value !== confirmPassword.value) {
              passwordError.value = "Passwords do not match";
              return;
            }

            passwordSaving.value = true;

            try {
              const response = await api.put("/auth/change-password", {
                currentPassword: currentPassword.value,
                newPassword: newPassword.value,
              });

              if (response.success) {
                passwordSuccess.value = "Password changed successfully";
                currentPassword.value = "";
                newPassword.value = "";
                confirmPassword.value = "";
              } else {
                passwordError.value =
                  response.error?.message || "Failed to change password";
              }
            } catch {
              passwordError.value = "Network error. Please try again.";
            } finally {
              passwordSaving.value = false;
            }
          }}
        >
          <div class="form-group">
            <label for="currentPassword">Current Password</label>
            <input
              type="password"
              id="currentPassword"
              value={currentPassword.value}
              onInput$={(e) => {
                currentPassword.value = (e.target as HTMLInputElement).value;
              }}
              placeholder="Enter current password"
              required
              disabled={passwordSaving.value}
            />
          </div>

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
              disabled={passwordSaving.value}
            />
          </div>

          <div class="form-group">
            <label for="confirmPwd">Confirm New Password</label>
            <input
              type="password"
              id="confirmPwd"
              value={confirmPassword.value}
              onInput$={(e) => {
                confirmPassword.value = (e.target as HTMLInputElement).value;
              }}
              placeholder="Confirm new password"
              required
              disabled={passwordSaving.value}
            />
          </div>

          <button
            type="submit"
            class="btn-primary"
            disabled={passwordSaving.value}
          >
            {passwordSaving.value ? "Changing..." : "Change Password"}
          </button>
        </form>
      </section>

      <div class="profile-back">
        <a href="/dashboard">&larr; Back to Dashboard</a>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Profile | App",
};
