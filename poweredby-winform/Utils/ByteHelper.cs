namespace poweredby_winform.Utils
{
    public class ByteHelper
    {

        public static byte[] Bytes4(int i)
        {
            byte[] bytes = new byte[4];
            bytes[0] = (byte)((i >> 24) & 0xFF);
            bytes[1] = (byte)((i >> 16) & 0xFF);
            bytes[2] = (byte)((i >> 8) & 0xFF);
            bytes[3] = (byte)(i & 0xFF);
            return bytes;
        }

        public static int Byte4ToInt(byte[] bytes)
        {
            return (bytes[0] << 24) | (bytes[1] << 16) | (bytes[2] << 8) | bytes[3];
        }

    }
}
