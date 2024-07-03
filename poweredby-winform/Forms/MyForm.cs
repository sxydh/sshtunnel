using System.Windows.Forms;

namespace poweredby_winform.Forms
{
    public class MyForm : Form
    {

        protected override CreateParams CreateParams
        {
            get
            {
                CreateParams cp = base.CreateParams;
                cp.ExStyle |= 0x02000000;
                return cp;
            }
        }

    }
}
