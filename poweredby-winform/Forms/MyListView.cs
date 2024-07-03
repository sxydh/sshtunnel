using System.Windows.Forms;

namespace poweredby_winform.Forms
{
    public class MyListView : ListView
    {

        public MyListView()
        {
            DoubleBuffered = true;
        }

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
