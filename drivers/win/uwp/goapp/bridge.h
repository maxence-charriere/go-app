#pragma once

#include <cstddef>

using namespace Platform;
using namespace Windows::Foundation;
using namespace Windows::ApplicationModel::AppService;

typedef void (*funcWinReturn)(char *retID, char *ret, char *err);
typedef void (*funcGoCall)(char *call);

IAsyncAction ^ BridgeConnectAsync();
void Bridge_WinCallReturn(String ^ retID, String ^ ret, String ^ err);
void Bridge_GoCall(String ^ method, String ^ input);

void BridgeRequestReceived(AppServiceConnection ^ connection,
                           AppServiceRequestReceivedEventArgs ^ args);
void BridgeClosed(AppServiceConnection ^ connection,
                  AppServiceClosedEventArgs ^ args);

char *cString(String ^ s);
String ^ winString(char *str);

extern "C"
{
    __declspec(dllexport) void Bridge_Init();
    __declspec(dllexport) void Bridge_Call(char *call);
    __declspec(dllexport) void Bridge_SetWinCallReturn(funcWinReturn f);
    __declspec(dllexport) void Bridge_SetGoCall(funcGoCall f);
}
