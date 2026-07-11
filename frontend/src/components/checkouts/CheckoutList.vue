<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Checkouts</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ checkoutsStore.checkouts.length }}
          total</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor"
            viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search checkouts..."
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
      <div class="min-w-3xl">
<div class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
  <div class="col-span-1">ID</div>
  <div class="col-span-1">Type</div>
  <div class="col-span-2">Item / Customer</div>
  <div class="col-span-1">Qty / Wt</div>
  <div class="col-span-2">Route / Receiver</div>
  <div class="col-span-2">Financial</div>
  <div class="col-span-3 text-right">Actions</div>
</div>

        <div v-if="loading" class="flex items-center justify-center py-12">
          <svg class="animate-spin w-10 h-10 text-accent" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
        </div>

        <div v-else-if="filteredCheckouts.length === 0" class="flex items-center justify-center py-12">
          <p class="text-sm text-muted">No checkouts found</p>
        </div>

        <div v-else>
          <CheckoutRow v-for="c in filteredCheckouts" :key="c.id" :checkout="c"
            @checkout-updated="handleCheckoutUpdated" @view-checkout="openDetailDialog"
            @checkout-deleted="handleCheckoutUpdated" />
        </div>
      </div>
    </div>

    <div v-if="!loading && filteredCheckouts.length > 0"
      class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredCheckouts.length }} of {{ checkoutsStore.checkouts.length }}</p>
    </div>

    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateCheckoutForm @checkout-created="handleCheckoutCreated" />
    </BaseDialog>

    <!-- Full Detail Dialog -->
    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedCheckout" class="space-y-5">
        <!-- Header with Image -->
        <div class="flex items-center gap-4">
          <div v-if="selectedCheckout.image_url"
            class="shrink-0 w-20 h-20 rounded-lg overflow-hidden border border-default">
            <img :src="getImageUrl(selectedCheckout.image_url) || undefined" alt="Checkout"
              class="w-full h-full object-cover" />
          </div>
          <div>
            <h2 class="text-xl font-bold text-primary">Checkout #{{ selectedCheckout.id }}</h2>
            <div class="flex items-center gap-2 mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                :class="selectedCheckout.type === 'pickup' ? 'status-info' : 'status-success'">{{ selectedCheckout.type
                }}</span>
            </div>
          </div>
        </div>

        <!-- Item Card -->
        <div v-if="checkoutItem" class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Item Details</h3>
          <div class="flex items-start gap-3">
            <div v-if="checkoutItem.image_url"
              class="shrink-0 w-14 h-14 rounded-lg overflow-hidden border border-default">
              <img :src="getImageUrl(checkoutItem.image_url) || undefined" :alt="checkoutItem.name"
                class="w-full h-full object-cover" />
            </div>
            <div v-else
              class="shrink-0 w-14 h-14 rounded-lg bg-surface border border-default flex items-center justify-center">
              <span class="text-xs font-semibold text-muted">#{{ checkoutItem.id }}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-primary">{{ checkoutItem.name }}</p>
              <div class="flex items-center gap-3 mt-1 flex-wrap">
                <span class="text-xs text-secondary">{{ checkoutItem.customer_phone }}</span>
                <span v-if="checkoutItem.customer_email" class="text-xs text-muted">{{ checkoutItem.customer_email
                  }}</span>
              </div>
              <div class="flex items-center gap-3 mt-1">
                <span class="text-xs text-muted">{{ checkoutItem.quantity }} {{ checkoutItem.quantity_unit }}
                  available</span>
                <span v-if="checkoutItem.category" class="text-xs text-muted">· {{ checkoutItem.category }}</span>
              </div>
              <div v-if="checkoutItem.weight" class="text-xs text-muted mt-0.5">Weight: {{ checkoutItem.weight }} {{
                checkoutItem.weight_unit || 'kg' }}</div>
            </div>
          </div>
        </div>
        <div v-else class="p-4 rounded-xl border border-default bg-surface-alt">
          <p class="text-sm text-muted">Item #{{ selectedCheckout.item_id }} (not found)</p>
        </div>

        <!-- Checkout Info -->
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Checkout Date</p>
            <p class="text-sm text-primary">{{ formatDate(selectedCheckout.checkout_date) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Quantity</p>
            <p class="text-sm text-primary">{{ selectedCheckout.quantity }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Receiver Name</p>
            <p class="text-sm text-primary">{{ selectedCheckout.receiver_name || '—' }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Receiver Phone</p>
            <p class="text-sm text-primary">{{ selectedCheckout.receiver_phone || '—' }}</p>
          </div>
        </div>

        <!-- Delivery Details -->
        <div v-if="selectedCheckout.type === 'delivery'" class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Delivery Details</h3>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">From Location</p>
              <p class="text-xs text-primary">{{ selectedCheckout.from_location || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">To Location</p>
              <p class="text-xs text-primary">{{ selectedCheckout.to_location || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Vehicle Number</p>
              <p class="text-xs text-primary">{{ selectedCheckout.vehicle_number || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Driver Name</p>
              <p class="text-xs text-primary">{{ selectedCheckout.driver_name || '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Driver Phone</p>
              <p class="text-xs text-primary">{{ selectedCheckout.driver_phone || '—' }}</p>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3 mt-3 pt-3 border-t border-divider">
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Delivery Charge</p>
              <p class="text-xs text-primary">{{ selectedCheckout.delivery_charge ? '৳' +
                selectedCheckout.delivery_charge : '—' }}</p>
            </div>
            <div class="space-y-1">
              <p class="text-[10px] text-muted uppercase">Delivery Cost</p>
              <p class="text-xs text-primary">{{ selectedCheckout.delivery_cost ? '৳' + selectedCheckout.delivery_cost :
                '—' }}</p>
            </div>
          </div>
        </div>

        <!-- Created/Updated -->
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Created</p>
            <p class="text-sm text-primary">{{ formatDate(selectedCheckout.created_at) }}</p>
          </div>
          <div class="space-y-1">
            <p class="text-xs text-muted uppercase">Updated</p>
            <p class="text-sm text-primary">{{ formatDate(selectedCheckout.updated_at) }}</p>
          </div>
        </div>

        <!-- Notes -->
        <div v-if="selectedCheckout.notes" class="space-y-1">
          <p class="text-xs text-muted uppercase">Notes</p>
          <p class="text-sm text-primary whitespace-pre-wrap">{{ selectedCheckout.notes }}</p>
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
import { useCheckoutsStore } from '@/stores/checkout'
import { useItemsStore } from '@/stores/items'
import type { Checkout } from '@/types/checkout'
import CheckoutRow from './CheckoutRow.vue'
import CreateCheckoutForm from './CreateCheckoutForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const checkoutsStore = useCheckoutsStore()
const itemsStore = useItemsStore()
const loading = ref(true)
const searchQuery = ref('')
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedCheckout = ref<Checkout | null>(null)
const showSortMenu = ref(false)
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
]

const checkoutItem = computed(() => {
  if (!selectedCheckout.value) return null
  return itemsStore.items.find(i => i.id === selectedCheckout.value!.item_id) || null
})

onMounted(async () => {
  if (itemsStore.items.length === 0) await itemsStore.fetchItems()
  await fetchCheckouts()
})

const fetchCheckouts = async () => { loading.value = true; try { await checkoutsStore.fetchCheckouts() } finally { loading.value = false } }
const handleCheckoutUpdated = async () => { await checkoutsStore.fetchCheckouts() }
const handleCheckoutCreated = async () => { showCreateDialog.value = false; await checkoutsStore.fetchCheckouts() }
const openCreateDialog = () => { showCreateDialog.value = true }

const openDetailDialog = (checkout: Checkout) => {
  selectedCheckout.value = checkout
  showDetailDialog.value = true
}

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }
const applySort = (key: string) => { currentSort.value = key; showSortMenu.value = false }

const copyToClipboard = async () => {
  const checkouts = filteredCheckouts.value
  if (checkouts.length === 0) return
 const headers = [
  'ID', 'Type', 'Item ID', 'Quantity', 'Weight', 
  'Receiver', 'Receiver Phone', 'From', 'To', 
  'Vehicle', 'Driver', 'Driver Phone',
  'Charge', 'Cost', 'Paid', 'Notes'
]
  const rows = checkouts.map(c => [
  c.id, 
  c.type, 
  c.item_id, 
  c.quantity, 
  c.weight || '', 
  c.receiver_name || '', 
  c.receiver_phone || '', 
  c.from_location || '', 
  c.to_location || '', 
  c.vehicle_number || '', 
  c.driver_name || '', 
  c.driver_phone || '',
  c.delivery_charge || '', 
  c.delivery_cost || '', 
  c.customer_paid || '', 
  c.notes || ''
])
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

const filteredCheckouts = computed(() => {
  let result = [...checkoutsStore.checkouts]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(c =>
String(c.id).includes(q) || 
c.type.includes(q) || 
String(c.item_id).includes(q) ||
String(c.quantity).includes(q) ||
(c.weight && String(c.weight).includes(q)) ||
(c.receiver_name && c.receiver_name.toLowerCase().includes(q)) || 
(c.receiver_phone && c.receiver_phone.includes(q)) ||
(c.from_location && c.from_location.toLowerCase().includes(q)) || 
(c.to_location && c.to_location.toLowerCase().includes(q)) ||
(c.vehicle_number && c.vehicle_number.toLowerCase().includes(q)) || 
(c.driver_name && c.driver_name.toLowerCase().includes(q)) ||
(c.driver_phone && c.driver_phone.includes(q)) ||
(c.notes && c.notes.toLowerCase().includes(q))
    )
  }
  if (currentSort.value === 'newest') result.sort((a, b) => b.id - a.id)
  else result.sort((a, b) => a.id - b.id)
  return result
})

const getStatusClass = (s: string): string => ({
  pending: 'border-yellow-400 bg-yellow-50 text-yellow-700',
  in_transit: 'border-blue-400 bg-blue-50 text-blue-700',
  complete: 'border-green-400 bg-green-50 text-green-700'
}[s] || '')

const formatDate = (d: string): string => d ? new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : '—'
</script>
