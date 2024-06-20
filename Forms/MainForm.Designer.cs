using Newtonsoft.Json;
using sshtunnel.Models;
using sshtunnel.Utils;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Reflection;
using System.Windows.Forms;

namespace sshtunnel.Forms
{
    partial class MainForm
    {
        private System.ComponentModel.IContainer components = null;

        /* Button */
        private Button execButton;
        private Button stopButton;

        /* Table */
        private DataGridView tunnelTable;
        private BindingList<Tunnel> tunnelList;

        /* LogView */
        private ListView logView;

        /* Other */
        private FileHelper fileHelper = FileHelper.New("./sshtunnel.config");

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
            /* 
             * 先 SuspendLayout 之后再 ResumeLayout 
             * 该步骤可以减少重绘以优化性能
             */
            SuspendLayout();

            /* Panel */
            var panel = new System.Windows.Forms.TableLayoutPanel
            {
                Name = "Panel",
                Dock = DockStyle.Fill,
                ColumnCount = 1,
                RowCount = 3,
                ColumnStyles = {
                    new ColumnStyle(SizeType.Percent, 100)
                },
                RowStyles = {
                    new RowStyle(SizeType.Absolute, 50),
                    new RowStyle(SizeType.Absolute, 400),
                    new RowStyle(SizeType.Percent, 100)
                },
                // BackColor = System.Drawing.Color.Blue
            };

            var buttonFlowPanel = new System.Windows.Forms.FlowLayoutPanel
            {
                Name = "ButtonFlowPanel",
                Dock = DockStyle.Fill,
                // BackColor = System.Drawing.Color.Orange,
            };
            panel.Controls.Add(buttonFlowPanel, 0, 0);

            /* Button */
            execButton = new Button
            {
                Name = "ExecButton",
                Text = "Run",
                Width = 100,
                Height = buttonFlowPanel.Height
            };
            execButton.Click += new System.EventHandler(HandleExecButtonClick);
            buttonFlowPanel.Controls.Add(execButton);

            stopButton = new Button
            {
                Name = "StopButton",
                Text = "Stop",
                Width = 100,
                Height = buttonFlowPanel.Height,
                Visible = false,
            };
            stopButton.Click += new System.EventHandler(HandleStopButtonClick);
            buttonFlowPanel.Controls.Add(stopButton);

            /* Table */
            tunnelTable = new DataGridView
            {
                Name = "TunnelTable",
                Dock = DockStyle.Fill,
                AutoSizeColumnsMode = DataGridViewAutoSizeColumnsMode.Fill,
                ColumnHeadersHeight = 50,
                RowTemplate = new DataGridViewRow
                {
                    Height = 40,
                },
                // BackgroundColor = System.Drawing.Color.Beige
            };
            tunnelList = new BindingList<Tunnel>();
            InitTunnelList(tunnelList);
            tunnelTable.DataSource = tunnelList;
            panel.Controls.Add(tunnelTable, 0, 1);

            /* Log */
            logView = new System.Windows.Forms.ListView
            {
                Name = "LogView",
                Dock = DockStyle.Fill,
                // BackColor = System.Drawing.Color.Chocolate
            };
            logView.View = View.Details;
            logView.Columns.Add("Log", 3000, HorizontalAlignment.Left);
            Type type = logView.GetType();
            PropertyInfo pi = type.GetProperty("DoubleBuffered", BindingFlags.Instance | BindingFlags.NonPublic);
            pi.SetValue(logView, true, null);
            panel.Controls.Add(logView, 0, 2);

            AutoScaleDimensions = new System.Drawing.SizeF(9F, 18F);
            AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            ClientSize = new System.Drawing.Size(1200, 800);
            Name = "MainForm";
            Text = "SSH Tunnel";
            Controls.Add(panel);

            ResumeLayout(false);
        }

        private void InitTunnelList(BindingList<Tunnel> tlb)
        {
            string str = fileHelper.R();
            if (string.IsNullOrEmpty(str))
            {
                return;
            }
            List<Tunnel> tl = JsonConvert.DeserializeObject<List<Tunnel>>(str);
            foreach (var t in tl)
            {
                tlb.Add(t);
            }
        }

        private void InitLogSource()
        {
            Action<Msg> action = msg =>
            {
                if (msg.Flag != "Log") { return; }
                this.Invoke((MethodInvoker)delegate
                {
                    logView.BeginUpdate();
                    logView.Items.Insert(
                        0,
                        new ListViewItem
                        {
                            Text = msg.Body
                        }
                    );
                    if (logView.Items.Count > 100)
                    {
                        logView.Items.RemoveAt(logView.Items.Count - 1);
                    }
                    logView.EndUpdate();
                });
            };
            tcpHelper.OnMsg(new Utils.TcpHelper.MsgHandler(action));
        }

        private void HandleExecButtonClick(object sender, EventArgs e)
        {
            if (tunnelList.Count == 0)
            {
                return;
            }
            List<Tunnel> tl = tunnelList.ToList();
            /* 正向 */
            List<Tunnel> tlp = tl.FindAll(ele => ele.Direction == 1);
            if (tlp.Count > 0) {
                Msg msg = new Msg
                {
                    Flag = "NewTunnel",
                    Body = JsonConvert.SerializeObject(tlp)
                };
                tcpHelper.Send(JsonConvert.SerializeObject(msg));
            }
            /* 反向 */
            List<Tunnel> tln = tl.FindAll(ele => ele.Direction == -1);
            if (tln.Count > 0)
            {
                Msg msg = new Msg
                {
                    Flag = "NewTunnel",
                    Body = JsonConvert.SerializeObject(tln)
                };
                tcpHelper.Send(JsonConvert.SerializeObject(msg));
            }
            execButton.Visible = false;
            stopButton.Visible = true;
        }

        private void HandleStopButtonClick(object sender, EventArgs e)
        {
            Msg msg = new Msg
            {
                Flag = "StopTunnel",
            };
            tcpHelper.Send(JsonConvert.SerializeObject(msg));
            execButton.Visible = true;
            stopButton.Visible = false;
        }

        private void OnFormClosing(object sender, FormClosingEventArgs e)
        {
            List<Tunnel> tl = tunnelList.ToList();
            fileHelper.W(JsonConvert.SerializeObject(tl));
        }

    }
}