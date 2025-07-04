class PixelZine {
  constructor() {
    this.contentDisplay = document.getElementById("content-display");
    this.welcomePage = document.getElementById("welcome-page");
    this.inputContainer = document.getElementById("input-container");
    this.inputField = document.getElementById("input-field");
    this.submitBtn = document.getElementById("submit-btn");

    this.currentNavItem = null;
    this.initZine();
  }

  initZine() {
    this.setupNavListeners();
    this.setupInputListeners();
    this.animateWelcomeQuote();
  }

  animateWelcomeQuote() {
    const quote = document.querySelector(".quote");
    if (quote) {
      quote.classList.add("content-entry");
    }
  }

  setupNavListeners() {
    const navItems = document.querySelectorAll(".nav-item");

    navItems.forEach((item) => {
      item.addEventListener("click", () => {
        this.handleNavItemClick(item);
      });
    });
  }

  setupInputListeners() {
    this.submitBtn.addEventListener("click", () => {
      const input = this.inputField.value.trim();
      if (input) {
        this.executeRequest(input);
      }
    });

    this.inputField.addEventListener("keydown", (event) => {
      if (event.key === "Enter") {
        const input = this.inputField.value.trim();
        if (input) {
          this.executeRequest(input);
        }
      }
    });
  }

  /**
   * navigation item click
   * @param {HTMLElement} navItem
   */
  handleNavItemClick(navItem) {
    document.querySelectorAll(".nav-item").forEach((item) => {
      item.classList.remove("active");
    });

    navItem.classList.add("active");

    this.currentNavItem = navItem;
    const endpoint = navItem.dataset.endpoint;
    const needsInput = navItem.dataset.needsInput === "true";
    const inputParam = navItem.dataset.inputParam;

    this.welcomePage.classList.add("page-exit");

    setTimeout(() => {
      this.welcomePage.classList.remove("current");
      this.welcomePage.classList.remove("page-exit");

      this.contentDisplay.classList.add("current");
      this.contentDisplay.classList.add("page-enter");

      if (needsInput) {
        this.showInputField(inputParam);
        this.contentDisplay.innerHTML = `
                    <div class="page-title">${this.getNavItemTitle(navItem)}</div>
                    <div class="page-text">
                        <p>Enter ${inputParam === "id" ? "a fragment ID" : "your search term"} to continue.</p>
                    </div>
                `;
      } else {
        this.hideInputField();
        this.executeEndpoint(endpoint);
      }
    }, 300);
  }

  /**
   * Execute request with input
   * @param {string} inputValue
   */
  async executeRequest(inputValue) {
    if (!this.currentNavItem) return;

    const endpoint = this.currentNavItem.dataset.endpoint;
    const inputParam = this.currentNavItem.dataset.inputParam;

    const params = {};
    params[inputParam] = inputValue;

    try {
      this.contentDisplay.innerHTML = `
                <div class="loading">Loading ${this.getNavItemTitle(this.currentNavItem)}</div>
            `;

      const result = await api.get(endpoint, params);
      this.displayResult(result, endpoint);
    } catch (error) {
      this.showError(error.message);
    }
  }

  /**
   * Execute endpoint without input
   * @param {string} endpoint
   */
  async executeEndpoint(endpoint) {
    try {
      this.contentDisplay.innerHTML = `
                <div class="loading">Loading</div>
            `;

      const result = await api.get(endpoint);
      this.displayResult(result, endpoint);
    } catch (error) {
      this.showError(error.message);
    }
  }

  /**
   * @param {string} paramName
   */
  showInputField(paramName) {
    this.inputContainer.classList.remove("hidden");
    this.inputField.placeholder =
      paramName === "id" ? "Enter fragment ID..." : "Enter search term...";
    this.inputField.focus();
  }

  hideInputField() {
    this.inputContainer.classList.add("hidden");
    this.inputField.value = "";
  }

  /**
   * @param {HTMLElement} navItem
   * @returns {string} - title
   */
  getNavItemTitle(navItem) {
    const navText = navItem.querySelector(".nav-text").textContent;
    return navText;
  }

  /**
   * API result
   * @param {Object} data
   * @param {string} endpoint
   */
  displayResult(data, endpoint) {
    let title, content;

    switch (endpoint) {
      case "/fragment":
        title = `Fragment #${data.id || "Unknown"}`;
        content = this.formatFragment(data);
        break;
      case "/search":
        title = `Search Results`;
        content = this.formatSearchResults(data);
        break;
      case "/random":
        title = `Random Fragment`;
        content = this.formatFragment(data);
        break;
      case "/quote":
        title = `Quote`;
        content = this.formatQuote(data);
        break;
      case "/info":
        title = `API Information`;
        content = this.formatInfo(data);
        break;
      case "/status":
        title = `System Status`;
        content = this.formatStatus(data);
        break;
      default:
        title = `API Response`;
        content = `<pre>${JSON.stringify(data, null, 2)}</pre>`;
    }

    this.contentDisplay.classList.remove("page-enter");
    this.contentDisplay.classList.add("page-exit");

    setTimeout(() => {
      this.contentDisplay.classList.remove("page-exit");
      this.contentDisplay.classList.add("page-enter");

      this.contentDisplay.innerHTML = `
                <div class="page-title">${title}</div>
                <div class="page-text content-entry">
                    ${content}
                </div>
            `;
    }, 300);
  }

  formatFragment(fragment) {
    return `
            <div class="fragment">
                <div class="fragment-text">${fragment.text || "No text available"}</div>
                ${fragment.notes ? `<div class="fragment-notes">${fragment.notes}</div>` : ""}
                <div class="page-decoration"><span>·</span><span>·</span><span>·</span></div>
            </div>
        `;
  }

  formatSearchResults(results) {
    if (!results || !results.length) {
      return "<p>No results found.</p>";
    }

    let html = `<p>Found ${results.length} results:</p><div class="page-decoration"><span>·</span><span>·</span><span>·</span></div>`;

    results.forEach((result) => {
      html += `
                <div class="search-result">
                    <div class="result-id">Fragment #${result.id}</div>
                    <div class="result-text">${result.text.substring(0, 150)}${result.text.length > 150 ? "..." : ""}</div>
                </div>
                <div class="page-decoration"><span>·</span></div>
            `;
    });

    return html;
  }

  formatQuote(quote) {
    return `
            <div class="quote">
                <p>"${quote.text || "No text available"}"</p>
            </div>
            ${quote.author ? `<p class="author">— ${quote.author}</p>` : ""}
        `;
  }

  formatInfo(info) {
    const description =
      info.project_info?.description || "No description available";
    const heteronymDesc = info.heteronyms_info?.description || "";

    return `
          <div class="info-section">
              <p>${description}</p>
              <p class="info-detail">API Version ${info.project_info?.api_version || "1.0.0"}</p>
          </div>

          <div class="page-decoration"><span>·</span><span>·</span><span>·</span></div>

          <div class="info-section">
              <p class="info-highlight">${info.fragments_count?.toLocaleString() || 0} fragments</p>
              <p class="info-detail">from ${info.authors_count || 0} authors across ${info.categories_count || 0} categories</p>
          </div>

          <div class="page-decoration"><span>·</span><span>·</span><span>·</span></div>

          <div class="info-section">
              <p class="info-title">The Heteronyms</p>
              <p>${heteronymDesc}</p>

              <div class="main-heteronyms">
                  ${
                    info.heteronyms_info?.main_heteronyms
                      ?.map(
                        (name) => `<span class="heteronym-name">${name}</span>`,
                      )
                      .join("") || ""
                  }
              </div>
          </div>

          <div class="page-decoration"><span>·</span><span>·</span><span>·</span></div>

          <div class="info-section">
              <p class="info-detail">Source: <a href="http://${info.project_info?.source || ""}" target="_blank">${info.project_info?.source || ""}</a></p>
              <p class="info-detail small">${info.project_info?.license || ""}</p>
          </div>
      `;
  }

  formatStatus(status) {
    let formattedUptime = status.uptime || "Unknown";
    if (formattedUptime.includes("m") && formattedUptime.includes("s")) {
    } else if (typeof formattedUptime === "number") {
      const minutes = Math.floor(formattedUptime / 60);
      const seconds = Math.floor(formattedUptime % 60);
      formattedUptime = `${minutes}m ${seconds}s`;
    }

    return `
          <div class="status-section">
              <p class="status-indicator ${status.status === "operational" ? "status-operational" : "status-error"}">
                  Status: ${status.status || "Unknown"}
              </p>
          </div>

          <div class="page-decoration"><span>·</span></div>

          <div class="status-details">
              <p>Version: ${status.version || "1.0.0"}</p>
              <p>Uptime: ${formattedUptime}</p>
              <p>Total requests: ${status.request_count?.toLocaleString() || "0"}</p>
          </div>

          ${status.message ? `<p class="status-message">${status.message}</p>` : ""}
      `;
  }

  showError(message) {
    this.contentDisplay.innerHTML = `
            <div class="page-title">Error</div>
            <div class="page-text">
                <div class="error">${message}</div>
                <p>Please try again or select a different option.</p>
            </div>
        `;
  }
}

document.addEventListener("DOMContentLoaded", () => {
  window.pixelZine = new PixelZine();
  window.apiSettings = new ApiSettings();
});
