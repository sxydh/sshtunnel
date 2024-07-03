using Newtonsoft.Json;

namespace poweredby_winform.Models
{
    public class Msg
    {

        [JsonProperty("flag")]
        public string Flag { get; set; }

        [JsonProperty("body")]
        public string Body { get; set; }

    }
}
