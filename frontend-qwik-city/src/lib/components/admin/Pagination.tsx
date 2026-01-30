import { component$ } from "@builder.io/qwik";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  pageSize: number;
  onPageChange$: (page: number) => void;
}

export const Pagination = component$<PaginationProps>(
  ({ currentPage, totalPages, totalItems, pageSize, onPageChange$ }) => {
    const startItem = (currentPage - 1) * pageSize + 1;
    const endItem = Math.min(currentPage * pageSize, totalItems);

    const getVisiblePages = (): (number | string)[] => {
      const pages: (number | string)[] = [];
      const delta = 2;

      if (totalPages <= 7) {
        for (let i = 1; i <= totalPages; i++) pages.push(i);
      } else {
        pages.push(1);
        if (currentPage > delta + 2) pages.push("...");
        const start = Math.max(2, currentPage - delta);
        const end = Math.min(totalPages - 1, currentPage + delta);
        for (let i = start; i <= end; i++) pages.push(i);
        if (currentPage < totalPages - delta - 1) pages.push("...");
        pages.push(totalPages);
      }
      return pages;
    };

    if (totalPages <= 1) return null;

    return (
      <div class="pagination">
        <div class="pagination-info">
          Showing {startItem} to {endItem} of {totalItems} items
        </div>

        <div class="pagination-controls">
          <button
            class="page-btn"
            disabled={currentPage === 1}
            onClick$={() => onPageChange$(currentPage - 1)}
            aria-label="Previous page"
          >
            &larr;
          </button>

          {getVisiblePages().map((page, i) =>
            page === "..." ? (
              <span key={`ellipsis-${i}`} class="page-ellipsis">
                ...
              </span>
            ) : (
              <button
                key={`page-${page}`}
                class={`page-btn ${page === currentPage ? "active" : ""}`}
                onClick$={() => onPageChange$(page as number)}
                aria-label={`Page ${page}`}
                aria-current={page === currentPage ? "page" : undefined}
              >
                {page}
              </button>
            )
          )}

          <button
            class="page-btn"
            disabled={currentPage === totalPages}
            onClick$={() => onPageChange$(currentPage + 1)}
            aria-label="Next page"
          >
            &rarr;
          </button>
        </div>
      </div>
    );
  }
);
