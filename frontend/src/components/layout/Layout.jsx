import { useState, useEffect } from 'react'
import { useLocation } from 'react-router-dom'
import { Navigation } from './Navigation'
import { Notifications } from '../notifications/Notifications'

export const Layout = ({ children }) => {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const location = useLocation()

  // Check if current page should show navigation
  const showNavigation =
    !location.pathname.startsWith('/login') &&
    !location.pathname.startsWith('/register') &&
    location.pathname !== '/'

  // Close mobile menu when route changes
  useEffect(() => {
    setIsMobileMenuOpen(false)
  }, [location.pathname])

  // Close mobile menu when clicking outside
  useEffect(() => {
    if (!isMobileMenuOpen) return

    const handleClickOutside = (event) => {
      if (
        !event.target.closest('.mobile-nav') &&
        !event.target.closest('.hamburger-btn')
      ) {
        setIsMobileMenuOpen(false)
      }
    }

    document.addEventListener('click', handleClickOutside)
    return () => document.removeEventListener('click', handleClickOutside)
  }, [isMobileMenuOpen])

  const toggleMobileMenu = () => {
    setIsMobileMenuOpen((prev) => !prev)
  }

  if (!showNavigation) {
    return (
      <>
        <div className="min-h-screen bg-background">{children}</div>
        <Notifications />
      </>
    )
  }

  return (
    <>
      <div className="min-h-screen bg-background">
        {/* Mobile Header (visible on small screens) */}
        <header className="sticky top-0 z-40 flex items-center justify-between border-b border-border bg-card px-4 py-3 lg:hidden">
          <div className="flex items-center space-x-3">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
              <svg
                className="h-5 w-5 text-primary-foreground"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"
                />
              </svg>
            </div>
            <h1 className="text-xl font-bold text-card-foreground">DebtTracker</h1>
          </div>

          {/* Hamburger Menu Button */}
          <button
            className="hamburger-btn rounded-lg p-2 transition-colors hover:bg-secondary"
            onClick={toggleMobileMenu}
            aria-label="Toggle navigation menu"
          >
            <svg className="h-6 w-6 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              {isMobileMenuOpen ? (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              ) : (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M4 6h16M4 12h16M4 18h16"
                />
              )}
            </svg>
          </button>
        </header>

        {/* Desktop Navigation (hidden on small screens) */}
        <div className="hidden lg:block">
          <Navigation />
        </div>

        {/* Mobile Navigation Overlay */}
        {isMobileMenuOpen && (
          <div className="fixed inset-0 z-50 lg:hidden">
            {/* Backdrop */}
            <div
              className="fixed inset-0 bg-black/50"
              onClick={() => setIsMobileMenuOpen(false)}
            />

            {/* Mobile Navigation Panel */}
            <div className="mobile-nav fixed left-0 top-0 h-full max-w-[85vw] transform bg-card transition-transform duration-300 ease-in-out">
              <Navigation mobile onNavigate={() => setIsMobileMenuOpen(false)} />
            </div>
          </div>
        )}

        {/* Main Content */}
        <main className="min-h-screen lg:ml-64">
          <div className="p-4 lg:p-6">{children}</div>
        </main>
      </div>

      {/* Global Notifications */}
      <Notifications />
    </>
  )
}

