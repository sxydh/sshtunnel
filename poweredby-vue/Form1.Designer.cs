using Microsoft.Web.WebView2.WinForms;
using System.Windows.Forms;

namespace poweredby_vue
{
    partial class Form1
    {
        private System.ComponentModel.IContainer components = null;
        private WebView2 webView2;

        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        private void InitializeComponent()
        {
            webView2 = new WebView2();
            (webView2 as System.ComponentModel.ISupportInitialize).BeginInit();
            webView2.Dock = System.Windows.Forms.DockStyle.Fill;
            webView2.Source = new System.Uri($"http://localhost:{fsPort}/?serverPort={wsPort}", System.UriKind.Absolute);
            webView2.EnsureCoreWebView2Async().ContinueWith(p =>
            {
                Invoke((MethodInvoker)delegate
                {
                    webView2.CoreWebView2.Settings.AreDefaultContextMenusEnabled = false;
                    webView2.CoreWebView2.Settings.IsWebMessageEnabled = true;
                    webView2.CoreWebView2.WebMessageReceived += WebMessageReceived;
                });
            });
            this.Controls.Add(webView2);
            (webView2 as System.ComponentModel.ISupportInitialize).EndInit();

            this.components = new System.ComponentModel.Container();
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(1800, 900);
            this.Text = "Form1";
        }   

    }
}

