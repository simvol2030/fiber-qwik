import { component$, useSignal, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { useNavigate } from "@builder.io/qwik-city";
import { FormBuilder, type FormField } from "~/lib/components/admin/FormBuilder";
import "~/lib/components/admin/FormBuilder.css";
import { adminApi } from "~/lib/api/admin";
import { toast } from "~/lib/stores/admin";
import "./index.css";

export default component$(() => {
  const nav = useNavigate();
  const loading = useSignal(false);

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
      label: "Password",
      type: "password",
      required: true,
      placeholder: "Enter password",
      validation: {
        minLength: 8,
        message: "Password must be at least 8 characters",
      },
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

  const initialData = {
    role: "user",
    isActive: true,
  };

  const handleSubmit = $(async (data: Record<string, unknown>) => {
    loading.value = true;

    try {
      const response = await adminApi.createUser({
        email: data.email as string,
        password: data.password as string,
        name: (data.name as string) || undefined,
        role: data.role as "user" | "admin",
        isActive: data.isActive as boolean,
      });

      if (response.success) {
        toast.success("User created successfully");
        nav("/admin/users");
      } else {
        toast.error(response.error?.message || "Failed to create user");
      }
    } catch (e) {
      toast.error("Failed to create user");
    } finally {
      loading.value = false;
    }
  });

  const handleCancel = $(() => {
    nav("/admin/users");
  });

  return (
    <div class="create-user-page">
      <div class="page-header">
        <a href="/admin/users" class="back-link">
          {"\u2190"} Back to Users
        </a>
        <h2 class="page-title">Create User</h2>
      </div>

      <div class="admin-card">
        <FormBuilder
          schema={schema}
          initialData={initialData}
          loading={loading.value}
          submitLabel="Create User"
          onSubmit$={handleSubmit}
          onCancel$={handleCancel}
        />
      </div>
    </div>
  );
});

export const head: DocumentHead = {
  title: "Create User | Admin | Box App",
};
