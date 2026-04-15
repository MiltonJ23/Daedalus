/**
 * Daedalus — Shared Winston Logging Configuration
 * Use for frontend and any Node.js services.
 */

const LOG_LEVELS = {
  error: 0,
  warn: 1,
  info: 2,
  http: 3,
  debug: 4,
};

function createLogger(serviceName = "daedalus") {
  const level = process.env.LOG_LEVEL || "info";

  function formatMessage(logLevel, message, meta = {}) {
    return JSON.stringify({
      timestamp: new Date().toISOString(),
      level: logLevel,
      service: serviceName,
      message,
      ...meta,
    });
  }

  const logger = {};
  Object.keys(LOG_LEVELS).forEach((lvl) => {
    logger[lvl] = (message, meta) => {
      if (LOG_LEVELS[lvl] <= LOG_LEVELS[level]) {
        const formatted = formatMessage(lvl, message, meta);
        if (lvl === "error") {
          console.error(formatted);
        } else if (lvl === "warn") {
          console.warn(formatted);
        } else {
          console.log(formatted);
        }
      }
    };
  });

  return logger;
}

module.exports = { createLogger, LOG_LEVELS };
