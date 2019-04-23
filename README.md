# dialog-loader

The dialog-loader Lambda is designed to load dialog content from the [Adaptive Dialog repo](https://github.com/adaptiveteam/adaptive-dialog) and then store it into a DynamoDB database for access the [Dialog Retriever](https://github.com/adaptiveteam/dialog-retriever) package. This Lambda will run every time new dialog is tagged in the [Adaptive Dialog repo](https://github.com/adaptiveteam/adaptive-dialog) ensuring the the platform always has access to the most up to date dialog.
