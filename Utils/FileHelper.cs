using System.IO;

namespace sshtunnel.Utils
{
    public class FileHelper
    {

        private string path;

        private FileHelper() { }

        public static FileHelper New(string p)
        {
            using (File.AppendText(p)) { };
            return new FileHelper
            {
                path = p
            };
        }

        public void W(string content)
        {
            File.WriteAllText(path, content);
        }

        public string R()
        {
            return File.ReadAllText(path);
        }

    }
}
