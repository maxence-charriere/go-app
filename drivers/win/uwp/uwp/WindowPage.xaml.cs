using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Runtime.InteropServices.WindowsRuntime;
using System.Threading.Tasks;
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
                try
                {
                    var frame = new Frame();
                    frame.Navigate(typeof(WindowPage), input);

                    Window.Current.Content = frame;
                    Window.Current.Activate();

                    viewID = ApplicationView.GetForCurrentView().Id;
                    setupWindow(input);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });

            await currentWindow.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, async () =>
            {
                try
                {
                    await ApplicationViewSwitcher.TryShowAsStandaloneAsync(viewID);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });

            Bridge.Return(returnID, null, null);
        }

        static void setupWindow(JsonObject input)
        {
            var win = Window.Current;
            win.SizeChanged += OnResized;
            win.Activated += OnActivated;

            SystemNavigationManagerPreview.GetForCurrentView().CloseRequested += OnClose;

            var view = ApplicationView.GetForCurrentView();
            var coreView = CoreApplication.GetCurrentView();

            ApplicationViewTitleBar titleBar = view.TitleBar;
            titleBar.ButtonBackgroundColor = Colors.Transparent;
            titleBar.InactiveBackgroundColor = Colors.Transparent;
            titleBar.ButtonInactiveBackgroundColor = Colors.Transparent;


            CoreApplicationViewTitleBar coreTitleBar = coreView.TitleBar;
            coreTitleBar.ExtendViewIntoTitleBar = true;
        }

        static async void OnClose(object sender, SystemNavigationCloseRequestedPreviewEventArgs e)
        {
            var deferral = e.GetDeferral();

            var frame = Window.Current.Content as Frame;
            if (frame != null)
            {
                var win = frame.Content as WindowPage;

                var input = new JsonObject();
                input["ID"] = JsonValue.CreateStringValue(win.ID);

                await Bridge.GoCall("windows.OnClose", input, true);
                Bridge.DeleteElem(win.ID);
            }

            deferral.Complete();
        }

        public WindowPage()
        {
            InitializeComponent();

            this.fullScreen = ApplicationView.GetForCurrentView().IsFullScreenMode;
            this.Webview.ScriptNotify += this.Webview_ScriptNotify;
            this.Webview.LoadCompleted += this.OnLoad;
            this.Webview.NavigationStarting += OnNavigationStart;
            this.Webview.UnsupportedUriSchemeIdentified += OnUnsupportedUriSchemeIdentified;
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
                try
                {
                    w.loadReturnID = returnID;
                    w.Webview.NavigateToString(page);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        void OnLoad(object sender, NavigationEventArgs e)
        {
            var returnID = this.loadReturnID;
            this.loadReturnID = "";
            Bridge.Return(returnID, null, null);
        }

       async void OnUnsupportedUriSchemeIdentified(WebView sender, WebViewUnsupportedUriSchemeIdentifiedEventArgs args)
        {
            if (args.Uri != null)
            {
                args.Handled = true;
                var input = new JsonObject();
                input["ID"] = JsonValue.CreateStringValue(this.ID);
                input["URL"] = JsonValue.CreateStringValue(args.Uri.ToString());

                await Bridge.GoCall("windows.OnNavigate", input, true);
            }
        }

        async void OnNavigationStart(WebView sender, WebViewNavigationStartingEventArgs args)
        {
            if (args.Uri != null)
            {
                var input = new JsonObject();
                input["ID"] = JsonValue.CreateStringValue(this.ID);
                input["URL"] = JsonValue.CreateStringValue(args.Uri.ToString());

                await Bridge.GoCall("windows.OnNavigate", input, true);
                args.Cancel = true;
            }
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
                try
                {
                    await w.Webview.InvokeScriptAsync("eval", args);
                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        internal static async void Bounds(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                try
                {
                    var bounds = Window.Current.Bounds;

                    var output = new JsonObject();
                    output["ReturnID"] = JsonValue.CreateStringValue(returnID);
                    output["X"] = JsonValue.CreateNumberValue(bounds.X);
                    output["Y"] = JsonValue.CreateNumberValue(bounds.Y);
                    output["Width"] = JsonValue.CreateNumberValue(bounds.Width);
                    output["Heigth"] = JsonValue.CreateNumberValue(bounds.Height);

                    Bridge.Return(returnID, output, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
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
                try
                {
                    ApplicationView.GetForCurrentView().TryResizeView(size);
                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        static async void OnResized(object sender, WindowSizeChangedEventArgs e)
        {
            var frame = Window.Current.Content as Frame;
            var win = frame.Content as WindowPage;
            
            var input = new JsonObject();
            input["ID"] = JsonValue.CreateStringValue(win.ID);
            input["Width"] = JsonValue.CreateNumberValue(e.Size.Width);
            input["Heigth"] = JsonValue.CreateNumberValue(e.Size.Height);
            await Bridge.GoCall("windows.OnResize", input, true);

            var view = ApplicationView.GetForCurrentView();
            var fullScreen = view.IsFullScreenMode;

            if (fullScreen == win.fullScreen)
            {
                return;
            }

            win.fullScreen = fullScreen;

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
                try
                {
                    Window.Current.Activate();
                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        static async void OnActivated(object sender, WindowActivatedEventArgs e)
        {
            var frame = Window.Current.Content as Frame;
            var win = frame.Content as WindowPage;

            var input = new JsonObject();
            input["ID"] = JsonValue.CreateStringValue(win.ID);

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
                    Bridge.Log("unkown activation state: {0}", e.WindowActivationState);
                    break;
            }
        }

        internal static async void FullScreen(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                try
                {
                    ApplicationView.GetForCurrentView().TryEnterFullScreenMode();
                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        internal static async void ExitFullScreen(JsonObject input, string returnID)
        {
            var ID = input.GetNamedString("ID");
            var w = Bridge.GetElem<WindowPage>(ID);

            await w.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                try
                {
                    ApplicationView.GetForCurrentView().ExitFullScreenMode();
                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }
    }
}
