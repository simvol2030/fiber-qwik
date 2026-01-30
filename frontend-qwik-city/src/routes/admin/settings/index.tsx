import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { adminApi, type AppSetting } from "~/lib/api/admin";
import { toast } from "~/lib/stores/admin";
import "./index.css";

export default component$(() => {
  const settings = useSignal<AppSetting[]>([]);
  const loading = useSignal(true);
  const saving = useSignal(false);
  const modifiedValues = useSignal<Record<string, string>>({});

  const loadSettings = $(async () => {
    loading.value = true;

    try {
      const response = await adminApi.getSettings();
      if (response.success && response.data) {
        settings.value = response.data;
        const vals: Record<string, string> = {};
        response.data.forEach((s: AppSetting) => {
          vals[s.key] = s.value;
        });
        modifiedValues.value = vals;
      } else {
        toast.error(response.error?.message || "Failed to load settings");
      }
    } catch (e) {
      toast.error("Failed to load settings");
    } finally {
      loading.value = false;
    }
  });

  useVisibleTask$(() => {
    loadSettings();
  });

  const getGroupedSettings = (): Record<string, AppSetting[]> => {
    const groups: Record<string, AppSetting[]> = {};
    settings.value.forEach((setting) => {
      const group = setting.group || "other";
      if (!groups[group]) {
        groups[group] = [];
      }
      groups[group].push(setting);
    });
    return groups;
  };

  const hasChanges = (): boolean => {
    return settings.value.some(
      (s) => s.value !== modifiedValues.value[s.key]
    );
  };

  const saveSettings = $(async () => {
    saving.value = true;

    const changedSettings = settings.value
      .filter((s) => s.value !== modifiedValues.value[s.key])
      .map((s) => ({ key: s.key, value: modifiedValues.value[s.key] }));

    if (changedSettings.length === 0) {
      toast.info("No changes to save");
      saving.value = false;
      return;
    }

    try {
      const response = await adminApi.updateSettings(changedSettings);
      if (response.success) {
        toast.success("Settings saved successfully");
        loadSettings();
      } else {
        toast.error(response.error?.message || "Failed to save settings");
      }
    } catch (e) {
      toast.error("Failed to save settings");
    } finally {
      saving.value = false;
    }
  });

  const resetChanges = $(() => {
    const vals: Record<string, string> = {};
    settings.value.forEach((s) => {
      vals[s.key] = s.value;
    });
    modifiedValues.value = vals;
  });

  const getGroupTitle = (group: string): string => {
    const titles: Record<string, string> = {
      general: "General Settings",
      auth: "Authentication",
      other: "Other Settings",
    };
    return titles[group] || group.charAt(0).toUpperCase() + group.slice(1);
  };

  const getGroupIcon = (group: string): string => {
    const icons: Record<string, string> = {
      general: "\u2699\uFE0F",
      auth: "\uD83D\uDD10",
      other: "\uD83D\uDCCB",
    };
    return icons[group] || "\uD83D\uDCCB";
  };

  const groupedSettings = getGroupedSettings();

  return (
    <div class="settings-page">
      <div class="page-header">
        <h2 class="page-title">Settings</h2>
        {hasChanges() && (
          <div class="header-actions">
            <button
              class="btn-reset"
              onClick$={resetChanges}
              disabled={saving.value}
            >
              Reset
            </button>
            <button
              class="btn-save"
              onClick$={saveSettings}
              disabled={saving.value}
            >
              {saving.value && <span class="spinner" />}
              Save Changes
            </button>
          </div>
        )}
      </div>

      {loading.value ? (
        <div class="loading-state">
          <div class="spinner large" />
          <span>Loading settings...</span>
        </div>
      ) : (
        Object.entries(groupedSettings).map(([group, groupSettings]) => (
          <div key={group} class="settings-group">
            <h3 class="group-title">
              <span class="group-icon">{getGroupIcon(group)}</span>
              {getGroupTitle(group)}
            </h3>

            <div class="admin-card">
              {groupSettings.map((setting) => (
                <div key={setting.key} class="setting-item">
                  <div class="setting-info">
                    <label for={setting.key} class="setting-label">
                      {setting.label || setting.key}
                    </label>
                    <span class="setting-key">{setting.key}</span>
                  </div>

                  <div class="setting-control">
                    {setting.type === "boolean" ? (
                      <label class="toggle">
                        <input
                          type="checkbox"
                          id={setting.key}
                          checked={
                            modifiedValues.value[setting.key] === "true"
                          }
                          onChange$={(e: Event) => {
                            const target = e.target as HTMLInputElement;
                            modifiedValues.value = {
                              ...modifiedValues.value,
                              [setting.key]: target.checked
                                ? "true"
                                : "false",
                            };
                          }}
                        />
                        <span class="toggle-slider" />
                      </label>
                    ) : setting.type === "number" ? (
                      <input
                        type="number"
                        id={setting.key}
                        value={modifiedValues.value[setting.key]}
                        onInput$={(e: Event) => {
                          const target = e.target as HTMLInputElement;
                          modifiedValues.value = {
                            ...modifiedValues.value,
                            [setting.key]: target.value,
                          };
                        }}
                        class="setting-input"
                      />
                    ) : (
                      <input
                        type="text"
                        id={setting.key}
                        value={modifiedValues.value[setting.key]}
                        onInput$={(e: Event) => {
                          const target = e.target as HTMLInputElement;
                          modifiedValues.value = {
                            ...modifiedValues.value,
                            [setting.key]: target.value,
                          };
                        }}
                        class="setting-input"
                      />
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))
      )}
    </div>
  );
});

export const head: DocumentHead = {
  title: "Settings | Admin | Box App",
};
