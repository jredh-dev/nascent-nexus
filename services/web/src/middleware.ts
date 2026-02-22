import { defineMiddleware } from 'astro:middleware';
import { proxyToPortal } from './lib/proxy';
import { isLocale, defaultLocale } from './i18n/config';

// Proxy Connect RPC, API calls, and portal-owned GET routes to the Go backend.
// Also handles bare (non-locale-prefixed) page paths by redirecting to /{locale}/...
export const onRequest = defineMiddleware(async (context, next) => {
  const { pathname } = context.url;

  // Proxy GET routes that the portal backend owns (not Astro pages)
  const isProxiedGet = context.request.method === 'GET' &&
    (pathname === '/logout' || pathname.startsWith('/auth/'));

  // Proxy Connect RPC, legacy API calls, and portal-owned GETs
  if (pathname.startsWith('/portal.v1.') || pathname.startsWith('/api/') || isProxiedGet) {
    return proxyToPortal(context.request, pathname + context.url.search);
  }

  // Redirect bare paths (e.g. /login → /en/login) for backwards compatibility.
  // Skip paths that are locale-prefixed, static assets, or the root /.
  const segments = pathname.split('/').filter(Boolean);
  if (segments.length > 0 && !isLocale(segments[0])) {
    // Check if this looks like a page route (not a file with extension)
    const lastSegment = segments[segments.length - 1];
    if (!lastSegment.includes('.')) {
      // Detect preferred locale from Accept-Language header
      const acceptLang = context.request.headers.get('Accept-Language') || '';
      let locale = defaultLocale;
      if (acceptLang.includes('es')) {
        locale = 'es';
      }
      const target = `/${locale}${pathname}${context.url.search}`;
      return context.redirect(target, 302);
    }
  }

  // Redirect bare / to /{locale}/
  if (pathname === '/') {
    const acceptLang = context.request.headers.get('Accept-Language') || '';
    let locale = defaultLocale;
    if (acceptLang.includes('es')) {
      locale = 'es';
    }
    return context.redirect(`/${locale}/`, 302);
  }

  return next();
});
