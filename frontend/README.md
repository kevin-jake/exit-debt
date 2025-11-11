# Pay Your Dues - Frontend

A modern debt tracking application built with React, Vite, and TailwindCSS.

## Tech Stack

- **Framework**: React 18.x
- **Build Tool**: Vite 7.x
- **Styling**: TailwindCSS 3.4.x
- **State Management**: Zustand
- **Routing**: React Router v6
- **Forms**: React Hook Form
- **Testing**: Vitest + React Testing Library
- **UI Components**: Custom components with TailwindCSS

## Getting Started

### Prerequisites

- Node.js >= 18.x
- npm or yarn

### Installation

```bash
# Install dependencies
npm install

# Copy environment variables
cp .env.example .env.development

# Start development server
npm run dev
```

The app will be available at `http://localhost:5173`

## Available Scripts

### Development

```bash
npm run dev              # Start development server
npm run build            # Build for production
npm run preview          # Preview production build
npm run lint             # Run ESLint
npm run format           # Run Prettier (if configured)
```

### Testing

```bash
npm run test             # Run tests
npm run test:ui          # Run tests with UI
npm run test:coverage    # Generate coverage report
```

## Project Structure

```
frontend/
├── src/
│   ├── api/              # API client and endpoints
│   │   └── client.js     # Main API client
│   ├── components/       # React components
│   │   ├── layout/       # Layout components (Navigation, Layout)
│   │   ├── contacts/     # Contact management components
│   │   ├── debts/        # Debt management components
│   │   ├── common/       # Shared components (LoadingSpinner, EmptyState, etc.)
│   │   └── notifications/ # Toast notifications
│   ├── pages/            # Page components
│   │   ├── LandingPage.jsx
│   │   ├── LoginPage.jsx
│   │   ├── RegisterPage.jsx
│   │   ├── DashboardPage.jsx
│   │   ├── ContactsPage.jsx
│   │   ├── DebtsPage.jsx
│   │   └── SettingsPage.jsx
│   ├── routes/           # Routing configuration
│   │   ├── index.jsx     # Main routes
│   │   ├── ProtectedRoute.jsx
│   │   ├── PublicRoute.jsx
│   │   └── routes.js     # Route constants
│   ├── stores/           # Zustand stores
│   │   ├── authStore.js
│   │   ├── contactsStore.js
│   │   ├── debtsStore.js
│   │   ├── paymentsStore.js
│   │   ├── notificationsStore.js
│   │   ├── settingsStore.js
│   │   └── themeStore.js
│   ├── hooks/            # Custom React hooks
│   │   ├── useAuth.js
│   │   ├── useDebounce.js
│   │   ├── useLocalStorage.js
│   │   ├── useMediaQuery.js
│   │   ├── useClickOutside.js
│   │   └── useKeyPress.js
│   ├── utils/            # Utility functions
│   │   ├── cn.js         # Class name utilities
│   │   └── formatters.js # Formatting utilities
│   ├── App.jsx           # Main App component
│   ├── main.jsx          # Entry point
│   └── index.css         # Global styles
├── public/               # Static assets
├── .env.example          # Environment variables template
├── vite.config.js        # Vite configuration
├── tailwind.config.js    # TailwindCSS configuration
├── vitest.config.js      # Vitest configuration
└── package.json          # Dependencies and scripts
```

## Features

### Authentication
- User registration and login
- JWT-based authentication
- Protected routes
- Automatic token management

### Contact Management
- Create, read, update, and delete contacts
- Search and filter contacts
- Contact details modal
- Pagination support

### Debt Tracking
- Track debts (I Owe / Owed to Me)
- Add payments to debts
- Track payment history
- Due date tracking with reminders
- Debt details with calculations

### Dashboard
- Overview of total debts
- Recent debts
- Upcoming due dates
- Quick actions

### Settings
- Theme selection (Light/Dark/System)
- Currency preferences
- Language settings
- Notification preferences

### UI/UX
- Responsive design (mobile, tablet, desktop)
- Dark mode support
- Toast notifications
- Loading states
- Empty states
- Smooth animations

## Environment Variables

Create a `.env.development` file based on `.env.example`:

```env
VITE_API_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Pay Your Dues
VITE_APP_VERSION=2.0.0
```

## API Integration

The frontend communicates with the backend API using a centralized API client (`src/api/client.js`).

### API Base URL

Configure the API URL in your environment file:
- Development: `http://localhost:8080/api/v1`
- Production: Set via `VITE_API_URL` environment variable

### Authentication

All authenticated requests include the JWT token in the `Authorization` header:
```
Authorization: Bearer <token>
```

## State Management

The app uses Zustand for state management. Stores are organized by domain:

- `authStore` - User authentication state
- `contactsStore` - Contact management state
- `debtsStore` - Debt management state
- `paymentsStore` - Payment management state
- `notificationsStore` - Toast notifications
- `settingsStore` - User preferences
- `themeStore` - Theme (dark/light mode)

Example usage:
```javascript
import { useAuthStore } from '@stores/authStore'

function MyComponent() {
  const user = useAuthStore((state) => state.user)
  const login = useAuthStore((state) => state.login)
  
  // Use user and login...
}
```

## Custom Hooks

The app includes several custom hooks for common functionality:

- `useAuth()` - Authentication helpers
- `useDebounce(value, delay)` - Debounce values
- `useLocalStorage(key, initialValue)` - LocalStorage with state
- `useMediaQuery(query)` - Responsive design helpers
- `useClickOutside(ref, handler)` - Detect clicks outside element
- `useKeyPress(key)` - Keyboard event handling

## Styling

### TailwindCSS

The app uses TailwindCSS for styling with custom theme variables defined in `src/index.css`.

### Custom CSS Classes

Common utility classes:
- `.btn`, `.btn-primary`, `.btn-secondary`, `.btn-destructive` - Button styles
- `.card` - Card container
- `.input` - Input field
- `.label` - Form label

### Theme

Colors are defined using CSS variables that support both light and dark modes:
- `--background`, `--foreground` - Main background and text
- `--primary`, `--secondary` - Primary and secondary colors
- `--destructive`, `--success`, `--warning` - Status colors
- `--border`, `--input`, `--muted` - UI element colors

## Code Splitting

The app uses React's `lazy()` and `Suspense` for code splitting, automatically splitting routes into separate chunks for optimal loading performance.

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Performance

- Initial load time: < 2 seconds
- Time to interactive: < 3 seconds
- Code splitting for optimal bundle sizes
- Lazy loading of routes
- Optimized images and assets

## Contributing

1. Create a feature branch from `main`
2. Make your changes
3. Run tests and linting
4. Submit a pull request

## License

[Your License Here]

## Support

For issues and questions, please create an issue in the repository.
