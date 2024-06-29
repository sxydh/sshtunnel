<template>
  <v-card>
    <v-data-table
        :headers="headers"
        :items="tunnels"
        hide-default-footer
    >
      <template v-slot:item="{ item, index }">
        <tr @input="handleTrInputEvent(index)">
          <td>
            <v-select
                v-model="item.direction"
                :items="directions"
                item-title="name"
                item-value="value"
                density="compact"
                variant="solo"
                hide-details
            />
          </td>
          <td>
            <v-text-field
                v-model="item.sshIp"
                density="comfortable"
                variant="solo"
                hide-details
                clearable
            />
          </td>
          <td>
            <v-text-field
                v-model="item.sshPort"
                type="number"
                hide-spin-buttons
                density="comfortable"
                variant="solo"
                hide-details
            />
          </td>
          <td>
            <v-text-field
                v-model="item.sshUser"
                density="comfortable"
                variant="solo"
                hide-details
                clearable
            />
          </td>
          <td>
            <v-text-field
                v-model="item.listenPort"
                type="number"
                hide-spin-buttons
                density="comfortable"
                variant="solo"
                hide-details
            />
          </td>
          <td>
            <v-text-field
                v-model="item.targetIp"
                density="comfortable"
                variant="solo"
                hide-details
                clearable
            />
          </td>
          <td>
            <v-text-field
                v-model="item.targetPort"
                type="number"
                hide-spin-buttons
                density="comfortable"
                variant="solo"
                hide-details
            />
          </td>
          <td>
            <v-text-field
                v-model="item.lastAlive"
                density="comfortable"
                variant="solo"
                hide-details
                disabled
            />
          </td>
          <td v-if="index != (tunnels.length - 1)">
            <v-icon
                icon="mdi-delete"
                color="#ff3a3a"
                @click="handleTrDeleteEvent(index)"
            />
          </td>
        </tr>
      </template>
    </v-data-table>
  </v-card>
</template>

<script
    setup
    lang="ts"
>
import {onMounted, ref} from 'vue'
import {Tunnel} from '@/models/Tunnel'
import {ifEqual} from '@/utils/obj_util'

/* 变量 */
const headers = [
  {
    title: 'Direction',
    align: 'center',
    key: 'direction',
    width: 150,
  },
  {
    title: 'SSH IP',
    align: 'center',
    key: 'sshIp',
    width: 200,
  },
  {
    title: 'SSH Port',
    align: 'center',
    key: 'sshPort',
    width: 150,
  },
  {
    title: 'SSH User',
    align: 'center',
    key: 'sshUser',
    width: 150,
  },
  {
    title: 'Listen Port',
    align: 'center',
    key: 'listenPort',
    width: 150,
  },
  {
    title: 'Target IP',
    align: 'center',
    key: 'targetIp',
    width: 150,
  },
  {
    title: 'Target Port',
    align: 'center',
    key: 'targetPort',
    width: 150,
  },
  {
    title: 'Last Alive',
    align: 'center',
    key: 'lastAlive',
    width: 200,
  },
  {
    title: 'Action',
    align: 'center',
    key: 'action',
    width: 50,
    sortable: false,
  },
]
const tunnels = ref<Tunnel[]>([])
const newTunnel = {
  id: null,
  direction: 1,
  sshIp: null,
  sshPort: null,
  sshUser: null,
  listenPort: null,
  targetIp: null,
  targetPort: null,
  status: null,
  delete: null,
  lastAlive: null,
}
const directions = ref([
  {
    name: 'Tunnel',
    value: 1,
  },
  {
    name: 'Reverse Tunnel',
    value: -1,
  },
])
// WebSocket
const params = new URLSearchParams(window.location.search)
const port = params.get('serverPort')
const webSocket = new WebSocket(`ws://localhost:${port}`)
webSocket.onopen = () => {
  console.debug(`WebSocket onopen: port=${port}`)
}
webSocket.onmessage = e => {
  console.debug(`WebSocket onmessage`, e)
}
webSocket.onerror = e => {
  console.debug(`WebSocket onerror`, e)
}
webSocket.onclose = e => {
  console.debug(`WebSocket onclose`, e)
}

/* 回调 */
onMounted(() => {
  initLastTunnel(tunnels.value)
})

/* 函数 */
const initLastTunnel = (tunnels: Tunnel[]) => {
  if (tunnels.length != 0 && ifEqual(tunnels[tunnels.length - 1], newTunnel)) {
    return
  }
  tunnels.push(JSON.parse(JSON.stringify(newTunnel)))
}
const handleTrInputEvent = (p: any) => {
  if (p == tunnels.value.length - 1) {
    initLastTunnel(tunnels.value)
  }
}
const handleTrDeleteEvent = (p: any) => {
  tunnels.value.splice(p, 1)
}
</script>

<style scoped>
td {
  text-align: center;
  padding: 5px 10px !important;
}
</style>