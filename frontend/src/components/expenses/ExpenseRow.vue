<template>
  <div
    class="grid grid-cols-12 items-center py-3.5 px-3 border-b border-default transition-all duration-200 relative text-sm group hover:border-b-accent/50 hover:shadow-[0_1px_0_0_rgba(232,33,39,0.3)]">

    <!-- Category -->
    <div class="col-span-3">
      <div class="flex items-center gap-2">
        <div class="min-w-0">
          <div class="font-medium text-primary truncate">{{ expense.category || '—' }}</div>
          <div class="text-[10px] text-muted">#{{ expense.id }}</div>
        </div>
      </div>
    </div>

    <!-- Amount -->
    <div class="col-span-2">
      <span class="text-sm font-semibold text-warning-text">৳{{ expense.amount.toFixed(2) }}</span>
    </div>

    <!-- Date -->
    <div class="col-span-2">
      <span class="text-xs text-secondary">{{ formatDate(expense.expense_date) }}</span>
    </div>

    <!-- Type (Salary/Regular) -->
    <div class="col-span-2">
      <span v-if="expense.is_salary"
        class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium bg-purple-50 text-purple-700 border border-purple-200">
        Salary
      </span>
      <span v-else class="text-xs text-muted">Regular</span>
    </div>

    <!-- Notes -->
    <div class="col-span-2">
      <span v-if="expense.notes" class="text-xs text-muted truncate block max-w-full">{{ expense.notes }}</span>
      <span v-else class="text-xs text-muted">—</span>
    </div>

    <!-- Actions -->
    <div class="col-span-1 flex justify-end pr-0.5 relative" @click.stop>
      <button @click="open = !open"
        class="w-7 h-7 flex items-center justify-center border border-default rounded-md hover:bg-surface-alt transition-all duration-200 opacity-0 group-hover:opacity-100"
        :class="{ 'opacity-100': open }">
        <svg class="w-3.5 h-3.5 text-secondary" fill="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="5" r="1.5" />
          <circle cx="12" cy="12" r="1.5" />
          <circle cx="12" cy="19" r="1.5" />
        </svg>
      </button>

      <Transition enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 scale-95 translate-y-1" enter-to-class="opacity-100 scale-100 translate-y-0"
        leave-active-class="transition ease-in duration-150" leave-from-class="opacity-100 scale-100 translate-y-0"
        leave-to-class="opacity-0 scale-95 translate-y-1">
        <div v-if="open"
          class="absolute right-0 top-9 w-44 bg-surface border border-default rounded-xl shadow-lg overflow-hidden z-50 py-1">
          <div class="px-3 py-1.5 border-b border-divider">
            <p class="text-[11px] font-semibold text-muted uppercase">Actions</p>
          </div>

          <button @click="openEditDialog"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-secondary hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>Edit
          </button>

          <div class="border-t border-divider my-0.5"></div>
          <button @click="openDeleteConfirm"
            class="w-full flex items-center gap-2.5 px-3 py-2 text-xs text-warning-text hover:bg-surface-alt transition-colors">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>Delete
          </button>
        </div>
      </Transition>
      <div v-if="open" class="fixed inset-0 z-40" @click="open = false"></div>
    </div>

    <!-- Edit Dialog -->
    <BaseDialog v-model="showEditDialog" size="lg" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-success-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-success-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Edit Expense</h2>
            <p class="text-xs text-muted">#{{ expense.id }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Category</label>
            <input v-model="editForm.category" type="text"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Amount</label>
            <div class="relative">
              <span class="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-muted">৳</span>
              <input v-model.number="editForm.amount" type="number" step="0.01" min="0"
                class="input w-full pl-7 pr-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
            </div>
          </div>
          <div class="space-y-1.5">
            <label class="text-sm font-medium text-primary">Date</label>
            <input v-model="editForm.expense_date" type="date"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent" />
          </div>
          <div class="space-y-1.5 md:col-span-2">
            <label class="text-sm font-medium text-primary">Notes</label>
            <textarea v-model="editForm.notes" rows="2"
              class="input w-full px-3 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-accent focus:border-transparent resize-none"></textarea>
          </div>
        </div>

        <!-- Image Upload in Edit -->
        <div class="border-t border-divider pt-4">
          <h3 class="text-sm font-semibold text-muted uppercase tracking-wider mb-3">Receipt</h3>
          <div class="flex items-start gap-4">
            <div class="shrink-0 relative group/edit-img">
              <div v-if="editImagePreview || expense.image_url"
                class="w-20 h-20 rounded-lg overflow-hidden border border-default">
                <img :src="editImagePreview || getImageUrl(expense.image_url) || undefined" alt="Receipt"
                  class="w-full h-full object-cover" />
              </div>
              <div v-else
                class="w-20 h-20 rounded-lg border-2 border-dashed border-default flex items-center justify-center bg-surface-alt">
                <svg class="w-6 h-6 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <button @click="triggerEditImageInput"
                class="absolute inset-0 rounded-lg bg-black/50 flex items-center justify-center opacity-0 group-hover/edit-img:opacity-100 transition-opacity cursor-pointer">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </button>
              <input ref="editImageInputRef" type="file" accept="image/*" class="hidden"
                @change="handleEditImageSelect" />
            </div>
            <div class="flex-1">
              <p class="text-xs text-muted">Click on image to change</p>
              <button v-if="expense.image_url && !editImageFile" @click="removeExistingImage" type="button"
                class="text-xs text-warning-text hover:underline mt-1">Remove current image</button>
            </div>
          </div>
        </div>

        <div class="flex gap-2 justify-end pt-2">
          <button @click="showEditDialog = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="submitEdit"
            class="px-4 py-2 text-sm font-semibold bg-accent text-accent-foreground rounded-lg">Save</button>
        </div>
      </div>
    </BaseDialog>

    <!-- Delete Confirm -->
    <BaseDialog v-model="showDeleteConfirm" @click.stop>
      <div class="space-y-4">
        <div class="flex items-center gap-3 mb-2">
          <div class="w-10 h-10 rounded-full bg-warning-bg flex items-center justify-center">
            <svg class="w-5 h-5 text-warning-text" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.34 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <div>
            <h2 class="text-lg font-semibold text-primary">Delete Expense</h2>
            <p class="text-xs text-muted">৳{{ expense.amount.toFixed(2) }}</p>
          </div>
        </div>
        <p class="text-sm text-secondary">Delete this expense?</p>
        <div class="flex gap-2 justify-end pt-2">
          <button @click="showDeleteConfirm = false"
            class="button px-4 py-2 text-sm rounded-lg hover-surface">Cancel</button>
          <button @click="deleteExpense"
            class="px-4 py-2 text-sm font-semibold bg-warning-text text-white rounded-lg">Delete</button>
        </div>
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { Expense } from '@/types/expense'
import { useExpensesStore } from '@/stores/expenses'
import { uploadImage, deleteImage, getImageUrl } from '@/utils/image'
import BaseDialog from '@/components/ui/BaseDialog.vue'
import { push } from 'notivue'

const props = defineProps<{ expense: Expense }>()
const emit = defineEmits<{ 'expense-updated': [] }>()

const expensesStore = useExpensesStore()
const open = ref(false)
const showEditDialog = ref(false)
const showDeleteConfirm = ref(false)

const editImageInputRef = ref<HTMLInputElement | null>(null)
const editImageFile = ref<File | null>(null)
const editImagePreview = ref<string | null>(null)

const editForm = ref({
  category: props.expense.category || '',
  amount: props.expense.amount,
  expense_date: props.expense.expense_date,
  notes: props.expense.notes || '',
})

const openEditDialog = () => {
  editForm.value = {
    category: props.expense.category || '',
    amount: props.expense.amount,
    expense_date: props.expense.expense_date,
    notes: props.expense.notes || '',
  }
  editImageFile.value = null
  editImagePreview.value = null
  showEditDialog.value = true
  open.value = false
}

const openDeleteConfirm = () => { showDeleteConfirm.value = true; open.value = false }

const triggerEditImageInput = () => editImageInputRef.value?.click()

const handleEditImageSelect = (event: Event) => {
  const input = event.target as HTMLInputElement; const file = input.files?.[0]
  if (!file) return
  if (!file.type.startsWith('image/')) { push.error('Please select an image file'); return }
  if (file.size > 5 * 1024 * 1024) { push.error('Image size should be less than 5MB'); return }
  editImageFile.value = file
  const reader = new FileReader(); reader.onload = (e) => { editImagePreview.value = e.target?.result as string }; reader.readAsDataURL(file)
}

const removeExistingImage = async () => {
  if (!props.expense.image_url) return
  try { await deleteImage('expenses', props.expense.id, props.expense.image_url); editImagePreview.value = null; emit('expense-updated') }
  catch { push.error('Failed to remove image') }
}

const submitEdit = async () => {
  try {
    if (editImageFile.value) { await uploadImage('expenses', props.expense.id, editImageFile.value) }
    const success = await expensesStore.updateExpense(props.expense.id, {
      category: editForm.value.category || null,
      amount: editForm.value.amount,
      expense_date: editForm.value.expense_date,
      notes: editForm.value.notes || null,
    })
    if (success) { showEditDialog.value = false; await expensesStore.fetchExpenses(); emit('expense-updated') }
  } catch { push.error('Failed to update expense') }
}

const deleteExpense = async () => {
  try { const success = await expensesStore.deleteExpense(props.expense.id); if (success) showDeleteConfirm.value = false }
  catch { push.error('Failed to delete expense') }
}

const formatDate = (d: string): string => new Date(d).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })
</script>
