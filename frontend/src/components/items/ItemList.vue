<template>
  <div class="flex flex-col h-full min-h-0 overflow-hidden">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-4 mt-1 shrink-0">
      <div class="flex items-center gap-3">
        <h2 class="text-lg font-semibold text-primary">Items</h2>
        <span class="text-sm text-muted bg-surface-alt px-2 py-0.5 rounded-md">{{ itemsStore.items.length }} total</span>
      </div>
      <div class="flex items-center gap-2 flex-wrap">
        <div class="relative flex-1 sm:flex-none w-full sm:w-auto">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search items..."
            class="input w-full sm:w-80 pl-9 pr-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-1 focus:ring-accent focus:border-transparent placeholder:text-muted" />
          <button v-if="searchQuery" @click="searchQuery = ''" type="button" class="absolute right-3 top-1/2 -translate-y-1/2 text-muted hover:text-secondary transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>

        <div class="flex items-center gap-2">
          <!-- Sort -->
          <div class="relative">
            <button @click="toggleSortDropdown" class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors" title="Sort">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" /></svg>
            </button>
            <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95">
              <div v-if="showSortMenu" class="absolute right-0 top-10 w-48 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
                <div class="px-3 py-2 border-b border-divider"><p class="text-xs font-semibold text-muted uppercase">Sort By</p></div>
                <button v-for="option in sortOptions" :key="option.key" @click="applySort(option.key)" class="w-full flex items-center justify-between px-3 py-2 text-sm text-secondary hover:bg-surface-alt transition-colors">
                  {{ option.label }}
                  <svg v-if="currentSort === option.key" class="w-4 h-4 text-accent" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
                </button>
              </div>
            </Transition>
          </div>
          <div v-if="showSortMenu" class="fixed inset-0 z-40" @click="showSortMenu = false"></div>

          <!-- Copy -->
          <button @click="copyToClipboard" class="w-9 h-9 flex items-center justify-center border border-default rounded-lg text-secondary hover-surface transition-colors relative" title="Copy to clipboard">
            <svg v-if="!copied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" /></svg>
            <svg v-else class="w-4 h-4 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
          </button>

          <button @click="openCreateDialog" class="h-9 px-4 flex items-center gap-2 bg-accent text-accent-foreground rounded-lg text-sm font-semibold hover:opacity-90 transition-all duration-200 active:scale-95 whitespace-nowrap">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" /></svg>
            Add
          </button>
        </div>
      </div>
    </div>

    <!-- Totals Bar -->
    <div class="flex items-center gap-4 mb-3 px-1 flex-wrap">
      <div class="flex items-center gap-1.5"><span class="text-xs text-muted">Amount:</span><span class="text-sm font-semibold text-primary">৳{{ totalAmount.toFixed(2) }}</span></div>
      <div class="flex items-center gap-1.5"><span class="text-xs text-muted">Deposit:</span><span class="text-sm font-semibold text-info-text">৳{{ totalDeposit.toFixed(2) }}</span></div>
      <div class="flex items-center gap-1.5"><span class="text-xs text-muted">Paid:</span><span class="text-sm font-semibold text-success-text">৳{{ totalCustomerPaid.toFixed(2) }}</span></div>
      <div class="flex items-center gap-1.5"><span class="text-xs text-muted">Due:</span><span class="text-sm font-semibold" :class="totalDue > 0 ? 'text-warning-text' : 'text-success-text'">৳{{ totalDue.toFixed(2) }}</span></div>
    </div>

    <!-- Table Header -->
    <div class="flex-1 min-h-0 overflow-auto">
      <div class="min-w-3xl">
        <div class="grid grid-cols-12 items-center py-3 px-3 border-b border-divider text-xs font-semibold text-muted uppercase tracking-wider shrink-0 bg-surface-alt rounded-t-lg">
          <div class="col-span-3">Product / Customer</div>
          <div class="col-span-2">Category</div>
          <div class="col-span-2">Storage / Account</div>
          <div class="col-span-1">Qty</div>
          <div class="col-span-2">Amount / Due</div>
          <div class="col-span-2 text-right">Actions</div>
        </div>

        <div v-if="loading" class="flex items-center justify-center py-12">
          <svg class="animate-spin w-10 h-10 text-accent" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" /><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" /></svg>
        </div>

        <div v-else-if="filteredItems.length === 0" class="flex items-center justify-center py-12">
          <p class="text-sm text-muted">No items found</p>
        </div>

        <div v-else>
          <ItemRow v-for="item in filteredItems" :key="item.id" :item="item" @item-updated="handleItemUpdated" @view-item="openDetailDialog" />
        </div>
      </div>
    </div>

    <div v-if="!loading && filteredItems.length > 0" class="flex items-center justify-between py-3 px-1 border-t border-divider shrink-0">
      <p class="text-xs text-muted">Showing {{ filteredItems.length }} of {{ itemsStore.items.length }}</p>
    </div>

    <!-- Create Dialog -->
    <BaseDialog v-model="showCreateDialog" size="lg">
      <CreateItemForm @item-created="handleItemCreated" />
    </BaseDialog>

    <!-- Detail Dialog -->
    <BaseDialog v-model="showDetailDialog" size="lg">
      <div v-if="selectedItem" class="space-y-5">
        <div class="flex items-center gap-4">
          <div v-if="selectedItem.image_url" class="shrink-0 w-20 h-20 rounded-lg overflow-hidden border border-default">
            <img :src="getImageUrl(selectedItem.image_url) || undefined" :alt="selectedItem.product_name" class="w-full h-full object-cover" />
          </div>
          <div>
            <h2 class="text-xl font-bold text-primary">{{ selectedItem.product_name }}</h2>
            <p class="text-sm text-muted">{{ selectedItem.customer_phone }}</p>
            <p v-if="selectedItem.customer_email" class="text-xs text-muted">{{ selectedItem.customer_email }}</p>
            <div v-if="selectedItem.lot_number" class="mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-surface-alt text-secondary">Lot: {{ selectedItem.lot_number }}</span>
            </div>
          </div>
        </div>

        <!-- Storage Info -->
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Godown Name</p><p class="text-sm text-primary">{{ selectedItem.storage_name || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Account Name</p><p class="text-sm text-primary">{{ selectedItem.account_name || '—' }}</p></div>
        </div>

        <!-- Category -->
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Category</p><p class="text-sm text-primary">{{ selectedItem.category || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Subcategory</p><p class="text-sm text-primary">{{ selectedItem.subcategory || '—' }}</p></div>
        </div>

        <!-- Quantity -->
        <div class="space-y-1"><p class="text-xs text-muted uppercase">Quantity</p><p class="text-sm text-primary">{{ selectedItem.quantity }} {{ selectedItem.quantity_unit }}</p></div>

        <!-- Financial Summary -->
        <div class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Financial</h3>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Amount</p><p class="text-sm font-semibold text-primary">৳{{ selectedItem.amount.toFixed(2) }}</p></div>
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Deposit</p><p class="text-sm font-semibold text-info-text">৳{{ selectedItem.deposit.toFixed(2) }}</p></div>
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Paid</p><p class="text-sm font-semibold text-success-text">৳{{ (selectedItem.customer_paid || 0).toFixed(2) }}</p></div>
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Due</p><p class="text-sm font-semibold" :class="getItemDue(selectedItem) > 0 ? 'text-warning-text' : 'text-success-text'">৳{{ getItemDue(selectedItem).toFixed(2) }}</p></div>
          </div>
        </div>

        <div v-if="selectedItem.notes" class="space-y-1"><p class="text-xs text-muted uppercase">Notes</p><p class="text-sm text-primary whitespace-pre-wrap">{{ selectedItem.notes }}</p></div>

        <div class="flex justify-end pt-2 border-t border-divider">
          <button @click="showDetailDialog = false" class="button px-4 py-2 text-sm rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useItemsStore } from '@/stores/items'
import type { Item } from '@/types/item'
import ItemRow from './ItemRow.vue'
import CreateItemForm from './CreateItemForm.vue'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { getImageUrl } from '@/utils/image'
import { push } from 'notivue'

const itemsStore = useItemsStore()

const loading = ref(true)
const searchQuery = ref('')
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedItem = ref<Item | null>(null)
const showSortMenu = ref(false)
const currentSort = ref('newest')
const copied = ref(false)

const sortOptions = [
  { key: 'newest', label: 'Newest First' },
  { key: 'oldest', label: 'Oldest First' },
  { key: 'name_asc', label: 'Name A-Z' },
  { key: 'name_desc', label: 'Name Z-A' },
]

// Totals
const totalAmount = computed(() => itemsStore.items.reduce((sum, i) => sum + (i.amount || 0), 0))
const totalDeposit = computed(() => itemsStore.items.reduce((sum, i) => sum + (i.deposit || 0), 0))
const totalCustomerPaid = computed(() => itemsStore.items.reduce((sum, i) => sum + (i.customer_paid || 0), 0))
const totalDue = computed(() => itemsStore.items.reduce((sum, i) => sum + getItemDue(i), 0))

const getItemDue = (item: Item): number => (item.amount || 0) - ((item.customer_paid || 0) + (item.deposit || 0))

onMounted(async () => { await fetchItems() })

const fetchItems = async () => { loading.value = true; try { await itemsStore.fetchItems() } finally { loading.value = false } }
const handleItemUpdated = async () => { await itemsStore.fetchItems() }
const handleItemCreated = async () => { showCreateDialog.value = false; await itemsStore.fetchItems() }
const openCreateDialog = () => { showCreateDialog.value = true }
const openDetailDialog = (item: Item) => { selectedItem.value = item; showDetailDialog.value = true }

const toggleSortDropdown = () => { showSortMenu.value = !showSortMenu.value }
const applySort = (key: string) => { currentSort.value = key; showSortMenu.value = false }

const copyToClipboard = async () => {
  const items = filteredItems.value
  if (items.length === 0) return
  const headers = ['ID', 'Product Name', 'Phone', 'Lot', 'Storage', 'Account', 'Category', 'Subcategory', 'Qty', 'Unit', 'Amount', 'Deposit', 'Paid', 'Due', 'Notes']
  const rows = items.map(i => [i.id, i.product_name, i.customer_phone, i.lot_number || '', i.storage_name || '', i.account_name || '', i.category || '', i.subcategory || '', i.quantity, i.quantity_unit, i.amount || '', i.deposit || '', i.customer_paid || 0, getItemDue(i).toFixed(2), i.notes || ''])
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

const filteredItems = computed(() => {
  let result = [...itemsStore.items]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(i =>
      // Product & Storage
      (i.product_name && i.product_name.toLowerCase().includes(q)) ||
      (i.storage_name && i.storage_name.toLowerCase().includes(q)) ||
      (i.account_name && i.account_name.toLowerCase().includes(q)) ||
      (i.lot_number && i.lot_number.toLowerCase().includes(q)) ||
      // Customer
      (i.customer_phone && i.customer_phone.includes(q)) ||
      (i.customer_email && i.customer_email.toLowerCase().includes(q)) ||
      // Category
      (i.category && i.category.toLowerCase().includes(q)) ||
      (i.subcategory && i.subcategory.toLowerCase().includes(q)) ||
      // Details
      String(i.quantity).includes(q) ||
      (i.quantity_unit && i.quantity_unit.toLowerCase().includes(q)) ||
      // Revenue
      (i.amount && String(i.amount).includes(q)) ||
      (i.deposit && String(i.deposit).includes(q)) ||
      (i.customer_paid && String(i.customer_paid).includes(q)) ||
      // Notes
      (i.notes && i.notes.toLowerCase().includes(q)) ||
      // ID
      String(i.id).includes(q)
    )
  }
  if (currentSort.value === 'newest') result.sort((a, b) => b.id - a.id)
  else if (currentSort.value === 'oldest') result.sort((a, b) => a.id - b.id)
  else if (currentSort.value === 'name_asc') result.sort((a, b) => a.product_name.localeCompare(b.product_name))
  else if (currentSort.value === 'name_desc') result.sort((a, b) => b.product_name.localeCompare(a.product_name))
  return result
})
</script>