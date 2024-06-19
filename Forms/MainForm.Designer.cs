using Newtonsoft.Json;
using sshtunnel.Models;
using System;
using System.ComponentModel;
using System.Windows.Forms;

namespace sshtunnel.Forms
{
    partial class MainForm
    {
        private System.ComponentModel.IContainer components = null;

        /* Panel */
        private System.Windows.Forms.Panel buttonPanel;
        private System.Windows.Forms.Panel tablePanel;

        /* Button */
        private Button directionButton;
        private string tunnelText = "Local To Remote";
        private string reverseTunnelText = "Remote To Local";
        private Button execButton;

        /* Table */
        private DataGridView tunnelTable;
        private BindingList<Tunnel> tunnelList;

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
            buttonPanel = new System.Windows.Forms.FlowLayoutPanel
            {
                Name = "ButtonPanel",
                Dock = DockStyle.Top,
                Height = 50,
            };

            tablePanel = new System.Windows.Forms.Panel
            {
                Name = "TablePanel",
                Dock = DockStyle.Fill,
                AutoScroll = true,
            };

            /* Button */
            directionButton = new Button
            {
                Name = "SwitchButton",
                Text = tunnelText,
                Width = 200,
                Height = 30,
            };
            directionButton.Click += new EventHandler(HandleSwitchButtonClick);
            buttonPanel.Controls.Add(directionButton);

            execButton = new Button
            {
                Name = "ExecButton",
                Text = "Run",
                Width = 100,
                Height = 30,
            };
            execButton.Click += new System.EventHandler(HandleExecButtonClick);
            buttonPanel.Controls.Add(execButton);

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
                }
            };
            tunnelList = new BindingList<Tunnel>();
            tunnelTable.DataSource = tunnelList;
            tablePanel.Controls.Add(tunnelTable);

            AutoScaleDimensions = new System.Drawing.SizeF(9F, 18F);
            AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            ClientSize = new System.Drawing.Size(1200, 600);
            Name = "MainForm";
            Text = "SSH Tunnel";
            Controls.Add(tablePanel);
            Controls.Add(buttonPanel);

            ResumeLayout(false);
        }

        private void HandleSwitchButtonClick(object sender, EventArgs e)
        {
            if (directionButton.Text == tunnelText)
            {
                directionButton.Text = reverseTunnelText;
            }
            else
            {
                directionButton.Text = tunnelText;
            }
        }

        private void HandleExecButtonClick(object sender, EventArgs e)
        {
            if (tunnelList.Count == 0)
            {
                return;
            }
            Msg msg = new Msg
            {
                Flag = directionButton.Text == tunnelText ? "NewTunnel" : "NewReverseTunnel",
                Body = JsonConvert.SerializeObject(tunnelList)
            };
            tcpHelper.Send(JsonConvert.SerializeObject(msg));
        }

    }
}