import { redirect } from "@sveltejs/kit";

// Define public routes that don't require authentication
const publicRoutes = [
  "/login",
  "/register",
  "/forgot-password",
  "/terms",
  "/privacy",
  "/test-api",
];

export function load({ url }) {
  // Check if the current route is public
  const isPublicRoute = publicRoutes.some((route) =>
    url.pathname.startsWith(route)
  );

  // For now, we'll let the client handle authentication checks
  // This avoids SSR issues with localStorage
  return {
    isPublicRoute,
    currentPath: url.pathname,
  };
}
