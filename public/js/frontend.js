(function bootstrapFrontend() {
  var installForm = document.querySelector('.install-form');
  if (installForm) {
    installForm.addEventListener('submit', function(event) {
      event.preventDefault();

      var externalKeyField = installForm.querySelector('input[name="external_key"]');
      if (externalKeyField) {
        externalKeyField.value = Date.now() + '-' + crypto.randomUUID();
      }

      var startURL = new URL(installForm.action, window.location.origin);
      var data = new FormData(installForm);
      data.forEach(function(value, key) { startURL.searchParams.set(key, value); });

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
