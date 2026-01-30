import { component$, Slot, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { Sidebar } from "./Sidebar";
import { Header } from "./Header";
import { Toast } from "./Toast";
import {
  initTheme,
  initSidebar,
  getSidebarCollapsed,
  subscribe as adminSubscribe,
} from "~/lib/stores/admin";

interface AdminLayoutProps {
  title?: string;
}

export const AdminLayout = component$<AdminLayoutProps>(
  ({ title = "Dashboard" }) => {
    const sidebarCollapsed = useSignal(false);

    useVisibleTask$(() => {
      initTheme();
      initSidebar();
      sidebarCollapsed.value = getSidebarCollapsed();

      const unsub = adminSubscribe(() => {
        sidebarCollapsed.value = getSidebarCollapsed();
      });
      return unsub;
    });

    return (
      <div
        class={`admin-layout ${sidebarCollapsed.value ? "sidebar-collapsed" : ""}`}
      >
        <Sidebar />
        <Header title={title} sidebarCollapsed={sidebarCollapsed.value} />

        <main class="admin-main">
          <div class="admin-content">
            <Slot />
          </div>
        </main>

        <Toast />
      </div>
    );
  }
);
