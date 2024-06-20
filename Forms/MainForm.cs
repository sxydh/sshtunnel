using sshtunnel.Utils;
using System.Runtime.InteropServices;
using System.Windows.Forms;

namespace sshtunnel.Forms
{
    public partial class MainForm : Form
    {

        private TcpHelper tcpHelper;

        public MainForm()
        {
            /* 初始化 UI 组件 */
            InitializeComponent();

            /* 初始化 GO 服务 */
            int goServerPort = InitGoServer();

            /* 初始化 TCP 客户端 */
            tcpHelper = TcpHelper.New(goServerPort);

            /* 初始化日志数据源 */
            InitLogSource();
        }

        [DllImport("go_export.dll")]
        public static extern int InitGoServer();
    }
}
