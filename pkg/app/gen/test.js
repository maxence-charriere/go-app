function goappNewNotification(notification) {
  goappShowNotification((title, options) => {
    try {
      const notification = new Notification(title, options);

      notification.onclick = (e) => {
        let target = options.target;
        if (!target) {
          target = "/";
        }

        window.location.href = target;
      };
    } catch (err) {
      console.log(err);
    }
  }, notification);
}

function goappShowNotification(showNotification, notification) {
  console.log(notification);

  const title = notification.title;
  delete notification.title;

  for (let action in notification.actions) {
    delete action.target;
  }

  showNotification(title, notification);
}
