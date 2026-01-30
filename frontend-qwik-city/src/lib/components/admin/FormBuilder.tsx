import { component$, useSignal, useVisibleTask$, $ } from "@builder.io/qwik";

export interface FormField {
  name: string;
  label: string;
  type:
    | "text"
    | "email"
    | "password"
    | "number"
    | "textarea"
    | "select"
    | "checkbox"
    | "date";
  required?: boolean;
  placeholder?: string;
  options?: { value: string; label: string }[];
  validation?: {
    min?: number;
    max?: number;
    minLength?: number;
    maxLength?: number;
    pattern?: string;
    message?: string;
  };
  disabled?: boolean;
  hint?: string;
}

interface FormBuilderProps {
  schema: FormField[];
  initialData?: Record<string, unknown>;
  submitLabel?: string;
  cancelLabel?: string;
  loading?: boolean;
  onSubmit$: (data: Record<string, unknown>) => void;
  onCancel$?: () => void;
}

function validateField(field: FormField, value: unknown): string {
  const strValue = String(value ?? "");

  if (field.required && !value && value !== false && value !== 0) {
    return `${field.label} is required`;
  }

  if (!value && value !== false && value !== 0) return "";

  if (field.validation) {
    const v = field.validation;
    if (v.minLength && strValue.length < v.minLength) {
      return (
        v.message ||
        `${field.label} must be at least ${v.minLength} characters`
      );
    }
    if (v.maxLength && strValue.length > v.maxLength) {
      return (
        v.message ||
        `${field.label} must be at most ${v.maxLength} characters`
      );
    }
    if (field.type === "number") {
      const numValue = Number(value);
      if (v.min !== undefined && numValue < v.min) {
        return v.message || `${field.label} must be at least ${v.min}`;
      }
      if (v.max !== undefined && numValue > v.max) {
        return v.message || `${field.label} must be at most ${v.max}`;
      }
    }
    if (v.pattern) {
      const regex = new RegExp(v.pattern);
      if (!regex.test(strValue)) {
        return v.message || `${field.label} format is invalid`;
      }
    }
  }

  if (field.type === "email" && strValue) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(strValue)) {
      return "Please enter a valid email address";
    }
  }

  return "";
}

export const FormBuilder = component$<FormBuilderProps>(
  ({
    schema,
    initialData = {},
    submitLabel = "Submit",
    cancelLabel = "Cancel",
    loading = false,
    onSubmit$,
    onCancel$,
  }) => {
    const formData = useSignal<Record<string, unknown>>({ ...initialData });
    const errors = useSignal<Record<string, string>>({});
    const touched = useSignal<Record<string, boolean>>({});

    useVisibleTask$(({ track }) => {
      track(() => initialData);
      // Initialize form data with defaults
      const data = { ...initialData };
      schema.forEach((field) => {
        if (data[field.name] === undefined) {
          if (field.type === "checkbox") {
            data[field.name] = false;
          } else {
            data[field.name] = "";
          }
        }
      });
      formData.value = data;
    });

    const handleBlur = $((field: FormField) => {
      touched.value = { ...touched.value, [field.name]: true };
      const error = validateField(field, formData.value[field.name]);
      errors.value = { ...errors.value, [field.name]: error };
    });

    const handleInput = $((field: FormField, event: Event) => {
      const target = event.target as HTMLInputElement;
      let newVal: unknown;

      if (field.type === "checkbox") {
        newVal = target.checked;
      } else if (field.type === "number") {
        newVal = target.value ? Number(target.value) : "";
      } else {
        newVal = target.value;
      }

      formData.value = { ...formData.value, [field.name]: newVal };

      if (touched.value[field.name]) {
        const error = validateField(field, newVal);
        errors.value = { ...errors.value, [field.name]: error };
      }
    });

    const handleSubmit = $(() => {
      let hasErrors = false;
      const newTouched: Record<string, boolean> = {};
      const newErrors: Record<string, string> = {};

      schema.forEach((field) => {
        newTouched[field.name] = true;
        const error = validateField(field, formData.value[field.name]);
        newErrors[field.name] = error;
        if (error) hasErrors = true;
      });

      touched.value = newTouched;
      errors.value = newErrors;

      if (!hasErrors) {
        onSubmit$(formData.value);
      }
    });

    return (
      <form
        class="form-builder"
        preventdefault:submit
        onSubmit$={handleSubmit}
      >
        {schema.map((field) => (
          <div
            key={field.name}
            class={`form-field ${touched.value[field.name] && errors.value[field.name] ? "has-error" : ""}`}
          >
            {field.type !== "checkbox" && (
              <label for={field.name} class="field-label">
                {field.label}
                {field.required && <span class="required">*</span>}
              </label>
            )}

            {field.type === "textarea" ? (
              <textarea
                id={field.name}
                name={field.name}
                placeholder={field.placeholder}
                disabled={field.disabled || loading}
                value={String(formData.value[field.name] ?? "")}
                onInput$={(e) => handleInput(field, e)}
                onBlur$={() => handleBlur(field)}
                class="field-input"
                rows={4}
              />
            ) : field.type === "select" ? (
              <select
                id={field.name}
                name={field.name}
                disabled={field.disabled || loading}
                value={String(formData.value[field.name] ?? "")}
                onChange$={(e) => handleInput(field, e)}
                onBlur$={() => handleBlur(field)}
                class="field-input"
              >
                <option value="">
                  {field.placeholder || "Select..."}
                </option>
                {(field.options || []).map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
            ) : field.type === "checkbox" ? (
              <label class="checkbox-label">
                <input
                  type="checkbox"
                  id={field.name}
                  name={field.name}
                  disabled={field.disabled || loading}
                  checked={Boolean(formData.value[field.name])}
                  onChange$={(e) => handleInput(field, e)}
                  onBlur$={() => handleBlur(field)}
                />
                <span class="checkbox-text">{field.label}</span>
                {field.required && <span class="required">*</span>}
              </label>
            ) : (
              <input
                type={field.type}
                id={field.name}
                name={field.name}
                placeholder={field.placeholder}
                disabled={field.disabled || loading}
                value={String(formData.value[field.name] ?? "")}
                onInput$={(e) => handleInput(field, e)}
                onBlur$={() => handleBlur(field)}
                class="field-input"
                min={field.validation?.min}
                max={field.validation?.max}
              />
            )}

            {field.hint && !errors.value[field.name] && (
              <span class="field-hint">{field.hint}</span>
            )}

            {touched.value[field.name] && errors.value[field.name] && (
              <span class="field-error">{errors.value[field.name]}</span>
            )}
          </div>
        ))}

        <div class="form-actions">
          {onCancel$ && (
            <button
              type="button"
              class="btn-cancel"
              onClick$={onCancel$}
              disabled={loading}
            >
              {cancelLabel}
            </button>
          )}
          <button type="submit" class="btn-submit" disabled={loading}>
            {loading && <span class="form-spinner"></span>}
            {submitLabel}
          </button>
        </div>
      </form>
    );
  }
);
