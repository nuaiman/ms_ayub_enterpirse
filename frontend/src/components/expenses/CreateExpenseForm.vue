<template>
  <div class="max-w-2xl mx-auto">
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-primary mb-2">Add Expense</h2>
      <p class="text-sm text-secondary">Record a new expense</p>
    </div>

    <div class="card rounded-2xl shadow-sm p-6">
      <form @submit.prevent="submit" class="space-y-6">

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Title <span class="text-warning-text">*</span></label>
            <input v-model="title" type="text" placeholder="Expense title" required
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Category <span class="text-warning-text">*</span></label>
            <select v-model="category" required
              class="input w-full px-3 py-2 rounded-lg cursor-pointer focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent">
              <option value="rent" class="bg-surface text-primary">Rent</option>
              <option value="salary" class="bg-surface text-primary">Salary</option>
              <option value="bill" class="bg-surface text-primary">Bill</option>
              <option value="other" class="bg-surface text-primary">Other</option>
            </select>
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Amount <span class="text-warning-text">*</span></label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">$</span>
              <input v-model.number="amount" type="number" step="0.01" min="0" placeholder="0.00" required
                class="input w-full pl-7 pr-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>

          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Date <span class="text-warning-text">*</span></label>
            <input v-model="expense_date" type="date" required
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>

          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Notes <span
                class="text-xs font-normal text-muted">(Optional)</span></label>
            <textarea v-model="notes" rows="2" placeholder="Add notes..."
              class="input w-full px-3 py-2 rounded-lg placeholder:text-muted focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
          </div>
        </div>

        <!-- Image Upload -->
        <div class="pb-6 border-b border-divider">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-4">
            Receipt / Image <span class="normal-case text-xs font-normal">(Optional)</span>
          </h3>
          <div class="flex items-start gap-4">
            <div class="shrink-0">
              <div v-if="imagePreview" class="relative w-24 h-24 rounded-lg overflow-hidden border border-default">
                <img :src="imagePreview" alt="Preview" class="w-full h-full object-cover" />
                <button type="button" @click="removeImage"
                  class="absolute top-1 right-1 w-5 h-5 bg-warning-text text-white rounded-full flex items-center justify-center text-xs hover:opacity-90">×</button>
              </div>
              <div v-else
                class="w-24 h-24 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt">
                <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
            </div>
            <div class="flex-1 space-y-3">
              <label class="text-sm font-medium text-primary block">Upload Receipt</label>
              <div class="flex gap-2">
                <label
                  class="button px-4 py-2 text-sm font-medium rounded-lg cursor-pointer hover-surface transition-all duration-200 inline-flex items-center gap-2"
                  :class="{ 'opacity-50 cursor-not-allowed': uploading }">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                  </svg>
                  {{ uploading ? 'Uploading...' : 'Choose File' }}
                  <input type="file" accept="image/*" class="hidden" @change="handleFileSelect" :disabled="uploading" />
                </label>
              </div>
              <div v-if="uploading" class="w-full bg-surface-alt rounded-full h-1.5 overflow-hidden">
                <div class="bg-accent h-full rounded-full transition-all duration-300"
                  :style="{ width: uploadProgress + '%' }"></div>
              </div>
              <p v-if="uploadError" class="text-xs text-warning-text">{{ uploadError }}</p>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-4">
          <button type="button" @click="resetForm"
            class="button px-4 py-2 text-sm font-medium rounded-lg hover-surface transition-all duration-200">Clear</button>
          <button type="submit" :disabled="submitting || !title || !amount || amount <= 0"
            class="px-6 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg shadow-sm hover:shadow-md transition-all duration-200 active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="submitting" class="inline-flex items-center gap-2">
              <svg class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
              </svg>
              Creating...
            </span>
            <span v-else>Add Expense</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { ExpenseCategory } from '@/types/expense'
import { useExpensesStore } from '@/stores/expenses'
import { uploadImage } from '@/utils/image'
import { push } from 'notivue'

const emit = defineEmits<{ 'expense-created': [] }>()

const expensesStore = useExpensesStore()

const submitting = ref(false)
const title = ref('')
const category = ref<ExpenseCategory>('other')
const amount = ref<number | null>(null)
const expense_date = ref(new Date().toISOString().split('T')[0] as string)
const notes = ref('')

// Image
const imageFile = ref<File | null>(null)
const imagePreview = ref<string | null>(null)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadError = ref<string | null>(null)

const handleFileSelect = (event: Event) => {
  const input = event.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Please select an image file'
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = 'Image size should be less than 5MB'
    return
  }
  uploadError.value = null
  imageFile.value = file
  const reader = new FileReader()
  reader.onload = (e) => {
    imagePreview.value = e.target?.result as string
  }
  reader.readAsDataURL(file)
}

const removeImage = () => {
  imageFile.value = null
  imagePreview.value = null
  uploadError.value = null
  uploadProgress.value = 0
  const fileInput = document.querySelector('input[type="file"]') as HTMLInputElement
  if (fileInput) fileInput.value = ''
}

const resetForm = () => {
  title.value = ''
  category.value = 'other'
  amount.value = null
  expense_date.value = new Date().toISOString().split('T')[0] as string
  notes.value = ''
  removeImage()
}

const submit = async () => {
  if (!title.value || !amount.value || amount.value <= 0 || !expense_date.value) return

  submitting.value = true

  try {
    // Step 1: Create expense WITHOUT image
    const newExpense = await expensesStore.createExpense({
      title: title.value,
      category: category.value,
      amount: Number(amount.value),
      expense_date: expense_date.value,
      notes: notes.value || null,
    })

    if (!newExpense) {
      submitting.value = false
      return
    }

    // Step 2: Upload image for the NEWLY CREATED expense
    if (imageFile.value) {
      uploading.value = true
      uploadProgress.value = 0

      const imageUrl = await uploadImage('expenses', newExpense.id, imageFile.value, (progress: number) => {
        uploadProgress.value = progress
      })

      uploading.value = false

      if (!imageUrl) {
        push.warning('Expense created but image upload failed')
      }
    }

    push.success('Expense added successfully!')
    resetForm()
    emit('expense-created')
  } catch (error) {
    console.error('Error creating expense:', error)
    push.error('Failed to create expense')
  } finally {
    submitting.value = false
  }
}
</script>
