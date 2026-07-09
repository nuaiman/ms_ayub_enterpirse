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
          <button v-if="searchQuery" @click="searchQuery = ''" type="button"
            class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="flex items-center gap-2">
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
            Add
          </button>
        </div>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-[768px]">
        <div
          class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
          <div class="col-span-3">Expense</div>
          <div class="col-span-1">Type</div>
          <div class="col-span-2">Amount</div>
          <div class="col-span-2">Date</div>
          <div class="col-span-2">Category</div>
          <div class="col-span-1">Notes</div>
          <div class="col-span-1 text-right">Actions</div>
        </div>

        <div v-if="loading" class="flex items-center justify-center py-12">
          <svg class="animate-spin w-10 h-10 text-accent" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        </div>

        <div v-else-if="filteredExpenses.length === 0" class="flex items-center justify-center py-12">
          <p class="text-sm text-muted">No expenses found</p>
        </div>

        <div v-else>
          <ExpenseRow v-for="e in filteredExpenses" :key="e.id" :expense="e" @expense-updated="handleExpenseUpdated" />
        </div>
      </div>
    </div>

    <div v-if="!loading && filteredExpenses.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredExpenses.length }} of {{ expensesStore.expenses.length }} ·
        Total: ৳{{ filteredTotal.toFixed(2) }}</p>
    </div>

    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateExpenseForm @expense-created="handleExpenseCreated" />
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useExpensesStore } from '@/stores/expenses'
import ExpenseRow from './ExpenseRow.vue'
import CreateExpenseForm from './CreateExpenseForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const expensesStore = useExpensesStore()
const loading = ref(true)
const searchQuery = ref('')
const showCreateDialog = ref(false)
const showSortMenu = ref(false)
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'amount_desc', label: 'Amount High-Low' },
  { key: 'amount_asc', label: 'Amount Low-High' },
]

onMounted(async () => { await fetchExpenses() })
const fetchExpenses = async () => { loading.value = true; try { await expensesStore.fetchExpenses() } finally { loading.value = false } }
const handleExpenseUpdated = async () => { await expensesStore.fetchExpenses() }
const handleExpenseCreated = async () => { showCreateDialog.value = false; await expensesStore.fetchExpenses() }
const openCreateDialog = () => { showCreateDialog.value = true }

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }
const applySort = (key: string) => { currentSort.value = key; showSortMenu.value = false }

const copyToClipboard = async () => {
  const expenses = filteredExpenses.value
  if (expenses.length === 0) return
  const headers = ['ID', 'Title', 'Type', 'Amount', 'Date', 'Category', 'Notes']
  const rows = expenses.map(e => [e.id, e.title, e.is_salary ? 'Salary' : 'Regular', e.amount, e.expense_date, e.category || '', e.notes || ''])
  const tsv = [headers, ...rows].map(r => r.map(c => String(c).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t')).join('\n')
  try {
    await navigator.clipboard.writeText(tsv); copied.value = true; push.success('Copied!')
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    const ta = document.createElement('textarea'); ta.value = tsv; ta.style.position = 'fixed'; ta.style.opacity = '0'
    document.body.appendChild(ta); ta.select(); document.execCommand('copy'); document.body.removeChild(ta)
    copied.value = true; push.success('Copied!'); setTimeout(() => { copied.value = false }, 2000)
  }
}

const filteredExpenses = computed(() => {
  let result = [...expensesStore.expenses]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(e => e.title.toLowerCase().includes(q) || (e.category && e.category.toLowerCase().includes(q)) || (e.notes && e.notes.toLowerCase().includes(q)))
  }
  if (currentSort.value === 'newest') result.sort((a, b) => new Date(b.expense_date).getTime() - new Date(a.expense_date).getTime())
  else if (currentSort.value === 'oldest') result.sort((a, b) => new Date(a.expense_date).getTime() - new Date(b.expense_date).getTime())
  else if (currentSort.value === 'amount_desc') result.sort((a, b) => b.amount - a.amount)
  else if (currentSort.value === 'amount_asc') result.sort((a, b) => a.amount - b.amount)
  return result
})

const filteredTotal = computed(() => filteredExpenses.value.reduce((sum, e) => sum + e.amount, 0))
</script>
