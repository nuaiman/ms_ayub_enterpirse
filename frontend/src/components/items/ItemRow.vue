<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)] cursor-pointer"
    @click="emit('view-item', item)">

    <!-- Product / Customer -->
    <div class="col-span-3">
      <div class="flex items-center gap-3">
        <div class="shrink-0">
          <div v-if="item.image_url" class="w-9 h-9 rounded-lg overflow-hidden border border-default">
            <img :src="getImageUrl(item.image_url) || undefined" :alt="item.product_name" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-9 h-9 rounded-lg bg-surface-alt border border-default flex items-center justify-center">
            <svg class="w-4 h-4 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>
          </div>
        </div>
        <div class="min-w-0">
          <div class="font-medium text-primary truncate">{{ item.product_name }}</div>
          <div class="flex items-center gap-2 text-xs text-muted">
            <span v-if="item.customer_phone">{{ item.customer_phone }}</span>
            <span v-if="item.lot_number" class="text-[10px] bg-surface-alt px-1.5 py-0.5 rounded">Lot: {{ item.lot_number }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Category -->
    <div class="col-span-2">
      <div v-if="item.category">
        <span class="text-xs text-secondary">{{ item.category }}</span>
        <div v-if="item.subcategory" class="text-[10px] text-muted mt-0.5">{{ item.subcategory }}</div>
      </div>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Storage Info -->
    <div class="col-span-2">
      <div v-if="item.storage_name || item.account_name" class="space-y-0.5">
        <div v-if="item.storage_name" class="text-xs text-secondary">{{ item.storage_name }}</div>
        <div v-if="item.account_name" class="text-[10px] text-muted">{{ item.account_name }}</div>
      </div>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Quantity -->
    <div class="col-span-1">
      <span class="text-sm font-medium text-primary">{{ item.quantity }}</span>
      <span class="text-xs text-muted ml-1">{{ item.quantity_unit }}</span>
    </div>

    <!-- Amount / Due -->
    <div class="col-span-2">
      <div v-if="item.amount" class="space-y-0.5">
        <span class="text-sm font-semibold text-warning-text">৳{{ item.amount.toFixed(2) }}</span>
        <div class="flex items-center gap-2 text-[10px]">
          <span v-if="item.deposit" class="text-info-text">Dep: ৳{{ item.deposit.toFixed(0) }}</span>
          <span v-if="item.customer_paid" class="text-success-text">Paid: ৳{{ item.customer_paid.toFixed(0) }}</span>
          <span :class="getItemDue(item) > 0 ? 'text-warning-text' : 'text-success-text'" class="font-medium">Due: ৳{{ getItemDue(item).toFixed(0) }}</span>
        </div>
      </div>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Actions -->
    <div class="col-span-2 flex justify-end pr-0.5 relative" @click.stop>
      <button @click="open = !open" class="w-7 h-7 flex items-center justify-center border border-default rounded-md hover:bg-surface-alt transition-all duration-200 opacity-0 group-hover:opacity-100" :class="{ 'opacity-100': open }">
        <svg class="w-3.5 h-3.5 text-secondary" fill="currentColor" viewBox="0 0 24 24"><circle cx="12" cy="5" r="1.5" /><circle cx="12" cy="12" r="1.5" /><circle cx="12" cy="19" r="1.5" /></svg>
      </button>
      <Transition enter-active-class="transition ease-out duration-200" enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0" leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 scale-100 translate-y-0" leave-to-class="opacity-0 scale-95 translate-y-1">
        <div v-if="open" class="absolute right-0 top-9 w-44 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider"><p class="text-[11px] font-semibold text-muted uppercase">Actions</p></div>
          <button @click="openEditDialog" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors"><svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>Edit</button>
          <button @click="openNewLotDialog" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors"><svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>Add New Lot</button>
          <div class="border-t border-divider my-0.5"></div>
          <button @click="openDeleteConfirm" class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors"><svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>Delete</button>
        </div>
      </Transition>
      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- Edit Dialog -->
    <BaseDialog v-model="showEditDialog" size="lg" @click.stop>
      <div class="space-y-4">
        <h2 class="text-lg font-semibold text-primary">Edit Item</h2>
        <div class="flex items-center gap-4 pb-4 border-b border-divider">
          <div class="shrink-0 relative group/edit-img">
            <div v-if="editImagePreview || item.image_url" class="w-16 h-16 rounded-lg overflow-hidden border border-default"><img :src="editImagePreview || getImageUrl(item.image_url) || undefined" :alt="item.product_name" class="w-full h-full object-cover" /></div>
            <div v-else class="w-16 h-16 rounded-lg bg-surface-alt border border-default flex items-center justify-center"><svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg></div>
            <button @click="triggerEditImageInput" class="absolute inset-0 rounded-lg bg-black/50 flex items-center justify-center opacity-0 group-hover/edit-img:opacity-100 transition-opacity cursor-pointer"><svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg></button>
            <input ref="editImageInputRef" type="file" accept="image/*" class="hidden" @change="handleEditImageSelect" />
          </div>
          <div><p class="text-xs text-muted">Click on image to change</p><button v-if="item.image_url && !editImageFile" @click="removeExistingImage" type="button" class="text-xs text-warning-text hover:underline mt-1">Remove current image</button><button v-if="editImageFile" @click="clearEditImage" type="button" class="text-xs text-muted hover:underline mt-1">Clear selection</button></div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5 md:col-span-2"><label class="text-sm font-medium text-primary">Product Name</label><input v-model="editForm.product_name" type="text" class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Storage Name</label><input v-model="editForm.storage_name" type="text" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Account Name</label><input v-model="editForm.account_name" type="text" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Lot Number</label><input v-model="editForm.lot_number" type="text" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Phone</label><input v-model="editForm.customer_phone" type="tel" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Email</label><input v-model="editForm.customer_email" type="email" class="input w-full px-3 py-2 rounded-lg" /></div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Category</label><input v-model="editForm.category" type="text" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Subcategory</label><input v-model="editForm.subcategory" type="text" class="input w-full px-3 py-2 rounded-lg" /></div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Quantity</label><div class="flex items-center gap-2"><button type="button" @click="editForm.quantity > 1 ? editForm.quantity-- : null" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt shrink-0">−</button><input v-model.number="editForm.quantity" type="number" min="1" class="input flex-1 px-3 py-2 rounded-lg text-center" /><button type="button" @click="editForm.quantity++" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt shrink-0">+</button></div></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Unit</label><select v-model="editForm.quantity_unit" class="input w-full px-3 py-2 rounded-lg"><option value="bag">Bag</option><option value="bottle">Bottle</option><option value="box">Box</option><option value="can">Can</option><option value="carton">Carton</option><option value="cup">Cup</option><option value="dozen">Dozen</option><option value="gallon">Gallon</option><option value="pack">Pack</option><option value="pair">Pair</option><option value="pcs">Pieces</option><option value="roll">Roll</option><option value="set">Set</option><option value="sheet">Sheet</option><option value="unit">Unit</option></select></div>
        </div>
        <div class="grid grid-cols-3 gap-4">
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Amount</label><input v-model.number="editForm.amount" type="number" step="0.01" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Deposit</label><input v-model.number="editForm.deposit" type="number" step="0.01" class="input w-full px-3 py-2 rounded-lg" /></div>
          <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Customer Paid</label><input v-model.number="editForm.customer_paid" type="number" step="0.01" class="input w-full px-3 py-2 rounded-lg" /></div>
        </div>
        <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Notes</label><textarea v-model="editForm.notes" rows="2" class="input w-full px-3 py-2 rounded-lg resize-none"></textarea></div>
        <div class="flex justify-end gap-3 pt-4 border-t border-divider"><button @click="showEditDialog = false" class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button><button @click="handleUpdate" :disabled="saving" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">{{ saving ? 'Saving...' : 'Save' }}</button></div>
      </div>
    </BaseDialog>

    <!-- New Lot Dialog -->
    <BaseDialog v-model="showNewLotDialog" size="lg" @click.stop>
      <CreateItemForm :prefill="newLotPrefill" @item-created="handleNewLotCreated" />
    </BaseDialog>

    <!-- Detail Dialog -->
    <BaseDialog v-model="showDetailDialog" size="lg" @click.stop>
      <div v-if="item" class="space-y-5">
        <div class="flex items-center gap-4">
          <div v-if="item.image_url" class="shrink-0 w-20 h-20 rounded-lg overflow-hidden border border-default"><img :src="getImageUrl(item.image_url) || undefined" :alt="item.product_name" class="w-full h-full object-cover" /></div>
          <div>
            <h2 class="text-xl font-bold text-primary">{{ item.product_name }}</h2>
            <p v-if="item.customer_phone" class="text-sm text-muted">{{ item.customer_phone }}</p>
            <p v-if="item.customer_email" class="text-xs text-muted">{{ item.customer_email }}</p>
            <div v-if="item.lot_number" class="mt-1"><span class="inline-flex items-center px-2 py-0.5 rounded-md text-xs font-medium bg-surface-alt text-secondary">Lot: {{ item.lot_number }}</span></div>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Storage Name</p><p class="text-sm text-primary">{{ item.storage_name || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Account Name</p><p class="text-sm text-primary">{{ item.account_name || '—' }}</p></div>
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Category</p><p class="text-sm text-primary">{{ item.category || '—' }}</p></div>
          <div class="space-y-1"><p class="text-xs text-muted uppercase">Subcategory</p><p class="text-sm text-primary">{{ item.subcategory || '—' }}</p></div>
        </div>
        <div class="space-y-1"><p class="text-xs text-muted uppercase">Quantity</p><p class="text-sm text-primary">{{ item.quantity }} {{ item.quantity_unit }}</p></div>
        <div class="p-4 rounded-xl border border-default bg-surface-alt">
          <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-3">Financial</h3>
          <div class="grid grid-cols-2 gap-3">
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Amount</p><p class="text-sm font-semibold text-primary">৳{{ item.amount.toFixed(2) }}</p></div>
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Deposit</p><p class="text-sm font-semibold text-info-text">৳{{ item.deposit.toFixed(2) }}</p></div>
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Paid</p><p class="text-sm font-semibold text-success-text">৳{{ (item.customer_paid || 0).toFixed(2) }}</p></div>
            <div class="space-y-1"><p class="text-[10px] text-muted uppercase">Due</p><p class="text-sm font-semibold" :class="getItemDue(item) > 0 ? 'text-warning-text' : 'text-success-text'">৳{{ getItemDue(item).toFixed(2) }}</p></div>
          </div>
        </div>
        <div v-if="item.notes" class="space-y-1"><p class="text-xs text-muted uppercase">Notes</p><p class="text-sm text-primary whitespace-pre-wrap">{{ item.notes }}</p></div>
        <div class="flex justify-end gap-3 pt-2 border-t border-divider">
          <button @click="openNewLotFromDetail" class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200 inline-flex items-center gap-2"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>Add New Lot</button>
          <button @click="showDetailDialog = false" class="button px-4 py-2 text-sm rounded-lg hover-surface">Close</button>
        </div>
      </div>
    </BaseDialog>

    <!-- Delete Dialog -->
    <BaseDialog v-model="showDeleteDialog" size="sm" @click.stop>
      <div class="space-y-4"><h2 class="text-lg font-semibold text-primary">Delete Item</h2><p class="text-sm text-secondary">Delete <strong>{{ item.product_name }}</strong>?</p><div class="flex justify-end gap-2 pt-2"><button @click="showDeleteDialog = false" class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button><button @click="handleDelete" :disabled="deleting" class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg">{{ deleting ? 'Deleting...' : 'Delete' }}</button></div></div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import type { Item } from '@/types/item'
import { useItemsStore } from '@/stores/items'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import CreateItemForm from './CreateItemForm.vue'
import { push } from 'notivue'
import { getImageUrl, uploadImage, deleteImage } from '@/utils/image'

const props = defineProps<{ item: Item }>()
const emit = defineEmits<{ 'item-updated': []; 'view-item': [item: Item] }>()

const itemsStore = useItemsStore()
const open = ref(false)
const showEditDialog = ref(false)
const showNewLotDialog = ref(false)
const showDetailDialog = ref(false)
const showDeleteDialog = ref(false)
const saving = ref(false)
const deleting = ref(false)
const editImageInputRef = ref<HTMLInputElement | null>(null)
const editImageFile = ref<File | null>(null)
const editImagePreview = ref<string | null>(null)

const editForm = reactive({
  product_name: '', storage_name: '', account_name: '', lot_number: '',
  customer_phone: '', customer_email: null as string | null,
  category: null as string | null, subcategory: null as string | null,
  quantity: 1, quantity_unit: 'pcs',
  amount: null as number | null, deposit: null as number | null, customer_paid: null as number | null,
  notes: null as string | null,
})

const getItemDue = (item: Item): number => (item.amount || 0) - ((item.customer_paid || 0) + (item.deposit || 0))

const newLotPrefill = computed(() => ({
  product_name: props.item.product_name,
  storage_name: props.item.storage_name || undefined,
  account_name: props.item.account_name || undefined,
  customer_phone: props.item.customer_phone || undefined,
  customer_email: props.item.customer_email || undefined,
  quantity_unit: props.item.quantity_unit,
}))

const openEditDialog = () => {
  open.value = false
  editForm.product_name = props.item.product_name
  editForm.storage_name = props.item.storage_name || ''
  editForm.account_name = props.item.account_name || ''
  editForm.lot_number = props.item.lot_number || ''
  editForm.customer_phone = props.item.customer_phone || ''
  editForm.customer_email = props.item.customer_email || null
  editForm.category = props.item.category || null
  editForm.subcategory = props.item.subcategory || null
  editForm.quantity = props.item.quantity
  editForm.quantity_unit = props.item.quantity_unit
  editForm.amount = props.item.amount || null
  editForm.deposit = props.item.deposit || null
  editForm.customer_paid = props.item.customer_paid || null
  editForm.notes = props.item.notes || null
  editImageFile.value = null; editImagePreview.value = null
  showEditDialog.value = true
}

const openNewLotDialog = () => { open.value = false; showNewLotDialog.value = true }
const openNewLotFromDetail = () => { showDetailDialog.value = false; showNewLotDialog.value = true }
const handleNewLotCreated = () => { showNewLotDialog.value = false; emit('item-updated') }

const triggerEditImageInput = () => editImageInputRef.value?.click()
const handleEditImageSelect = (event: Event) => { const input = event.target as HTMLInputElement; const file = input.files?.[0]; if (!file) return; if (!file.type.startsWith('image/')) { push.error('Please select an image file'); return }; if (file.size > 5 * 1024 * 1024) { push.error('Image size should be less than 5MB'); return }; editImageFile.value = file; const reader = new FileReader(); reader.onload = (e) => { editImagePreview.value = e.target?.result as string }; reader.readAsDataURL(file) }
const clearEditImage = () => { editImageFile.value = null; editImagePreview.value = null; if (editImageInputRef.value) editImageInputRef.value.value = '' }
const removeExistingImage = async () => { if (!props.item.image_url) return; try { await deleteImage('items', props.item.id, props.item.image_url); editImagePreview.value = null; push.success('Image removed'); emit('item-updated') } catch { push.error('Failed to remove image') } }

const handleUpdate = async () => {
  saving.value = true
  try {
    if (editImageFile.value) { await uploadImage('items', props.item.id, editImageFile.value) }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const payload: any = {
      product_name: editForm.product_name, storage_name: editForm.storage_name || null,
      account_name: editForm.account_name || null, lot_number: editForm.lot_number || null,
      customer_phone: editForm.customer_phone || null, customer_email: editForm.customer_email,
      category: editForm.category, subcategory: editForm.subcategory,
      quantity: editForm.quantity, quantity_unit: editForm.quantity_unit,
      amount: editForm.amount, deposit: editForm.deposit, customer_paid: editForm.customer_paid,
      notes: editForm.notes,
    }
    await itemsStore.updateItem(props.item.id, payload)
    push.success('Item updated'); showEditDialog.value = false; emit('item-updated')
  } catch { push.error('Failed to update item') }
  finally { saving.value = false }
}

const openDeleteConfirm = () => { open.value = false; showDeleteDialog.value = true }
const handleDelete = async () => { deleting.value = true; try { await itemsStore.deleteItem(props.item.id); push.success('Item deleted'); showDeleteDialog.value = false; emit('item-updated') } catch { push.error('Failed to delete item') } finally { deleting.value = false } }
</script>