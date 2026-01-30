import { component$, useSignal, $ } from "@builder.io/qwik";
import { SearchInput } from "./SearchInput";
import { Pagination } from "./Pagination";

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  format?: "text" | "date" | "datetime" | "currency" | "badge" | "boolean";
  width?: string;
  badgeColors?: Record<string, string>;
}

export interface Action {
  label: string;
  icon?: string;
  variant?: "default" | "primary" | "danger";
  onClick: (item: Record<string, unknown>) => void;
}

interface DataTableProps {
  data: Record<string, unknown>[];
  columns: Column[];
  searchable?: boolean;
  sortable?: boolean;
  pageSize?: number;
  totalItems?: number;
  currentPage?: number;
  totalPages?: number;
  actions?: Action[];
  loading?: boolean;
  emptyMessage?: string;
  onSearch$?: (value: string) => void;
  onSort$?: (key: string, dir: "asc" | "desc") => void;
  onPageChange$?: (page: number) => void;
  onExport$?: () => void;
}

function formatValue(value: unknown, column: Column): string {
  if (value === null || value === undefined) return "-";
  switch (column.format) {
    case "date":
      return new Date(value as string).toLocaleDateString();
    case "datetime":
      return new Date(value as string).toLocaleString();
    case "currency":
      return new Intl.NumberFormat("en-US", {
        style: "currency",
        currency: "USD",
      }).format(value as number);
    case "boolean":
      return (value as boolean) ? "Yes" : "No";
    default:
      return String(value);
  }
}

function getBadgeClass(value: unknown, column: Column): string {
  if (!column.badgeColors) {
    const defaults: Record<string, string> = {
      admin: "badge-admin",
      user: "badge-user",
      active: "badge-active",
      inactive: "badge-inactive",
      true: "badge-active",
      false: "badge-inactive",
    };
    return defaults[String(value).toLowerCase()] || "badge-default";
  }
  return column.badgeColors[String(value)] || "badge-default";
}

export const DataTable = component$<DataTableProps>(
  ({
    data,
    columns,
    searchable = true,
    sortable = true,
    pageSize = 10,
    totalItems = 0,
    currentPage = 1,
    totalPages = 1,
    actions = [],
    loading = false,
    emptyMessage = "No data found",
    onSearch$,
    onSort$,
    onPageChange$,
    onExport$,
  }) => {
    const sortKey = useSignal("");
    const sortDir = useSignal<"asc" | "desc">("asc");
    const searchValue = useSignal("");

    const handleSort = $((key: string) => {
      if (!sortable) return;
      if (sortKey.value === key) {
        sortDir.value = sortDir.value === "asc" ? "desc" : "asc";
      } else {
        sortKey.value = key;
        sortDir.value = "asc";
      }
      onSort$?.(sortKey.value, sortDir.value);
    });

    const handleSearch = $((value: string) => {
      searchValue.value = value;
      onSearch$?.(value);
    });

    const colSpan = columns.length + (actions.length > 0 ? 1 : 0);

    return (
      <div class="data-table-wrapper">
        {(searchable || onExport$) && (
          <div class="table-toolbar">
            {searchable && (
              <div class="search-wrapper">
                <SearchInput
                  value={searchValue.value}
                  onSearch$={handleSearch}
                  placeholder="Search..."
                />
              </div>
            )}
            {onExport$ && (
              <button class="btn-export" onClick$={onExport$}>
                &#128229; Export CSV
              </button>
            )}
          </div>
        )}

        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                {columns.map((column) => (
                  <th
                    key={column.key}
                    style={column.width ? { width: column.width } : {}}
                    class={
                      sortable && column.sortable !== false ? "sortable" : ""
                    }
                    onClick$={() => {
                      if (column.sortable !== false) handleSort(column.key);
                    }}
                  >
                    <span class="th-content">
                      {column.label}
                      {sortable && column.sortable !== false && (
                        <span
                          class={`sort-indicator ${sortKey.value === column.key ? "active" : ""}`}
                        >
                          {sortKey.value === column.key
                            ? sortDir.value === "asc"
                              ? "\u2191"
                              : "\u2193"
                            : "\u2195"}
                        </span>
                      )}
                    </span>
                  </th>
                ))}
                {actions.length > 0 && (
                  <th class="actions-column">Actions</th>
                )}
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={colSpan} class="loading-cell">
                    <div class="loading-content">
                      <div class="loading-spinner"></div>
                      <span>Loading...</span>
                    </div>
                  </td>
                </tr>
              ) : data.length === 0 ? (
                <tr>
                  <td colSpan={colSpan} class="empty-cell">
                    {emptyMessage}
                  </td>
                </tr>
              ) : (
                data.map((item, rowIdx) => (
                  <tr key={rowIdx}>
                    {columns.map((column) => (
                      <td key={column.key}>
                        {column.format === "badge" ? (
                          <span
                            class={`badge ${getBadgeClass(item[column.key], column)}`}
                          >
                            {formatValue(item[column.key], column)}
                          </span>
                        ) : column.format === "boolean" ? (
                          <span
                            class={`badge ${item[column.key] ? "badge-active" : "badge-inactive"}`}
                          >
                            {formatValue(item[column.key], column)}
                          </span>
                        ) : (
                          formatValue(item[column.key], column)
                        )}
                      </td>
                    ))}
                    {actions.length > 0 && (
                      <td class="actions-cell">
                        {actions.map((action, aIdx) => (
                          <button
                            key={aIdx}
                            class={`action-btn action-${action.variant || "default"}`}
                            onClick$={() => action.onClick(item)}
                            title={action.label}
                          >
                            {action.icon && <span>{action.icon}</span>}
                            <span class="action-label">{action.label}</span>
                          </button>
                        ))}
                      </td>
                    )}
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>

        {onPageChange$ && totalPages > 1 && (
          <Pagination
            currentPage={currentPage}
            totalPages={totalPages}
            totalItems={totalItems}
            pageSize={pageSize}
            onPageChange$={onPageChange$}
          />
        )}
      </div>
    );
  }
);
