using System;
using Windows.UI.Xaml.Media;

namespace uwp
{
    public class color
    {
        public static SolidColorBrush GetSolidColorBrush(string c)
        {
            c = c.Replace("#", string.Empty);
            byte r = (byte)(Convert.ToUInt32(c.Substring(0, 2), 16));
            byte g = (byte)(Convert.ToUInt32(c.Substring(2, 2), 16));
            byte b = (byte)(Convert.ToUInt32(c.Substring(4, 2), 16));

            return new SolidColorBrush(Windows.UI.Color.FromArgb(255, r, g, b));
        }
    }
}
