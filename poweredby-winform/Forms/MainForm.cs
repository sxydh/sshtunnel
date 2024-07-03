using Newtonsoft.Json;
using poweredby_winform.Models;
using poweredby_winform.Utils;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Runtime.InteropServices;
using System.Windows.Forms;

namespace poweredby_winform.Forms
{
    public partial class MainForm : MyForm
    {

        private TcpHelper tcpHelper;
        private FileHelper fileHelper = FileHelper.New("./sshtunnel.config");

        public MainForm()
        {
            /* 初始化 GO 服务 */
            int port = InitGoServer();

            /* 初始化 TCP 客户端 */
            tcpHelper = TcpHelper.New(port);

            /* 初始化 UI 组件 */
            InitializeComponent();

            /* 防止闪退 */
            AppDomain.CurrentDomain.UnhandledException += new UnhandledExceptionEventHandler(HandleUnknownException);

            /* 关闭回调 */
            FormClosing += new FormClosingEventHandler(OnFormClosing);
        }

        public void HandleUnknownException(object sender, UnhandledExceptionEventArgs e)
        {
            Exception ex = e.ExceptionObject as Exception;
            MessageBox.Show(ex.Message, "Error", MessageBoxButtons.OK, MessageBoxIcon.Error);
        }

        private void OnFormClosing(object sender, FormClosingEventArgs e)
        {
            List<Tunnel> tl = tunnelList.ToList();
            fileHelper.W(JsonConvert.SerializeObject(tl));
        }

        [DllImport("go_export.dll")]
        public static extern int InitGoServer();
    }
}
