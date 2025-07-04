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
    const version = info.project_info?.api_version || "";

    return `
          <p>${description}</p>
          ${version ? `<p>API Version: ${version}</p>` : ""}

          <div class="page-decoration"><span>·</span><span>·</span><span>·</span></div>

          <p><strong>Collection</strong>: ${info.fragments_count.toLocaleString()} fragments</p>
          <p><strong>Authors</strong>: ${info.authors_count} heteronyms and personas</p>
          <p><strong>Categories</strong>: ${info.categories_count} categories</p>

          <div class="page-decoration"><span>·</span></div>

          <p><strong>Main Heteronyms</strong>:</p>
          <ul>
              ${info.heteronyms_info?.main_heteronyms?.map((name) => `<li>${name}</li>`).join("") || ""}
          </ul>
      `;
  }

  formatStatus(status) {
    return `
            <p>Uptime: ${status.uptime || "Unknown"}</p>
            <p>Total requests: ${status.requests || "0"}</p>
            <p>Fragments available: ${status.fragments || "0"}</p>
            ${status.message ? `<p>${status.message}</p>` : ""}
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
