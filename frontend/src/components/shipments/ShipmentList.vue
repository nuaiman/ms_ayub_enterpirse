<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Shipments</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ shipmentsStore.shipments.length }}
          total</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search by ID, item, driver, location..."
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
            <span class="hidden sm:inline">New Shipment</span>
          </button>
        </div>
      </div>
    </div>

    <div
      class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
      <div class="col-span-1 cursor-pointer hover:text-primary transition-colors flex items-center gap-1"
        @click="toggleSort('id')">ID<svg v-if="sortField === 'id'" class="w-3 h-3"
          :class="{ 'rotate-180': sortDirection === 'desc' }" fill="currentColor" viewBox="0 0 24 24">
          <path d="M7 10l5 5 5-5z" />
        </svg></div>
      <div class="col-span-2">Type / Status</div>
      <div class="col-span-3">Item</div>
      <div class="col-span-2">Qty / Route</div>
      <div class="col-span-2">Financial</div>
      <div class="col-span-2 text-right pr-0.5">Actions</div>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <svg class="animate-spin w-10 h-10 text-accent" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
    </div>

    <div v-else-if="filteredShipments.length === 0" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-3">
        <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
          <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M13 16V6a1 1 0 00-1-1H4a1 1 0 00-1 1v10a1 1 0 001 1h8zm0 0v3a1 1 0 01-1 1H4a1 1 0 01-1-1v-3m5-10h9a1 1 0 011 1v10a1 1 0 01-1 1h-3m-5-5h.01" />
          </svg>
        </div>
        <p class="text-sm font-medium text-primary">No shipments found</p>
      </div>
    </div>

    <div v-else class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-max">
        <ShipmentRow v-for="s in filteredShipments" :key="s.id" :shipment="s" @shipment-updated="handleShipmentUpdated"
          @view-shipment="openDetailDialog" />
      </div>
    </div>

    <div v-if="!loading && filteredShipments.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredShipments.length }} of {{ shipmentsStore.shipments.length }}</p>
      <button @click="refreshShipments"
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
      <CreateShipmentForm @shipment-created="handleShipmentCreated" />
    </BaseDialog>

    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedShipment" class="space-y-5">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-bold text-primary">Shipment #{{ selectedShipment.id }}</h2>
            <div class="flex items-center gap-2 mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="selectedShipment.shipment_type === 'pickup' ? 'status-info' : 'status-success'">{{
                  selectedShipment.shipment_type }}</span>
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium border"
                :class="getStatusClass(selectedShipment.status)">{{ selectedShipment.status }}</span>
            </div>
          </div>
          <span class="text-xs text-muted">Qty: {{ selectedShipment.quantity }}</span>
        </div>

        <!-- Item Info Card -->
        <div class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Item</h3>
          <div class="flex items-start gap-3">
            <div
              class="w-10 h-10 rounded-lg bg-surface border border-default flex items-center justify-center shrink-0">
              <span class="text-xs font-semibold text-muted">#{{ selectedShipment.item_id }}</span>
            </div>
            <div>
              <p class="text-sm font-medium text-primary">{{ getItemName(selectedShipment.item_id) }}</p>
              <p class="text-xs text-muted">{{ getItemQuantity(selectedShipment.item_id) }} {{
                getItemUnit(selectedShipment.item_id) }} available</p>
              <p v-if="getItemWeight(selectedShipment.item_id)" class="text-xs text-muted">{{
                getItemWeight(selectedShipment.item_id) }} kg</p>
            </div>
          </div>
        </div>

        <!-- Delivery Details -->
        <div v-if="selectedShipment.shipment_type === 'delivery'"
          class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Delivery Details</h3>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Vehicle</p>
              <p class="text-xs text-primary">{{ selectedShipment.vehicle_number || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Driver</p>
              <p class="text-xs text-primary">{{ selectedShipment.driver_name || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Driver Phone</p>
              <p class="text-xs text-primary">{{ selectedShipment.driver_phone || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Company</p>
              <p class="text-xs text-primary">{{ selectedShipment.company_name || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">From</p>
              <p class="text-xs text-primary">{{ selectedShipment.from_location || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">To</p>
              <p class="text-xs text-primary">{{ selectedShipment.to_location || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Receiver</p>
              <p class="text-xs text-primary">{{ selectedShipment.receiver_name || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Receiver Phone</p>
              <p class="text-xs text-primary">{{ selectedShipment.receiver_phone || '—' }}</p>
            </div>
          </div>
        </div>

        <!-- Financial -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Cust. Charge</p>
            <p class="text-sm font-medium text-primary">{{ selectedShipment.customer_charge ?
              '$' + selectedShipment.customer_charge : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Co. Cost</p>
            <p class="text-sm font-medium text-primary">{{ selectedShipment.company_cost ?
              '$' + selectedShipment.company_cost : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Co. Paid</p>
            <p class="text-sm font-medium text-primary">{{ selectedShipment.company_paid ?
              '$' + selectedShipment.company_paid : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Cust. Paid</p>
            <p class="text-sm font-medium text-primary">{{ selectedShipment.customer_paid ?
              '$' + selectedShipment.customer_paid : '—' }}</p>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-3">
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Created</p>
            <p class="text-xs text-primary">{{ formatDateTime(selectedShipment.created_at) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Shipped</p>
            <p class="text-xs text-primary">{{ selectedShipment.shipped_at ? formatDateTime(selectedShipment.shipped_at)
              : '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-[10px] text-muted uppercase">Received</p>
            <p class="text-xs text-primary">{{ selectedShipment.received_at ?
              formatDateTime(selectedShipment.received_at) : '—' }}</p>
          </div>
        </div>

        <div v-if="selectedShipment.notes" class="space-y-1">
          <p class="text-[10px] text-muted uppercase">Notes</p>
          <p class="text-xs text-primary whitespace-pre-wrap">{{ selectedShipment.notes }}</p>
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
import { useShipmentsStore } from '@/stores/shipments'
import { useItemsStore } from '@/stores/items'
import type { Shipment } from '@/types/shipment'
import ShipmentRow from './ShipmentRow.vue'
import CreateShipmentForm from './CreateShipmentForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const shipmentsStore = useShipmentsStore()
const itemsStore = useItemsStore()
const loading = ref(true)
const refreshing = ref(false)
const searchQuery = ref('')
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedShipment = ref<Shipment | null>(null)
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

const getItem = (itemId: number) => itemsStore.items.find(i => i.id === itemId)
const getItemName = (id: number) => getItem(id)?.name || `Item #${id}`
const getItemQuantity = (id: number) => getItem(id)?.quantity ?? '?'
const getItemUnit = (id: number) => getItem(id)?.quantity_unit ?? ''
const getItemWeight = (id: number) => getItem(id)?.weight ?? null

onMounted(async () => {
  if (itemsStore.items.length === 0) await itemsStore.fetchItems()
  await fetchShipments()
})
const fetchShipments = async () => { loading.value = true; try { await shipmentsStore.fetchShipments() } finally { loading.value = false } }
const refreshShipments = async () => { refreshing.value = true; try { await shipmentsStore.fetchShipments() } finally { refreshing.value = false } }
const handleShipmentUpdated = async () => { await shipmentsStore.fetchShipments() }
const handleShipmentCreated = async () => { showCreateDialog.value = false; await shipmentsStore.fetchShipments() }
const openCreateDialog = () => { showCreateDialog.value = true }
const openDetailDialog = (s: Shipment) => { selectedShipment.value = s; showDetailDialog.value = true }

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
  const shipments = filteredShipments.value
  if (shipments.length === 0) return
  const headers = ['ID', 'Type', 'Status', 'Item ID', 'Item Name', 'Qty', 'Vehicle', 'Driver', 'Driver Phone', 'Company', 'From', 'To', 'Receiver', 'Receiver Phone', 'Cust. Charge', 'Co. Cost', 'Co. Paid', 'Cust. Paid', 'Shipped', 'Received', 'Notes']
  const rows = shipments.map(s => [s.id, s.shipment_type, s.status, s.item_id, getItemName(s.item_id), s.quantity, s.vehicle_number || '', s.driver_name || '', s.driver_phone || '', s.company_name || '', s.from_location || '', s.to_location || '', s.receiver_name || '', s.receiver_phone || '', s.customer_charge ?? '', s.company_cost ?? '', s.company_paid ?? '', s.customer_paid ?? '', s.shipped_at ? formatDate(s.shipped_at) : '', s.received_at ? formatDate(s.received_at) : '', s.notes || ''])
  const tsv = [headers, ...rows].map(r => r.map(c => String(c).replace(/\t/g, ' ').replace(/\n/g, ' ')).join('\t')).join('\n')
  try { await navigator.clipboard.writeText(tsv); copied.value = true; push.success('Copied!'); setTimeout(() => { copied.value = false }, 2000) }
  catch { const ta = document.createElement('textarea'); ta.value = tsv; ta.style.position = 'fixed'; ta.style.opacity = '0'; document.body.appendChild(ta); ta.select(); document.execCommand('copy'); document.body.removeChild(ta); copied.value = true; push.success('Copied!'); setTimeout(() => { copied.value = false }, 2000) }
}

const filteredShipments = computed(() => {
  const shipments = [...shipmentsStore.shipments]
  let result = shipments
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = shipments.filter(s => {
      if (String(s.id).includes(q) || s.shipment_type.includes(q) || s.status.includes(q) || String(s.item_id).includes(q) || (s.notes && s.notes.toLowerCase().includes(q))) return true
      if (s.driver_name && s.driver_name.toLowerCase().includes(q)) return true
      if (s.vehicle_number && s.vehicle_number.toLowerCase().includes(q)) return true
      if (s.from_location && s.from_location.toLowerCase().includes(q)) return true
      if (s.to_location && s.to_location.toLowerCase().includes(q)) return true
      if (s.receiver_name && s.receiver_name.toLowerCase().includes(q)) return true
      if (s.company_name && s.company_name.toLowerCase().includes(q)) return true
      const item = getItem(s.item_id)
      if (item && item.name.toLowerCase().includes(q)) return true
      return false
    })
  }
  result.sort((a, b) => { let c = 0; if (sortField.value === 'id') c = a.id - b.id; else c = a.status.localeCompare(b.status); return sortDirection.value === 'desc' ? -c : c })
  return result
})

const getStatusClass = (s: string): string => ({ pending: 'status-warning', in_transit: 'status-info', completed: 'status-success', cancelled: 'status-warning' }[s] || '')
const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
const formatDateTime = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' })
</script>
