<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">

    <!-- HEADER -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Customers</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ customersStore.customers.length }} total</span>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search customers..."
            class="input w-full sm:w-80 pl-9 pr-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent focus:border-transparent placeholder:text-muted"
          />
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
          <div class="relative">
            <button @click="toggleSortDropdown"
              class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors" title="Sort">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
              </svg>
            </button>
            <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95"
              enter-to-class="opacity-100 scale-100" leave-active-class="transition ease-in duration-150"
              leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
              <div v-if="showSortMenu"
                class="absolute right-0 top-10 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
                <div class="px-3 py-2 border-b border-divider"><p class="text-xs font-semibold text-muted uppercase">Sort By</p></div>
                <button v-for="option in sortOptions" :key="option.key" @click="applySort(option.key)"
                  class="w-full flex items-center justify-between px-3 py-2 text-sm text-secondary hover:bg-surface-alt transition-colors">
                  {{ option.label }}
                  <svg v-if="currentSort === option.key" class="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                </button>
              </div>
            </Transition>
          </div>

          <button @click="copyToClipboard"
            class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors relative" title="Copy to clipboard">
            <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>

          <button @click="openCreateDialog"
            class="h-9 px-4 flex items-center gap-2 bg-accent text-accent-foreground rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" /></svg>
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
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">ID Type</label>
            <select v-model="filters.id_type" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Types</option>
              <option value="nid">NID</option>
              <option value="passport">Passport</option>
              <option value="driving_license">Driving License</option>
              <option value="birth_certificate">Birth Certificate</option>
              <option value="trade_license">Trade License</option>
              <option value="other">Other</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Has Items</label>
            <select v-model="filters.has_items" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Customers</option>
              <option value="true">Has Items</option>
              <option value="false">No Items</option>
            </select>
          </div>
          <div class="space-y-1">
            <label class="text-xs font-medium text-muted">Has Contract</label>
            <select v-model="filters.has_contract" class="input w-full px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent">
              <option value="">All Customers</option>
              <option value="true">Has Contract Items</option>
              <option value="false">No Contract Items</option>
            </select>
          </div>
          <div class="space-y-1 md:col-span-3">
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

    <div v-if="showSortMenu" class="fixed inset-0 z-40" @click="showSortMenu = false"></div>

    <!-- TABLE - With horizontal scroll for mobile -->
    <div class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-3xl">
        <!-- TABLE HEADER -->
        <div class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
          <div class="col-span-3 flex items-center gap-1 cursor-pointer hover:text-primary transition-colors" @click="toggleSort('name')">
            Name
            <svg v-if="sortField==='name'" class="w-3 h-3" :class="{'rotate-180':sortDirection==='desc'}" fill="currentColor" viewBox="0 0 24 24"><path d="M7 10l5 5 5-5z"/></svg>
          </div>
          <div class="col-span-2">Phone</div>
          <div class="col-span-2">Company</div>
          <div class="col-span-2">Items</div>
          <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1" @click="toggleSort('created_at')">
            Created
            <svg v-if="sortField==='created_at'" class="w-3 h-3" :class="{'rotate-180':sortDirection==='desc'}" fill="currentColor" viewBox="0 0 24 24"><path d="M7 10l5 5 5-5z"/></svg>
          </div>
          <div class="col-span-1 text-right pr-0.5">Actions</div>
        </div>

        <!-- LOADING -->
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="text-center space-y-4">
            <svg class="animate-spin w-10 h-10 text-accent mx-auto" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            <p class="text-sm text-muted">Loading customers...</p>
          </div>
        </div>

        <!-- EMPTY -->
        <div v-else-if="filteredCustomers.length === 0" class="flex items-center justify-center py-12">
          <div class="text-center space-y-3">
            <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
              <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <p class="text-sm font-medium text-primary">No customers found</p>
            <p class="text-xs text-muted mt-1">{{ searchQuery || hasActiveFilters ? 'adjust your filters' : 'add new customer' }}</p>
            <button v-if="hasActiveFilters" @click="clearAllFilters" class="text-xs text-accent hover:underline">
              Clear all filters
            </button>
          </div>
        </div>

        <!-- CUSTOMER LIST -->
        <div v-else>
          <CustomerRow v-for="customer in filteredCustomers" :key="customer.id" :customer="customer"
            @customer-updated="handleCustomerUpdated" @view-customer="openDetailDialog" />
        </div>
      </div>
    </div>

    <!-- FOOTER -->
    <div v-if="!loading && filteredCustomers.length > 0" class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredCustomers.length }} of {{ customersStore.customers.length }}</p>
      <div class="flex items-center gap-2">
        <span v-if="hasActiveFilters" class="text-xs text-accent">(filtered)</span>
        <button @click="refreshCustomers" class="button px-3 py-1.5 text-xs rounded-lg hover-surface transition-colors inline-flex items-center gap-1.5" :disabled="refreshing">
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': refreshing }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- CREATE CUSTOMER DIALOG -->
    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateCustomerForm @customer-created="handleCustomerCreated" />
    </BaseDialog>

    <!-- CUSTOMER DETAIL DIALOG -->
    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedCustomer" class="space-y-5">
        <!-- Customer Header -->
        <div class="flex items-center gap-4">
          <div class="shrink-0">
            <div class="w-16 h-16 rounded-full bg-surface-alt border border-default flex items-center justify-center">
              <span class="text-xl font-semibold text-muted">{{ getInitials(selectedCustomer.name) }}</span>
            </div>
          </div>
          <div>
            <h2 class="text-xl font-bold text-primary">{{ selectedCustomer.name }}</h2>
            <p class="text-sm text-muted">{{ selectedCustomer.phone }}</p>
            <p v-if="selectedCustomer.company" class="text-sm text-secondary mt-0.5">{{ selectedCustomer.company }}</p>
          </div>
        </div>

        <div class="flex items-center gap-4 pb-4 border-b border-divider flex-wrap">
          <span class="text-xs text-muted">ID: {{ selectedCustomer.id }}</span>
          <span v-if="selectedCustomer.user_id" class="text-xs text-muted">Managed by: #{{ selectedCustomer.user_id }}</span>
          <span v-if="selectedCustomer.id_type" class="text-xs text-muted">ID Type: {{ selectedCustomer.id_type }}</span>
          <span v-if="selectedCustomer.id_number" class="text-xs text-muted">ID: {{ selectedCustomer.id_number }}</span>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1"><p class="text-xs text-muted uppercase tracking-wider">Email</p><p class="text-sm text-primary">{{ selectedCustomer.email || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase tracking-wider">Phone</p><p class="text-sm text-primary">{{ selectedCustomer.phone }}</p></div>
          <div class="space-y-1 md:col-span-2"><p class="text-xs text-muted uppercase tracking-wider">Address</p><p class="text-sm text-primary">{{ selectedCustomer.address || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase tracking-wider">ID Type</p><p class="text-sm text-primary capitalize">{{ selectedCustomer.id_type || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase tracking-wider">ID Number</p><p class="text-sm text-primary">{{ selectedCustomer.id_number || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase tracking-wider">Created</p><p class="text-sm text-primary">{{ formatDate(selectedCustomer.created_at) }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase tracking-wider">Updated</p><p class="text-sm text-primary">{{ formatDate(selectedCustomer.updated_at) }}</p></div>
          <div class="space-y-1 md:col-span-2"><p class="text-xs text-muted uppercase tracking-wider">Notes</p><p class="text-sm text-primary whitespace-pre-wrap">{{ selectedCustomer.notes || '—' }}</p></div>
        </div>

        <!-- Items Section -->
        <div v-if="customerItems.length > 0" class="border-t border-divider pt-4">
          <h3 class="text-sm font-semibold text-primary mb-3">Items ({{ customerItems.length }})</h3>
          <div class="grid grid-cols-1 gap-2">
            <div v-for="item in customerItems" :key="item.id"
              class="flex items-center gap-3 p-3 rounded-lg border border-default bg-surface-alt hover:bg-surface transition-colors cursor-pointer"
              @click="viewItem(item)">
              <div class="w-10 h-10 rounded-lg bg-surface border border-default flex items-center justify-center shrink-0 overflow-hidden">
                <img v-if="item.image_url" :src="getImageUrl(item.image_url) || undefined" :alt="item.name" class="w-full h-full object-cover" />
                <svg v-else class="w-5 h-5 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium text-primary truncate">{{ item.name }}</p>
                <div class="flex items-center gap-2 mt-0.5">
                  <span class="text-xs text-muted">{{ item.quantity }} {{ item.quantity_unit }}</span>
                  <span v-if="item.weight" class="text-xs text-muted">· {{ item.weight }} kg</span>
                  <span v-if="item.duration_type" class="text-xs text-accent">· Contract: {{ item.duration_type }}</span>
                </div>
                <div v-if="item.status" class="mt-1">
                  <span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium border"
                    :class="item.status === 'active' ? 'status-success' : item.status === 'completed' ? 'status-info' : 'status-warning'">
                    {{ item.status }}
                  </span>
                </div>
              </div>
              <svg class="w-4 h-4 text-muted shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </div>
        </div>

        <!-- No Items Message -->
        <div v-else class="border-t border-divider pt-4">
          <p class="text-sm text-muted text-center py-4">No items found for this customer</p>
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
import { useCustomersStore } from '@/stores/customers'
import { useItemsStore } from '@/stores/items'
import type { Customer } from '@/types/customer'
import type { Item } from '@/types/item'
import CustomerRow from './CustomerRow.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import CreateCustomerForm from './CreateCustomerForm.vue'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'
// import { useRouter } from 'vue-router'

// const router = useRouter()
const customersStore = useCustomersStore()
const itemsStore = useItemsStore()

const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const showAdvancedFilters = ref(false)

// Filter state
const filters = ref({
  id_type: '',
  has_items: '',
  has_contract: '',
  date_from: '',
  date_to: ''
})

const sortField = ref<'name' | 'created_at'>('created_at')
const sortDirection = ref<'asc' | 'desc'>('desc')
const showSortMenu = ref(false)
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedCustomer = ref<Customer | null>(null)
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'name_asc', label: 'Name A-Z' },
  { key: 'name_desc', label: 'Name Z-A' },
]

const hasActiveFilters = computed(() => {
  return !!searchQuery.value ||
         !!filters.value.id_type ||
         !!filters.value.has_items ||
         !!filters.value.has_contract ||
         !!filters.value.date_from ||
         !!filters.value.date_to
})

const activeFilterCount = computed(() => {
  let count = 0
  if (filters.value.id_type) count++
  if (filters.value.has_items) count++
  if (filters.value.has_contract) count++
  if (filters.value.date_from) count++
  if (filters.value.date_to) count++
  return count
})

const customerItems = computed(() => {
  if (!selectedCustomer.value) return []
  return itemsStore.items.filter(i => i.customer_id === selectedCustomer.value!.id)
})

// Debounced search
let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(searchQuery, () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    // Search will be applied via computed property
  }, 300)
})

onMounted(async () => {
  await fetchCustomers()
  await itemsStore.fetchItems()
})

const fetchCustomers = async () => { loading.value = true; try { await customersStore.fetchCustomers() } finally { loading.value = false } }
const refreshCustomers = async () => { refreshing.value = true; try { await customersStore.fetchCustomers() } finally { refreshing.value = false } }
const handleCustomerUpdated = async () => { await customersStore.fetchCustomers() }
const handleCustomerCreated = async () => { showCreateDialog.value = false; await customersStore.fetchCustomers() }
const openCreateDialog = () => { showCreateDialog.value = true }

const openDetailDialog = async (customer: Customer) => {
  selectedCustomer.value = customer
  if (itemsStore.items.length === 0) await itemsStore.fetchItems()
  showDetailDialog.value = true
}

// const viewItem = (item: Item) => {
//   showDetailDialog.value = false
//   router.push(`/items/${item.id}`)
// }

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const viewItem = (item: Item) => {
  showDetailDialog.value = false
  // router.push(`/items/${item.id}`)
}

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }

const applySort = (key: string) => {
  currentSort.value = key; showSortMenu.value = false
  switch (key) {
    case 'newest': sortField.value = 'created_at'; sortDirection.value = 'desc'; break
    case 'oldest': sortField.value = 'created_at'; sortDirection.value = 'asc'; break
    case 'name_asc': sortField.value = 'name'; sortDirection.value = 'asc'; break
    case 'name_desc': sortField.value = 'name'; sortDirection.value = 'desc'; break
  }
}

const clearFilters = () => {
  filters.value = {
    id_type: '',
    has_items: '',
    has_contract: '',
    date_from: '',
    date_to: ''
  }
}

const clearAllFilters = () => {
  searchQuery.value = ''
  clearFilters()
}

const copyToClipboard = async () => {
  const customers = filteredCustomers.value
  if (customers.length === 0) return
  const headers = ['ID', 'Name', 'Phone', 'Email', 'Company', 'Address', 'ID Type', 'ID Number', 'Notes', 'Created']
  const rows = customers.map(c => [c.id, c.name, c.phone, c.email || '', c.company || '', c.address || '', c.id_type || '', c.id_number || '', c.notes || '', formatDate(c.created_at)])
  const tsv = [headers, ...rows].map(row => row.map(cell => String(cell).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t')).join('\n')
  try { await navigator.clipboard.writeText(tsv); copied.value = true; push.success('Copied!'); setTimeout(() => { copied.value = false }, 2000) }
  catch { const ta = document.createElement('textarea'); ta.value = tsv; ta.style.position = 'fixed'; ta.style.opacity = '0'; document.body.appendChild(ta); ta.select(); document.execCommand('copy'); document.body.removeChild(ta); copied.value = true; push.success('Copied!'); setTimeout(() => { copied.value = false }, 2000) }
}

const filteredCustomers = computed(() => {
  const customers = [...customersStore.customers]
  let result = customers

  // Search across ALL fields
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase().trim()
    result = customers.filter(c => {
      // Search in all string fields
      const searchableFields = [
        c.name,
        c.phone,
        c.email,
        c.company,
        c.address,
        c.id_type,
        c.id_number,
        c.notes,
        String(c.id) // Search by ID too
      ]

      // Also search in nested fields
      const allText = searchableFields
        .filter(f => f != null)
        .join(' ')
        .toLowerCase()

      return allText.includes(query)
    })
  }

  // Filter by ID Type
  if (filters.value.id_type) {
    result = result.filter(c => c.id_type === filters.value.id_type)
  }

  // Filter by Has Items
  if (filters.value.has_items) {
    const hasItems = filters.value.has_items === 'true'
    result = result.filter(c => {
      const itemCount = itemsStore.items.filter(i => i.customer_id === c.id).length
      return hasItems ? itemCount > 0 : itemCount === 0
    })
  }

  // Filter by Has Contract
  if (filters.value.has_contract) {
    const hasContract = filters.value.has_contract === 'true'
    result = result.filter(c => {
      const contractItems = itemsStore.items.filter(i =>
        i.customer_id === c.id && i.duration_type && i.status
      ).length
      return hasContract ? contractItems > 0 : contractItems === 0
    })
  }

  // Filter by Date Range
  if (filters.value.date_from) {
    const from = new Date(filters.value.date_from)
    result = result.filter(c => new Date(c.created_at) >= from)
  }
  if (filters.value.date_to) {
    const to = new Date(filters.value.date_to)
    to.setHours(23, 59, 59, 999)
    result = result.filter(c => new Date(c.created_at) <= to)
  }

  // Sort
  result.sort((a, b) => {
    let comparison = 0
    if (sortField.value === 'name') comparison = a.name.localeCompare(b.name)
    else if (sortField.value === 'created_at') comparison = new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
    return sortDirection.value === 'desc' ? -comparison : comparison
  })

  return result
})

const toggleSort = (field: 'name' | 'created_at') => {
  if (sortField.value === field) sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  else { sortField.value = field; sortDirection.value = 'asc' }
}

const getInitials = (name: string): string => name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
</script>
