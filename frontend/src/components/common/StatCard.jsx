import { cn } from '@utils/cn'

export const StatCard = ({
  title,
  value,
  icon,
  iconBgColor = 'bg-primary/10',
  iconColor = 'text-primary',
  valueColor = 'text-foreground',
  trend,
  trendLabel,
}) => {
  return (
    <div className="card p-6">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <h3 className="text-sm font-medium text-muted-foreground">{title}</h3>
          <p className={cn('mt-2 text-3xl font-bold', valueColor)}>{value}</p>
          {trend && trendLabel && (
            <p className="mt-2 flex items-center text-sm">
              {trend > 0 ? (
                <svg
                  className="mr-1 h-4 w-4 text-success"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                  />
                </svg>
              ) : (
                <svg
                  className="mr-1 h-4 w-4 text-destructive"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 17h8m0 0v-8m0 8l-8-8-4 4-6-6"
                  />
                </svg>
              )}
              <span className="text-muted-foreground">{trendLabel}</span>
            </p>
          )}
        </div>
        {icon && (
          <div className={cn('flex h-12 w-12 items-center justify-center rounded-lg', iconBgColor)}>
            <div className={cn('h-6 w-6', iconColor)}>{icon}</div>
          </div>
        )}
      </div>
    </div>
  )
}

