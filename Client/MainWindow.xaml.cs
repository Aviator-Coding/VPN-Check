using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Navigation;
using System.Windows.Shapes;

namespace VPN_Check
{
    /// <summary>
    /// Interaction logic for MainWindow.xaml
    /// </summary>
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();

        }


        private void CheckButton_OnClick(object sender, RoutedEventArgs e)
        {
            try
            {
                var rsp = Task.Run(() => VPNCheck.CheckVpn());
                rsp.Wait();

                if (rsp.Result != null)
                {
                    Ip.Content = rsp.Result.Data.IP;
                    countryCode.Content = rsp.Result.Data.CountryCode;
                    countryName.Content = rsp.Result.Data.CountryName;
                    Asn.Content = rsp.Result.Data.Asn;
                    isp.Content = rsp.Result.Data.Isp;
                    block.Content = rsp.Result.Data.Block;
                    if (!rsp.Result.Data.Block.isProxy())
                    {
                        App.Background = Brushes.Crimson;
                    }
                    else App.Background = Brushes.Lime;
                }
            }
            catch (Exception exception)
            {
                // ignored
            }
        }
    }
}