<template>
  <div class="min-h-screen bg-app flex flex-col">

    <div class="w-full flex-1 flex flex-col min-h-0">

      <!-- MAIN CONTAINER -->
      <div class="bg-surface border border-default flex flex-col flex-1 min-h-0">

        <!-- HEADER -->
        <div class="flex items-center gap-4 px-6 py-3 border-b border-divider shrink-0">

          <!-- Tabs -->
          <div class="flex-1 min-w-0 flex items-center">
            <AppTabs v-model="activeTab" :tabs="filteredTabs" />
          </div>

          <!-- Settings -->
          <div class="flex items-center shrink-0">
            <AppSettings />
          </div>

        </div>

        <!-- CONTENT -->
        <div class="flex-1 min-h-0 flex flex-col overflow-hidden">

          <section v-if="activeTab === 'customers'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <CustomerList />
          </section>

          <section v-else-if="activeTab === 'contracts'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <ContractList />
          </section>

          <section v-else-if="activeTab === 'items'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <ItemList />
          </section>

          <section v-else-if="activeTab === 'shipments'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <ShipmentList />
          </section>

          <section v-else-if="activeTab === 'users'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <UserList />
          </section>

          <section v-else-if="activeTab === 'expenses'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <ExpenseList />
          </section>

          <section v-else-if="activeTab === 'audit'" class="flex-1 min-h-0 flex flex-col overflow-hidden p-6">
            <LogList />
          </section>

          <section v-else class="flex-1 flex items-center justify-center">
            <div class="text-center space-y-3">
              <div class="w-16 h-16 mx-auto rounded-full bg-surface-alt flex items-center justify-center">
                <svg class="w-8 h-8 text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                    d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium text-primary">{{ formattedActiveTab }} coming soon</p>
                <p class="text-xs text-muted mt-1">This section is under development</p>
              </div>
            </div>
          </section>

        </div>

      </div>

    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

import AppSettings from '@/components/AppSettings.vue'
import AppTabs from '@/components/ui/AppTabs.vue'

import UserList from '@/components/users/UserList.vue'
import CustomerList from '@/components/customers/CustomerList.vue'
import ContractList from '@/components/contracts/ContractList.vue'
import ItemList from '@/components/items/ItemList.vue'
import ShipmentList from '@/components/shipments/ShipmentList.vue'
import ExpenseList from '@/components/expenses/ExpenseList.vue'
import LogList from '@/components/logs/LogList.vue'

const auth = useAuthStore()

const activeTab = ref('customers')

const formattedActiveTab = computed(() =>
  activeTab.value.charAt(0).toUpperCase() + activeTab.value.slice(1)
)

const allTabs = [
  { key: 'customers', label: 'Customers', roles: ['admin', 'manager'] },
  { key: 'contracts', label: 'Contracts', roles: ['admin', 'manager'] },
  { key: 'items', label: 'Items', roles: ['admin', 'manager', 'staff'] },
  { key: 'shipments', label: 'Shipments', roles: ['admin', 'manager'] },
  { key: 'divider1', label: '', divider: true, roles: ['admin', 'manager'] },
  { key: 'expenses', label: 'Expenses', roles: ['admin', 'manager', 'accounts'] },
  { key: 'users', label: 'Users', roles: ['admin', 'manager'] },
  { key: 'divider2', label: '', divider: true, roles: ['admin'] },
  { key: 'audit', label: 'Audits', roles: ['admin'] },
]

const filteredTabs = computed(() => {
  const role = auth.user?.role || 'staff'
  return allTabs.filter(tab => tab.roles.includes(role))
})

// Set initial active tab based on role
if (filteredTabs.value.length > 0) {
  const firstTab = filteredTabs.value[0]
  if (firstTab) {
    activeTab.value = firstTab.key
  }
}
</script>
