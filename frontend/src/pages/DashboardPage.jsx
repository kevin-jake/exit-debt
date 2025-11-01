import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useDebtsStore } from '@stores/debtsStore'
import { useContactsStore } from '@stores/contactsStore'
import { LoadingSpinner } from '@components/common/LoadingSpinner'
import { StatCard } from '@components/common/StatCard'
import { EmptyState } from '@components/common/EmptyState'
import { CreateDebtModal } from '@components/debts/CreateDebtModal'
import { CreateContactModal } from '@components/contacts/CreateContactModal'
import {
  formatCurrency,
  formatRelativeTime,
  getDaysUntilDue,
  getDueDateColor,
  getInitials,
} from '@utils/formatters'
import { ROUTES } from '@/routes/routes'

export const DashboardPage = () => {
  const navigate = useNavigate()
  const { debts, isLoading: debtsLoading, fetchDebts } = useDebtsStore()
  const { contacts, isLoading: contactsLoading, fetchContacts } = useContactsStore()
  const [isLoading, setIsLoading] = useState(true)
  const [showCreateDebtModal, setShowCreateDebtModal] = useState(false)
  const [showCreateContactModal, setShowCreateContactModal] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setIsLoading(true)
    try {
      await Promise.all([fetchDebts(), fetchContacts()])
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleDebtCreated = async () => {
    setShowCreateDebtModal(false)
    await fetchDebts()
  }

  const handleContactCreated = async () => {
    setShowCreateContactModal(false)
    await fetchContacts()
  }

  // Calculate totals
  const totalIOwe = debts
    .filter((debt) => debt.debt_type === 'i_owe')
    .reduce((sum, debt) => sum + parseFloat(debt.total_amount || 0), 0)

  const totalOwedToMe = debts
    .filter((debt) => debt.debt_type === 'owed_to_me')
    .reduce((sum, debt) => sum + parseFloat(debt.total_amount || 0), 0)

  // Get recent debts (last 6)
  const recentDebts = [...debts]
    .sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
    .slice(0, 6)

  // Get upcoming due dates
  const upcomingDueDates = debts
    .filter((debt) => debt.due_date)
    .map((debt) => ({
      ...debt,
      daysUntilDue: getDaysUntilDue(debt.due_date),
    }))
    .filter((debt) => debt.daysUntilDue !== null && debt.daysUntilDue <= 30)
    .sort((a, b) => a.daysUntilDue - b.daysUntilDue)
    .slice(0, 5)

  if (isLoading) {
    return <LoadingSpinner size="lg" message="Loading dashboard..." className="py-12" />
  }

  return (
    <div className="mx-auto max-w-7xl space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
          <p className="mt-1 text-muted-foreground">Overview of your debt tracking</p>
        </div>
      </div>

      {/* Quick Overview Section */}
      <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
        <StatCard
          title="Total Amount I Owe"
          value={formatCurrency(totalIOwe)}
          valueColor="text-destructive"
          iconBgColor="bg-destructive/10"
          iconColor="text-destructive"
          icon={
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"
              />
            </svg>
          }
        />

        <StatCard
          title="Total Amount Owed to Me"
          value={formatCurrency(totalOwedToMe)}
          valueColor="text-success"
          iconBgColor="bg-success/10"
          iconColor="text-success"
          icon={
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"
              />
            </svg>
          }
        />

        <StatCard
          title="Total Debts"
          value={debts.length}
          iconBgColor="bg-primary/10"
          iconColor="text-primary"
          icon={
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
          }
        />

        <StatCard
          title="Total Contacts"
          value={contacts.length}
          iconBgColor="bg-secondary/10"
          iconColor="text-secondary"
          icon={
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
          }
        />
      </div>

      {/* Quick Actions Section */}
      <div className="card p-6">
        <h2 className="mb-4 text-xl font-semibold text-foreground">Quick Actions</h2>
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <button
            onClick={() => setShowCreateDebtModal(true)}
            className="flex items-center space-x-3 rounded-lg border border-border bg-card p-4 text-left transition-colors hover:bg-muted/50"
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
              <svg
                className="h-5 w-5 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 4v16m8-8H4"
                />
              </svg>
            </div>
            <div>
              <h3 className="font-medium text-foreground">Add New Debt</h3>
              <p className="text-sm text-muted-foreground">Create a new debt entry</p>
            </div>
          </button>

          <button
            onClick={() => setShowCreateContactModal(true)}
            className="flex items-center space-x-3 rounded-lg border border-border bg-card p-4 text-left transition-colors hover:bg-muted/50"
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-success/10">
              <svg
                className="h-5 w-5 text-success"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"
                />
              </svg>
            </div>
            <div>
              <h3 className="font-medium text-foreground">Add Contact</h3>
              <p className="text-sm text-muted-foreground">Create a new contact</p>
            </div>
          </button>

          <button
            onClick={() => navigate(ROUTES.DEBTS)}
            className="flex items-center space-x-3 rounded-lg border border-border bg-card p-4 text-left transition-colors hover:bg-muted/50"
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-warning/10">
              <svg
                className="h-5 w-5 text-warning"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
                />
              </svg>
            </div>
            <div>
              <h3 className="font-medium text-foreground">View All Debts</h3>
              <p className="text-sm text-muted-foreground">Manage your debts</p>
            </div>
          </button>
        </div>
      </div>

      {/* Recent Debts Section */}
      <div className="card p-6">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-xl font-semibold text-foreground">Recent Debts</h2>
          <button
            onClick={() => navigate(ROUTES.DEBTS)}
            className="text-sm text-primary hover:text-primary/80"
          >
            View All →
          </button>
        </div>

        {recentDebts.length === 0 ? (
          <EmptyState
            icon="debts"
            title="No debts yet"
            description="Get started by adding your first debt entry."
            action={() => setShowCreateDebtModal(true)}
            actionLabel="Add Debt"
          />
        ) : (
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
            {recentDebts.map((debt) => (
              <div
                key={debt.id}
                className="cursor-pointer rounded-lg border border-border p-4 transition-shadow hover:shadow-md"
                onClick={() => navigate(ROUTES.DEBTS)}
              >
                <div className="mb-3 flex items-start justify-between">
                  <div className="flex items-center space-x-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
                      <span className="font-medium text-primary">
                        {getInitials(debt.contact?.name || 'Unknown')}
                      </span>
                    </div>
                    <div>
                      <h3 className="font-medium text-foreground">
                        {debt.contact?.name || 'Unknown Contact'}
                      </h3>
                      <p className="text-sm text-muted-foreground">
                        {debt.description || 'No description'}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="flex items-center justify-between">
                  <p
                    className={`font-semibold ${debt.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'}`}
                  >
                    {formatCurrency(parseFloat(debt.total_amount || 0))}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {formatRelativeTime(debt.updated_at)}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* Upcoming Due Dates Section */}
      {upcomingDueDates.length > 0 && (
        <div className="card p-6">
          <h2 className="mb-4 text-xl font-semibold text-foreground">Upcoming Due Dates</h2>
          <div className="space-y-3">
            {upcomingDueDates.map((debt) => (
              <div
                key={debt.id}
                className="flex cursor-pointer items-center justify-between rounded-lg border border-border p-3 transition-colors hover:bg-muted/50"
                onClick={() => navigate(ROUTES.DEBTS)}
              >
                <div className="flex items-center space-x-3">
                  <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10">
                    <span className="text-sm font-medium text-primary">
                      {getInitials(debt.contact?.name || 'Unknown')}
                    </span>
                  </div>
                  <div>
                    <h4 className="font-medium text-foreground">
                      {debt.contact?.name || 'Unknown Contact'}
                    </h4>
                    <p className="text-sm text-muted-foreground">
                      {debt.description || 'No description'}
                    </p>
                  </div>
                </div>

                <div className="text-right">
                  <p
                    className={`font-semibold ${debt.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'}`}
                  >
                    {formatCurrency(parseFloat(debt.total_amount || 0))}
                  </p>
                  <p className={`text-sm ${getDueDateColor(debt.daysUntilDue)}`}>
                    {debt.daysUntilDue < 0
                      ? `Overdue by ${Math.abs(debt.daysUntilDue)} days`
                      : debt.daysUntilDue === 0
                        ? 'Due today'
                        : `Due in ${debt.daysUntilDue} days`}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Create Debt Modal */}
      {showCreateDebtModal && (
        <CreateDebtModal
          onDebtCreated={handleDebtCreated}
          onClose={() => setShowCreateDebtModal(false)}
        />
      )}

      {/* Create Contact Modal */}
      {showCreateContactModal && (
        <CreateContactModal
          onContactCreated={handleContactCreated}
          onClose={() => setShowCreateContactModal(false)}
        />
      )}
    </div>
  )
}
