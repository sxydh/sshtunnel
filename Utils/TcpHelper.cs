using System.Net.Sockets;
using System.Text;

namespace sshtunnel.Utils
{

    public class TcpHelper
    {

        private TcpClient client;
        private NetworkStream stream;

        private TcpHelper() { }

        public static TcpHelper New(int port)
        {
            TcpHelper tcpHelper = new TcpHelper
            {
                client = new TcpClient("127.0.0.1", port)
            };
            tcpHelper.stream = tcpHelper.client.GetStream();
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

    }
}
