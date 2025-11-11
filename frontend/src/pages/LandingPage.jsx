import { Link } from 'react-router-dom'
import { ROUTES } from '@/routes/routes'

export const LandingPage = () => {
  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto px-4 py-16">
        <div className="text-center">
          <h1 className="mb-4 text-5xl font-bold text-primary">Pay Your Dues</h1>
          <p className="mb-8 text-xl text-muted-foreground">
            Manage your debts and payments efficiently
          </p>
          <div className="flex justify-center gap-4">
            <Link to={ROUTES.LOGIN} className="btn-primary">
              Sign In
            </Link>
            <Link to={ROUTES.REGISTER} className="btn-outline">
              Get Started
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

