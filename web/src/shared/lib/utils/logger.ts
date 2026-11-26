import { isProdEnv } from 'shared/config';

/**
 * Tagged logger. Silent in production unless explicitly enabled. No PII: never log tokens or e-mails.
 */
export class Logger {
  constructor(
    private _tag: string,
    private _isActive = true,
    private _runOnProduction = false,
  ) {}

  log(message: string, ...args: unknown[]) {
    if (this._isActive && isProdEnv === this._runOnProduction)
      // eslint-disable-next-line no-console
      console.log(`[${this._tag}]: ${message}`, ...args);
  }
}
