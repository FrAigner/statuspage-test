<template>
  <div class="space-y-8">
    <!-- Services Management -->
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <div class="px-4 py-5 sm:px-6 flex justify-between items-center">
        <div>
          <h2 class="text-lg font-medium leading-6 text-gray-900">Services verwalten</h2>
          <p class="mt-1 max-w-2xl text-sm text-gray-500">Fügen Sie neue Services hinzu oder aktualisieren Sie bestehende</p>
        </div>
        <button @click="showServiceModal = true" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium">
          Service hinzufügen
        </button>
      </div>
      <div class="border-t border-gray-200">
        <ul class="divide-y divide-gray-200">
          <li v-for="service in services" :key="service.id" class="px-4 py-4">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-sm font-medium">{{ service.name }}</h3>
                <p class="text-sm text-gray-500">{{ service.description }}</p>
              </div>
              <div class="flex space-x-4">
                <select 
                  v-model="service.status"
                  @change="updateService(service)"
                  class="rounded-md border-gray-300 text-sm"
                >
                  <option value="operational">Operational</option>
                  <option value="degraded">Degraded</option>
                  <option value="outage">Outage</option>
                </select>
                <button 
                  @click="deleteService(service.id)"
                  class="text-red-600 hover:text-red-800"
                >
                  Löschen
                </button>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>

    <!-- Incidents Management -->
    <div class="bg-white shadow overflow-hidden sm:rounded-lg">
      <div class="px-4 py-5 sm:px-6 flex justify-between items-center">
        <div>
          <h2 class="text-lg font-medium leading-6 text-gray-900">Vorfälle verwalten</h2>
          <p class="mt-1 max-w-2xl text-sm text-gray-500">Erstellen und aktualisieren Sie Vorfälle</p>
        </div>
        <button @click="showIncidentModal = true" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium">
          Vorfall erstellen
        </button>
      </div>
      <div class="border-t border-gray-200">
        <ul class="divide-y divide-gray-200">
          <li v-for="incident in incidents" :key="incident.id" class="px-4 py-4">
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-sm font-medium">{{ incident.title }}</h3>
                <p class="text-sm text-gray-500">{{ incident.description }}</p>
              </div>
              <div class="flex space-x-4">
                <select 
                  v-model="incident.status"
                  @change="updateIncident(incident)"
                  class="rounded-md border-gray-300 text-sm"
                >
                  <option value="investigating">Investigating</option>
                  <option value="identified">Identified</option>
                  <option value="monitoring">Monitoring</option>
                  <option value="resolved">Resolved</option>
                </select>
                <router-link 
                  :to="{ name: 'IncidentDetails', params: { id: incident.id }}"
                  class="text-blue-600 hover:text-blue-800"
                >
                  Details
                </router-link>
              </div>
            </div>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import type { Service, Incident } from '../types'

const services = ref<Service[]>([])
const incidents = ref<Incident[]>([])
const showServiceModal = ref(false)
const showIncidentModal = ref(false)

interface NewService {
  name: string;
  description: string;
  status: Service['status'];
}

interface NewIncident {
  title: string;
  description: string;
  status: Incident['status'];
  impact: Incident['impact'];
  service_id: string;
}

const fetchData = async () => {
  try {
    const [servicesRes, incidentsRes] = await Promise.all([
      axios.get('http://localhost:8080/api/services'),
      axios.get('http://localhost:8080/api/incidents')
    ])
    services.value = servicesRes.data
    incidents.value = incidentsRes.data
  } catch (error) {
    console.error('Error fetching data:', error)
  }
}

const updateService = async (service) => {
  try {
    await axios.put(`http://localhost:8080/api/services/${service.id}`, service)
  } catch (error) {
    console.error('Error updating service:', error)
  }
}

const deleteService = async (id) => {
  if (!confirm('Sind Sie sicher, dass Sie diesen Service löschen möchten?')) return
  
  try {
    await axios.delete(`http://localhost:8080/api/services/${id}`)
    services.value = services.value.filter(s => s.id !== id)
  } catch (error) {
    console.error('Error deleting service:', error)
  }
}

const updateIncident = async (incident) => {
  try {
    await axios.put(`http://localhost:8080/api/incidents/${incident.id}`, incident)
  } catch (error) {
    console.error('Error updating incident:', error)
  }
}

onMounted(fetchData)
</script>
