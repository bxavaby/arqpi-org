class ApiSettings {
  constructor() {
    this.apiKey = localStorage.getItem("arqpi_key") || "";
    this.createSettingsUI();
  }

  createSettingsUI() {
    const footerLinks = document.querySelector(".zine-footer p");
    const apiKeyLink = document.createElement("span");
    apiKeyLink.innerHTML = ' | <a href="#" class="api-key-link">API Key</a>';
    footerLinks.appendChild(apiKeyLink);

    const modal = document.createElement("div");
    modal.className = "settings-modal hidden";
    modal.innerHTML = `
            <div class="settings-content">
                <h2>API KEY</h2>
                <p>Enter your API key to access unlimited requests.</p>
                <div class="api-key-form">
                    <input type="text" id="api-key-input" placeholder="Paste your key here" value="${this.apiKey}" />
                    <div class="button-group">
                        <button id="save-api-key">Save</button>
                        <button id="clear-api-key">Clear</button>
                    </div>
                </div>
                <p class="settings-info">
                    <small>Don't have a key? <a href="https://ko-fi.com/bxav" target="_blank">Support on Ko-fi</a> to get unlimited API access.</small>
                </p>
                <button class="close-modal">Close</button>
            </div>
        `;
    document.body.appendChild(modal);

    this.setupEventListeners(apiKeyLink.querySelector("a"), modal);
  }

  setupEventListeners(settingsBtn, modal) {
    settingsBtn.addEventListener("click", () => {
      modal.classList.toggle("hidden");
    });

    modal.querySelector(".close-modal").addEventListener("click", () => {
      modal.classList.add("hidden");
    });

    modal.addEventListener("click", (e) => {
      if (e.target === modal) {
        modal.classList.add("hidden");
      }
    });

    modal.querySelector("#save-api-key").addEventListener("click", () => {
      const keyInput = modal.querySelector("#api-key-input");
      this.saveApiKey(keyInput.value.trim());
      modal.classList.add("hidden");
      this.showNotification("API key saved!");
    });

    modal.querySelector("#clear-api-key").addEventListener("click", () => {
      modal.querySelector("#api-key-input").value = "";
      this.saveApiKey("");
      modal.classList.add("hidden");
      this.showNotification("API key cleared");
    });
  }

  saveApiKey(key) {
    this.apiKey = key;
    localStorage.setItem("arqpi_key", key);
    if (window.api) {
      window.api.setApiKey(key);
    }
  }

  showNotification(message) {
    const notification = document.createElement("div");
    notification.className = "notification";
    notification.textContent = message;
    document.body.appendChild(notification);

    setTimeout(() => {
      notification.classList.add("show");
      setTimeout(() => {
        notification.classList.remove("show");
        setTimeout(() => {
          document.body.removeChild(notification);
        }, 300);
      }, 2000);
    }, 10);
  }
}
