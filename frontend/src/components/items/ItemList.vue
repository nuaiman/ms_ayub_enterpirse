<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Items</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ itemsStore.items.length }} total</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search items..."
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
          <!-- Sort Button -->
          <div class="relative">
            <button @click="toggleSortDropdown"
              class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors" title="Sort">
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

          <!-- Copy Button -->
          <button @click="copyToClipboard"
            class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors relative" title="Copy to clipboard">
            <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>

          <!-- Add Item Button -->
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
            <label class="text-xs font-medium text-muted">Status</label>
            <select v-model="filters.status" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Status</option>
              <option value="active">Active</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Has Contract</label>
            <select v-model="filters.has_contract" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Items</option>
              <option value="true">Has Contract</option>
              <option value="false">No Contract</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Has Customer</label>
            <select v-model="filters.has_customer" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Items</option>
              <option value="true">Has Customer</option>
              <option value="false">No Customer</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Quantity Unit</label>
            <select v-model="filters.quantity_unit" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Units</option>
              <option value="pcs">Pieces</option>
              <option value="g">Grams</option>
              <option value="kg">Kilograms</option>
              <option value="ton">Tons</option>
              <option value="ml">Milliliters</option>
              <option value="liter">Liters</option>
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
        <div class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
          <div class="col-span-3 cursor-pointer hover:text-primary transition-colors flex items-center gap-1" @click="toggleSort('name')">
            Name
            <svg v-if="sortField === 'name'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2">Customer</div>
          <div class="col-span-2">Contract</div>
          <div class="col-span-2">Dimensions</div>
          <div class="col-span-1 cursor-pointer hover:text-primary transition-colors flex items-center gap-1" @click="toggleSort('created_at')">
            Created
            <svg v-if="sortField === 'created_at'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor" viewBox="0 0 24 24">
              <path d="M7 10l5 5 5-5z" />
            </svg>
          </div>
          <div class="col-span-2 text-right pr-0.5">Actions</div>
        </div>

        <!-- Loading State -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="text-center space-y-4">
            <svg class="animate-spin w-10 h-10 text-accent mx-auto" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            <p class="text-sm text-muted">Loading items...</p>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else-if="filteredItems.length === 0" class="flex items-center justify-center py-12">
          <div class="text-center space-y-3">
            <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
              <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <p class="text-sm font-medium text-primary">No items found</p>
            <p class="text-xs text-muted mt-1">{{ searchQuery || hasActiveFilters ? 'adjust your filters' : 'add new item' }}</p>
            <button v-if="hasActiveFilters" @click="clearAllFilters" class="text-xs text-accent hover:underline">
              Clear all filters
            </button>
          </div>
        </div>

        <!-- Items List -->
        <div v-else>
          <ItemRow
            v-for="item in filteredItems"
            :key="item.id"
            :item="item"
            @item-updated="handleItemUpdated"
            @view-item="openDetailDialog"
          />
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div v-if="!loading && filteredItems.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredItems.length }} of {{ itemsStore.items.length }} items</p>
      <div class="flex items-center gap-2">
        <span v-if="hasActiveFilters" class="text-xs text-accent">(filtered)</span>
        <button @click="refreshItems"
          class="button px-3 py-1.5 text-xs rounded-lg hover-surface transition-colors inline-flex items-center gap-1.5"
          :disabled="refreshing">
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Dialogs -->
    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateItemForm @item-created="handleItemCreated" />
    </BaseDialog>

    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedItem" class="space-y-5">
        <!-- Item Header -->
        <div class="flex items-center gap-4">
          <div v-if="selectedItem.image_url" class="shrink-0 w-24 h-24 rounded-lg overflow-hidden border border-default">
            <img :src="getImageUrl(selectedItem.image_url) || undefined" :alt="selectedItem.name"
              class="w-full h-full object-cover" />
          </div>
          <div>
            <h2 class="text-xl font-bold text-primary">{{ selectedItem.name }}</h2>
            <p class="text-sm text-muted">{{ selectedItem.quantity }} {{ selectedItem.quantity_unit }}</p>
            <p v-if="selectedItem.customer_id" class="text-sm text-secondary mt-1">
              Customer: {{ getCustomerName(selectedItem.customer_id) }}
            </p>
          </div>
        </div>

        <!-- Contract Info Card (if contract exists) -->
        <div v-if="hasContract(selectedItem)" class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Contract Details</h3>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Duration</p>
              <p class="text-xs text-primary">{{ selectedItem.duration || '—' }} {{ selectedItem.duration_type }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Status</p>
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold border"
                :class="getStatusClass(selectedItem.status || 'active')">
                {{ selectedItem.status || 'active' }}
              </span>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Price</p>
              <p class="text-xs text-primary">{{ selectedItem.price ? '৳' + selectedItem.price : '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Security Deposit</p>
              <p class="text-xs text-primary">{{ selectedItem.security_deposit ? '৳' + selectedItem.security_deposit : '—' }}</p>
            </div>
            <div class="space-y-1 col-span-2">
              <p class="text-[10px] text-muted uppercase">Estimated Value</p>
              <p class="text-xs text-primary">{{ selectedItem.estimated_value ? '৳' + selectedItem.estimated_value : '—' }}</p>
            </div>
          </div>
        </div>

        <!-- Details Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Weight</p>
            <p class="text-sm text-primary">{{ selectedItem.weight ? selectedItem.weight + ' kg' : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Dimensions</p>
            <p class="text-sm text-primary">{{ getDimensions(selectedItem) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Created</p>
            <p class="text-sm text-primary">{{ formatDate(selectedItem.created_at) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Updated</p>
            <p class="text-sm text-primary">{{ formatDate(selectedItem.updated_at) }}</p>
          </div>
          <div class="space-y-1 md:col-span-2">
            <p class="text-xs text-muted uppercase tracking-wider">Notes</p>
            <p class="text-sm text-primary whitespace-pre-wrap">{{ selectedItem.notes || '—' }}</p>
          </div>
        </div>

        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showDetailDialog = false" class="button px-4 py-2 text-sm rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useItemsStore } from '@/stores/items'
import { useCustomersStore } from '@/stores/customers'
import type { Item } from '@/types/item'
import { hasContract } from '@/types/item'
import ItemRow from './ItemRow.vue'
import CreateItemForm from './CreateItemForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const itemsStore = useItemsStore()
const customersStore = useCustomersStore()

const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const showAdvancedFilters = ref(false)
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedItem = ref<Item | null>(null)
const sortField = ref<'name' | 'created_at'>('created_at')
const sortDirection = ref<'asc' | 'desc'>('desc')
const showSortMenu = ref(false)
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'name_asc', label: 'Name A-Z' },
  { key: 'name_desc', label: 'Name Z-A' },
]

// Filter state
const filters = ref({
  status: '',
  has_contract: '',
  has_customer: '',
  quantity_unit: '',
  date_from: '',
  date_to: ''
})

const hasActiveFilters = computed(() => {
  return !!searchQuery.value ||
         !!filters.value.status ||
         !!filters.value.has_contract ||
         !!filters.value.has_customer ||
         !!filters.value.quantity_unit ||
         !!filters.value.date_from ||
         !!filters.value.date_to
})

const activeFilterCount = computed(() => {
  let count = 0
  if (filters.value.status) count++
  if (filters.value.has_contract) count++
  if (filters.value.has_customer) count++
  if (filters.value.quantity_unit) count++
  if (filters.value.date_from) count++
  if (filters.value.date_to) count++
  return count
})

const getCustomerName = (id: number): string => {
  const customer = customersStore.customers.find(c => c.id === id)
  return customer?.name || `#${id}`
}

const getStatusClass = (status: string) => {
  const classes = {
    active: 'border-green-400 bg-green-50 text-green-700',
    completed: 'border-blue-400 bg-blue-50 text-blue-700',
    cancelled: 'border-red-400 bg-red-50 text-red-700'
  }
  return classes[status as keyof typeof classes] || classes.active
}

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    // Search will be applied via computed property
  }, 300)
})

onMounted(async () => {
  if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
  await fetchItems()
})

const fetchItems = async () => { loading.value = true; try { await itemsStore.fetchItems() } finally { loading.value = false } }
const refreshItems = async () => { refreshing.value = true; try { await itemsStore.fetchItems() } finally { refreshing.value = false } }
const handleItemUpdated = async () => { await itemsStore.fetchItems() }
const handleItemCreated = async () => { showCreateDialog.value = false; await itemsStore.fetchItems() }
const openCreateDialog = () => { showCreateDialog.value = true }
const openDetailDialog = (item: Item) => { selectedItem.value = item; showDetailDialog.value = true }

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }

const applySort = (key: string) => {
  currentSort.value = key
  showSortMenu.value = false
  switch (key) {
    case 'newest':
      sortField.value = 'created_at'
      sortDirection.value = 'desc'
      break
    case 'oldest':
      sortField.value = 'created_at'
      sortDirection.value = 'asc'
      break
    case 'name_asc':
      sortField.value = 'name'
      sortDirection.value = 'asc'
      break
    case 'name_desc':
      sortField.value = 'name'
      sortDirection.value = 'desc'
      break
  }
}

const copyToClipboard = async () => {
  const items = filteredItems.value
  if (items.length === 0) return

  const headers = ['ID', 'Name', 'Quantity', 'Unit', 'Weight (kg)', 'Customer', 'Contract', 'Status', 'Length', 'Width', 'Height', 'Notes', 'Created']

  const rows = items.map(i => {
    const customer = i.customer_id ? getCustomerName(i.customer_id) : ''
    const contract = i.duration_type ? `${i.duration_type}${i.duration ? ` (${i.duration})` : ''}` : ''
    return [
      i.id,
      i.name,
      i.quantity,
      i.quantity_unit,
      i.weight ?? '',
      customer,
      contract,
      i.status || '',
      i.length ?? '',
      i.width ?? '',
      i.height ?? '',
      i.notes || '',
      formatDate(i.created_at),
    ]
  })

  const tsv = [headers, ...rows].map(row => row.map(cell => String(cell).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t')).join('\n')

  try {
    await navigator.clipboard.writeText(tsv)
    copied.value = true
    push.success('Copied to clipboard! Paste into any spreadsheet')
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    const textarea = document.createElement('textarea')
    textarea.value = tsv
    textarea.style.position = 'fixed'
    textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    push.success('Copied to clipboard!')
    setTimeout(() => { copied.value = false }, 2000)
  }
}

const clearFilters = () => {
  filters.value = {
    status: '',
    has_contract: '',
    has_customer: '',
    quantity_unit: '',
    date_from: '',
    date_to: ''
  }
}

const clearAllFilters = () => {
  searchQuery.value = ''
  clearFilters()
}

const toggleSort = (field: 'name' | 'created_at') => {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'asc'
  }
  // Update currentSort to match
  if (field === 'name') {
    currentSort.value = sortDirection.value === 'asc' ? 'name_asc' : 'name_desc'
  } else {
    currentSort.value = sortDirection.value === 'asc' ? 'oldest' : 'newest'
  }
}

const filteredItems = computed(() => {
  const items = [...itemsStore.items]
  let result = items

  // Search across ALL fields
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase().trim()
    result = items.filter(i => {
      // Search in all string fields
      const searchableFields = [
        i.name,
        i.notes,
        String(i.id),
        i.quantity_unit,
        i.duration_type,
        i.status,
        String(i.quantity),
        i.weight ? String(i.weight) : '',
        i.length ? String(i.length) : '',
        i.width ? String(i.width) : '',
        i.height ? String(i.height) : '',
        i.price ? String(i.price) : '',
        i.security_deposit ? String(i.security_deposit) : '',
        i.estimated_value ? String(i.estimated_value) : '',
      ]

      // Search in customer fields
      if (i.customer_id) {
        const customer = customersStore.customers.find(c => c.id === i.customer_id)
        if (customer) {
          searchableFields.push(customer.name)
          searchableFields.push(customer.phone)
          searchableFields.push(customer.email || '')
          searchableFields.push(customer.company || '')
          searchableFields.push(customer.address || '')
          searchableFields.push(customer.id_type || '')
          searchableFields.push(customer.id_number || '')
        }
        searchableFields.push(String(i.customer_id))
      }

      const allText = searchableFields
        .filter(f => f != null)
        .join(' ')
        .toLowerCase()

      return allText.includes(q)
    })
  }

  // Filter by Status
  if (filters.value.status) {
    result = result.filter(i => i.status === filters.value.status)
  }

  // Filter by Has Contract
  if (filters.value.has_contract) {
    const hasContractFilter = filters.value.has_contract === 'true'
    result = result.filter(i => {
      const has = !!(i.duration_type && i.status)
      return has === hasContractFilter
    })
  }

  // Filter by Has Customer
  if (filters.value.has_customer) {
    const hasCustomer = filters.value.has_customer === 'true'
    result = result.filter(i => {
      const has = !!i.customer_id
      return has === hasCustomer
    })
  }

  // Filter by Quantity Unit
  if (filters.value.quantity_unit) {
    result = result.filter(i => i.quantity_unit === filters.value.quantity_unit)
  }

  // Filter by Date Range
  if (filters.value.date_from) {
    const from = new Date(filters.value.date_from)
    result = result.filter(i => new Date(i.created_at) >= from)
  }
  if (filters.value.date_to) {
    const to = new Date(filters.value.date_to)
    to.setHours(23, 59, 59, 999)
    result = result.filter(i => new Date(i.created_at) <= to)
  }

  // Sort
  result.sort((a, b) => {
    let comparison = 0
    if (sortField.value === 'name') {
      comparison = a.name.localeCompare(b.name)
    } else if (sortField.value === 'created_at') {
      comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    }
    return sortDirection.value === 'desc' ? -comparison : comparison
  })

  return result
})

const getDimensions = (item: Item): string => {
  const parts = []
  if (item.length) parts.push(`L: ${item.length}`)
  if (item.width) parts.push(`W: ${item.width}`)
  if (item.height) parts.push(`H: ${item.height}`)
  return parts.length > 0 ? parts.join(' × ') : '—'
}

const formatDate = (d: string): string => {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
}
</script>
