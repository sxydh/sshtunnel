using Newtonsoft.Json;
using poweredby_winform.Models;
using System.Collections.Generic;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;

namespace poweredby_winform.Utils
{
    public class TcpHelper
    {

        private TcpClient client;
        private NetworkStream stream;
        private List<MsgHandler> msgHandlers = new List<MsgHandler>();
        public delegate void MsgHandler(Msg msg);

        private TcpHelper() { }

        public static TcpHelper New(int port)
        {
            TcpHelper tcpHelper = new TcpHelper
            {
                client = new TcpClient("127.0.0.1", port)
            };
            tcpHelper.stream = tcpHelper.client.GetStream();

            Task.Run(() =>
            {
                while (true)
                {
                    byte[] bytes = new byte[4];
                    tcpHelper.stream.Read(bytes, 0, bytes.Length);
                    bytes = new byte[ByteHelper.Byte4ToInt(bytes)];
                    tcpHelper.stream.Read(bytes, 0, bytes.Length);
                    Msg msg = JsonConvert.DeserializeObject<Msg>(Encoding.UTF8.GetString(bytes));
                    foreach (var msgHandler in tcpHelper.msgHandlers)
                    {
                        msgHandler.Invoke(msg);
                    }
                }
            });

            return tcpHelper;
        }

        public void Send(string msg)
        {
            byte[] bytes = Encoding.UTF8.GetBytes(msg);
            /* 解决粘包问题 */
            stream.Write(ByteHelper.Bytes4(bytes.Length), 0, 4);
            bytes = Encoding.UTF8.GetBytes(msg);
            stream.Write(bytes, 0, bytes.Length);
        }

        public void OnMsg(MsgHandler msgHandler)
        {
            msgHandlers.Add(msgHandler);
        }

    }
}
