<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add New Item</h2>
      <p class="text-sm text-secondary">Register a new item in inventory</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <!-- Product & Storage Information -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Product & Storage <span class="text-warning-text">*</span>
          </h3>
          <div class="space-y-4">
            <AutocompleteField v-model="product_name" label="Product Name" placeholder="Product name / description" required :suggestions="productSuggestions" />
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <AutocompleteField v-model="storage_name" label="Godown Name" placeholder="Warehouse/facility" :suggestions="storageSuggestions" />
              <AutocompleteField v-model="account_name" label="Account Name" placeholder="Account reference" :suggestions="accountSuggestions" />
              <AutocompleteField v-model="lot_number" label="Lot Number" placeholder="Lot/batch number" :suggestions="lotSuggestions" />
            </div>
          </div>
        </div>

        <!-- Customer Information -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Customer Information <span class="normal-case text-xs font-normal">(Optional)</span></h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <AutocompleteField v-model="customer_phone" label="Phone" placeholder="+880 1XXX-XXXXXX" :suggestions="phoneSuggestions" />
            <AutocompleteField v-model="customer_email" label="Email" placeholder="customer@example.com" :suggestions="emailSuggestions" />
          </div>
        </div>

        <!-- Item Details -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Item Details</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <AutocompleteField v-model="category" label="Category" placeholder="e.g. electronics, furniture..." :suggestions="categorySuggestions" />
            <AutocompleteField v-model="subcategory" label="Subcategory" placeholder="e.g. laptop, sofa..." :suggestions="subcategorySuggestions" />
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Quantity</label>
              <div class="flex items-center gap-2">
                <button type="button" @click="quantity > 1 ? quantity-- : null" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" /></svg></button>
                <input v-model.number="quantity" type="number" min="1" placeholder="1" class="input flex-1 px-3 py-2 rounded-lg text-center placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
                <button type="button" @click="quantity++" class="w-9 h-9 flex items-center justify-center rounded-lg border border-default text-secondary hover:bg-surface-alt transition-colors shrink-0"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v14M5 12h14" /></svg></button>
              </div>
            </div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Unit</label><select v-model="quantity_unit" class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent"><option value="bag">Bag</option><option value="bottle">Bottle</option><option value="box">Box</option><option value="can">Can</option><option value="carton">Carton</option><option value="cup">Cup</option><option value="dozen">Dozen</option><option value="gallon">Gallon</option><option value="pack">Pack</option><option value="pair">Pair</option><option value="pcs">Pieces</option><option value="roll">Roll</option><option value="set">Set</option><option value="sheet">Sheet</option><option value="unit">Unit</option></select></div>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Weight</label>
              <input v-model.number="weight" type="number" step="0.001" min="0" placeholder="0.000"
                class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
            <div class="space-y-1.5">
              <label class="text-sm font-medium text-primary">Weight Unit</label>
              <select v-model="weight_unit" class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
                <option value="kg">kg</option><option value="g">g</option><option value="mg">mg</option><option value="oz">oz</option><option value="lb">lb</option><option value="ton">ton</option>
              </select>
            </div>
          </div>
        </div>

        <!-- Revenue Details -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Revenue Details <span class="normal-case text-xs font-normal">(Optional)</span></h3>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Amount</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="amount" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Security Deposit</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="deposit" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
            <div class="space-y-1.5"><label class="text-sm font-medium text-primary">Customer Paid</label><div class="relative"><span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span><input v-model.number="customer_paid" type="number" step="0.01" min="0" placeholder="0.00" class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" /></div></div>
          </div>
        </div>

        <!-- Image & Notes -->
        <div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Image <span class="normal-case text-xs font-normal">(Optional)</span></h3>
              <div class="flex items-start gap-4">
                <div class="shrink-0">
                  <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default"><img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" /><button type="button" @click="removeImage" class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full flex items-center justify-center text-xs hover:opacity-90 transition-opacity">×</button></div>
                  <div v-else class="w-24 h-24 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt"><svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg></div>
                </div>
                <div class="flex-1 space-y-3"><label class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer hover-surface transition-all duration-200 inline-flex items-center gap-2" :class="{ 'opacity-50 cursor-not-allowed': uploading }"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>{{ uploading ? 'Uploading...' : 'Choose File' }}<input type="file" accept="image/*" class="hidden" @change="handleFileSelect" :disabled="uploading" /></label></div>
              </div>
            </div>
            <div>
              <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">Notes <span class="normal-case text-xs font-normal">(Optional)</span></h3>
              <textarea v-model="notes" rows="4" placeholder="Add notes..." class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm" class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Clear</button>
          <button type="submit" :disabled="isSubmitDisabled" class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2"><svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>Creating...</span>
            <span v-else>Create Item</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import type { CreateItemPayload } from '@/types/item'
import { useItemsStore } from '@/stores/items'
import { uploadImage } from '@/utils/image'
import { push } from 'notivue'
import AutocompleteField from '@/components/ui/AutocompleteField.vue'

const props = defineProps<{
  prefill?: {
    product_name?: string; storage_name?: string | null; account_name?: string | null
    lot_number?: string | null; customer_phone?: string | null; customer_email?: string | null
    quantity_unit?: string
  }
}>()

const emit = defineEmits<{ 'item-created': [] }>()
const itemsStore = useItemsStore()
const submitting = ref(false)

const product_name = ref('')
const storage_name = ref('')
const account_name = ref('')
const lot_number = ref('')
const customer_phone = ref('')
const customer_email = ref('')
const category = ref('')
const subcategory = ref('')
const quantity = ref(1)
const quantity_unit = ref('pcs')
const weight = ref<number | null>(null)
const weight_unit = ref('kg')
const amount = ref<number | null>(null)
const deposit = ref<number | null>(null)
const customer_paid = ref<number | null>(null)
const notes = ref('')
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)

const productSuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.product_name) s.add(i.product_name.toLowerCase()) }); return Array.from(s).sort() })
const storageSuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.storage_name) s.add(i.storage_name.toLowerCase()) }); return Array.from(s).sort() })
const accountSuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.account_name) s.add(i.account_name.toLowerCase()) }); return Array.from(s).sort() })
const lotSuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.lot_number) s.add(i.lot_number.toLowerCase()) }); return Array.from(s).sort() })
const phoneSuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.customer_phone) s.add(i.customer_phone) }); return Array.from(s).sort() })
const emailSuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.customer_email) s.add(i.customer_email.toLowerCase()) }); return Array.from(s).sort() })
const categorySuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.category) s.add(i.category.toLowerCase()) }); return Array.from(s).sort() })
const subcategorySuggestions = computed(() => { const s = new Set<string>(); itemsStore.items.forEach(i => { if (i.subcategory) s.add(i.subcategory.toLowerCase()) }); return Array.from(s).sort() })

const isSubmitDisabled = computed(() => submitting.value || !product_name.value)

onMounted(() => {
  if (props.prefill) {
    if (props.prefill.product_name) product_name.value = props.prefill.product_name
    if (props.prefill.storage_name) storage_name.value = props.prefill.storage_name || ''
    if (props.prefill.account_name) account_name.value = props.prefill.account_name || ''
    if (props.prefill.lot_number) lot_number.value = props.prefill.lot_number || ''
    if (props.prefill.customer_phone) customer_phone.value = props.prefill.customer_phone || ''
    if (props.prefill.customer_email) customer_email.value = props.prefill.customer_email || ''
    if (props.prefill.quantity_unit) quantity_unit.value = props.prefill.quantity_unit
  }
})

const handleFileSelect = (event: Event) => { const input = event.target as HTMLInputElement; const file = input.files?.[0]; if (!file) return; if (!file.type.startsWith('image/')) { push.error('Please select an image file'); return }; if (file.size > 5 * 1024 * 1024) { push.error('Image size should be less than 5MB'); return }; imageFile.value = file; const reader = new FileReader(); reader.onload = (e) => { imagePreview.value = e.target?.result as string }; reader.readAsDataURL(file) }
const removeImage = () => { imageFile.value = null; imagePreview.value = null }

const resetForm = () => {
  product_name.value = ''; storage_name.value = ''; account_name.value = ''; lot_number.value = ''
  customer_phone.value = ''; customer_email.value = ''
  category.value = ''; subcategory.value = ''
  quantity.value = 1; quantity_unit.value = 'pcs'
  weight.value = null; weight_unit.value = 'kg'
  amount.value = null; deposit.value = null; customer_paid.value = null
  notes.value = ''; removeImage()
  if (props.prefill) {
    if (props.prefill.product_name) product_name.value = props.prefill.product_name
    if (props.prefill.storage_name) storage_name.value = props.prefill.storage_name || ''
    if (props.prefill.account_name) account_name.value = props.prefill.account_name || ''
    if (props.prefill.customer_phone) customer_phone.value = props.prefill.customer_phone || ''
    if (props.prefill.customer_email) customer_email.value = props.prefill.customer_email || ''
    if (props.prefill.quantity_unit) quantity_unit.value = props.prefill.quantity_unit
  }
}

const submit = async () => {
  if (!product_name.value) return
  submitting.value = true
  try {
    const payload: CreateItemPayload = {
      product_name: product_name.value.trim(),
      storage_name: storage_name.value.trim() || null,
      account_name: account_name.value.trim() || null,
      lot_number: lot_number.value.trim() || null,
      customer_phone: customer_phone.value.trim() || null,
      customer_email: customer_email.value?.trim() || null,
      category: category.value.trim().toLowerCase() || null,
      subcategory: subcategory.value.trim().toLowerCase() || null,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      quantity_unit: quantity_unit.value as any, quantity: quantity.value,
      weight: weight.value || null,
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      weight_unit: weight_unit.value as any,
      amount: amount.value || 0, deposit: deposit.value || 0, customer_paid: customer_paid.value || 0,
      notes: notes.value || null,
    }
    const newItem = await itemsStore.createItem(payload)
    if (!newItem) { submitting.value = false; return }
    if (imageFile.value) { uploading.value = true; await uploadImage('items', newItem.id, imageFile.value); uploading.value = false }
    push.success('Item created successfully!')
    resetForm()
    emit('item-created')
  } catch { push.error('Failed to create item') }
  finally { submitting.value = false }
}
</script>