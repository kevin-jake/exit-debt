import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useDebtsStore } from '@stores/debtsStore'
import { useContactsStore } from '@stores/contactsStore'
import { LoadingSpinner } from '@components/common/LoadingSpinner'
import { StatCard } from '@components/common/StatCard'
import { EmptyState } from '@components/common/EmptyState'
import { CreateDebtModal } from '@components/debts/CreateDebtModal'
import { CreateContactModal } from '@components/contacts/CreateContactModal'
import { DebtDetailsModal } from '@components/debts/DebtDetailsModal'
import { EditDebtModal } from '@components/debts/EditDebtModal'
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
  const [selectedDebt, setSelectedDebt] = useState(null)
  const [editingDebt, setEditingDebt] = useState(null)
  const [upcomingPayments, setUpcomingPayments] = useState([])

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setIsLoading(true)
    try {
      await Promise.all([fetchDebts(), fetchContacts(), fetchUpcomingPayments()])
    } catch (error) {
      console.error('Failed to load dashboard data:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const fetchUpcomingPayments = async () => {
    try {
      const { apiClient } = await import('@/api/client')
      const payments = await apiClient.getUpcomingPayments(7) // Get upcoming payments for next 7 days
      setUpcomingPayments(payments || [])
    } catch (error) {
      console.error('Failed to fetch upcoming payments:', error)
      setUpcomingPayments([])
    }
  }

  const handleDebtCreated = async () => {
    setShowCreateDebtModal(false)
    await Promise.all([fetchDebts(), fetchUpcomingPayments()])
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

  // Split upcoming payments into overdue and upcoming based on days_until_due
  const overduePayments = upcomingPayments
    .filter((payment) => payment.days_until_due < 0)
    .sort((a, b) => a.days_until_due - b.days_until_due) // Most overdue first
    .slice(0, 5)

  const upcomingDueDates = upcomingPayments
    .filter((payment) => payment.days_until_due >= 0)
    .sort((a, b) => a.days_until_due - b.days_until_due) // Soonest first
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
      <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
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
      </div>
      {(overduePayments.length > 0 || upcomingDueDates.length > 0) && (
        <div
          className={`grid grid-cols-1 gap-6 ${upcomingDueDates.length > 0 ? 'md:grid-cols-2' : ''}`}
        >
          {/* Overdue Payments Section */}
          {overduePayments.length > 0 && (
            <div className="card p-6">
              <div className="mb-4 flex items-center justify-between">
                <h2 className="text-xl font-semibold text-foreground">Overdue Payments</h2>
                <span className="rounded-full bg-destructive/10 px-3 py-1 text-sm font-medium text-destructive">
                  {overduePayments.length} Overdue
                </span>
              </div>
              <div className="space-y-3">
                {overduePayments.map((payment) => {
                  // Find the debt to navigate to it on click
                  const debt = debts.find((d) => d.id === payment.debt_list_id)
                  return (
                    <div
                      key={payment.debt_list_id}
                      className="flex cursor-pointer items-center justify-between rounded-lg border p-3 transition-colors hover:bg-muted/50"
                      onClick={() => debt && setSelectedDebt(debt)}
                    >
                      <div className="flex items-center space-x-3">
                        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10">
                          <span className="text-sm font-medium text-primary">
                            {getInitials(payment.contact_name || 'Unknown')}
                          </span>
                        </div>
                        <div>
                          <h4 className="font-medium text-foreground">
                            {payment.contact_name || 'Unknown Contact'}
                          </h4>
                          <p className="text-sm text-muted-foreground">
                            {payment.debt_type === 'i_owe' ? 'I Owe' : 'Owed to Me'}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p
                          className={`font-semibold ${payment.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'}`}
                        >
                          {formatCurrency(parseFloat(payment.amount || 0))}
                        </p>
                        <p className="text-sm font-medium text-destructive">
                          {Math.abs(payment.days_until_due)} day
                          {Math.abs(payment.days_until_due) === 1 ? '' : 's'} overdue
                        </p>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
          )}

          {/* Upcoming Installment Dates Section */}
          {upcomingDueDates.length > 0 && (
            <div className="card p-6">
              <h2 className="mb-4 text-xl font-semibold text-foreground">Upcoming Due Dates</h2>
              <div className="space-y-3">
                {upcomingDueDates.map((payment) => {
                  // Find the debt to navigate to it on click
                  const debt = debts.find((d) => d.id === payment.debt_list_id)
                  return (
                    <div
                      key={payment.debt_list_id}
                      className="flex cursor-pointer items-center justify-between rounded-lg border border-border p-3 transition-colors hover:bg-muted/50"
                      onClick={() => debt && setSelectedDebt(debt)}
                    >
                      <div className="flex items-center space-x-3">
                        <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10">
                          <span className="text-sm font-medium text-primary">
                            {getInitials(payment.contact_name || 'Unknown')}
                          </span>
                        </div>
                        <div>
                          <h4 className="font-medium text-foreground">
                            {payment.contact_name || 'Unknown Contact'}
                          </h4>
                          <p className="text-sm text-muted-foreground">
                            {payment.debt_type === 'i_owe' ? 'I Owe' : 'Owed to Me'}
                          </p>
                        </div>
                      </div>

                      <div className="text-right">
                        <p
                          className={`font-semibold ${payment.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'}`}
                        >
                          {formatCurrency(parseFloat(payment.amount || 0))}
                        </p>
                        <p className={`text-sm ${getDueDateColor(payment.days_until_due)}`}>
                          {payment.days_until_due === 0
                            ? 'Due today'
                            : `Due in ${payment.days_until_due} day${payment.days_until_due === 1 ? '' : 's'}`}
                        </p>
                      </div>
                    </div>
                  )
                })}
              </div>
            </div>
          )}
        </div>
      )}
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
                onClick={() => setSelectedDebt(debt)}
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
                        {debt.debt_type === 'i_owe' ? 'I Owe' : 'Owed to Me'}
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

      {/* Debt Details Modal */}
      {selectedDebt && !editingDebt && (
        <DebtDetailsModal
          debt={selectedDebt}
          onClose={() => setSelectedDebt(null)}
          onEdit={() => {
            setEditingDebt(selectedDebt)
            setSelectedDebt(null)
          }}
          onDelete={async () => {
            if (window.confirm('Are you sure you want to delete this debt?')) {
              try {
                await useDebtsStore.getState().deleteDebt(selectedDebt.id)
                setSelectedDebt(null)
                await Promise.all([fetchDebts(), fetchUpcomingPayments()])
              } catch (error) {
                console.error('Failed to delete debt:', error)
                alert('Failed to delete debt. Please try again.')
              }
            }
          }}
        />
      )}

      {/* Edit Debt Modal */}
      {editingDebt && (
        <EditDebtModal
          debt={editingDebt}
          onClose={() => setEditingDebt(null)}
          onDebtUpdated={async () => {
            setEditingDebt(null)
            await Promise.all([fetchDebts(), fetchUpcomingPayments()])
          }}
        />
      )}
    </div>
  )
}
