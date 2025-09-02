// Error handling utilities

export interface ApiError {
  error: string;
  details?: string;
  request_id: string;
  timestamp: string;
}

export function isApiError(error: any): error is ApiError {
  return (
    typeof error === "object" &&
    error !== null &&
    "error" in error &&
    "request_id" in error &&
    "timestamp" in error
  );
}

export function extractErrorMessage(error: any): string {
  if (isApiError(error)) {
    return error.error || "An API error occurred";
  }

  if (error instanceof Error) {
    return error.message;
  }

  if (typeof error === "string") {
    return error;
  }

  return "An unexpected error occurred";
}

export function extractErrorDetails(error: any): string | undefined {
  if (isApiError(error)) {
    return error.details;
  }

  return undefined;
}

export function getRequestId(error: any): string | undefined {
  if (isApiError(error)) {
    return error.request_id;
  }

  return undefined;
}

export function handleApiError(error: any, context: string): string {
  const message = extractErrorMessage(error);
  const details = extractErrorDetails(error);
  const requestId = getRequestId(error);

  console.error(`[${context}] API Error:`, {
    message,
    details,
    requestId,
    originalError: error,
  });

  // Return user-friendly error message
  if (message.includes("401") || message.includes("Unauthorized")) {
    return "Your session has expired. Please log in again.";
  }

  if (message.includes("403") || message.includes("Forbidden")) {
    return "You do not have permission to perform this action.";
  }

  if (message.includes("404") || message.includes("Not Found")) {
    return "The requested resource was not found.";
  }

  if (message.includes("422") || message.includes("Validation")) {
    return "Please check your input and try again.";
  }

  if (message.includes("500") || message.includes("Internal Server Error")) {
    return "A server error occurred. Please try again later.";
  }

  if (message.includes("Network") || message.includes("fetch")) {
    return "Network error. Please check your connection and try again.";
  }

  return message || "An unexpected error occurred. Please try again.";
}

export function isNetworkError(error: any): boolean {
  if (error instanceof TypeError && error.message.includes("fetch")) {
    return true;
  }

  if (error instanceof Error) {
    return error.message.includes("Network") || error.message.includes("fetch");
  }

  return false;
}

export function isAuthError(error: any): boolean {
  const message = extractErrorMessage(error);
  return (
    message.includes("401") ||
    message.includes("Unauthorized") ||
    message.includes("403") ||
    message.includes("Forbidden")
  );
}

export function isValidationError(error: any): boolean {
  const message = extractErrorMessage(error);
  return (
    message.includes("422") ||
    message.includes("Validation") ||
    message.includes("Invalid")
  );
}

export function isServerError(error: any): boolean {
  const message = extractErrorMessage(error);
  return message.includes("500") || message.includes("Internal Server Error");
}

export function shouldRetry(error: any): boolean {
  // Don't retry auth errors, validation errors, or client errors
  if (isAuthError(error) || isValidationError(error)) {
    return false;
  }

  // Retry network errors and server errors
  return isNetworkError(error) || isServerError(error);
}

export function getRetryDelay(attempt: number, baseDelay = 1000): number {
  // Exponential backoff with jitter
  const exponentialDelay = baseDelay * Math.pow(2, attempt - 1);
  const jitter = Math.random() * 0.1 * exponentialDelay;
  return Math.min(exponentialDelay + jitter, 30000); // Max 30 seconds
}
