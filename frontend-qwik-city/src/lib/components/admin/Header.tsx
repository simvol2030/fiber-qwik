import { component$, useSignal, useVisibleTask$ } from "@builder.io/qwik";
import { useNavigate } from "@builder.io/qwik-city";
import {
  getTheme,
  toggleTheme,
  toggleMobileSidebar,
  subscribe as adminSubscribe,
} from "~/lib/stores/admin";
import {
  getUser,
  logout,
  subscribe as authSubscribe,
} from "~/lib/stores/auth";

interface HeaderProps {
  title?: string;
  sidebarCollapsed?: boolean;
}

export const Header = component$<HeaderProps>(
  ({ title = "Dashboard", sidebarCollapsed = false }) => {
    const nav = useNavigate();
    const theme = useSignal(getTheme());
    const user = useSignal(getUser());
    const showUserMenu = useSignal(false);

    useVisibleTask$(() => {
      const unsub1 = adminSubscribe(() => {
        theme.value = getTheme();
      });
      const unsub2 = authSubscribe(() => {
        user.value = getUser();
      });
      return () => {
        unsub1();
        unsub2();
      };
    });

    return (
      <header
        class={`admin-header ${sidebarCollapsed ? "sidebar-collapsed" : ""}`}
      >
        <div class="header-left">
          <button
            class="menu-btn mobile-only"
            onClick$={() => toggleMobileSidebar()}
            aria-label="Toggle menu"
          >
            &#9776;
          </button>
          <h1 class="page-title">{title}</h1>
        </div>

        <div class="header-right">
          <button
            class="theme-toggle"
            onClick$={() => toggleTheme()}
            aria-label="Toggle theme"
          >
            {theme.value === "light" ? "\uD83C\uDF19" : "\u2600\uFE0F"}
          </button>

          <div class="user-menu-wrapper">
            <button
              class="user-menu-trigger"
              onClick$={() => {
                showUserMenu.value = !showUserMenu.value;
              }}
            >
              <span class="user-avatar">
                {user.value?.name?.[0]?.toUpperCase() ||
                  user.value?.email?.[0]?.toUpperCase() ||
                  "A"}
              </span>
              <span class="user-name desktop-only">
                {user.value?.name || user.value?.email}
              </span>
              <span class="dropdown-arrow">{"\u25BC"}</span>
            </button>

            {showUserMenu.value && (
              <>
                <div class="user-menu">
                  <a
                    href="/admin/profile"
                    class="menu-item"
                    onClick$={() => {
                      showUserMenu.value = false;
                    }}
                  >
                    <span>{"\uD83D\uDC64"}</span> Profile
                  </a>
                  <a
                    href="/admin/settings"
                    class="menu-item"
                    onClick$={() => {
                      showUserMenu.value = false;
                    }}
                  >
                    <span>{"\u2699\uFE0F"}</span> Settings
                  </a>
                  <hr />
                  <button
                    class="menu-item logout"
                    onClick$={async () => {
                      await logout();
                      nav("/login");
                    }}
                  >
                    <span>{"\uD83D\uDEAA"}</span> Logout
                  </button>
                </div>
                <div
                  class="menu-overlay"
                  onClick$={() => {
                    showUserMenu.value = false;
                  }}
                  role="presentation"
                />
              </>
            )}
          </div>
        </div>
      </header>
    );
  }
);
