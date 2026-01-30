import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { useNavigate } from "@builder.io/qwik-city";
import { DataTable } from "~/lib/components/admin/DataTable";
import type { Column, Action } from "~/lib/components/admin/DataTable";
import { ConfirmDialog } from "~/lib/components/admin/ConfirmDialog";
import "~/lib/components/admin/DataTable.css";
import "~/lib/components/admin/SearchInput.css";
import "~/lib/components/admin/Pagination.css";
import "~/lib/components/admin/ConfirmDialog.css";
import { adminApi, type AdminUser, type ListParams } from "~/lib/api/admin";
import { toast } from "~/lib/stores/admin";
import "./index.css";

export default component$(() => {
  const nav = useNavigate();
  const users = useSignal<AdminUser[]>([]);
  const loading = useSignal(true);
  const totalItems = useSignal(0);
  const totalPages = useSignal(1);
  const currentPage = useSignal(1);
  const pageSize = useSignal(10);
  const searchValue = useSignal("");
  const sortBy = useSignal("created_at");
  const sortDir = useSignal<"asc" | "desc">("desc");

  const showDeleteConfirm = useSignal(false);
  const userToDelete = useSignal<AdminUser | null>(null);

  const columns: Column[] = [
    { key: "email", label: "Email", sortable: true },
    { key: "name", label: "Name", sortable: true },
    { key: "role", label: "Role", sortable: true, format: "badge" },
    { key: "isActive", label: "Status", format: "boolean" },
    { key: "createdAt", label: "Created", sortable: true, format: "date" },
    { key: "lastLoginAt", label: "Last Login", format: "datetime" },
  ];

  const actions: Action[] = [
    {
      label: "Edit",
      icon: "\u270F\uFE0F",
      variant: "primary",
      onClick: (item: Record<string, unknown>) =>
        nav(`/admin/users/${item.id}`),
    },
    {
      label: "Delete",
      icon: "\uD83D\uDDD1\uFE0F",
      variant: "danger",
      onClick: (item: Record<string, unknown>) => {
        userToDelete.value = item as unknown as AdminUser;
        showDeleteConfirm.value = true;
      },
    },
  ];

  const loadUsers = $(async () => {
    loading.value = true;

    const params: ListParams = {
      page: currentPage.value,
      pageSize: pageSize.value,
      search: searchValue.value || undefined,
      sortBy: sortBy.value,
      sortDir: sortDir.value,
    };

    try {
      const response = await adminApi.getUsers(params);
      if (response.success && response.data) {
        users.value = response.data.items;
        totalItems.value = response.data.total;
        totalPages.value = response.data.totalPages;
      } else {
        toast.error(response.error?.message || "Failed to load users");
      }
    } catch (e) {
      toast.error("Failed to load users");
    } finally {
      loading.value = false;
    }
  });

  useVisibleTask$(() => {
    loadUsers();
  });

  const handleSearch = $((value: string) => {
    searchValue.value = value;
    currentPage.value = 1;
    loadUsers();
  });

  const handleSort = $((key: string, dir: "asc" | "desc") => {
    sortBy.value = key;
    sortDir.value = dir;
    loadUsers();
  });

  const handlePageChange = $((page: number) => {
    currentPage.value = page;
    loadUsers();
  });

  const handleDelete = $(async () => {
    if (!userToDelete.value) return;

    try {
      const response = await adminApi.deleteUser(userToDelete.value.id);
      if (response.success) {
        toast.success("User deleted successfully");
        showDeleteConfirm.value = false;
        userToDelete.value = null;
        loadUsers();
      } else {
        toast.error(response.error?.message || "Failed to delete user");
      }
    } catch (e) {
      toast.error("Failed to delete user");
    }
  });

  const exportCSV = $(() => {
    const headers = [
      "Email",
      "Name",
      "Role",
      "Status",
      "Created",
      "Last Login",
    ];
    const rows = users.value.map((u) => [
      u.email,
      u.name || "",
      u.role,
      u.isActive ? "Active" : "Inactive",
      new Date(u.createdAt).toLocaleDateString(),
      u.lastLoginAt ? new Date(u.lastLoginAt).toLocaleDateString() : "",
    ]);

    const csv = [headers.join(","), ...rows.map((r) => r.join(","))].join("\n");

    const blob = new Blob([csv], { type: "text/csv" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `users-${new Date().toISOString().split("T")[0]}.csv`;
    a.click();
    URL.revokeObjectURL(url);

    toast.success("CSV exported successfully");
  });

  const deleteMessage = userToDelete.value
    ? `Are you sure you want to delete ${userToDelete.value.email || "this user"}? This action cannot be undone.`
    : "";

  return (
    <div class="users-page">
      <div class="page-header">
        <h2 class="page-title">Users</h2>
        <a href="/admin/users/new" class="btn-create">
          <span>{"\u2795"}</span> Add User
        </a>
      </div>

      <DataTable
        data={users.value as unknown as Record<string, unknown>[]}
        columns={columns}
        actions={actions}
        loading={loading.value}
        totalItems={totalItems.value}
        currentPage={currentPage.value}
        totalPages={totalPages.value}
        pageSize={pageSize.value}
        searchable
        sortable
        emptyMessage="No users found"
        onSearch$={handleSearch}
        onSort$={handleSort}
        onPageChange$={handlePageChange}
        onExport$={exportCSV}
      />

      <ConfirmDialog
        open={showDeleteConfirm.value}
        title="Delete User"
        message={deleteMessage}
        confirmLabel="Delete"
        variant="danger"
        onConfirm$={handleDelete}
        onCancel$={() => {
          showDeleteConfirm.value = false;
          userToDelete.value = null;
        }}
      />
    </div>
  );
});

export const head: DocumentHead = {
  title: "Users | Admin | Box App",
};
