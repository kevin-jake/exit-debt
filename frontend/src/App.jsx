import { useEffect } from 'react'
import { BrowserRouter } from 'react-router-dom'
import { AppRoutes } from './routes'
import { useAuthStore } from '@stores/authStore'
import { initializeTheme } from '@stores/themeStore'
import { Layout } from '@components/layout/Layout'

function App() {
  const initAuth = useAuthStore((state) => state.init)
  useEffect(() => {
    // Initialize theme from localStorage
    initializeTheme()

    // Initialize auth state
    initAuth()
  }, [initAuth])

  return (
    <BrowserRouter>
      <Layout>
        <AppRoutes />
      </Layout>
    </BrowserRouter>
  )
}

export default App
