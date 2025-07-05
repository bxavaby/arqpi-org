function checkLocalStorage() {
  try {
    localStorage.setItem("test", "test");
    localStorage.removeItem("test");
    return true;
  } catch (e) {
    /* console.error("localStorage not available:", e); */
    return false;
  }
}

class ApiClient {
  constructor(baseUrl = "https://arqpi-org.onrender.com") {
    this.baseUrl = baseUrl;
    this.localStorageAvailable = checkLocalStorage();
    this.loadApiKey();

    /* console.log("Storage available:", this.localStorageAvailable); */
    /* console.log("API key loaded:", this.apiKey ? "Yes" : "No"); */
  }

  loadApiKey() {
    if (this.localStorageAvailable) {
      this.apiKey = localStorage.getItem("arqpi_key") || null;
    } else {
      this.apiKey = null;
    }
  }

  /**
   * GET request
   * @param {string} endpoint
   * @param {Object} params
   * @returns {Promise}
   */
  async get(endpoint, params = {}) {
    this.loadApiKey();

    try {
      const url = new URL(`${this.baseUrl}${endpoint}`);

      /* console.log("Making request to:", endpoint); */
      /* console.log("Current API key:", this.apiKey ? "Key exists" : "No key"); */

      Object.keys(params).forEach((key) => {
        if (params[key] !== undefined && params[key] !== null) {
          url.searchParams.append(key, params[key]);
        }
      });

      if (this.apiKey) {
        url.searchParams.append("key", this.apiKey);
        /* console.log("Added API key to request"); */
      }

      /* console.log("Final request URL:", url.toString()); */

      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      /* console.error("API request failed:", error); */
      throw error;
    }
  }

  setApiKey(key) {
    console.log(
      "ApiClient.setApiKey called with:",
      key ? "Key provided" : "Empty key",
    );

    if (key && key.trim() !== "") {
      this.apiKey = key;
      if (this.localStorageAvailable) {
        localStorage.setItem("arqpi_key", key);
        /* console.log("API key saved to localStorage"); */
      }
    } else {
      this.clearApiKey();
      /* console.log("API key cleared (empty value provided)"); */
    }
  }

  clearApiKey() {
    /* console.log("Clearing API key"); */
    this.apiKey = null;
    if (this.localStorageAvailable) {
      localStorage.removeItem("arqpi_key");
    }
  }
}

const api = new ApiClient();
