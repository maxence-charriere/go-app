using System;
using System.Collections.Generic;
using System.Threading.Tasks;
using Windows.ApplicationModel;
using Windows.ApplicationModel.AppService;
using Windows.ApplicationModel.Background;
using Windows.ApplicationModel.Core;
using Windows.Data.Json;
using Windows.Foundation.Collections;
using Windows.UI.Core;

namespace uwp
{
    class Bridge
    {
        static AppServiceConnection conn = null;
        static bool launched = false;
        static bool connected = false;
        static object locker = new object();
        static BackgroundTaskDeferral deferral = null;
        static Dictionary<string, Action<JsonObject, string>> handlers = new Dictionary<string, Action<JsonObject, string>>();
        static Dictionary<string, object> elems = new Dictionary<string, object>();
        static Queue<deferredGoCall> deferredGoCalls = new Queue<deferredGoCall>();

        public static async void TryLaunchGoApp()
        {
            if (launched)
            {
                return;
            }

            await FullTrustProcessLauncher.LaunchFullTrustProcessForCurrentAppAsync();
            launched = true;
        }

        public static async void NewConn(IBackgroundTaskInstance task)
        {
            AppServiceTriggerDetails appService = task.TriggerDetails as AppServiceTriggerDetails;
            if (appService == null)
            {
                return;
            }

            deferral = task.GetDeferral();
            task.Canceled += Task_Canceled;

            conn = appService.AppServiceConnection;
            conn.RequestReceived += Conn_RequestReceived;
            conn.ServiceClosed += Conn_ServiceClosed;

            lock (locker)
            {
                connected = true;
            }

            while (deferredGoCalls.Count != 0)
            {
                var call = deferredGoCalls.Dequeue();
                await GoCall(call.Method, call.Input, call.UI);
            }
        }

        private static async void Conn_RequestReceived(AppServiceConnection sender, AppServiceRequestReceivedEventArgs args)
        {
            AppServiceDeferral msgDeferral = args.GetDeferral();

            await CoreApplication.MainView.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                string value = args.Request.Message["Value"].ToString();
                var req = JsonObject.Parse(value);
                string returnID = req.GetNamedValue("ReturnID").GetString();

                try
                {
                    string method = req.GetNamedValue("Method").GetString();
                    JsonObject input = req.GetNamedObject("Input");

                    var handler = handlers[method];
                    if (handler == null)
                    {
                        throw new Exception(string.Format("{0} is not handled", method));
                    }

                    handler(input, returnID);
                }
                catch (Exception e)
                {
                    Return(returnID, null, e.Message);
                }
                finally
                {
                    msgDeferral.Complete();
                }
            });
        }

        private static void Conn_ServiceClosed(AppServiceConnection sender, AppServiceClosedEventArgs args)
        {
            deferral.Complete();
        }

        private static void Task_Canceled(IBackgroundTaskInstance sender, BackgroundTaskCancellationReason reason)
        {
            deferral.Complete();
        }

        public static void Handle(string method, Action<JsonObject, string> handler)
        {
            handlers[method] = handler;
        }

        public static async void Return(string returnID, JsonObject input, string err)
        {
            var data = new ValueSet();
            data["Operation"] = "Return";
            data["ReturnID"] = returnID;
            data["Input"] = "";
            data["Err"] = err;

            if (input != null)
            {
                data["Input"] = input.ToString();
            }

            await conn.SendMessageAsync(data);
        }

        public static async Task<JsonObject> GoCall(string method, JsonObject input, bool ui)
        {
            lock (locker)
            {
                if (!connected)
                {
                    deferredGoCalls.Enqueue(new deferredGoCall() {
                        Method = method,
                        Input = input,
                        UI = ui,
                    });

                    return null;
                }
            }


            var msg = new ValueSet();
            msg["Operation"] = "Call";
            msg["Method"] = method;
            msg["Input"] = "";
            msg["Ui"] = ui.ToString();

            if (input != null)
            {
                msg["Input"] = input.ToString();
            }

            var res = await conn.SendMessageAsync(msg);
            if (res.Status != AppServiceResponseStatus.Success)
            {
                return null;
            }

            var resMsg = res.Message;
            if (!resMsg.ContainsKey("Value")) {
                return null;
            }

            return JsonObject.Parse(resMsg["Value"] as string);
        }

        public static async void Log(string format, params object[] v)
        {
            var msg = string.Format(format, v);
            JsonObject input = new JsonObject();
            input["Msg"] = JsonValue.CreateStringValue(msg);
            await GoCall("driver.Log", input, false);
        }

        public static void PutElem(string ID, object elem)
        {
            lock (locker)
            {
                elems.Add(ID, elem);
            }
        }

        public static void DeleteElem(string ID)
        {
            lock (locker)
            {
                elems.Remove(ID);
            }
        }

        public static T GetElem<T>(string ID) where T : class
        {
            lock (locker)
            {
                var elem = elems[ID];
                if (elem == null)
                {
                    throw new Exception(string.Format("elem {0} is not found", ID));
                }

                var tElem = elem as T;



                if (!(elem is T))
                {
                    throw new Exception(string.Format("elem {0} is not a {1}", ID, elem.GetType().ToString()));
                }

                return tElem;
            }
        }
    }

    struct deferredGoCall
    {
        public string Method;
        public JsonObject Input;
        public bool UI;
    }
}
