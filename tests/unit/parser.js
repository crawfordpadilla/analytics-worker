// parser.js

/**
 * Parses a user agent string and extracts relevant information.
 * @param {string} userAgent The user agent string to parse.
 * @returns {object} An object containing the parsed user agent information.
 */
function parseUserAgent(userAgent) {
  if (!userAgent || typeof userAgent !== 'string') {
    return {
      browser: 'Unknown',
      os: 'Unknown',
      device: 'Unknown'
    };
  }

  const lowerUserAgent = userAgent.toLowerCase();

  let browser = 'Unknown';
  let os = 'Unknown';
  let device = 'Unknown';

  // Browser detection
  if (lowerUserAgent.includes('chrome')) {
    browser = 'Chrome';
  } else if (lowerUserAgent.includes('firefox')) {
    browser = 'Firefox';
  } else if (lowerUserAgent.includes('safari')) {
    browser = 'Safari';
  } else if (lowerUserAgent.includes('edge')) {
    browser = 'Edge';
  } else if (lowerUserAgent.includes('opera') || lowerUserAgent.includes('opr')) {
    browser = 'Opera';
  } else if (lowerUserAgent.includes('msie') || lowerUserAgent.includes('trident')) {
    browser = 'Internet Explorer';
  }

  // OS detection
  if (lowerUserAgent.includes('windows')) {
    os = 'Windows';
  } else if (lowerUserAgent.includes('mac os x')) {
    os = 'macOS';
  } else if (lowerUserAgent.includes('android')) {
    os = 'Android';
  } else if (lowerUserAgent.includes('ios')) {
    os = 'iOS';
  } else if (lowerUserAgent.includes('linux')) {
    os = 'Linux';
  }

  // Device detection
  if (lowerUserAgent.includes('mobile')) {
    device = 'Mobile';
  } else if (lowerUserAgent.includes('tablet')) {
    device = 'Tablet';
  } else {
    device = 'Desktop';
  }

  return {
    browser,
    os,
    device
  };
}

/**
 * Parses a URL and extracts relevant information.
 * @param {string} url The URL string to parse.
 * @returns {object} An object containing the parsed URL information.
 */
function parseUrl(url) {
  try {
    const parsedUrl = new URL(url);
    return {
      hostname: parsedUrl.hostname,
      pathname: parsedUrl.pathname,
      searchParams: Object.fromEntries(parsedUrl.searchParams)
    };
  } catch (error) {
    // Invalid URL
    return {
      hostname: null,
      pathname: null,
      searchParams: {}
    };
  }
}

/**
 * Parses a JSON string. Returns null if parsing fails.
 * @param {string} jsonString The JSON string to parse.
 * @returns {object|null} The parsed JSON object or null if parsing fails.
 */
function parseJson(jsonString) {
  try {
    return JSON.parse(jsonString);
  } catch (error) {
    return null;
  }
}

module.exports = {
  parseUserAgent,
  parseUrl,
  parseJson
};