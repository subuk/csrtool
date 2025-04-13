import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div className="text-center">
      <h1 className="text-4xl font-bold tracking-tight text-gray-900 dark:text-gray-100">
        404
      </h1>
      <p className="mt-4 text-lg text-gray-600 dark:text-gray-400">
        Page not found
      </p>
      <div className="mt-6">
        <Link
          to="/"
          className="btn btn-primary"
        >
          Go back home
        </Link>
      </div>
    </div>
  )
}
