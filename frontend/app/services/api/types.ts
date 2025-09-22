/**
 * These types indicate the shape of the data you expect to receive from your
 * API endpoint, assuming it's a JSON object like we have.
 */
export interface EpisodeItem {
  title: string
  pubDate: string
  link: string
  guid: string
  author: string
  thumbnail: string
  description: string
  content: string
  enclosure: {
    link: string
    type: string
    length: number
    duration: number
    rating: { scheme: string; value: string }
  }
  categories: string[]
}

export interface ApiFeedResponse {
  status: string
  feed: {
    url: string
    title: string
    link: string
    author: string
    description: string
    image: string
  }
  items: EpisodeItem[]
}

export interface Course {
  className?: string
  title: string
  courseType?: string
  schedule: string
  description: string
  level: "beginner" | "intermediate" | "advanced"
  price: number
  capacity: boolean
}


export interface RegisterRequest {
  name: string
  phone: string
  gender?: "male" | "female"
  password: string
  role?: "student" | "instructor"
}

export interface LoginRequest {
  phone: string
  password: string
}

export interface RefreshTokenRequest {
  refreshToken: string
}

export interface LoginResponse {
  token: string
  refresh: string
  role: "admin" | "instructor" | "user"
}

export interface UserProfileResponse {
  phone: string
  name: string
  gender: "male" | "female"
  avatar: string
}

/**
 * The options used to configure apisauce.
 */
export interface ApiConfig {
  /**
   * The URL of the api.
   */
  url: string

  /**
   * Milliseconds before we timeout the request.
   */
  timeout: number
}
