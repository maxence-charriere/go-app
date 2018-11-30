using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.InteropServices.WindowsRuntime;
using Windows.ApplicationModel.Core;
using Windows.Data.Json;
using Windows.Foundation;
using Windows.Foundation.Collections;
using Windows.UI;
using Windows.UI.Core;
using Windows.UI.ViewManagement;
using Windows.UI.Xaml;
using Windows.UI.Xaml.Controls;
using Windows.UI.Xaml.Controls.Primitives;
using Windows.UI.Xaml.Data;
using Windows.UI.Xaml.Input;
using Windows.UI.Xaml.Media;
using Windows.UI.Xaml.Navigation;

// The Blank Page item template is documented at https://go.microsoft.com/fwlink/?LinkId=234238

namespace uwp
{
    /// <summary>
    /// An empty page that can be used on its own or navigated to within a Frame.
    /// </summary>
    public sealed partial class WindowPage : Page
    {
        private string ID = "";
        private string loadReturnID = "";
        private object locker = new object();
        private bool fullScreen = false;

        public WindowPage()
        {
            InitializeComponent();

            var win = Window.Current;
            var view = ApplicationView.GetForCurrentView();
            var coreView = CoreApplication.GetCurrentView();


            ApplicationViewTitleBar titleBar = view.TitleBar;
            titleBar.ButtonBackgroundColor = Colors.Transparent;
            titleBar.InactiveBackgroundColor = Colors.Transparent;
            titleBar.ButtonInactiveBackgroundColor = Colors.Transparent;

            CoreApplicationViewTitleBar coreTitleBar = coreView.TitleBar;
            coreTitleBar.ExtendViewIntoTitleBar = true;

            this.fullScreen = view.IsFullScreenMode;


            this.Webview.ScriptNotify += this.Webview_ScriptNotify;
            win.SizeChanged += this.OnResized;
            this.Webview.LoadCompleted += this.OnLoad;
            win.Activated += this.OnActivated;
            this.Unloaded += this.OnUnload;
        }

        private void OnUnload(object sender, RoutedEventArgs e)
        {
            var win = Window.Current;
            win.SizeChanged -= this.OnResized;
            this.Webview.LoadCompleted -= this.OnLoad;
            win.Activated -= this.OnActivated;
            this.Unloaded -= this.OnUnload;
        }

        protected override void OnNavigatedTo(NavigationEventArgs e)
        {
            base.OnNavigatedTo(e);

            JsonObject input = e.Parameter as JsonObject;
            this.ID = input.GetNamedString("ID");

            Bridge.PutElem(this.ID, this);
        }

        internal static void Load(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var page = input.GetNamedString("Page");
            var w = Bridge.GetElem<WindowPage>(ID);

            lock (w.locker)
            {
                w.loadReturnID = returnID;
            }

            w.Webview.NavigateToString(page);
        }

        private void OnLoad(object sender, NavigationEventArgs e)
        {
            lock (this.locker)
            {
                var returnID = this.loadReturnID;
                this.loadReturnID = "";
                Bridge.Return(returnID, null, "");
            }
        }

        private void Webview_ScriptNotify(object sender, NotifyEventArgs e)
        {
            var input = new JsonObject();
            input["ID"] = JsonValue.CreateStringValue(this.ID);
            input["Mapping"] = JsonValue.CreateStringValue(e.Value);

            Bridge.GoCall("windows.OnCallback", input, true);
        }

        internal static async void Render(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            var changes = input.GetNamedString("Changes");
            changes = string.Format("render({0})", changes);
            var args = new string[] { changes };
            await w.Webview.InvokeScriptAsync("eval", args);


            Bridge.Return(returnID, null, "");
        }

        internal static void Bounds(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            var window = Window.Current;
            var bounds = window.Bounds;

            var output = new JsonObject();
            output["ReturnID"] = JsonValue.CreateStringValue(returnID);
            output["X"] = JsonValue.CreateNumberValue(bounds.X);
            output["Y"] = JsonValue.CreateNumberValue(bounds.Y);
            output["Width"] = JsonValue.CreateNumberValue(bounds.Width);
            output["Heigth"] = JsonValue.CreateNumberValue(bounds.Height);

            Bridge.Return(returnID, output, "");
        }

        internal static void Resize(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            var x = input.GetNamedNumber("Width");
            var y = input.GetNamedNumber("Height");
            var size = new Size(x, y);

            ApplicationView.GetForCurrentView().TryResizeView(size);
            Bridge.Return(returnID, null, "");
        }

        private void OnResized(object sender, WindowSizeChangedEventArgs e)
        {
            var input = new JsonObject();
            input["Width"] = JsonValue.CreateNumberValue(e.Size.Width);
            input["Heigth"] = JsonValue.CreateNumberValue(e.Size.Height);
            Bridge.GoCall("windows.OnResize", input, true);

            var view = ApplicationView.GetForCurrentView();

            lock(this.locker)
            {
                var fullScreen = view.IsFullScreenMode;
                if (fullScreen == this.fullScreen)
                {
                    return;
                }
                this.fullScreen = fullScreen;

                if (fullScreen)
                {
                    Bridge.GoCall("windows.OnFullScreen", input, true);
                }
                else
                {
                    Bridge.GoCall("windows.OnExitFullScreen", input, true);
                }
            }
        }

        internal static void Focus(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);
            Window.Current.Activate();
        }

        private void OnActivated(object sender, WindowActivatedEventArgs e)
        {
            switch (e.WindowActivationState)
            {
                case CoreWindowActivationState.CodeActivated:
                case CoreWindowActivationState.PointerActivated:
                    Bridge.GoCall("windows.OnFocus", null, true);
                    break;

                case CoreWindowActivationState.Deactivated:
                    Bridge.GoCall("windows.OnBlur", null, true);
                    break;

                default:
                    throw new Exception(string.Format("unkown activation state: {0}", e.WindowActivationState));
            }
        }

        internal static void FullScreen(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            ApplicationView.GetForCurrentView().TryEnterFullScreenMode();
            Bridge.Return(returnID, null, "");
        }

        internal static void ExitFullScreen(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            ApplicationView.GetForCurrentView().ExitFullScreenMode();
            Bridge.Return(returnID, null, "");
        }
    }
}
