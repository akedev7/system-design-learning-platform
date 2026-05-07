'use client'

import { useUser } from '@auth0/nextjs-auth0/client'

export default function AuthButton() {
  const { user, isLoading } = useUser()

  if (isLoading) {
    return <span className="text-sm text-gray-500">Loading...</span>
  }

  if (user) {
    return (
      <div className="flex items-center gap-4">
        <span className="text-sm text-gray-700">Welcome, {user.name}</span>
        <a
          href="/auth/logout"
          className="px-4 py-2 text-sm bg-gray-200 rounded hover:bg-gray-300 transition-colors"
        >
          Logout
        </a>
      </div>
    )
  }

  return (
    <a
      href="/auth/login"
      className="px-4 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
    >
      Log in
    </a>
  )
}