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
          class="tunnel-table"
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
                  v-model.number="item.sshPort"
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
                  v-model.number="item.listenPort"
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
                  v-model.number="item.targetPort"
                  :rules="[rules.required(index)]"
                  type="number"
                  hide-spin-buttons
                  density="comfortable"
                  variant="solo"
              />
            </td>
            <td>
              <v-text-field
                  :id="`${item.id}_last_alive`"
                  density="comfortable"
                  variant="solo"
                  readonly
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

    <v-data-table
        :items="logs"
        hide-default-header
        hide-default-footer
        class="log-table"
    />
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
import {stringify} from '@/utils/json_util'

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
    return (value: any) => index === (tunnels.value.length - 1) || /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(value) || value == 'localhost' || 'Please input IP'
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
const job = ref<boolean>(true)
const params = new URLSearchParams(window.location.search)
const port = params.get('serverPort')
const webSocket = ref()
const logs = ref<{ text: string }[]>([])
const messages = ref<number>(0)
const wd = window as any

/* 回调 */
onMounted(() => {
  initWs()
  initLastTunnel(tunnels.value)
})

/* 函数 */
const initWs = () => {
  const ws = new WebSocket(`ws://localhost:${port}`)
  ws.onopen = () => {
    onOpen()
  }
  ws.onmessage = e => {
    onMessage(e)
  }
  ws.onerror = e => {
    onError(e)
  }
  ws.onclose = e => {
    onClose(e)
  }
  webSocket.value = ws
}
const onOpen = () => {
  console.debug(`WebSocket onopen: port=${port}`)
  initTunnels()
  initJob()
}
const onMessage = (e: any) => {
  try {
    messages.value += 1
    const msg: Msg = JSON.parse(e.data)
    switch (msg.flag) {
      case 'Log':
        logs.value.unshift({text: msg.body})
        if (logs.value.length > 100) {
          logs.value.splice(-1, 1)
        }
        break
      case 'ListSavedTunnel':
        const savedTunnels: Tunnel[] = JSON.parse(msg.body)
        if (savedTunnels.length > 0) {
          tunnels.value.unshift(...savedTunnels)
        }
        break
      case 'ListTunnel':
        const serverTunnels: Tunnel[] = JSON.parse(msg.body)
        if (serverTunnels.length > 0) {
          if (btnCase.value !== stopBtn) {
            btnCase.value = stopBtn
          }
          if (tunnels.value.length > 1) {
            for (const tunnel of tunnels.value) {
              const serverTunnel = serverTunnels.find(se => se.id === tunnel.id)
              if (serverTunnel) {
                handleLastAlive(`${tunnel.id}_last_alive`, serverTunnel.lastAlive)
              }
            }
          }
        }
        break
      case 'SaveTunnel':
        wd.chrome.webview.postMessage('1')
        break
    }
  } finally {
    messages.value -= 1
  }
}
const onError = (e: any) => {
  console.debug(`WebSocket onerror`, e)
}
const onClose = (e: any) => {
  console.debug(`WebSocket onclose`, e)
  initWs()
}
const send = (msg: Msg) => {
  const ws = webSocket.value as any
  if (ws.readyState !== WebSocket.OPEN) {
    return
  }
  webSocket.value.send(stringify(msg))
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
  const newTunnel = JSON.parse(stringify(tunnelTemplate))
  newTunnel.id = new Date().getTime().toString()
  tunnels.push(newTunnel)
}
const initJob = () => {
  handleJob()
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
  if (btnCase.value.icon === pushBtn.icon) {
    const pushed = await handlePushEvent()
    if (pushed) {
      btnCase.value = stopBtn
    }
  } else if (btnCase.value.icon === stopBtn.icon) {
    const stopped = handleStopEvent()
    if (stopped) {
      btnCase.value = pushBtn
    }
  }
}
const handlePushEvent = async (): Promise<boolean> => {
  // noinspection ES6RedundantAwait
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
      body: stringify(targetList),
    }
    send(msg)
  }
  /* 反向 */
  targetList = list.filter(ele => ele.direction === -1)
  if (targetList.length > 0) {
    const msg: Msg = {
      flag: 'NewReverseTunnel',
      body: stringify(targetList),
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
const handleSaveEvent = (tunnels: Tunnel[]) => {
  const targetList = tunnels.slice(0, -1)
  const msg: Msg = {
    flag: 'SaveTunnel',
    body: stringify(targetList),
  }
  send(msg)
}
const handleLastAlive = (id: string, val: string | null) => {
  const input = document.getElementById(id) as any
  input.value = val
  const diff = new Date().getTime() - new Date(input.value).getTime()
  const diffMi = diff / 1000
  if (diffMi > 20) {
    input.className = `v-input-warn`
  } else {
    input.className = `v-input-success`
  }
}
const handleJob = () => {
  if (messages.value === 0) {
    // 用于更新最后存活时间
    const msg: Msg = {
      flag: 'ListTunnel',
      body: '',
    }
    send(msg)
  }
  if (!job.value) {
    return
  }
  setTimeout(handleJob, 10)
}
wd.beforeClose = () => {
  handleSaveEvent(tunnels.value)
  job.value = false
  console.debug(`beforeClose going...`)
}
</script>

<style scoped>
.tunnel-table {
  :deep(td) {
    text-align: center;
    padding: 5px 10px !important;
  }

  :deep(.v-input-success) {
    color: #06fc91 !important;
  }

  :deep(.v-input-warn) {
    color: #ff3a3a !important;
  }
}

.log-table :deep(td) {
  color: #adadad;
}
</style>