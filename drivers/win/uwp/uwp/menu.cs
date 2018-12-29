using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using Windows.Data.Json;
using Windows.UI.Core;
using Windows.UI.Xaml;
using Windows.UI.Xaml.Controls;

namespace uwp
{
    public class Menu
    {
        public string ID { get; set; }
        Dictionary<string, object> Nodes { get; set; }
        public CompoNode Root { get; set; }


        public Menu(string ID)
        {
            this.ID = ID;
            this.Nodes = new Dictionary<string, object>();
        }

        public static void New(JsonObject input, string returnID)
        {
            var menu = new Menu(input.GetNamedString("ID"));
            Bridge.PutElem(menu.ID, menu);
            Bridge.Return(returnID, null, null);
        }

        public static void Load(JsonObject input, string returnID)
        {
            var menu = Bridge.GetElem<Menu>(input.GetNamedString("ID"));
            menu.Root = null;
            Bridge.Return(returnID, null, null);
        }

        public static async void Render(JsonObject input, string returnID)
        {
            await Window.Current.Dispatcher.RunAsync(CoreDispatcherPriority.Normal, () =>
            {
                try
                {
                    var menu = Bridge.GetElem<Menu>(input.GetNamedString("ID"));
                    var changes = JsonArray.Parse(input.GetNamedString("Changes"));

                    foreach (var c in changes)
                    {
                        var change = c.GetObject();
                        var action = change.GetNamedNumber("Action");

                        switch (action)
                        {
                            case 0:
                                menu.setRoot(change);
                                break;

                            case 1:
                                menu.newNode(change);
                                break;

                            case 2:
                                menu.delNode(change);
                                break;

                            case 3:
                                menu.setAttr(change);
                                break;

                            case 4:
                                menu.delAttr(change);
                                break;

                            case 6:
                                menu.appendChild(change);
                                break;

                            default:
                                throw new Exception(string.Format("{0} change is not supported", action));
                        }
                    }

                    Bridge.Return(returnID, null, null);
                }
                catch (Exception e)
                {
                    Bridge.Return(returnID, null, e.Message);
                }
            });
        }

        public void setRoot(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var c = this.Nodes[nodeID] as CompoNode;
            c.isRootCompo = true;
            this.Root = c;
        }

        public void newNode(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var compoID = change.GetNamedString("CompoID", "");
            var type = change.GetNamedString("Type");
            var isCompo = change.GetNamedBoolean("IsCompo", false);

            if (isCompo)
            {
                var c = new CompoNode()
                {
                    ID = nodeID,
                    type = type,
                    isRootCompo = false,
                };

                this.Nodes[nodeID] = c;
                return;
            }

            if (type == "menu")
            {
                var container = new MenuContainer()
                {
                    ID = nodeID,
                    compoID = compoID,
                    elemID = this.ID,
                    item = new MenuFlyoutSubItem(),
                };

                this.Nodes[nodeID] = container;
                return;
            }

            if (type == "menuitem")
            {
                var item = new MenuItem()
                {
                    ID = nodeID,
                    compoID = compoID,
                    elemID = this.ID,
                    item = new MenuFlyoutItem(),
                };

                this.Nodes[nodeID] = item;
                return;
            }

            throw new Exception(string.Format("menu does not support {0} tag", type));
        }

        public void delNode(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            this.Nodes.Remove(nodeID);
        }

        public void setAttr(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var key = change.GetNamedString("Key");
            var value = change.GetNamedString("Value", "");

            var node = this.Nodes[nodeID] as IMenuWithAttr;
            node.setAttr(key, value);
        }

        public void delAttr(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var key = change.GetNamedString("Key");

            var node = this.Nodes[nodeID] as IMenuWithAttr;
            node.delAttr(key);
        }

        public void appendChild(JsonObject change)
        {
            var nodeID = change.GetNamedString("NodeID");
            var childID = change.GetNamedString("ChildID");
            var node = this.Nodes[nodeID];

            if (node is CompoNode)
            {
                var cnode = node as CompoNode;
                cnode.rootID = childID;
                return;
            }

            var parent = node as MenuContainer;
            var child = this.Nodes[childID];
            var childRoot = this.CompoRoot(child);
            parent.appendChild(childRoot);
        }

        public object CompoRoot(object node)
        {
            if (node == null || !(node is CompoNode))
            {
                return node;
            }

            var c = node as CompoNode;
            return this.CompoRoot(this.Nodes[c.rootID]);
        }
    }

    public class CompoNode
    {
        public string ID;
        public string rootID;
        public string type;
        public bool isRootCompo;
    }

    public interface IMenuWithAttr
    {
        void setAttr(string key, string value);
        void delAttr(string key);
    }

    public class MenuContainer : IMenuWithAttr
    {
        public string ID { get; set; }
        public string compoID { get; set; }
        public string elemID { get; set; }
        public MenuFlyoutSubItem item { get; set; }

        public void setAttr(string key, string value)
        {
            switch (key)
            {
                case "label":
                    item.Text = value;
                    break;
            }
        }

        public void delAttr(string key)
        {
            switch (key)
            {
                case "label":
                    item.Text = "";
                    break;
            }
        }

        public void appendChild(object child)
        {
            if (child is MenuContainer)
            {
                var container = child as MenuContainer;
                this.item.Items.Add(container.item);
                return;
            }

            if (child is MenuItem)
            {
                var item = child as MenuItem;
                this.item.Items.Add(item.item);
                return;
            }

            throw new Exception("unknow child node type: " + child.GetType().ToString());
        }
    }

    public class MenuItem : IMenuWithAttr
    {
        public string ID { get; set; }
        public string compoID { get; set; }
        public string elemID { get; set; }
        public MenuFlyoutItem item { get; set; }

        public void setAttr(string key, string value)
        {
            switch (key)
            {
                case "label":
                    item.Text = value;
                    break;
            }
        }

        public void delAttr(string key)
        {
            switch (key)
            {
                case "label":
                    item.Text = "";
                    break;
            }
        }
    }
}
