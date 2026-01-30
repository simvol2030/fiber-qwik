import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { useLocation } from "@builder.io/qwik-city";
import {
  getSidebarCollapsed,
  getSidebarMobileOpen,
  toggleSidebar,
  closeMobileSidebar,
  subscribe as adminSubscribe,
} from "~/lib/stores/admin";
import { getUser, subscribe as authSubscribe } from "~/lib/stores/auth";

interface MenuItem {
  icon: string;
  label: string;
  href: string;
}

const menuItems: MenuItem[] = [
  { icon: "\uD83D\uDCCA", label: "Dashboard", href: "/admin" },
  { icon: "\uD83D\uDC65", label: "Users", href: "/admin/users" },
  { icon: "\uD83D\uDCC1", label: "Files", href: "/admin/files" },
  { icon: "\u2699\uFE0F", label: "Settings", href: "/admin/settings" },
  { icon: "\uD83D\uDC64", label: "Profile", href: "/admin/profile" },
];

export const Sidebar = component$(() => {
  const loc = useLocation();
  const collapsed = useSignal(getSidebarCollapsed());
  const mobileOpen = useSignal(getSidebarMobileOpen());
  const userEmail = useSignal(getUser()?.email || "");

  useVisibleTask$(() => {
    const unsub1 = adminSubscribe(() => {
      collapsed.value = getSidebarCollapsed();
      mobileOpen.value = getSidebarMobileOpen();
    });
    const unsub2 = authSubscribe(() => {
      userEmail.value = getUser()?.email || "";
    });
    return () => {
      unsub1();
      unsub2();
    };
  });

  const currentPath = loc.url.pathname;

  function isActive(href: string): boolean {
    if (href === "/admin") {
      return currentPath === "/admin" || currentPath === "/admin/";
    }
    return currentPath.startsWith(href);
  }

  return (
    <>
      <aside
        class={`sidebar ${collapsed.value ? "collapsed" : ""} ${mobileOpen.value ? "mobile-open" : ""}`}
      >
        <div class="sidebar-header">
          <a
            href="/admin"
            class="logo"
            onClick$={() => closeMobileSidebar()}
          >
            {collapsed.value ? (
              <span class="logo-icon">A</span>
            ) : (
              <span class="logo-text">Admin Panel</span>
            )}
          </a>
          <button
            class="toggle-btn desktop-only"
            onClick$={() => toggleSidebar()}
            aria-label="Toggle sidebar"
          >
            {collapsed.value ? "\u2192" : "\u2190"}
          </button>
          <button
            class="toggle-btn mobile-only"
            onClick$={() => closeMobileSidebar()}
            aria-label="Close sidebar"
          >
            &#10005;
          </button>
        </div>

        <nav class="sidebar-nav">
          <ul>
            {menuItems.map((item) => (
              <li key={item.href}>
                <a
                  href={item.href}
                  class={`nav-item ${isActive(item.href) ? "active" : ""}`}
                  onClick$={() => closeMobileSidebar()}
                  title={collapsed.value ? item.label : undefined}
                >
                  <span class="nav-icon">{item.icon}</span>
                  {!collapsed.value && (
                    <span class="nav-label">{item.label}</span>
                  )}
                </a>
              </li>
            ))}
          </ul>
        </nav>

        <div class="sidebar-footer">
          {!collapsed.value && (
            <div class="user-info">
              <span class="user-email">{userEmail.value}</span>
              <span class="user-role">Administrator</span>
            </div>
          )}
          <a href="/" class="back-link" title="Back to site">
            <span class="nav-icon">{"\uD83C\uDFE0"}</span>
            {!collapsed.value && (
              <span class="nav-label">Back to Site</span>
            )}
          </a>
        </div>
      </aside>

      {mobileOpen.value && (
        <div
          class="sidebar-overlay"
          onClick$={() => closeMobileSidebar()}
          role="presentation"
        />
      )}
    </>
  );
});
