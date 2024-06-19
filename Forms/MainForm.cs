using System.Runtime.InteropServices;
using System.Windows.Forms;

namespace sshtunnel.Forms
{
    public partial class MainForm : Form
    {

        private int goServerPort;

        public MainForm()
        {
            /* 初始化 UI 组件 */
            InitializeComponent();

            /* 初始化 GO 服务 */
            goServerPort = InitGoServer();
        }

        [DllImport("go_export.dll")]
        public static extern int InitGoServer();
    }
}
