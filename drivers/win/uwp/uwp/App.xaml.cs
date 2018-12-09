using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Runtime.InteropServices.WindowsRuntime;
using Windows.ApplicationModel;
using Windows.ApplicationModel.Activation;
using Windows.ApplicationModel.AppService;
using Windows.ApplicationModel.Background;
using Windows.ApplicationModel.Core;
using Windows.Data.Json;
using Windows.Foundation;
using Windows.Foundation.Collections;
using Windows.Storage;
using Windows.UI.Core;
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
    /// Provides application-specific behavior to supplement the default Application class.
    /// </summary>
    sealed partial class App : Application
    {
        /// <summary>
        /// Initializes the singleton application object.  This is the first line of authored code
        /// executed, and as such is the logical equivalent of main() or WinMain().
        /// </summary>
        public App()
        {
            this.InitializeComponent();
            this.Suspending += OnSuspending;

            Bridge.Handle("windows.New", this.newWindow);
            Bridge.Handle("windows.Load", WindowPage.Load);
            Bridge.Handle("windows.Render", WindowPage.Render);
            Bridge.Handle("windows.Position", WindowPage.Bounds);
            Bridge.Handle("windows.Size", WindowPage.Bounds);
            Bridge.Handle("windows.Resize", WindowPage.Resize);
            Bridge.Handle("windows.Focus", WindowPage.Focus);
            Bridge.Handle("windows.FullScreen", WindowPage.FullScreen);
            Bridge.Handle("windows.ExitFullScreen", WindowPage.ExitFullScreen);
        }

        protected override void OnBackgroundActivated(BackgroundActivatedEventArgs args)
        {
            base.OnBackgroundActivated(args);
            Bridge.NewConn(args.TaskInstance);
        }

        /// <summary>
        /// Invoked when the application is launched normally by the end user.  Other entry points
        /// will be used suc h as when the application is launched to open a specific file.
        /// </summary>
        /// <param name = "e">Details about the launch request and process.</param>
        protected override void OnLaunched(LaunchActivatedEventArgs e)
        {
            Bridge.TryLaunchGoApp();
        }

        protected override async void OnActivated(IActivatedEventArgs e)
        {
            base.OnActivated(e);
            Bridge.TryLaunchGoApp();
          

            switch (e.Kind)
            {
                case ActivationKind.Protocol:
                    var pe = e as ProtocolActivatedEventArgs;

                    var input = new JsonObject();
                    input["URL"] = JsonValue.CreateStringValue(pe.Uri.ToString());

                    await Bridge.GoCall("driver.OnURLOpen", input, true);
                    break;
            }
        }

        protected override async void OnFileActivated(FileActivatedEventArgs args)
        {
            base.OnFileActivated(args);
            Bridge.TryLaunchGoApp();

            var filenames = new JsonArray();

            foreach (var f in args.Files)
            {
                filenames.Add(JsonValue.CreateStringValue(f.Path));
            }

            var input = new JsonObject();
            input["Filenames"] = filenames;

            await Bridge.GoCall("driver.OnFilesOpen", input, true);
        }

        /// <summary>
        /// Invoked when Navigation to a certain page fails
        /// </summary>
        /// <param name="sender">The Frame which failed navigation</param>
        /// <param name="e">Details about the navigation failure</param>
        void OnNavigationFailed(object sender, NavigationFailedEventArgs e)
        {
            Bridge.Log("uwp => navigation failed: {0}", e.SourcePageType.FullName);
            throw new Exception("Failed to load Page " + e.SourcePageType.FullName);
        }

        /// <summary>
        /// Invoked when application execution is being suspended.  Application state is saved
        /// without knowing  whether the application will be terminated or resumed with the contents
        /// of memory still intact.
        /// </summary>
        /// <param name="sender">The source of the suspend request.</param>
        /// <param name="e">Details about the suspend request.</param>
        private void OnSuspending(object sender, SuspendingEventArgs e)
        {
            var deferral = e.SuspendingOperation.GetDeferral();
            //TODO: Save application state and stop any background activity
            deferral.Complete();
        }

        private void newWindow(JsonObject input, string returnID)
        {
            Frame frame = Window.Current.Content as Frame;
            if (frame == null)
            {
                frame = new Frame();

                frame.NavigationFailed += OnNavigationFailed;
                Window.Current.Content = frame;
            }


            frame.Navigate(typeof(WindowPage), input);
            Window.Current.Activate();

            Bridge.Return(returnID, null, "");
            return;

            //CoreApplicationView view = CoreApplication.CreateNewView();
            //int viewID = 0;

            //await view.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            //{
            //    Frame frame = new Frame();
            //    frame.Navigate(typeof(WindowPage));
            //    Window.Current.Content = frame;
            //    Window.Current.Activate();
            //    viewID = ApplicationView.GetForCurrentView().Id;
            //});

            //bool viewShown = await ApplicationViewSwitcher.TryShowAsStandaloneAsync(viewID);
        }
    }
}