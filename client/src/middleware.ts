import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// For local testing without Auth0, set DISABLE_AUTH=true in docker-compose environment
const disableAuth = process.env.DISABLE_AUTH === 'true'

export async function middleware(request: NextRequest) {
  // Bypass auth for local development
  if (disableAuth) {
    return NextResponse.next()
  }

  // Dynamic import to avoid loading Auth0 when disabled
  const { auth0 } = await import('./lib/auth0')

  if (request.nextUrl.pathname.startsWith('/admin')) {
    const session = await auth0.getSession(request)

    if (!session?.user) {
      return NextResponse.redirect(new URL('/', request.url))
    }

    const userRole = session.user['https://ralph-learning.com/role']

    if (!userRole || userRole !== 'Admin') {
      return NextResponse.redirect(new URL('/', request.url))
    }
  }
  return auth0.middleware(request)
}

export const config = {
  matcher: [
    '/((?!_next/static|_next/image|favicon.ico|sitemap.xml|robots.txt).*)'
  ]
}