import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

import AuthView from '@/views/AuthView.vue'
import HomeView from '@/views/HomeView.vue'
import NotFoundView from '@/views/NotFoundView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      component: AuthView,
    },

    {
      path: '/home',
      component: HomeView,
      meta: {
        requiresAuth: true,
        roles: ['admin', 'manager', 'accounts', 'staff'],
      },
    },

    {
      path: '/:pathMatch(.*)*',
      component: NotFoundView,
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  const isAuthRoute = to.path === '/auth'

  // not logged in → redirect to auth
  if (!auth.isAuthenticated && !isAuthRoute) {
    return '/auth'
  }

  // role guard (only if route defines roles)
  const roles = to.meta.roles as string[] | undefined

  if (roles && auth.user && !roles.includes(auth.user.role)) {
    return '/auth'
  }

  // prevent going back to /auth if already logged in
  if (auth.isAuthenticated && isAuthRoute) {
    return '/home'
  }
})

export default router
