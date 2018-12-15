using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.InteropServices.WindowsRuntime;
using Windows.ApplicationModel.Core;
using Windows.Data.Json;
using Windows.Foundation;
using Windows.Foundation.Collections;
using Windows.Foundation.Metadata;
using Windows.UI;
using Windows.UI.Core;
using Windows.UI.Core.Preview;
using Windows.UI.ViewManagement;
using Windows.UI.Xaml;
using Windows.UI.Xaml.Controls;
using Windows.UI.Xaml.Controls.Primitives;
using Windows.UI.Xaml.Data;
using Windows.UI.Xaml.Input;
using Windows.UI.Xaml.Media;
using Windows.UI.Xaml.Navigation;

namespace uwp
{
    /// <summary>
    /// The page that describes a window content.
    /// </summary>
    public sealed partial class WindowPage : Page
    {
        string ID = "";
        string loadReturnID = "";
        bool fullScreen = false;

        static Window currentWindow = null;

        internal static async void NewWindow(JsonObject input, string returnID)
        {
            CoreApplicationView view = null;
            var viewID = 0;

            if (currentWindow == null)
            {
                view = CoreApplication.MainView;
                currentWindow = Window.Current;
            }
            else
            {
                view = CoreApplication.CreateNewView();
            }

            await view.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                var frame = new Frame();
                frame.NavigationFailed += OnNavigationFailed;
                frame.Navigate(typeof(WindowPage), input);

                Window.Current.Content = frame;
                Window.Current.Activate();

                viewID = ApplicationView.GetForCurrentView().Id;
                SystemNavigationManagerPreview.GetForCurrentView().CloseRequested += OnClose;
            });

            await currentWindow.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, async () =>
            {
                await ApplicationViewSwitcher.TryShowAsStandaloneAsync(viewID);
            });

            Bridge.Return(returnID, null, null);
        }

        static void OnClose(object sender, SystemNavigationCloseRequestedPreviewEventArgs e)
        {
            SystemNavigationManagerPreview.GetForCurrentView().CloseRequested -= OnClose;
        }

        static void OnNavigationFailed(object sender, NavigationFailedEventArgs e)
        {
            throw new Exception("Failed to load Page " + e.SourcePageType.FullName);
        }

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



        void OnUnload(object sender, RoutedEventArgs e)
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
            var frosted = input.GetNamedBoolean("FrostedBackground");


            Color bg = Color.FromArgb(255, 50, 52, 54);

            if (Application.Current.RequestedTheme == ApplicationTheme.Light)
            {
                bg = Color.FromArgb(255, 236, 236, 236);
            }

            if (frosted && ApiInformation.IsTypePresent("Windows.UI.Xaml.Media.XamlCompositionBrushBase"))
            {
                AcrylicBrush frostedBrush = new AcrylicBrush();
                frostedBrush.BackgroundSource = Windows.UI.Xaml.Media.AcrylicBackgroundSource.HostBackdrop;
                frostedBrush.TintColor = bg;
                frostedBrush.FallbackColor = bg;
                frostedBrush.TintOpacity = 0.80;
                this.Root.Background = frostedBrush;
            }
            else
            {
                this.Root.Background = new SolidColorBrush(bg);
            }

            Bridge.PutElem(this.ID, this);
        }

        internal static async void Load(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var page = input.GetNamedString("Page");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                w.loadReturnID = returnID;
                w.Webview.NavigateToString(page);
            });
        }

        void OnLoad(object sender, NavigationEventArgs e)
        {
            var returnID = this.loadReturnID;
            this.loadReturnID = "";
            Bridge.Return(returnID, null, null);
        }

        async void Webview_ScriptNotify(object sender, NotifyEventArgs e)
        {
            var input = new JsonObject();
            input["ID"] = JsonValue.CreateStringValue(this.ID);
            input["Mapping"] = JsonValue.CreateStringValue(e.Value);

            await Bridge.GoCall("windows.OnCallback", input, true);
        }

        internal static async void Render(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            var changes = input.GetNamedString("Changes");
            changes = string.Format("render({0})", changes);
            var args = new string[] { changes };

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, async () =>
            {
                await w.Webview.InvokeScriptAsync("eval", args);
                Bridge.Return(returnID, null, null);
            });
        }

        internal static async void Bounds(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                var bounds = Window.Current.Bounds;

                var output = new JsonObject();
                output["ReturnID"] = JsonValue.CreateStringValue(returnID);
                output["X"] = JsonValue.CreateNumberValue(bounds.X);
                output["Y"] = JsonValue.CreateNumberValue(bounds.Y);
                output["Width"] = JsonValue.CreateNumberValue(bounds.Width);
                output["Heigth"] = JsonValue.CreateNumberValue(bounds.Height);

                Bridge.Return(returnID, output, null);
            });
        }

        internal static async void Resize(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            var x = input.GetNamedNumber("Width");
            var y = input.GetNamedNumber("Height");
            var size = new Size(x, y);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                ApplicationView.GetForCurrentView().TryResizeView(size);
                Bridge.Return(returnID, null, null);
            });
        }

        async void OnResized(object sender, WindowSizeChangedEventArgs e)
        {
            var input = new JsonObject();
            input["ID"] = JsonValue.CreateStringValue(this.ID);
            input["Width"] = JsonValue.CreateNumberValue(e.Size.Width);
            input["Heigth"] = JsonValue.CreateNumberValue(e.Size.Height);
            await Bridge.GoCall("windows.OnResize", input, true);

            var view = ApplicationView.GetForCurrentView();
            var fullScreen = view.IsFullScreenMode;

            if (fullScreen == this.fullScreen)
            {
                return;
            }

            this.fullScreen = fullScreen;

            if (fullScreen)
            {
                await Bridge.GoCall("windows.OnFullScreen", input, true);
            }
            else
            {
                await Bridge.GoCall("windows.OnExitFullScreen", input, true);
            }
        }

        internal static async void Focus(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                Window.Current.Activate();
                Bridge.Return(returnID, null, null);
            });
        }

        async void OnActivated(object sender, WindowActivatedEventArgs e)
        {
            var input = new JsonObject();
            input["ID"] = JsonValue.CreateStringValue(this.ID);

            switch (e.WindowActivationState)
            {
                case CoreWindowActivationState.CodeActivated:
                case CoreWindowActivationState.PointerActivated:
                    var w = Window.Current;

                    await CoreApplication.MainView.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
                    {
                        currentWindow = w;
                    });

                    await Bridge.GoCall("windows.OnFocus", input, true);
                    break;

                case CoreWindowActivationState.Deactivated:
                    await Bridge.GoCall("windows.OnBlur", input, true);
                    break;

                default:
                    throw new Exception(string.Format("unkown activation state: {0}", e.WindowActivationState));
            }
        }

        internal static async void FullScreen(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                ApplicationView.GetForCurrentView().TryEnterFullScreenMode();
                Bridge.Return(returnID, null, null);
            });
        }

        internal static async void ExitFullScreen(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                ApplicationView.GetForCurrentView().ExitFullScreenMode();
                Bridge.Return(returnID, null, null);
            });
        }
    }
}
