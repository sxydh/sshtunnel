using Newtonsoft.Json;

namespace sshtunnel.Models
{
    public class Tunnel
    {

        [JsonProperty("sshIp")]
        public string SSHIp { get; set; }

        [JsonProperty("sshPort")]
        public string SSHPort { get; set; }

        [JsonProperty("sshUser")]
        public string SSHUser { get; set; }

        [JsonProperty("listenPort")]
        public string ListenPort { get; set; }

        [JsonProperty("targetIp")]
        public string TargetIp { get; set; }

        [JsonProperty("targetPort")]
        public string TargetPort { get; set; }

        [JsonProperty("status")]
        public string Status { get; set; }

        [JsonProperty("delete")]
        public string Delete { get; set; }

    }
}
