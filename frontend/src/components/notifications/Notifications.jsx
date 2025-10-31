import { useNotificationsStore } from '@stores/notificationsStore'
import { cn } from '@utils/cn'

const NotificationIcon = ({ type }) => {
  const icons = {
    success: (
      <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 13l4 4L19 7" />
      </svg>
    ),
    error: (
      <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M6 18L18 6M6 6l12 12"
        />
      </svg>
    ),
    warning: (
      <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
        />
      </svg>
    ),
    info: (
      <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth="2"
          d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
    ),
  }

  return icons[type] || null
}

const getTypeClasses = (type) => {
  const classes = {
    success: 'bg-success border-success dark:bg-success dark:text-success-foreground text-white',
    error:
      'bg-destructive border-destructive dark:bg-destructive dark:text-destructive-foreground text-white',
    warning: 'bg-warning border-warning dark:bg-warning dark:text-warning-foreground text-white',
    info: 'bg-primary border-primary dark:bg-primary dark:text-primary-foreground text-white',
  }

  return classes[type] || 'bg-muted border-border text-foreground'
}

export const Notifications = () => {
  const notifications = useNotificationsStore((state) => state.notifications)
  const removeNotification = useNotificationsStore((state) => state.removeNotification)

  return (
    <div className="fixed right-4 top-4 z-50 max-w-sm space-y-2">
      {notifications.map((notification) => (
        <div key={notification.id} className="animate-slide-up" role="alert">
          <div
            className={cn(
              'flex items-start rounded-lg border p-4 shadow-lg',
              getTypeClasses(notification.type)
            )}
          >
            <div className="mr-3 flex-shrink-0">
              <NotificationIcon type={notification.type} />
            </div>

            <div className="min-w-0 flex-1">
              <p className="text-sm">{notification.message}</p>
            </div>

            <button
              onClick={() => removeNotification(notification.id)}
              className="ml-3 flex-shrink-0 opacity-60 transition-opacity hover:opacity-100"
              aria-label="Close notification"
            >
              <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>
      ))}
    </div>
  )
}
