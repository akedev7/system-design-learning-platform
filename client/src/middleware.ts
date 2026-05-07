import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { auth0 } from './lib/auth0'

export async function middleware(request: NextRequest) {
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