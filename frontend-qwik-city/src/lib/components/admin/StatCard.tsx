import { component$ } from "@builder.io/qwik";

interface StatCardProps {
  title: string;
  value: string | number;
  icon?: string;
  change?: number;
  changeLabel?: string;
  variant?: "default" | "primary" | "success" | "warning" | "error";
}

export const StatCard = component$<StatCardProps>(
  ({ title, value, icon, change, changeLabel, variant = "default" }) => {
    const isPositiveChange = change !== undefined && change >= 0;

    return (
      <div class={`stat-card stat-card-${variant}`}>
        {icon && <div class="stat-icon">{icon}</div>}
        <div class="stat-content">
          <span class="stat-title">{title}</span>
          <span class="stat-value">{value}</span>
          {change !== undefined && (
            <span
              class={`stat-change ${isPositiveChange ? "positive" : "negative"}`}
            >
              {isPositiveChange ? "\u2191" : "\u2193"} {Math.abs(change)}%
              {changeLabel && (
                <span class="change-label">{changeLabel}</span>
              )}
            </span>
          )}
        </div>
      </div>
    );
  }
);
