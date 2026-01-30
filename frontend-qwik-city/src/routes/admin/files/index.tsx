import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";
import type { DocumentHead } from "@builder.io/qwik-city";
import { ConfirmDialog } from "~/lib/components/admin/ConfirmDialog";
import "~/lib/components/admin/ConfirmDialog.css";
import { adminApi, type FileInfo } from "~/lib/api/admin";
import { api } from "~/lib/api/client";
import { toast } from "~/lib/stores/admin";
import "./index.css";

export default component$(() => {
  const files = useSignal<FileInfo[]>([]);
  const loading = useSignal(true);
  const currentDir = useSignal("");
  const totalSize = useSignal(0);

  const showDeleteConfirm = useSignal(false);
  const fileToDelete = useSignal<FileInfo | null>(null);

  const showUploadPanel = useSignal(false);
  const selectedFiles = useSignal<File[]>([]);
  const uploading = useSignal(false);
  const isDragging = useSignal(false);

  const allowedTypes = [
    "image/jpeg",
    "image/png",
    "image/gif",
    "image/webp",
    "application/pdf",
  ];
  const maxFileSize = 10 * 1024 * 1024;

  const loadFiles = $(async (dir: string = "") => {
    loading.value = true;

    try {
      const response = await adminApi.getFiles(dir);
      if (response.success && response.data) {
        files.value = response.data.files;
        currentDir.value = response.data.currentDir;
        totalSize.value = response.data.totalSize;
      } else {
        toast.error(response.error?.message || "Failed to load files");
      }
    } catch (e) {
      toast.error("Failed to load files");
    } finally {
      loading.value = false;
    }
  });

  useVisibleTask$(() => {
    loadFiles("");
  });

  const navigateToDir = $((path: string) => {
    loadFiles(path);
  });

  const navigateUp = $(() => {
    if (!currentDir.value) return;
    const parts = currentDir.value.split("/");
    parts.pop();
    loadFiles(parts.join("/"));
  });

  const handleDelete = $(async () => {
    if (!fileToDelete.value) return;

    try {
      const response = await adminApi.deleteFile(fileToDelete.value.path);
      if (response.success) {
        toast.success(
          `${fileToDelete.value.isDir ? "Folder" : "File"} deleted successfully`
        );
        showDeleteConfirm.value = false;
        fileToDelete.value = null;
        loadFiles(currentDir.value);
      } else {
        toast.error(response.error?.message || "Failed to delete");
      }
    } catch (e) {
      toast.error("Failed to delete");
    }
  });

  const addFiles = $((newFiles: File[]) => {
    const validated: File[] = [];
    for (const file of newFiles) {
      if (!allowedTypes.includes(file.type)) {
        toast.error(
          `${file.name}: File type not allowed. Use JPG, PNG, GIF, WebP, or PDF.`
        );
        continue;
      }
      if (file.size > maxFileSize) {
        toast.error(`${file.name}: File too large. Maximum size is 10MB.`);
        continue;
      }
      if (
        selectedFiles.value.some(
          (f) => f.name === file.name && f.size === file.size
        )
      ) {
        continue;
      }
      validated.push(file);
    }
    if (validated.length + selectedFiles.value.length > 10) {
      toast.error("Maximum 10 files per upload.");
      return;
    }
    selectedFiles.value = [...selectedFiles.value, ...validated];
  });

  const handleUpload = $(async () => {
    if (selectedFiles.value.length === 0 || uploading.value) return;
    uploading.value = true;

    const token = api.getAccessToken();
    if (!token) {
      const refreshed = await api.refreshToken();
      if (!refreshed) {
        toast.error("Authentication required. Please log in again.");
        uploading.value = false;
        return;
      }
    }

    const formData = new FormData();
    for (const file of selectedFiles.value) {
      formData.append("files", file);
    }

    try {
      const response = await fetch("/api/upload/multiple", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${api.getAccessToken()}`,
        },
        body: formData,
        credentials: "include",
      });

      const data = await response.json();

      if (data.success) {
        const count = selectedFiles.value.length;
        toast.success(
          `${count} file${count > 1 ? "s" : ""} uploaded successfully`
        );
        showUploadPanel.value = false;
        selectedFiles.value = [];
        loadFiles(currentDir.value);
      } else {
        toast.error(data.error?.message || "Upload failed");
      }
    } catch {
      toast.error("Upload failed. Please try again.");
    } finally {
      uploading.value = false;
    }
  });

  const formatSize = (bytes: number): string => {
    if (bytes === 0) return "0 B";
    const k = 1024;
    const sizes = ["B", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString();
  };

  const getFileIcon = (file: FileInfo): string => {
    if (file.isDir) return "\uD83D\uDCC1";
    const icons: Record<string, string> = {
      jpg: "\uD83D\uDDBC\uFE0F",
      jpeg: "\uD83D\uDDBC\uFE0F",
      png: "\uD83D\uDDBC\uFE0F",
      gif: "\uD83D\uDDBC\uFE0F",
      webp: "\uD83D\uDDBC\uFE0F",
      svg: "\uD83D\uDDBC\uFE0F",
      pdf: "\uD83D\uDCC4",
      doc: "\uD83D\uDCDD",
      docx: "\uD83D\uDCDD",
      xls: "\uD83D\uDCCA",
      xlsx: "\uD83D\uDCCA",
      zip: "\uD83D\uDCE6",
      mp4: "\uD83C\uDFAC",
      webm: "\uD83C\uDFAC",
      mp3: "\uD83C\uDFB5",
      txt: "\uD83D\uDCC3",
      json: "\uD83D\uDCCB",
    };
    return icons[file.extension.toLowerCase()] || "\uD83D\uDCC4";
  };

  const deleteTitle = fileToDelete.value
    ? `Delete ${fileToDelete.value.isDir ? "Folder" : "File"}`
    : "Delete";

  const deleteMessage = fileToDelete.value
    ? `Are you sure you want to delete '${fileToDelete.value.name}'? ${fileToDelete.value.isDir ? "All contents will be deleted." : ""} This action cannot be undone.`
    : "";

  return (
    <div class="files-page">
      <div class="page-header">
        <h2 class="page-title">Files</h2>
        <div class="page-header-actions">
          <span class="total-size">Total: {formatSize(totalSize.value)}</span>
          <button
            class="btn-upload"
            onClick$={() => {
              showUploadPanel.value = true;
              selectedFiles.value = [];
            }}
          >
            + Upload Files
          </button>
        </div>
      </div>

      {/* Upload Panel */}
      {showUploadPanel.value && (
        <div class="upload-panel admin-card">
          <div class="upload-panel-header">
            <h3>Upload Files</h3>
            <button
              class="upload-close"
              onClick$={() => {
                if (!uploading.value) {
                  showUploadPanel.value = false;
                  selectedFiles.value = [];
                }
              }}
              disabled={uploading.value}
            >
              &times;
            </button>
          </div>

          <div
            class={`drop-zone ${isDragging.value ? "dragging" : ""}`}
            role="button"
            tabIndex={0}
            preventdefault:dragover
            preventdefault:drop
            onDragOver$={() => {
              isDragging.value = true;
            }}
            onDragLeave$={() => {
              isDragging.value = false;
            }}
            onDrop$={(e: DragEvent) => {
              isDragging.value = false;
              if (e.dataTransfer?.files) {
                addFiles(Array.from(e.dataTransfer.files));
              }
            }}
            onClick$={() => {
              const input = document.getElementById(
                "file-upload-input"
              ) as HTMLInputElement;
              input?.click();
            }}
            onKeyDown$={(e: KeyboardEvent) => {
              if (e.key === "Enter") {
                const input = document.getElementById(
                  "file-upload-input"
                ) as HTMLInputElement;
                input?.click();
              }
            }}
          >
            <span class="drop-icon">{"\uD83D\uDCE4"}</span>
            <p class="drop-text">Drop files here or click to browse</p>
            <p class="drop-hint">
              JPG, PNG, GIF, WebP, PDF -- max 10MB each, up to 10 files
            </p>
          </div>

          <input
            type="file"
            id="file-upload-input"
            onChange$={(e: Event) => {
              const input = e.target as HTMLInputElement;
              if (input.files) {
                addFiles(Array.from(input.files));
              }
              input.value = "";
            }}
            multiple
            accept=".jpg,.jpeg,.png,.gif,.webp,.pdf"
            hidden
          />

          {selectedFiles.value.length > 0 && (
            <>
              <div class="selected-files">
                {selectedFiles.value.map((file, i) => (
                  <div key={`${file.name}-${i}`} class="selected-file">
                    <span class="selected-file-name">{file.name}</span>
                    <span class="selected-file-size">
                      {formatSize(file.size)}
                    </span>
                    {!uploading.value && (
                      <button
                        class="selected-file-remove"
                        onClick$={() => {
                          selectedFiles.value = selectedFiles.value.filter(
                            (_, idx) => idx !== i
                          );
                        }}
                      >
                        &times;
                      </button>
                    )}
                  </div>
                ))}
              </div>

              <div class="upload-actions">
                <button
                  class="btn-upload-start"
                  onClick$={handleUpload}
                  disabled={uploading.value}
                >
                  {uploading.value
                    ? "Uploading..."
                    : `Upload ${selectedFiles.value.length} file${selectedFiles.value.length > 1 ? "s" : ""}`}
                </button>
              </div>
            </>
          )}
        </div>
      )}

      {/* Breadcrumb */}
      <div class="breadcrumb">
        <button
          class={`breadcrumb-item ${!currentDir.value ? "active" : ""}`}
          onClick$={() => loadFiles("")}
        >
          {"\uD83C\uDFE0"} Root
        </button>
        {currentDir.value &&
          currentDir.value.split("/").map((part, i) => {
            const parts = currentDir.value.split("/");
            return (
              <span key={`${part}-${i}`}>
                <span class="breadcrumb-separator">/</span>
                <button
                  class={`breadcrumb-item ${i === parts.length - 1 ? "active" : ""}`}
                  onClick$={() =>
                    navigateToDir(parts.slice(0, i + 1).join("/"))
                  }
                >
                  {part}
                </button>
              </span>
            );
          })}
      </div>

      <div class="admin-card">
        {loading.value ? (
          <div class="loading-state">
            <div class="spinner" />
            <span>Loading files...</span>
          </div>
        ) : files.value.length === 0 ? (
          <div class="empty-state">
            <span class="empty-icon">{"\uD83D\uDCC2"}</span>
            <p>No files found</p>
            {currentDir.value && (
              <button class="btn-back" onClick$={navigateUp}>
                {"\u2190"} Go Back
              </button>
            )}
          </div>
        ) : (
          <div class="files-grid">
            {currentDir.value && (
              <button class="file-item file-back" onClick$={navigateUp}>
                <span class="file-icon">{"\u2B06\uFE0F"}</span>
                <span class="file-name">..</span>
                <span class="file-meta">Parent Directory</span>
              </button>
            )}

            {files.value.map((file) => (
              <div
                key={file.path}
                class={`file-item ${file.isDir ? "is-dir" : ""}`}
              >
                {file.isDir ? (
                  <button
                    class="file-content"
                    onClick$={() => navigateToDir(file.path)}
                  >
                    <span class="file-icon">{getFileIcon(file)}</span>
                    <span class="file-name">{file.name}</span>
                    <span class="file-meta">Folder</span>
                  </button>
                ) : (
                  <a
                    href={`/uploads/${file.path}`}
                    target="_blank"
                    class="file-content"
                    rel="noopener noreferrer"
                  >
                    <span class="file-icon">{getFileIcon(file)}</span>
                    <span class="file-name">{file.name}</span>
                    <span class="file-meta">
                      {formatSize(file.size)} {"\u2022"}{" "}
                      {formatDate(file.modTime)}
                    </span>
                  </a>
                )}
                <button
                  class="file-delete"
                  onClick$={() => {
                    fileToDelete.value = file;
                    showDeleteConfirm.value = true;
                  }}
                  title="Delete"
                >
                  {"\uD83D\uDDD1\uFE0F"}
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      <ConfirmDialog
        open={showDeleteConfirm.value}
        title={deleteTitle}
        message={deleteMessage}
        confirmLabel="Delete"
        variant="danger"
        onConfirm$={handleDelete}
        onCancel$={() => {
          showDeleteConfirm.value = false;
          fileToDelete.value = null;
        }}
      />
    </div>
  );
});

export const head: DocumentHead = {
  title: "Files | Admin | Box App",
};
