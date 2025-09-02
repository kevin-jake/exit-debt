// Form validation utilities

export interface ValidationRule {
  required?: boolean;
  minLength?: number;
  maxLength?: number;
  pattern?: RegExp;
  custom?: (value: any) => string | null;
}

export interface ValidationResult {
  isValid: boolean;
  errors: Record<string, string>;
}

export function validateField(
  value: any,
  rules: ValidationRule,
  fieldName: string
): string | null {
  // Required validation
  if (
    rules.required &&
    (!value || (typeof value === "string" && !value.trim()))
  ) {
    return `${fieldName} is required`;
  }

  // Skip other validations if value is empty and not required
  if (!value || (typeof value === "string" && !value.trim())) {
    return null;
  }

  // Type-specific validations
  if (typeof value === "string") {
    // Min length validation
    if (rules.minLength && value.length < rules.minLength) {
      return `${fieldName} must be at least ${rules.minLength} characters`;
    }

    // Max length validation
    if (rules.maxLength && value.length > rules.maxLength) {
      return `${fieldName} must be no more than ${rules.maxLength} characters`;
    }

    // Pattern validation
    if (rules.pattern && !rules.pattern.test(value)) {
      return `${fieldName} format is invalid`;
    }
  }

  // Custom validation
  if (rules.custom) {
    const customError = rules.custom(value);
    if (customError) {
      return customError;
    }
  }

  return null;
}

export function validateForm(
  data: Record<string, any>,
  rules: Record<string, ValidationRule>
): ValidationResult {
  const errors: Record<string, string> = {};

  for (const [fieldName, fieldRules] of Object.entries(rules)) {
    const value = data[fieldName];
    const error = validateField(value, fieldRules, fieldName);

    if (error) {
      errors[fieldName] = error;
    }
  }

  return {
    isValid: Object.keys(errors).length === 0,
    errors,
  };
}

// Common validation rules
export const commonRules = {
  email: {
    pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    custom: (value: string) => {
      if (value && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
        return "Please enter a valid email address";
      }
      return null;
    },
  },

  phone: {
    pattern: /^[\d\s\-\+\(\)]+$/,
    custom: (value: string) => {
      if (value && !/^[\d\s\-\+\(\)]+$/.test(value)) {
        return "Please enter a valid phone number";
      }
      return null;
    },
  },

  required: {
    required: true,
  },

  name: {
    required: true,
    minLength: 2,
    maxLength: 100,
  },

  amount: {
    required: true,
    custom: (value: string) => {
      if (value && !/^\d+(\.\d{1,2})?$/.test(value)) {
        return "Please enter a valid amount (e.g., 100.50)";
      }
      if (value && parseFloat(value) <= 0) {
        return "Amount must be greater than 0";
      }
      return null;
    },
  },

  date: {
    required: true,
    custom: (value: string) => {
      if (value) {
        const date = new Date(value);
        if (isNaN(date.getTime())) {
          return "Please enter a valid date";
        }
        if (date < new Date()) {
          return "Date cannot be in the past";
        }
      }
      return null;
    },
  },
};

// Contact validation rules
export const contactValidationRules = {
  name: commonRules.name,
  email: { ...commonRules.email },
  phone: { ...commonRules.phone },
  notes: { maxLength: 500 },
};

// Debt validation rules
export const debtValidationRules = {
  contact_id: commonRules.required,
  total_amount: commonRules.amount,
  currency: { required: true, minLength: 3, maxLength: 3 },
  debt_type: { required: true, minLength: 2, maxLength: 50 },
  installment_plan: { required: true, minLength: 2, maxLength: 50 },
  description: { maxLength: 1000 },
  notes: { maxLength: 1000 },
};

// Payment validation rules
export const paymentValidationRules = {
  amount: commonRules.amount,
  due_date: commonRules.date,
};
