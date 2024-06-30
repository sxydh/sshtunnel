<template>
  <v-card>
    <v-card-actions>
      <v-row>
        <v-col class="d-flex align-center justify-center">
          <v-text-field
              class="mx-3"
              v-model="filterValue"
              prepend-icon="mdi-magnify"
              variant="solo"
              density="comfortable"
              hide-details
              clearable
              max-width="300"
          />
          <v-btn
              class="mx-3"
              variant="tonal"
              density="comfortable"
              :icon="btnCase.icon"
              :color="btnCase.color"
              @click.prevent="handleBtnEvent"
          />
        </v-col>
      </v-row>
    </v-card-actions>

    <v-form ref="tunnelForm">
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
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.sshIp"
                  :rules="[rules.ip(index)]"
                  density="comfortable"
                  variant="solo"
                  clearable
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.sshPort"
                  :rules="[rules.required(index)]"
                  type="number"
                  hide-spin-buttons
                  density="comfortable"
                  variant="solo"
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.sshUser"
                  :rules="[rules.required(index)]"
                  density="comfortable"
                  variant="solo"
                  clearable
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.listenPort"
                  :rules="[rules.required(index)]"
                  type="number"
                  hide-spin-buttons
                  density="comfortable"
                  variant="solo"
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.targetIp"
                  :rules="[rules.ip(index)]"
                  density="comfortable"
                  variant="solo"
                  clearable
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.targetPort"
                  :rules="[rules.required(index)]"
                  type="number"
                  hide-spin-buttons
                  density="comfortable"
                  variant="solo"
              />
            </td>
            <td>
              <v-text-field
                  v-model="item.lastAlive"
                  density="comfortable"
                  variant="solo"
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
    </v-form>
  </v-card>
</template>

<script
    setup
    lang="ts"
>
import {onMounted, ref} from 'vue'
import {Tunnel} from '@/models/Tunnel'
import {Msg} from '@/models/Msg'
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
    width: 200,
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
const tunnelTemplate = {
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
const filterValue = ref()
const rules = ref({
  required(index: number) {
    return (value: any) => index === (tunnels.value.length - 1) || !!value || 'Please input'
  },
  ip(index: number) {
    return (value: any) => index === (tunnels.value.length - 1) || /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(value) || 'Please input IP'
  },
})
const tunnelForm = ref(null)
const pushBtn = {
  icon: 'mdi-send',
  color: '#2DC1DD',
}
const stopBtn = {
  icon: 'mdi-stop-circle',
  color: '#ff3a3a',
}
const btnCase = ref(pushBtn)
// WebSocket
const params = new URLSearchParams(window.location.search)
const port = params.get('serverPort')
const webSocket = new WebSocket(`ws://localhost:${port}`)
webSocket.onopen = () => {
  console.debug(`WebSocket onopen: port=${port}`)
  initTunnels()
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
window.onbeforeunload = () => {
  const targetList = tunnels.value.slice(0, -1)
  const msg: Msg = {
    flag: 'SaveTunnel',
    body: JSON.stringify(targetList),
  }
  send(msg)
}

/* 函数 */
const send = (msg: Msg) => {
  console.debug(`webSocket.send`, msg)
  webSocket.send(JSON.stringify(msg))
}
const initTunnels = () => {
  const msg: Msg = {
    flag: 'ListSavedTunnel',
    body: '',
  }
  send(msg)
}
const initLastTunnel = (tunnels: Tunnel[]) => {
  if (tunnels.length != 0 && ifEqual(tunnels[tunnels.length - 1], tunnelTemplate, ['id'])) {
    return
  }
  const newTunnel = JSON.parse(JSON.stringify(tunnelTemplate))
  newTunnel.id = new Date().getTime().toString()
  tunnels.push(newTunnel)
}
const handleTrInputEvent = (p: any) => {
  if (p == tunnels.value.length - 1) {
    initLastTunnel(tunnels.value)
  }
}
const handleTrDeleteEvent = (p: any) => {
  tunnels.value.splice(p, 1)
}
const handleBtnEvent = async () => {
  if (btnCase.value === pushBtn) {
    const pushed = await handlePushEvent()
    if (pushed) {
      btnCase.value = stopBtn
    }
  } else if (btnCase.value === stopBtn) {
    const stopped = handleStopEvent()
    if (stopped) {
      btnCase.value = pushBtn
    }
  }
}
const handlePushEvent = async (): Promise<boolean> => {
  const results = await (tunnelForm.value as any).validate()
  if (!results.valid) {
    return false
  }
  let list = tunnels.value.slice(0, -1)
  if (list.length === 0) {
    return false
  }
  /* 正向 */
  let targetList = list.filter(ele => ele.direction === 1)
  if (targetList.length > 0) {
    const msg: Msg = {
      flag: 'NewTunnel',
      body: JSON.stringify(targetList),
    }
    send(msg)
  }
  /* 反向 */
  targetList = list.filter(ele => ele.direction === -1)
  if (targetList.length > 0) {
    const msg: Msg = {
      flag: 'NewReverseTunnel',
      body: JSON.stringify(targetList),
    }
    send(msg)
  }
  return true
}
const handleStopEvent = (): boolean => {
  const msg: Msg = {
    flag: 'StopTunnel',
    body: '',
  }
  send(msg)
  return true
}
</script>

<style scoped>
td {
  text-align: center;
  padding: 5px 10px !important;
}
</style>