import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { useNavigate, useLocation } from "@builder.io/qwik-city";
import { FormBuilder, type FormField } from "~/lib/components/admin/FormBuilder";
import "~/lib/components/admin/FormBuilder.css";
import { adminApi, type AdminUser } from "~/lib/api/admin";
import { toast } from "~/lib/stores/admin";
import "./index.css";

export default component$(() => {
  const nav = useNavigate();
  const loc = useLocation();
  const userId = loc.params.id;

  const user = useSignal<AdminUser | null>(null);
  const loading = useSignal(true);
  const saving = useSignal(false);

  const schema: FormField[] = [
    {
      name: "email",
      label: "Email",
      type: "email",
      required: true,
      placeholder: "user@example.com",
    },
    {
      name: "password",
      label: "New Password",
      type: "password",
      placeholder: "Leave empty to keep current password",
      validation: {
        minLength: 8,
        message: "Password must be at least 8 characters",
      },
      hint: "Only fill if you want to change the password",
    },
    {
      name: "name",
      label: "Name",
      type: "text",
      placeholder: "Full name (optional)",
    },
    {
      name: "role",
      label: "Role",
      type: "select",
      required: true,
      options: [
        { value: "user", label: "User" },
        { value: "admin", label: "Admin" },
      ],
    },
    {
      name: "isActive",
      label: "Active",
      type: "checkbox",
      hint: "Allow user to log in",
    },
  ];

  useVisibleTask$(async () => {
    loading.value = true;

    try {
      const response = await adminApi.getUser(userId);
      if (response.success && response.data) {
        user.value = response.data;
      } else {
        toast.error(response.error?.message || "Failed to load user");
        nav("/admin/users");
      }
    } catch (e) {
      toast.error("Failed to load user");
      nav("/admin/users");
    } finally {
      loading.value = false;
    }
  });

  const handleSubmit = $(async (data: Record<string, unknown>) => {
    saving.value = true;

    try {
      const updateData: Record<string, unknown> = {
        email: data.email,
        name: data.name || null,
        role: data.role,
        isActive: data.isActive,
      };

      if (data.password) {
        updateData.password = data.password;
      }

      const response = await adminApi.updateUser(userId, updateData);

      if (response.success) {
        toast.success("User updated successfully");
        nav("/admin/users");
      } else {
        toast.error(response.error?.message || "Failed to update user");
      }
    } catch (e) {
      toast.error("Failed to update user");
    } finally {
      saving.value = false;
    }
  });

  const handleCancel = $(() => {
    nav("/admin/users");
  });

  const initialData = user.value
    ? {
        email: user.value.email,
        name: user.value.name || "",
        role: user.value.role,
        isActive: user.value.isActive,
        password: "",
      }
    : {};

  return (
    <div class="edit-user-page">
      <div class="page-header">
        <a href="/admin/users" class="back-link">
          {"\u2190"} Back to Users
        </a>
        <h2 class="page-title">Edit User</h2>
      </div>

      {loading.value ? (
        <div class="loading-state">
          <div class="spinner" />
          <span>Loading user...</span>
        </div>
      ) : user.value ? (
        <div class="admin-card">
          <div class="user-meta">
            <div class="meta-item">
              <span class="meta-label">ID:</span>
              <code class="meta-value">{user.value.id}</code>
            </div>
            <div class="meta-item">
              <span class="meta-label">Created:</span>
              <span class="meta-value">
                {new Date(user.value.createdAt).toLocaleString()}
              </span>
            </div>
            {user.value.lastLoginAt && (
              <div class="meta-item">
                <span class="meta-label">Last Login:</span>
                <span class="meta-value">
                  {new Date(user.value.lastLoginAt).toLocaleString()}
                </span>
              </div>
            )}
          </div>

          <hr class="divider" />

          <FormBuilder
            schema={schema}
            initialData={initialData}
            loading={saving.value}
            submitLabel="Save Changes"
            onSubmit$={handleSubmit}
            onCancel$={handleCancel}
          />
        </div>
      ) : null}
    </div>
  );
});

export const head: DocumentHead = {
  title: "Edit User | Admin | Box App",
};
