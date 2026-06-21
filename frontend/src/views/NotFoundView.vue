<template>
  <div class="min-h-screen bg-app flex items-center justify-center p-6">
    <div class="text-center max-w-md">
      <!-- 404 Illustration -->
      <div class="mb-8">
        <div
          class="w-24 h-24 mx-auto rounded-full bg-surface-alt border border-default flex items-center justify-center mb-6">
          <svg class="w-12 h-12 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h1 class="text-6xl font-bold text-primary mb-2">404</h1>
        <p class="text-lg font-medium text-primary mb-2">Page not found</p>
        <p class="text-sm text-muted mb-8">
          The page you're looking for doesn't exist or has been moved.
        </p>
      </div>

      <!-- Actions -->
      <div class="flex flex-col sm:flex-row items-center justify-center gap-3">
        <button @click="$router.push('/home')" class="w-full sm:w-auto px-6 py-2.5 bg-accent text-accent-foreground rounded-lg text-sm font-semibold
                 shadow-sm hover:shadow-md transition-all duration-200 active:scale-95
                 inline-flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-4 0a1 1 0 01-1-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 01-1 1" />
          </svg>
          Go to Home
        </button>
        <button @click="$router.go(-1)" class="w-full sm:w-auto px-6 py-2.5 border border-default rounded-lg text-sm font-medium text-secondary
                 hover-surface transition-all duration-200 active:scale-95
                 inline-flex items-center justify-center gap-2">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Go Back
        </button>
      </div>
    </div>
  </div>
</template>
