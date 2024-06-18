using System;
using System.Data;
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
        private Button switchButton;
        private Button execButton;

        /* Table */
        private DataGridView configGrid;
        // Local to remote config list
        private DataTable configList;
        // Remote to local config list
        private DataTable reverseConfigList;

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
            this.SuspendLayout();

            /* Panel */
            this.buttonPanel = new System.Windows.Forms.FlowLayoutPanel
            {
                Name = "ButtonPanel",
                Dock = DockStyle.Top,
                Height = 50,
            };
            this.tablePanel = new System.Windows.Forms.Panel
            {
                Name = "TablePanel",
                Dock = DockStyle.Fill,
                AutoScroll = true,
            };

            /* Button */
            this.switchButton = new Button
            {
                Name = "SwitchButton",
                Text = "点我切换方向",
                Width = 200,
                Height = 30,
            };
            this.buttonPanel.Controls.Add(this.switchButton);
            this.execButton = new Button
            {
                Name = "ExecButton",
                Text = "执行",
                Width = 100,
                Height = 30,
            };
            this.execButton.Click += new System.EventHandler(HandleExecButtonClick);
            this.buttonPanel.Controls.Add(this.execButton);

            /* Table */
            this.configGrid = new DataGridView
            {
                Name = "ConfigGrid",
                Dock = DockStyle.Fill,
                AutoSizeColumnsMode = DataGridViewAutoSizeColumnsMode.Fill,
            };
            
            configGrid.DataSource = BuildConfigList();
            this.tablePanel.Controls.Add(configGrid);

            this.AutoScaleDimensions = new System.Drawing.SizeF(9F, 18F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(1000, 600);
            this.Name = "MainForm";
            this.Text = "SSH Tunnel";
            this.Controls.Add(this.tablePanel);
            this.Controls.Add(this.buttonPanel);

            this.ResumeLayout(false);
        }

        private void HandleExecButtonClick(object sender, EventArgs e)
        {
        }

        private DataTable BuildConfigList()
        {
            DataTable dataTable = new DataTable();
            dataTable.Columns.Add("Local Port", typeof(string));
            dataTable.Columns.Add("Target IP", typeof(string));
            dataTable.Columns.Add("Target Port", typeof(string));
            this.configList = dataTable;
            return dataTable;
        }

        private DataTable BuildReverseConfigList()
        {
            DataTable dataTable = new DataTable();
            dataTable.Columns.Add("Remote Port", typeof(string));
            dataTable.Columns.Add("Target IP", typeof(string));
            dataTable.Columns.Add("Target Port", typeof(string));
            this.reverseConfigList = dataTable;
            return dataTable;
        }

    }
}