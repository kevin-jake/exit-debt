import { useAuthStore } from '@stores/authStore'

/**
 * Custom hook for authentication
 * Provides convenient access to auth state and methods
 */
export const useAuth = () => {
  const user = useAuthStore((state) => state.user)
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  const login = useAuthStore((state) => state.login)
  const register = useAuthStore((state) => state.register)
  const logout = useAuthStore((state) => state.logout)
  const init = useAuthStore((state) => state.init)

  return {
    user,
    isAuthenticated,
    login,
    register,
    logout,
    init,
  }
}

