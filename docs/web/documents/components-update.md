## Components update

The component now displays the username in its title and provides input for the user to type his/her name. When the user does so, an event handler is called and the name is stored in the component field named `name`.

The **[Update()](/reference#Composer) method call is what tells the component that its state changed and that its appearance must be updated**.

It internally triggers the `Render()` method and performs a diff with the current component state in order to define and process the changes. Here is how rendering diff behave:

| Diff                                                       | Modification                              |
| ---------------------------------------------------------- | ----------------------------------------- |
| Different types of nodes (Text, HTML element or Component) | Current node is replaced                  |
| Different texts                                            | Current node text value is updated        |
| Different HTML elements                                    | Current node is replaced                  |
| Different HTML element attributes                          | Current node attributes are updated       |
| Different HTML element event handlers                      | Current node event handlers are updated   |
| Different component types                                  | Current node is replaced                  |
| Different exported fields on a same component type         | Current component fields are updated      |
| Different non-exported fields on a same component type     | No modifications                          |
| Extra node in the new the tree                             | Node added to the current tree            |
| Missing node in the new tree                               | Extra node is the current tree is removed |
