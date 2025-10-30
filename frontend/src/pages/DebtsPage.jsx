import { useState } from 'react'
import { DebtsTable } from '@components/debts/DebtsTable'
import { CreateDebtModal } from '@components/debts/CreateDebtModal'

export const DebtsPage = () => {
  const [showCreateModal, setShowCreateModal] = useState(false)

  const handleDebtCreated = () => {
    setShowCreateModal(false)
  }

  return (
    <div className="mx-auto max-w-7xl space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-foreground">Debt Management</h1>
          <p className="mt-1 text-muted-foreground">
            Track and manage all your debts and payment obligations
          </p>
        </div>
        <button onClick={() => setShowCreateModal(true)} className="btn-primary">
          <svg className="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M12 6v6m0 0v6m0-6h6m-6 0H6"
            />
          </svg>
          Create Debt
        </button>
      </div>

      {/* Debts Table */}
      <DebtsTable />

      {/* Create Debt Modal */}
      {showCreateModal && (
        <CreateDebtModal
          onDebtCreated={handleDebtCreated}
          onClose={() => setShowCreateModal(false)}
        />
      )}
    </div>
  )
}

