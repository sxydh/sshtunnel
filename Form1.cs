using System.Windows.Forms;

namespace sshtunnel
{
    public partial class Form1 : Form
    {

        private int fsPort;
        private int wsPort;

        public Form1()
        {
            fsPort = InitFsServer();
            wsPort = InitWsServer();
            InitializeComponent();
        }

        [System.Runtime.InteropServices.DllImport("fs.dll", EntryPoint = "InitFsServer")]
        private static extern int InitFsServer();

        [System.Runtime.InteropServices.DllImport("ws.dll", EntryPoint = "InitWsServer")]
        private static extern int InitWsServer();

    }
}
