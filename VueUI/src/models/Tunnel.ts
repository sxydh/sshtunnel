export interface Tunnel {
    id: string | null
    direction: number | null
    sshIp: string | null
    sshPort: number | null
    sshUser: string | null
    listenPort: number | null
    targetIp: string | null
    targetPort: number | null
    status: number | null
    delete: number | null
    lastAlive: string | null
}