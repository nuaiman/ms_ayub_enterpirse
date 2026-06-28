<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Expenses</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ expensesStore.expenses.length }}
          total</span>
        <span class="text-sm font-semibold text-warning-text">৳{{ filteredTotal.toFixed(2) }}</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search expenses..."
            class="input w-full sm:w-80 pl-9 pr-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent focus:border-transparent placeholder:text-muted" />
          <button
            v-if="searchQuery"
            @click="searchQuery = ''"
            type="button"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors"
          >
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Advanced Filters Toggle -->
        <button
          @click="showAdvancedFilters = !showAdvancedFilters"
          class="h-9 px-3 flex items-center gap-1.5 border border-default rounded-lg text-secondary hover-surface transition-colors text-sm whitespace-nowrap"
          :class="{ 'border-accent text-accent': showAdvancedFilters }"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
          </svg>
          <span class="hidden sm:inline">Filters</span>
          <span v-if="activeFilterCount > 0" class="w-5 h-5 rounded-full bg-accent text-accent-foreground text-[10px] flex items-center justify-center font-bold">
            {{ activeFilterCount }}
          </span>
        </button>

        <div class="flex items-center gap-2">
          <!-- Month/Year Filter -->
          <select v-model="monthFilter"
            class="input px-3 py-2 rounded-lg text-sm cursor-pointer focus:outline-none focus:ring-1 focus:ring-accent focus:border-transparent max-w-[120px] sm:max-w-none">
            <option value="" class="bg-surface text-primary">All Months</option>
            <option v-for="m in availableMonths" :key="m.value" :value="m.value" class="bg-surface text-primary">{{
              m.label }}</option>
          </select>

          <div class="relative">
            <button @click="toggleSortDropdown"
              class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors"
              title="Sort">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
              </svg>
            </button>
            <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95"
              enter-to-class="opacity-100 scale-100" leave-active-class="transition ease-in duration-150"
              leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
              <div v-if="showSortMenu"
                class="absolute right-0 top-10 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
                <div class="px-3 py-2 border-b border-divider">
                  <p class="text-xs font-semibold text-muted uppercase">Sort By</p>
                </div>
                <button v-for="option in sortOptions" :key="option.key" @click="applySort(option.key)"
                  class="w-full flex items-center justify-between px-3 py-2 text-sm text-secondary hover:bg-surface-alt transition-colors">
                  {{ option.label }}
                  <svg v-if="currentSort === option.key" class="w-4 h-4 text-accent" fill="none" stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
              </div>
            </Transition>
          </div>
          <div v-if="showSortMenu" class="fixed inset-0 z-40" @click="showSortMenu = false"></div>

          <button @click="copyToClipboard"
            class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors relative"
            title="Copy to clipboard">
            <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>

          <button @click="openCreateDialog"
            class="h-9 px-4 flex items-center gap-2 bg-accent text-accent-foreground rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" />
            </svg>
            <span class="hidden sm:inline">Add</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Advanced Filters -->
    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="showAdvancedFilters" class="bg-surface-alt rounded-lg p-4 mb-4 border border-default">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Category</label>
            <select v-model="filters.category" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Categories</option>
              <option value="rent">Rent</option>
              <option value="salary">Salary</option>
              <option value="bill">Bill</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Min Amount</label>
            <input v-model.number="filters.min_amount" type="number" step="0.01" min="0" placeholder="0.00"
              class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent" />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Max Amount</label>
            <input v-model.number="filters.max_amount" type="number" step="0.01" min="0" placeholder="1000.00"
              class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent" />
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Has Image</label>
            <select v-model="filters.has_image" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Expenses</option>
              <option value="true">Has Receipt</option>
              <option value="false">No Receipt</option>
            </select>
          </div>
          <div class="space-y-1 md:col-span-4">
            <label class="text-xs font-medium text-muted">Date Range</label>
            <div class="flex gap-3">
              <input v-model="filters.date_from" type="date" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent" />
              <span class="text-xs text-muted self-center">to</span>
              <input v-model="filters.date_to" type="date" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent" />
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-2 mt-4 pt-3 border-t border-divider">
          <button @click="clearFilters" class="button px-3 py-1.5 text-xs rounded-lg hover-surface transition-colors">
            Clear Filters
          </button>
          <button @click="showAdvancedFilters = false" class="button px-3 py-1.5 text-xs rounded-lg bg-accent text-accent-foreground hover:opacity-90 transition-colors">
            Apply Filters
          </button>
        </div>
      </div>
    </Transition>

    <!-- Table Header -->
    <div class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-[768px]">
        <div
          class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
          <div class="col-span-3 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
            @click="toggleSort('title')">
            Expense
            <svg v-if="sortField === 'title'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
              fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
            @click="toggleSort('amount')">
            Amount
            <svg v-if="sortField === 'amount'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
              fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
            @click="toggleSort('category')">
            Category
            <svg v-if="sortField === 'category'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
              fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
            @click="toggleSort('expense_date')">
            Date
            <svg v-if="sortField === 'expense_date'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
              fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2">Notes</div>
          <div class="col-span-1 text-right pr-0.5">Actions</div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="text-center space-y-4">
            <svg class="animate-spin w-10 h-10 text-accent mx-auto" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            <p class="text-sm text-muted">Loading expenses...</p>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredExpenses.length === 0" class="flex items-center justify-center py-12">
          <div class="text-center space-y-3">
            <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
              <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-primary">No expenses found</p>
            <p class="text-xs text-muted mt-1">{{ searchQuery || hasActiveFilters ? 'adjust your filters' : 'add a new expense' }}</p>
            <button v-if="hasActiveFilters" @click="clearAllFilters" class="text-xs text-accent hover:underline">
              Clear all filters
            </button>
          </div>
        </div>

        <!-- Expenses List -->
        <div v-else>
          <ExpenseRow v-for="e in filteredExpenses" :key="e.id" :expense="e" @expense-updated="handleExpenseUpdated"
            @view-expense="openDetailDialog" />
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div v-if="!loading && filteredExpenses.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredExpenses.length }} of {{ expensesStore.expenses.length }} ·
        Total: ৳{{ filteredTotal.toFixed(2) }}</p>
      <div class="flex items-center gap-2">
        <span v-if="hasActiveFilters" class="text-xs text-accent">(filtered)</span>
        <button @click="refreshExpenses"
          class="button px-3 py-1.5 text-xs rounded-lg hover-surface transition-colors inline-flex items-center gap-1.5"
          :disabled="refreshing">
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Dialogs -->
    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateExpenseForm @expense-created="handleExpenseCreated" />
    </BaseDialog>

    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedExpense" class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-bold text-primary">{{ selectedExpense.title }}</h2>
            <div class="flex items-center gap-2 mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium capitalize"
                :class="getCategoryClass(selectedExpense.category)">{{ selectedExpense.category }}</span>
              <span class="text-xs text-muted">ID: #{{ selectedExpense.id }}</span>
            </div>
          </div>
          <span class="text-xl font-bold text-warning-text">৳{{ selectedExpense.amount.toFixed(2) }}</span>
        </div>
        <div v-if="selectedExpense.image_url" class="rounded-lg overflow-hidden border border-default">
          <img :src="getImageUrl(selectedExpense.image_url) || undefined" alt="Receipt"
            class="w-full max-h-64 object-contain bg-surface-alt" />
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Expense Date</p>
            <p class="text-sm text-primary">{{ formatDate(selectedExpense.expense_date) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Category</p>
            <p class="text-sm text-primary capitalize">{{ selectedExpense.category }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Created</p>
            <p class="text-sm text-primary">{{ formatDateTime(selectedExpense.created_at) }}</p>
          </div>
          <div class="space-y-1 md:col-span-3">
            <p class="text-xs text-muted uppercase">Notes</p>
            <p class="text-sm text-primary whitespace-pre-wrap">{{ selectedExpense.notes || '—' }}</p>
          </div>
        </div>
        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showDetailDialog = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useExpensesStore } from '@/stores/expenses'
import type { Expense } from '@/types/expense'
import ExpenseRow from './ExpenseRow.vue'
import CreateExpenseForm from './CreateExpenseForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const expensesStore = useExpensesStore()
const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const monthFilter = ref('')
const showAdvancedFilters = ref(false)
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedExpense = ref<Expense | null>(null)
const showSortMenu = ref(false)
const sortField = ref<'title' | 'amount' | 'category' | 'expense_date'>('expense_date')
const sortDirection = ref<'asc' | 'desc'>('desc')
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'amount_desc', label: 'Amount High-Low' },
  { key: 'amount_asc', label: 'Amount Low-High' },
  { key: 'title_asc', label: 'Title A-Z' },
  { key: 'title_desc', label: 'Title Z-A' },
  { key: 'category_asc', label: 'Category A-Z' },
  { key: 'category_desc', label: 'Category Z-A' },
]

// Filter state
const filters = ref({
  category: '',
  min_amount: null as number | null,
  max_amount: null as number | null,
  has_image: '',
  date_from: '',
  date_to: ''
})

const hasActiveFilters = computed(() => {
  return !!searchQuery.value ||
         !!monthFilter.value ||
         !!filters.value.category ||
         !!filters.value.min_amount ||
         !!filters.value.max_amount ||
         !!filters.value.has_image ||
         !!filters.value.date_from ||
         !!filters.value.date_to
})

const activeFilterCount = computed(() => {
  let count = 0
  if (monthFilter.value) count++
  if (filters.value.category) count++
  if (filters.value.min_amount) count++
  if (filters.value.max_amount) count++
  if (filters.value.has_image) count++
  if (filters.value.date_from) count++
  if (filters.value.date_to) count++
  return count
})

// Generate available months from actual expense data
const availableMonths = computed(() => {
  const months: Record<string, string> = {}
  expensesStore.expenses.forEach(e => {
    const d = new Date(e.expense_date)
    const key = `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`
    if (!months[key]) {
      months[key] = d.toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    }
  })
  return Object.entries(months)
    .map(([value, label]) => ({ value, label }))
    .sort((a, b) => b.value.localeCompare(a.value))
})

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    // Search will be applied via computed property
  }, 300)
})

onMounted(async () => { await fetchExpenses() })

const fetchExpenses = async () => { loading.value = true; try { await expensesStore.fetchExpenses() } finally { loading.value = false } }
const refreshExpenses = async () => { refreshing.value = true; try { await expensesStore.fetchExpenses() } finally { refreshing.value = false } }
const handleExpenseUpdated = async () => { await expensesStore.fetchExpenses() }
const handleExpenseCreated = async () => { showCreateDialog.value = false; await expensesStore.fetchExpenses() }
const openCreateDialog = () => { showCreateDialog.value = true }
const openDetailDialog = (e: Expense) => { selectedExpense.value = e; showDetailDialog.value = true }

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }

const applySort = (key: string) => {
  currentSort.value = key
  showSortMenu.value = false
  switch (key) {
    case 'newest': sortField.value = 'expense_date'; sortDirection.value = 'desc'; break
    case 'oldest': sortField.value = 'expense_date'; sortDirection.value = 'asc'; break
    case 'amount_desc': sortField.value = 'amount'; sortDirection.value = 'desc'; break
    case 'amount_asc': sortField.value = 'amount'; sortDirection.value = 'asc'; break
    case 'title_asc': sortField.value = 'title'; sortDirection.value = 'asc'; break
    case 'title_desc': sortField.value = 'title'; sortDirection.value = 'desc'; break
    case 'category_asc': sortField.value = 'category'; sortDirection.value = 'asc'; break
    case 'category_desc': sortField.value = 'category'; sortDirection.value = 'desc'; break
  }
}

const toggleSort = (field: 'title' | 'amount' | 'category' | 'expense_date') => {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'asc'
  }
  // Update currentSort to match
  if (field === 'title') {
    currentSort.value = sortDirection.value === 'asc' ? 'title_asc' : 'title_desc'
  } else if (field === 'amount') {
    currentSort.value = sortDirection.value === 'asc' ? 'amount_asc' : 'amount_desc'
  } else if (field === 'category') {
    currentSort.value = sortDirection.value === 'asc' ? 'category_asc' : 'category_desc'
  } else if (field === 'expense_date') {
    currentSort.value = sortDirection.value === 'asc' ? 'oldest' : 'newest'
  }
}

const copyToClipboard = async () => {
  const expenses = filteredExpenses.value
  if (expenses.length === 0) return
  const headers = ['ID', 'Title', 'Category', 'Amount', 'Date', 'Notes', 'Created']
  const rows = expenses.map(e => [
    e.id,
    e.title,
    e.category,
    e.amount,
    e.expense_date,
    e.notes || '',
    formatDateTime(e.created_at)
  ])
  const tsv = [headers, ...rows].map(r => r.map(c => String(c).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t')).join('\n')
  try {
    await navigator.clipboard.writeText(tsv)
    copied.value = true
    push.success('Copied!')
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    const ta = document.createElement('textarea')
    ta.value = tsv
    ta.style.position = 'fixed'
    ta.style.opacity = '0'
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    copied.value = true
    push.success('Copied!')
    setTimeout(() => { copied.value = false }, 2000)
  }
}

const clearFilters = () => {
  filters.value = {
    category: '',
    min_amount: null,
    max_amount: null,
    has_image: '',
    date_from: '',
    date_to: ''
  }
  monthFilter.value = ''
}

const clearAllFilters = () => {
  searchQuery.value = ''
  clearFilters()
}

const filteredExpenses = computed(() => {
  let result = [...expensesStore.expenses]

  // Search across ALL fields
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase().trim()
    result = result.filter(e => {
      const searchableFields = [
        e.title,
        e.category,
        e.notes,
        String(e.id),
        String(e.amount),
        e.expense_date,
      ]

      const allText = searchableFields
        .filter(f => f != null)
        .join(' ')
        .toLowerCase()

      return allText.includes(q)
    })
  }

  // Month/Year filter
  if (monthFilter.value) {
    const [year, month] = monthFilter.value.split('-')
    result = result.filter(e => {
      const d = new Date(e.expense_date)
      return d.getFullYear() === parseInt(year!) && d.getMonth() + 1 === parseInt(month!)
    })
  }

  // Filter by Category
  if (filters.value.category) {
    result = result.filter(e => e.category === filters.value.category)
  }

  // Filter by Min Amount
  if (filters.value.min_amount !== null) {
    result = result.filter(e => e.amount >= (filters.value.min_amount || 0))
  }

  // Filter by Max Amount
  if (filters.value.max_amount !== null) {
    result = result.filter(e => e.amount <= (filters.value.max_amount || Infinity))
  }

  // Filter by Has Image
  if (filters.value.has_image) {
    const hasImage = filters.value.has_image === 'true'
    result = result.filter(e => {
      const has = !!e.image_url
      return has === hasImage
    })
  }

  // Filter by Date Range
  if (filters.value.date_from) {
    const from = new Date(filters.value.date_from)
    result = result.filter(e => new Date(e.expense_date) >= from)
  }
  if (filters.value.date_to) {
    const to = new Date(filters.value.date_to)
    to.setHours(23, 59, 59, 999)
    result = result.filter(e => new Date(e.expense_date) <= to)
  }

  // Sort
  result.sort((a, b) => {
    let comparison = 0
    if (sortField.value === 'title') {
      comparison = a.title.localeCompare(b.title)
    } else if (sortField.value === 'amount') {
      comparison = a.amount - b.amount
    } else if (sortField.value === 'category') {
      comparison = a.category.localeCompare(b.category)
    } else if (sortField.value === 'expense_date') {
      comparison = new Date(a.expense_date).getTime() - new Date(b.expense_date).getTime()
    }
    return sortDirection.value === 'desc' ? -comparison : comparison
  })
  return result
})

const filteredTotal = computed(() => filteredExpenses.value.reduce((sum, e) => sum + e.amount, 0))

const getCategoryClass = (cat: string): string => {
  switch (cat) {
    case 'rent': return 'bg-blue-50 text-blue-700 border border-blue-200'
    case 'salary': return 'bg-green-50 text-green-700 border border-green-200'
    case 'bill': return 'bg-yellow-50 text-yellow-700 border border-yellow-200'
    default: return 'bg-gray-50 text-gray-700 border border-gray-200'
  }
}

const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
const formatDateTime = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' })
</script>
