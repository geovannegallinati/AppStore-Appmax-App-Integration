(function bootstrapFrontend() {
  var installForm = document.querySelector('.install-form');
  if (installForm) {
    var appIDField = installForm.querySelector('input[name="app_id"]');
    var externalKeyField = installForm.querySelector('input[name="external_key"]');
    if (externalKeyField && externalKeyField.value.trim() === '') {
      externalKeyField.value = 'install-' + Date.now();
    }

    installForm.addEventListener('submit', function(event) {
      event.preventDefault();

      var appID = appIDField ? appIDField.value.trim() : '';
      var externalKey = externalKeyField ? externalKeyField.value.trim() : '';
      if (!appID || !externalKey) {
        installForm.submit();
        return;
      }

      var startURL = new URL(installForm.action, window.location.origin);
      startURL.searchParams.set('app_id', appID);
      startURL.searchParams.set('external_key', externalKey);

      var newTab = window.open(startURL.toString(), '_blank', 'noopener,noreferrer');
      if (!newTab) {
        window.location.href = startURL.toString();
      }
    });
  }

  var copyButtons = document.querySelectorAll('[data-copy-target]');
  copyButtons.forEach(function(button) {
    button.addEventListener('click', function() {
      var targetSelector = button.getAttribute('data-copy-target');
      if (!targetSelector) {
        return;
      }

      var target = document.querySelector(targetSelector);
      if (!target) {
        return;
      }

      var text = target.textContent.trim();
      if (!text) {
        return;
      }

      navigator.clipboard.writeText(text).then(function() {
        var original = button.textContent;
        button.textContent = 'Copied';
        setTimeout(function() {
          button.textContent = original;
        }, 1200);
      }).catch(function() {
        button.textContent = 'Failed';
      });
    });
  });
})();
