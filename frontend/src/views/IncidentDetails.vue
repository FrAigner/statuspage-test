<template>
  <div v-if="incident" class="space-y-8">
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <div class="px-4 py-5 sm:px-6">
        <h2 class="text-lg font-medium leading-6 text-gray-900">{{ incident.title }}</h2>
        <p class="mt-1 max-w-2xl text-sm text-gray-500">
          Status: 
          <span :class="getIncidentStatusClass(incident.status)" class="px-2 py-1 rounded-full ml-2">
            {{ incident.status }}
          </span>
        </p>
      </div>
      <div class="border-t border-gray-200 px-4 py-5">
        <div class="prose max-w-none">
          {{ incident.description }}
        </div>
      </div>
    </div>

    <!-- Updates -->
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <div class="px-4 py-5 sm:px-6 flex justify-between items-center">
        <h3 class="text-lg font-medium leading-6 text-gray-900">Updates</h3>
        <button 
          @click="showUpdateModal = true"
          class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium"
        >
          Update hinzufügen
        </button>
      </div>
      <div class="border-t border-gray-200">
        <ul class="divide-y divide-gray-200">
          <li v-for="update in updates" :key="update.id" class="px-4 py-4">
            <div class="space-y-2">
              <div class="flex justify-between">
                <span :class="getIncidentStatusClass(update.status)" class="px-2 py-1 text-xs rounded-full">
                  {{ update.status }}
                </span>
                <span class="text-sm text-gray-500">
                  {{ new Date(update.created_at).toLocaleString() }}
                </span>
              </div>
              <p class="text-sm text-gray-700">{{ update.message }}</p>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
import type { Incident, IncidentUpdate } from '../types'

const incident = ref<Incident | null>(null)
const updates = ref<IncidentUpdate[]>([])
const showUpdateModal = ref(false)

const getIncidentStatusClass = (status: string) => {
  switch (status) {
    case 'investigating':
      return 'bg-yellow-100 text-yellow-800'
    case 'identified':
      return 'bg-blue-100 text-blue-800'
    case 'monitoring':
      return 'bg-purple-100 text-purple-800'
    case 'resolved':
      return 'bg-green-100 text-green-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const fetchIncidentDetails = async () => {
  try {
    const res = await axios.get(`http://localhost:8080/api/incidents/${route.params.id}`)
    incident.value = res.data.incident
    updates.value = res.data.updates
  } catch (error) {
    console.error('Error fetching incident details:', error)
  }
}

onMounted(fetchIncidentDetails)
</script>
