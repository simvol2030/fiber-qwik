import { component$, Slot, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { useLocation } from "@builder.io/qwik-city";
import { AdminLayout } from "~/lib/components/admin/AdminLayout";
import "~/lib/components/admin/AdminLayout.css";
import "~/lib/components/admin/Sidebar.css";
import "~/lib/components/admin/Header.css";
import "~/lib/components/admin/Toast.css";

export default component$(() => {
  const loc = useLocation();
  const title = useSignal("Admin");

  useVisibleTask$(({ track }) => {
    track(() => loc.url.pathname);

    const path = loc.url.pathname;

    if (path === "/admin" || path === "/admin/") {
      title.value = "Dashboard";
    } else if (path.includes("/admin/users/new")) {
      title.value = "Create User";
    } else if (path.match(/\/admin\/users\/[^/]+/)) {
      title.value = "Edit User";
    } else if (path.startsWith("/admin/users")) {
      title.value = "Users";
    } else if (path.startsWith("/admin/files")) {
      title.value = "Files";
    } else if (path.startsWith("/admin/settings")) {
      title.value = "Settings";
    } else if (path.startsWith("/admin/profile")) {
      title.value = "Profile";
    } else {
      title.value = "Admin";
    }
  });

  return (
    <AdminLayout title={title.value}>
      <Slot />
    </AdminLayout>
  );
});
