import { createRouter, createWebHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

import AuthView from "@/views/AuthView.vue";
import HomeView from "@/views/HomeView.vue";
import NotFoundView from "@/views/NotFoundView.vue";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // Auth route
    {
      path: "/auth",
      component: AuthView,
      meta: { requiresGuest: true },
    },

    // Home route
    {
      path: "/home",
      component: HomeView,
      meta: {
        requiresAuth: true,
        roles: ["admin", "manager", "accounts", "staff"],
      },
    },

    // Support route - accessible to everyone
    {
      path: "/support",
      name: "support",
      component: () => import("@/views/SupportView.vue"),
      meta: { requiresAuth: false }, // No authentication required
    },

    // Catch-all 404
    {
      path: "/:pathMatch(.*)*",
      component: NotFoundView,
    },
  ],
});

router.beforeEach((to) => {
  const auth = useAuthStore();

  const isAuthRoute = to.path === "/auth";
  const isSupportRoute = to.path === "/support";

  // Allow access to support page without authentication
  if (isSupportRoute) {
    return true;
  }

  // not logged in → redirect to auth
  if (!auth.isAuthenticated && !isAuthRoute) {
    return "/auth";
  }

  // role guard (only if route defines roles)
  const roles = to.meta.roles as string[] | undefined;

  if (roles && auth.user && !roles.includes(auth.user.role)) {
    return "/auth";
  }

  // prevent going back to /auth if already logged in
  if (auth.isAuthenticated && isAuthRoute) {
    return "/home";
  }
});

export default router;
