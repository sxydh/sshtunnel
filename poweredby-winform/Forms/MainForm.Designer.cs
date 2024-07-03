using Newtonsoft.Json;
using poweredby_winform.Models;
using poweredby_winform.Utils;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace poweredby_winform.Forms
{
    partial class MainForm
    {
        private System.ComponentModel.IContainer components = null;

        /* Button */
        private MyButton execButton;
        private MyButton stopButton;

        /* Table */
        private MyDataGridView tunnelTable;
        private BindingList<Tunnel> tunnelList;

        /* LogView */
        private MyListView logView;

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
            var panel = new MyTableLayoutPanel
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
            };

            /* Button */
            execButton = new MyButton
            {
                Name = "ExecButton",
                Text = "Run",
                Width = 100,
                Dock = DockStyle.Left,
            };
            execButton.Click += new System.EventHandler(HandleExecButtonClick);
            panel.Controls.Add(execButton, 0, 0);

            stopButton = new MyButton
            {
                Name = "StopButton",
                Text = "Stop",
                Width = 100,
                Dock = DockStyle.Left,
                Visible = false,
            };
            stopButton.Click += new System.EventHandler(HandleStopButtonClick);
            panel.Controls.Add(stopButton, 0, 0);

            /* Table */
            tunnelTable = new MyDataGridView
            {
                Name = "TunnelTable",
                Dock = DockStyle.Fill,
                AutoGenerateColumns = false,
                ColumnHeadersHeight = 50,
                Columns = {
                    new DataGridViewTextBoxColumn {
                        HeaderText = "Direction",
                        DataPropertyName = "Direction",
                        Width = 120
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "SSH IP",
                        DataPropertyName = "SSHIp",
                        Width = 200
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "SSH Port",
                        DataPropertyName = "SSHPort",
                        Width = 120
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "SSH User",
                        DataPropertyName = "SSHUser",
                        Width = 150
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "Listen Port",
                        DataPropertyName = "ListenPort",
                        Width = 150
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "Target IP",
                        DataPropertyName = "TargetIp",
                        Width = 200
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "Target Port",
                        DataPropertyName = "TargetPort",
                        Width = 150
                    },
                    new DataGridViewTextBoxColumn {
                        HeaderText = "Last Alive",
                        DataPropertyName = "LastAlive",
                        AutoSizeMode = DataGridViewAutoSizeColumnMode.Fill,
                    }
                },
                RowTemplate = new DataGridViewRow
                {
                    Height = 40,
                },
            };
            tunnelList = new BindingList<Tunnel>();
            InitTunnelList(tunnelList);
            InitTunnelListHeartBeat(tunnelList);
            tunnelTable.DataSource = tunnelList;
            panel.Controls.Add(tunnelTable, 0, 1);

            /* Log */
            logView = new MyListView
            {
                Name = "LogView",
                Dock = DockStyle.Fill,
                View = View.Details,
                HeaderStyle = ColumnHeaderStyle.None,
                Columns = {
                    new ColumnHeader {
                        Text = "Log",
                        Width = 3000,
                    }
                }
            };
            InitLogSource();
            panel.Controls.Add(logView, 0, 2);

            AutoScaleDimensions = new System.Drawing.SizeF(9F, 18F);
            AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            ClientSize = new System.Drawing.Size(1500, 800);
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

        private void InitTunnelListHeartBeat(BindingList<Tunnel> tlb)
        {
            Task.Run(() =>
            {
                while (true)
                {
                    Msg msg = new Msg
                    {
                        Flag = "ListTunnel"
                    };
                    tcpHelper.Send(JsonConvert.SerializeObject(msg));
                    Thread.Sleep(2000);
                }
            });
            Action<Msg> action = msg =>
            {
                if (msg.Flag != "ListTunnel") { return; }
                List<Tunnel> tl = JsonConvert.DeserializeObject<List<Tunnel>>(msg.Body);
                this.Invoke((MethodInvoker)delegate
                {
                    tunnelTable.SuspendLayout();
                    foreach (Tunnel tb in tlb)
                    {
                        foreach (Tunnel t in tl)
                        {
                            if (tb.Id == t.Id)
                            {
                                tb.LastAlive = t.LastAlive;
                            }
                        }
                    }
                    tunnelTable.Invalidate();
                    tunnelTable.ResumeLayout();
                });
            };
            tcpHelper.OnMsg(new TcpHelper.MsgHandler(action));
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
            foreach (Tunnel tb in tunnelList)
            {
                if (string.IsNullOrEmpty(tb.Id))
                {
                    tb.Id = Guid.NewGuid().ToString();
                }
            }
            List<Tunnel> tl = tunnelList.ToList();
            /* 正向 */
            List<Tunnel> tlp = tl.FindAll(ele => ele.Direction == 1);
            if (tlp.Count > 0)
            {
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
                    Flag = "NewReverseTunnel",
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

    }
}