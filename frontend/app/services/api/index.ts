/**
 * This Api class lets you define an API endpoint and methods to request
 * data and process it.
 *
 * See the [Backend API Integration](https://docs.infinite.red/ignite-cli/boilerplate/app/services/#backend-api-integration)
 * documentation for more details.
 */
import { ApiResponse, ApisauceInstance, create } from "apisauce"

import Config from "@/config"

import type { ApiConfig, LoginRequest, LoginResponse, RegisterRequest, UserProfileResponse } from "./types"
import { courses } from "./mockdata"
import { FormValidationError, GeneralApiProblem, getFormValidationProblem, getGeneralApiProblem } from "./apiProblem"

/**
 * Configuring the apisauce instance.
 */
export const DEFAULT_API_CONFIG: ApiConfig = {
  url: Config.API_URL,
  timeout: 10000,
}

/**
 * Manages all requests to the API. You can use this class to build out
 * various requests that you need to call from your backend API.
 */
export class Api {
  apisauce: ApisauceInstance
  config: ApiConfig

  /**
   * Set up our API instance. Keep this lightweight!
   */
  constructor(config: ApiConfig = DEFAULT_API_CONFIG) {
    console.debug(config)
    this.config = config
    this.apisauce = create({
      baseURL: this.config.url,
      timeout: this.config.timeout,
      headers: {
        Accept: "application/json",
      },
    })
  }


  setAuthHeader(token: string) {
    this.apisauce.setHeader("Authorization", `Bearer ${token}`)
    console.log("token set", token)
  }

  async Register(data: RegisterRequest): Promise<{ kind: string, done: boolean } | GeneralApiProblem | FormValidationError> {
    const response = await this.apisauce.post("register", data)
    console.log(response)
    try {

      if (!response.ok) {
        const validationProblem = getFormValidationProblem(response)
        if (validationProblem) return validationProblem
        const problem = getGeneralApiProblem(response)
        if (problem) return problem
      }


      return { kind: "ok", done: true }
    } catch (e) {
      if (__DEV__ && e instanceof Error) {
        console.error(`Bad data: ${e.message}\n${response.data}`, e.stack)
      }
      return { kind: "bad-data" }
    }
  }

  async Login(data: LoginRequest): Promise<LoginResponse | GeneralApiProblem> {
    const response: ApiResponse<LoginResponse> = await this.apisauce.post("/login", data)
    try {

      if (!response.ok) {
        const problem = getGeneralApiProblem(response)
        if (problem) return problem
      }


      return { kind: "ok", token: response.data?.token || "", refresh: response.data?.refresh || "" }
    } catch (e) {
      if (__DEV__ && e instanceof Error) {
        console.error(`Bad data: ${e.message}\n${response.data}`, e.stack)
      }
      return { kind: "bad-data" }
    }
  }

  async getProfile(): Promise<{ kind: "ok" } & UserProfileResponse | GeneralApiProblem> {
    const response: ApiResponse<UserProfileResponse> = await this.apisauce.get("/users/me")
    console.debug("getProfile", response)
    try {

      if (!response.ok) {
        const problem = getGeneralApiProblem(response)
        if (problem) return problem
      }

      if (response.data) {
        return { kind: "ok", ...response.data }
      }
      return { kind: "bad-data" }
    } catch (e) {
      if (__DEV__ && e instanceof Error) {
        console.error(`Bad data: ${e.message}\n${response.data}`, e.stack)
      }
      return { kind: "bad-data" }
    }

  }

  async getCourses() {
    return courses
  }
}

// Singleton instance of the API for convenience
export const api = new Api()
