using Newtonsoft.Json;

namespace poweredby_winform.Models
{
    public class Tunnel
    {

        [JsonProperty("id")]
        public string Id { get; set; }

        [JsonProperty("direction")]
        public int Direction { get; set; }

        [JsonProperty("sshIp")]
        public string SSHIp { get; set; }

        [JsonProperty("sshPort")]
        public int SSHPort { get; set; }

        [JsonProperty("sshUser")]
        public string SSHUser { get; set; }

        [JsonProperty("listenPort")]
        public int ListenPort { get; set; }

        [JsonProperty("targetIp")]
        public string TargetIp { get; set; }

        [JsonProperty("targetPort")]
        public int TargetPort { get; set; }

        [JsonProperty("status")]
        public int Status { get; set; }

        [JsonProperty("delete")]
        public int Delete { get; set; }

        [JsonProperty("lastAlive")]
        public string LastAlive { get; set; }

    }
}
