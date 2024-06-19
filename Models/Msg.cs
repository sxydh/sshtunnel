using Newtonsoft.Json;

namespace sshtunnel.Models
{
    public class Msg
    {

        [JsonProperty("flag")]
        public string Flag { get; set; }

        [JsonProperty("body")]
        public string Body { get; set; }

    }
}
