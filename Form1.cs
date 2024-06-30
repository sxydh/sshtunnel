using Microsoft.Web.WebView2.Core;
using System.Windows.Forms;

namespace sshtunnel
{
    public partial class Form1 : Form
    {

        private int fsPort;
        private int wsPort;
        private int closeFlag = -1;

        public Form1()
        {
            fsPort = InitFsServer();
            wsPort = InitWsServer();
            InitializeComponent();

            FormClosing += BeforeClose;
        }

        private void BeforeClose(object sender, FormClosingEventArgs e)
        {
            if (closeFlag <= 0)
            {
                closeFlag += 1;
                e.Cancel = true;
                webView2.ExecuteScriptAsync("window.beforeClose()");
            }
        }

        private void WebMessageReceived(object sender, CoreWebView2WebMessageReceivedEventArgs e)
        {
            var msg = e.TryGetWebMessageAsString();
            if (msg != null)
            {
                closeFlag = 1;
                Close();
            }
        }

        [System.Runtime.InteropServices.DllImport("fs.dll", EntryPoint = "InitFsServer")]
        private static extern int InitFsServer();

        [System.Runtime.InteropServices.DllImport("ws.dll", EntryPoint = "InitWsServer")]
        private static extern int InitWsServer();

    }
}
