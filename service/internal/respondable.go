// TODO: return respondible object
// new respondable(quote, error)
/**
  export interface Respondable<T> {
    data?: T;
    success: boolean;
    error?: ResponseError; // will be the `error` type.
  	// Go does not have Exceptions. error is simpler
  }

  export class ResponseOK<T> implements ServiceResponse<T> {
    success: boolean;
    data?: T | undefined;
    error?: ResponseError | undefined;
    constructor(data?: T, success = true) {
      this.success = success;
      this.data = data;
    }
  }

  export class ResponseFailed implements ServiceResponse<ErrorTypeName> {
    success: boolean;
    error?: ResponseError | undefined;
    params?: any[];
    constructor(type: ErrorTypeName, errorMessage?: string, params?: any[]) {
      this.success = false;
      this.error = this.errorFrom(type, errorMessage);
      this.params = params;
    }

    private errorFrom(level: ErrorTypeName, message?: string): ResponseError {
      switch (level) {
        case ErrorTypeName.unknown:
          return new ResponseError(ErrorTypeName.unknown, message);

        case ErrorTypeName.service:
          return new ResponseError(ErrorTypeName.service, message);

        case ErrorTypeName.notFound:
          return new ResponseError(ErrorTypeName.notFound, message);

        case ErrorTypeName.internal:
          return new ResponseError(ErrorTypeName.internal, message);

        case ErrorTypeName.bad:
          return new ResponseError(ErrorTypeName.bad, message);

        case ErrorTypeName.db:
          return new ResponseError(ErrorTypeName.db, message);

        case ErrorTypeName.notAcceptable:
          return new ResponseError(ErrorTypeName.notAcceptable, message);

        default:
          return new ResponseError(ErrorTypeName.unknown, message);
      }
    }
  }
*/