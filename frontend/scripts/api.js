class ApiClient {
  constructor(baseUrl = "https://arqpi-org.onrender.com") {
    this.baseUrl = baseUrl;
    this.apiKey = localStorage.getItem("arqpi_key") || null;
  }

  /**
   * GET request
   * @param {string} endpoint
   * @param {Object} params
   * @returns {Promise}
   */
  async get(endpoint, params = {}) {
    try {
      const url = new URL(`${this.baseUrl}${endpoint}`);

      Object.keys(params).forEach((key) => {
        if (params[key] !== undefined && params[key] !== null) {
          url.searchParams.append(key, params[key]);
        }
      });

      if (this.apiKey) {
        url.searchParams.append("key", this.apiKey);
      }

      const response = await fetch(url);

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error("API request failed:", error);
      throw error;
    }
  }

  setApiKey(key) {
    this.apiKey = key;
    localStorage.setItem("arqpi_key", key);
  }

  clearApiKey() {
    this.apiKey = null;
    localStorage.removeItem("arqpi_key");
  }
}

const api = new ApiClient();
