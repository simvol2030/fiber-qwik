import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { FormBuilder, type FormField } from "~/lib/components/admin/FormBuilder";
import "~/lib/components/admin/FormBuilder.css";
import { toast } from "~/lib/stores/admin";
import {
  getUser,
  subscribe as authSubscribe,
} from "~/lib/stores/auth";
import "./index.css";

export default component$(() => {
  const user = useSignal(getUser());
  const saving = useSignal(false);
  const activeTab = useSignal<"profile" | "password">("profile");

  useVisibleTask$(() => {
    const unsub = authSubscribe(() => {
      user.value = getUser();
    });
    return unsub;
  });

  const profileSchema: FormField[] = [
    {
      name: "email",
      label: "Email",
      type: "email",
      required: true,
      disabled: true,
      hint: "Email cannot be changed",
    },
    {
      name: "name",
      label: "Name",
      type: "text",
      placeholder: "Your name",
    },
  ];

  const passwordSchema: FormField[] = [
    {
      name: "currentPassword",
      label: "Current Password",
      type: "password",
      required: true,
      placeholder: "Enter current password",
    },
    {
      name: "newPassword",
      label: "New Password",
      type: "password",
      required: true,
      placeholder: "Enter new password",
      validation: {
        minLength: 8,
        message: "Password must be at least 8 characters",
      },
    },
    {
      name: "confirmPassword",
      label: "Confirm New Password",
      type: "password",
      required: true,
      placeholder: "Confirm new password",
    },
  ];

  const profileData = user.value
    ? {
        email: user.value.email,
        name: user.value.name || "",
      }
    : {};

  const handleProfileSubmit = $(async (_data: Record<string, unknown>) => {
    saving.value = true;
    try {
      toast.info("Profile update coming soon");
    } catch (e) {
      toast.error("Failed to update profile");
    } finally {
      saving.value = false;
    }
  });

  const handlePasswordSubmit = $(async (data: Record<string, unknown>) => {
    if (data.newPassword !== data.confirmPassword) {
      toast.error("Passwords do not match");
      return;
    }

    saving.value = true;
    try {
      toast.info("Password change coming soon");
    } catch (e) {
      toast.error("Failed to change password");
    } finally {
      saving.value = false;
    }
  });

  return (
    <div class="profile-page">
      <div class="page-header">
        <h2 class="page-title">Profile</h2>
      </div>

      {/* User Info Card */}
      <div class="admin-card user-card">
        <div class="user-avatar">
          {user.value?.name?.[0]?.toUpperCase() ||
            user.value?.email?.[0]?.toUpperCase() ||
            "A"}
        </div>
        <div class="user-info">
          <h3 class="user-name">{user.value?.name || "No name set"}</h3>
          <p class="user-email">{user.value?.email}</p>
          <span class="user-role badge-admin">Administrator</span>
        </div>
      </div>

      {/* Tabs */}
      <div class="tabs">
        <button
          class={`tab ${activeTab.value === "profile" ? "active" : ""}`}
          onClick$={() => {
            activeTab.value = "profile";
          }}
        >
          Profile Information
        </button>
        <button
          class={`tab ${activeTab.value === "password" ? "active" : ""}`}
          onClick$={() => {
            activeTab.value = "password";
          }}
        >
          Change Password
        </button>
      </div>

      {/* Tab Content */}
      <div class="admin-card">
        {activeTab.value === "profile" ? (
          <FormBuilder
            schema={profileSchema}
            initialData={profileData}
            loading={saving.value}
            submitLabel="Update Profile"
            onSubmit$={handleProfileSubmit}
          />
        ) : (
          <FormBuilder
            schema={passwordSchema}
            loading={saving.value}
            submitLabel="Change Password"
            onSubmit$={handlePasswordSubmit}
          />
        )}
      </div>

      {/* Additional Info */}
      <div class="admin-card info-card">
        <h4 class="info-title">Account Information</h4>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">User ID</span>
            <code class="info-value">{user.value?.id || "-"}</code>
          </div>
          <div class="info-item">
            <span class="info-label">Role</span>
            <span class="info-value">Administrator</span>
          </div>
          <div class="info-item">
            <span class="info-label">Member Since</span>
            <span class="info-value">
              {user.value?.createdAt
                ? new Date(user.value.createdAt).toLocaleDateString()
                : "-"}
            </span>
          </div>
        </div>
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Profile | Admin | Box App",
};
