import { useState, useEffect, useMemo } from 'react'
import { useDebtsStore } from '@stores/debtsStore'
import { useNotificationsStore } from '@stores/notificationsStore'
import { DebtDetailsModal } from './DebtDetailsModal'
import { EditDebtModal } from './EditDebtModal'
import { DeleteDebtModal } from './DeleteDebtModal'
import { LoadingSpinner } from '@components/common/LoadingSpinner'
import { EmptyState } from '@components/common/EmptyState'
import {
  formatCurrency,
  formatRelativeTime,
  getInitials,
  truncateText,
  getDebtStatus,
} from '@utils/formatters'

export const DebtsTable = () => {
  const { debts, isLoading, fetchDebts } = useDebtsStore()
  const showError = useNotificationsStore((state) => state.error)

  const [selectedDebt, setSelectedDebt] = useState(null)
  const [showDetailsModal, setShowDetailsModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [showDeleteDialog, setShowDeleteDialog] = useState(false)
  const [debtToDelete, setDebtToDelete] = useState(null)

  // Filter and search state
  const [searchQuery, setSearchQuery] = useState('')
  const [typeFilter, setTypeFilter] = useState('all') // all, i_owe, owed_to_me
  const [sortBy, setSortBy] = useState('updated_at')
  const [sortOrder, setSortOrder] = useState('desc')

  // Pagination
  const [currentPage, setCurrentPage] = useState(1)
  const itemsPerPage = 10

  useEffect(() => {
    loadDebts()
  }, [])

  const loadDebts = async () => {
    try {
      await fetchDebts()
    } catch (error) {
      showError('Failed to load debts. Please try again.')
    }
  }

  const filteredAndSortedDebts = useMemo(() => {
    let filtered = [...debts]

    // Apply type filter
    if (typeFilter !== 'all') {
      filtered = filtered.filter((debt) => debt.debt_type === typeFilter)
    }

    // Apply search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase()
      filtered = filtered.filter(
        (debt) =>
          debt.contact?.name?.toLowerCase().includes(query) ||
          debt.description?.toLowerCase().includes(query) ||
          debt.notes?.toLowerCase().includes(query)
      )
    }

    // Apply sorting
    filtered.sort((a, b) => {
      let aValue = a[sortBy]
      let bValue = b[sortBy]

      if (sortBy === 'total_amount') {
        aValue = parseFloat(aValue || 0)
        bValue = parseFloat(bValue || 0)
      } else if (sortBy === 'created_at' || sortBy === 'updated_at') {
        aValue = new Date(aValue).getTime()
        bValue = new Date(bValue).getTime()
      } else if (typeof aValue === 'string' && typeof bValue === 'string') {
        aValue = aValue.toLowerCase()
        bValue = bValue.toLowerCase()
      }

      if (aValue === null) aValue = ''
      if (bValue === null) bValue = ''

      if (sortOrder === 'asc') {
        return aValue < bValue ? -1 : aValue > bValue ? 1 : 0
      } else {
        return aValue > bValue ? -1 : aValue < bValue ? 1 : 0
      }
    })

    return filtered
  }, [debts, searchQuery, typeFilter, sortBy, sortOrder])

  const totalPages = Math.ceil(filteredAndSortedDebts.length / itemsPerPage)
  const paginatedDebts = filteredAndSortedDebts.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  )

  const handleSort = (column) => {
    if (sortBy === column) {
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')
    } else {
      setSortBy(column)
      setSortOrder('asc')
    }
  }

  const viewDebt = (debt) => {
    setSelectedDebt(debt)
    setShowDetailsModal(true)
  }

  const editDebt = (debt) => {
    setSelectedDebt(debt)
    setShowEditModal(true)
  }

  const confirmDeleteDebt = (debt) => {
    setDebtToDelete(debt)
    setShowDeleteDialog(true)
  }

  const handleDebtUpdated = () => {
    setShowEditModal(false)
    setSelectedDebt(null)
    loadDebts()
  }

  const handleDebtDeleted = () => {
    setShowDeleteDialog(false)
    setDebtToDelete(null)
    loadDebts()
  }

  if (isLoading) {
    return <LoadingSpinner size="lg" message="Loading debts..." className="py-12" />
  }

  return (
    <div className="space-y-6">
      {/* Header with Search and Filters */}
      <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
        <div className="max-w-md flex-1">
          <div className="relative">
            <svg
              className="absolute left-3 top-1/2 h-5 w-5 -translate-y-1/2 transform text-muted-foreground"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder="Search debts..."
              className="input pl-10"
            />
          </div>
        </div>

        {/* Type Filter */}
        <div className="flex items-center space-x-2">
          <button
            onClick={() => setTypeFilter('all')}
            className={`rounded-lg px-4 py-2 text-sm font-medium transition-colors ${
              typeFilter === 'all'
                ? 'bg-primary text-primary-foreground'
                : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
            }`}
          >
            All
          </button>
          <button
            onClick={() => setTypeFilter('i_owe')}
            className={`rounded-lg px-4 py-2 text-sm font-medium transition-colors ${
              typeFilter === 'i_owe'
                ? 'bg-destructive text-destructive-foreground'
                : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
            }`}
          >
            I Owe
          </button>
          <button
            onClick={() => setTypeFilter('owed_to_me')}
            className={`rounded-lg px-4 py-2 text-sm font-medium transition-colors ${
              typeFilter === 'owed_to_me'
                ? 'bg-success text-success-foreground'
                : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
            }`}
          >
            Owed to Me
          </button>
        </div>
      </div>

      {filteredAndSortedDebts.length === 0 ? (
        <EmptyState
          icon="debts"
          title="No debts found"
          description={
            searchQuery || typeFilter !== 'all'
              ? 'Try adjusting your search or filters.'
              : 'Get started by adding your first debt entry.'
          }
        />
      ) : (
        <>
          {/* Desktop Table */}
          <div className="card hidden overflow-hidden lg:block">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="border-b border-border bg-muted/50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Contact
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Description
                    </th>
                    <th
                      className="cursor-pointer px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground"
                      onClick={() => handleSort('total_amount')}
                    >
                      Amount {sortBy === 'total_amount' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Type
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Status
                    </th>
                    <th
                      className="cursor-pointer px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground"
                      onClick={() => handleSort('updated_at')}
                    >
                      Updated {sortBy === 'updated_at' && (sortOrder === 'asc' ? '↑' : '↓')}
                    </th>
                    <th className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border bg-card">
                  {paginatedDebts.map((debt) => (
                    <tr
                      key={debt.id}
                      className="cursor-pointer transition-colors duration-200 hover:bg-muted/30"
                      onClick={() => viewDebt(debt)}
                    >
                      <td className="px-6 py-4">
                        <div className="flex items-center">
                          <div className="mr-3 flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10">
                            <span className="text-sm font-medium text-primary">
                              {getInitials(debt.contact?.name || 'Unknown')}
                            </span>
                          </div>
                          <div
                            className="max-w-[150px] truncate text-sm font-medium text-foreground"
                            title={debt.contact?.name || 'Unknown Contact'}
                          >
                            {debt.contact?.name || 'Unknown Contact'}
                          </div>
                        </div>
                      </td>
                      <td className="px-6 py-4">
                        <div
                          className="max-w-[200px] truncate text-sm text-muted-foreground"
                          title={debt.description || 'No description'}
                        >
                          {debt.description || 'No description'}
                        </div>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        <span
                          className={`text-sm font-semibold ${
                            debt.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'
                          }`}
                        >
                          {formatCurrency(parseFloat(debt.total_amount || 0))}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        <span
                          className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                            debt.debt_type === 'i_owe'
                              ? 'bg-destructive/10 text-destructive'
                              : 'bg-success/10 text-success'
                          }`}
                        >
                          {debt.debt_type === 'i_owe' ? 'I Owe' : 'Owed to Me'}
                        </span>
                      </td>
                      <td className="whitespace-nowrap px-6 py-4">
                        {(() => {
                          var status = getDebtStatus(debt.due_date)
                          if (debt.status && debt.status !== 'active') {
                            switch (debt.status) {
                              case 'settled':
                                status = {
                                  label: 'Settled',
                                  color: 'text-success',
                                  bgColor: 'bg-success/10',
                                }
                                break
                              case 'overdue':
                                status = {
                                  label: 'Overdue',
                                  color: 'text-destructive',
                                  bgColor: 'bg-destructive/10',
                                }
                                break
                              case 'archived':
                                status = {
                                  label: 'Archived',
                                  color: 'text-muted-foreground',
                                  bgColor: 'bg-muted/50',
                                }
                                break
                              default:
                                status = {
                                  label: 'Active',
                                  color: 'text-muted-foreground',
                                  bgColor: 'bg-muted/50',
                                }
                                break
                            }
                          }
                          return (
                            <span
                              className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${status.color} ${status.bgColor}`}
                            >
                              {status.label}
                            </span>
                          )
                        })()}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm text-muted-foreground">
                        {formatRelativeTime(debt.updated_at)}
                      </td>
                      <td className="whitespace-nowrap px-6 py-4 text-sm font-medium">
                        <div
                          className="flex items-center space-x-2"
                          onClick={(e) => e.stopPropagation()}
                        >
                          <button
                            onClick={() => viewDebt(debt)}
                            className="p-1 text-primary hover:text-primary/80"
                            title="View Details"
                          >
                            <svg
                              className="h-4 w-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                              />
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                              />
                            </svg>
                          </button>
                          <button
                            onClick={() => editDebt(debt)}
                            className="p-1 text-primary hover:text-primary/80"
                            title="Edit"
                          >
                            <svg
                              className="h-4 w-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                              />
                            </svg>
                          </button>
                          <button
                            onClick={() => confirmDeleteDebt(debt)}
                            className="p-1 text-destructive hover:text-destructive/80"
                            title="Delete"
                          >
                            <svg
                              className="h-4 w-4"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                              />
                            </svg>
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Mobile Card Layout */}
          <div className="space-y-4 lg:hidden">
            {paginatedDebts.map((debt) => {
              const status = getDebtStatus(debt.due_date)
              return (
                <div
                  key={debt.id}
                  className="card cursor-pointer p-4"
                  onClick={() => viewDebt(debt)}
                >
                  <div className="mb-3 flex items-start justify-between">
                    <div className="flex items-center space-x-3">
                      <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-primary/10">
                        <span className="text-sm font-medium text-primary">
                          {getInitials(debt.contact?.name || 'Unknown')}
                        </span>
                      </div>
                      <div className="min-w-0 flex-1">
                        <div
                          className="truncate font-medium text-foreground"
                          title={debt.contact?.name || 'Unknown Contact'}
                        >
                          {debt.contact?.name || 'Unknown Contact'}
                        </div>
                        <div className="flex flex-wrap gap-1">
                          <span
                            className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${
                              debt.debt_type === 'i_owe'
                                ? 'bg-destructive/10 text-destructive'
                                : 'bg-success/10 text-success'
                            }`}
                          >
                            {debt.debt_type === 'i_owe' ? 'I Owe' : 'Owed to Me'}
                          </span>
                          <span
                            className={`inline-flex rounded-full px-2 py-1 text-xs font-medium ${status.color} ${status.bgColor}`}
                          >
                            {status.label}
                          </span>
                        </div>
                      </div>
                    </div>
                    <span className="ml-2 shrink-0 text-xs text-muted-foreground">
                      {formatRelativeTime(debt.updated_at)}
                    </span>
                  </div>

                  <p
                    className="mb-3 truncate text-sm text-muted-foreground"
                    title={debt.description || 'No description'}
                  >
                    {debt.description || 'No description'}
                  </p>

                  <div className="flex items-center justify-between">
                    <span
                      className={`text-lg font-semibold ${
                        debt.debt_type === 'i_owe' ? 'text-destructive' : 'text-success'
                      }`}
                    >
                      {formatCurrency(parseFloat(debt.total_amount || 0))}
                    </span>
                    <div
                      className="flex justify-end space-x-2"
                      onClick={(e) => e.stopPropagation()}
                    >
                      <button
                        onClick={() => viewDebt(debt)}
                        className="btn-secondary px-3 py-1 text-xs"
                      >
                        View
                      </button>
                      <button
                        onClick={() => editDebt(debt)}
                        className="btn-secondary px-3 py-1 text-xs"
                      >
                        Edit
                      </button>
                    </div>
                  </div>
                </div>
              )
            })}
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between">
              <div className="text-sm text-muted-foreground">
                Showing {(currentPage - 1) * itemsPerPage + 1} to{' '}
                {Math.min(currentPage * itemsPerPage, filteredAndSortedDebts.length)} of{' '}
                {filteredAndSortedDebts.length} debts
              </div>
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                  disabled={currentPage === 1}
                  className="btn-secondary disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Previous
                </button>
                <span className="text-sm text-muted-foreground">
                  Page {currentPage} of {totalPages}
                </span>
                <button
                  onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                  disabled={currentPage === totalPages}
                  className="btn-secondary disabled:cursor-not-allowed disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </>
      )}

      {/* Modals */}
      {showDetailsModal && selectedDebt && (
        <DebtDetailsModal
          debt={selectedDebt}
          onClose={() => {
            setShowDetailsModal(false)
            setSelectedDebt(null)
          }}
          onEdit={() => {
            setShowDetailsModal(false)
            editDebt(selectedDebt)
          }}
          onDelete={() => {
            setShowDetailsModal(false)
            confirmDeleteDebt(selectedDebt)
          }}
        />
      )}

      {showEditModal && selectedDebt && (
        <EditDebtModal
          debt={selectedDebt}
          onClose={() => {
            setShowEditModal(false)
            setSelectedDebt(null)
          }}
          onDebtUpdated={handleDebtUpdated}
        />
      )}

      {showDeleteDialog && debtToDelete && (
        <DeleteDebtModal
          debt={debtToDelete}
          onConfirm={handleDebtDeleted}
          onClose={() => {
            setShowDeleteDialog(false)
            setDebtToDelete(null)
          }}
        />
      )}
    </div>
  )
}
