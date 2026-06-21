<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">

    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Contracts</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ contractsStore.contracts.length }}
          total</span>
      </div>

      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search by ID, customer name, phone..."
            class="input w-full sm:w-64 pl-9 pr-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent focus:border-transparent placeholder:text-muted" />
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
            <span class="hidden sm:inline">New Contract</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Table -->
    <div
      class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
      <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
        @click="toggleSort('id')">
        ID
        <svg v-if="sortField === 'id'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
          fill="currentColor" viewBox="0 0 24 24">
          <path d="M7 10l5 5 5-5z" />
        </svg>
      </div>
      <div class="col-span-2">Customer</div>
      <div class="col-span-2">Duration</div>
      <div class="col-span-2">Value</div>
      <div class="col-span-2 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
        @click="toggleSort('status')">
        Status
        <svg v-if="sortField === 'status'" class="w-3 h-3" :class="{ 'rotate-180': sortDirection === 'desc' }"
          fill="currentColor" viewBox="0 0 24 24">
          <path d="M7 10l5 5 5-5z" />
        </svg>
      </div>
      <div class="col-span-2 text-right pr-0.5">Actions</div>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-4">
        <svg class="animate-spin w-10 h-10 text-accent mx-auto" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
          <path class="opacity-75" fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
        </svg>
        <p class="text-sm text-muted">Loading contracts...</p>
      </div>
    </div>

    <div v-else-if="filteredContracts.length === 0" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-3">
        <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
          <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <p class="text-sm font-medium text-primary">No contracts found</p>
        <p class="text-xs text-muted mt-1">{{ searchQuery ? 'adjust your search' : 'create new contract' }}</p>
      </div>
    </div>

    <div v-else class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-max">
        <ContractRow v-for="contract in filteredContracts" :key="contract.id" :contract="contract"
          @contract-updated="handleContractUpdated" @view-contract="openDetailDialog" />
      </div>
    </div>

    <div v-if="!loading && filteredContracts.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredContracts.length }} of {{ contractsStore.contracts.length }}
        contracts</p>
      <button @click="refreshContracts"
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

    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateContractForm @contract-created="handleContractCreated" />
    </BaseDialog>

    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedContract" class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-bold text-primary">Contract #{{ selectedContract.id }}</h2>
            <p class="text-sm text-muted capitalize">{{ selectedContract.status }}</p>
          </div>
          <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold border"
            :class="selectedContract.status === 'active' ? 'status-success' : selectedContract.status === 'completed' ? 'status-info' : 'status-warning'">
            {{ selectedContract.status }}
          </span>
        </div>

        <!-- Customer Info Card -->
        <div class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Customer</h3>
          <div class="flex items-start gap-4">
            <div
              class="w-12 h-12 rounded-full bg-surface border border-default flex items-center justify-center shrink-0">
              <span class="text-sm font-semibold text-muted">#{{ selectedContract.customer_id }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-primary">{{ getCustomerName(selectedContract.customer_id) }}</p>
              <p class="text-xs text-secondary">{{ getCustomer(selectedContract.customer_id)?.phone || '—' }}</p>
              <p v-if="getCustomer(selectedContract.customer_id)?.company" class="text-xs text-muted">{{
                getCustomer(selectedContract.customer_id)?.company }}</p>
              <div class="grid grid-cols-2 gap-2 mt-2">
                <div>
                  <p class="text-[10px] text-muted uppercase">Email</p>
                  <p class="text-xs text-primary">{{ getCustomer(selectedContract.customer_id)?.email || '—' }}</p>
                </div>
                <div>
                  <p class="text-[10px] text-muted uppercase">Address</p>
                  <p class="text-xs text-primary">{{ getCustomer(selectedContract.customer_id)?.address || '—' }}</p>
                </div>
                <div>
                  <p class="text-[10px] text-muted uppercase">ID Type</p>
                  <p class="text-xs text-primary capitalize">{{ getCustomer(selectedContract.customer_id)?.id_type ||
                    '—' }}</p>
                </div>
                <div>
                  <p class="text-[10px] text-muted uppercase">ID Number</p>
                  <p class="text-xs text-primary">{{ getCustomer(selectedContract.customer_id)?.id_number || '—' }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Duration</p>
            <p class="text-sm text-primary">{{ selectedContract.duration || '—' }} {{ selectedContract.duration_type }}
            </p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Price</p>
            <p class="text-sm text-primary">{{ selectedContract.price ? '$' + selectedContract.price : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Security Deposit</p>
            <p class="text-sm text-primary">{{ selectedContract.security_deposit ? '$' +
              selectedContract.security_deposit : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Est. Value</p>
            <p class="text-sm text-primary">{{ selectedContract.estimated_value ? '$' + selectedContract.estimated_value
              : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Created</p>
            <p class="text-sm text-primary">{{ formatDate(selectedContract.created_at) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase tracking-wider">Ended</p>
            <p class="text-sm text-primary">{{ selectedContract.ended_at ? formatDate(selectedContract.ended_at) : '—'
              }}</p>
          </div>
          <div class="space-y-1 md:col-span-2">
            <p class="text-xs text-muted uppercase tracking-wider">Notes</p>
            <p class="text-sm text-primary whitespace-pre-wrap">{{ selectedContract.notes || '—' }}</p>
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
import { ref, computed, onMounted } from 'vue'
import { useContractsStore } from '@/stores/contracts'
import { useCustomersStore } from '@/stores/customers'
import type { Contract } from '@/types/contract'
import ContractRow from './ContractRow.vue'
import CreateContractForm from './CreateContractForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const contractsStore = useContractsStore()
const customersStore = useCustomersStore()

const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedContract = ref<Contract | null>(null)
const showSortMenu = ref(false)
const sortField = ref<'id' | 'status'>('id')
const sortDirection = ref<'asc' | 'desc'>('desc')
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'status_asc', label: 'Status A-Z' },
  { key: 'status_desc', label: 'Status Z-A' },
]

const getCustomer = (customerId: number) => {
  return customersStore.customers.find(c => c.id === customerId)
}

const getCustomerName = (customerId: number): string => {
  const customer = getCustomer(customerId)
  return customer?.name || `Customer #${customerId}`
}

onMounted(async () => {
  if (customersStore.customers.length === 0) await customersStore.fetchCustomers()
  await fetchContracts()
})

const fetchContracts = async () => { loading.value = true; try { await contractsStore.fetchContracts() } finally { loading.value = false } }
const refreshContracts = async () => { refreshing.value = true; try { await contractsStore.fetchContracts() } finally { refreshing.value = false } }
const handleContractUpdated = async () => { await contractsStore.fetchContracts() }
const handleContractCreated = async () => { showCreateDialog.value = false; await contractsStore.fetchContracts() }
const openCreateDialog = () => { showCreateDialog.value = true }

const openDetailDialog = (contract: Contract) => {
  selectedContract.value = contract
  showDetailDialog.value = true
}

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }

const applySort = (key: string) => {
  currentSort.value = key; showSortMenu.value = false
  switch (key) {
    case 'newest': sortField.value = 'id'; sortDirection.value = 'desc'; break
    case 'oldest': sortField.value = 'id'; sortDirection.value = 'asc'; break
    case 'status_asc': sortField.value = 'status'; sortDirection.value = 'asc'; break
    case 'status_desc': sortField.value = 'status'; sortDirection.value = 'desc'; break
  }
}

const toggleSort = (field: 'id' | 'status') => {
  if (sortField.value === field) sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  else { sortField.value = field; sortDirection.value = 'asc' }
}

const copyToClipboard = async () => {
  const contracts = filteredContracts.value
  if (contracts.length === 0) return

  const headers = ['ID', 'Customer Name', 'Customer Phone', 'Duration Type', 'Duration', 'Price', 'Security Deposit', 'Est. Value', 'Status', 'Ended', 'Created', 'Notes']

  const rows = contracts.map(c => {
    const customer = getCustomer(c.customer_id)
    return [
      c.id,
      customer?.name || '',
      customer?.phone || '',
      c.duration_type,
      c.duration ?? '',
      c.price ?? '',
      c.security_deposit ?? '',
      c.estimated_value ?? '',
      c.status,
      c.ended_at ? formatDate(c.ended_at) : '',
      formatDate(c.created_at),
      c.notes || '',
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
    textarea.style.position = 'fixed'; textarea.style.opacity = '0'
    document.body.appendChild(textarea)
    textarea.select(); document.execCommand('copy')
    document.body.removeChild(textarea)
    copied.value = true
    push.success('Copied to clipboard!')
    setTimeout(() => { copied.value = false }, 2000)
  }
}

const filteredContracts = computed(() => {
  const contracts = [...contractsStore.contracts]
  let result = contracts
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = contracts.filter(c => {
      if (String(c.id).includes(q)) return true
      if (c.status.toLowerCase().includes(q)) return true
      if (c.notes && c.notes.toLowerCase().includes(q)) return true
      const customer = getCustomer(c.customer_id)
      if (customer) {
        if (customer.name.toLowerCase().includes(q)) return true
        if (customer.phone.includes(q)) return true
        if (customer.company && customer.company.toLowerCase().includes(q)) return true
      }
      return false
    })
  }
  result.sort((a, b) => {
    let comparison = 0
    if (sortField.value === 'id') comparison = a.id - b.id
    else if (sortField.value === 'status') comparison = a.status.localeCompare(b.status)
    return sortDirection.value === 'desc' ? -comparison : comparison
  })
  return result
})

const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
</script>
