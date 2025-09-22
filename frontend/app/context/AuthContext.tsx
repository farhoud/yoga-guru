import { supabase } from "@/services/supabase"
import { Session } from "@supabase/supabase-js"
import {
  createContext,
  FC,
  PropsWithChildren,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react"
import { useMMKVString, useMMKVObject } from "react-native-mmkv"
import { useForm, Controller } from "react-hook-form"
import { zodResolver } from "@hookform/resolvers/zod"
import { z } from "zod"
import { api } from "@/services/api"
import {
  useToast,
  Toast,
  ToastTitle,
  ToastDescription,
} from "@/components/ui/toast"
import { FormValidationError } from "@/services/api/apiProblem"

export type AuthContextType = {
  isAuthenticated: boolean
  logout: () => void
  login: (fromData: LoginFormData) => Promise<void>
  signup: (fromData: SignupFormData) => Promise<void>
  role?: string
  profile?: UserProfile
  loading: boolean
  error: string
}
// Zod schema for validation
export const signupSchema = z.object({
  name: z.string().min(3, "حداقل ۳  حرف"),
  phone: z
    .string()
    .regex(/^\+?[0-9]{10,14}$/, "شماره تلفن اشتباه است"),
  password: z.string().min(6, "حداقل ۶ حرف"),
  gender: z.literal(["male", "female"])
})

export type SignupFormData = z.infer<typeof signupSchema>

// Zod schema for validation
export const loginSchema = z.object({
  phone: z.e164(),
  password: z.string().min(6, "Password must be at least 6 characters"),
})


export interface UserProfile {
  name: string
  phone: string
  gender: "male" | "female"
  avatarUrl?: string
}

export type LoginFormData = z.infer<typeof loginSchema>

export const AuthContext = createContext<AuthContextType | null>(null)

export interface AuthProviderProps { }

export const AuthProvider: FC<PropsWithChildren<AuthProviderProps>> = ({ children }) => {
  const [token, setToken] = useMMKVString("authToken")
  const [refresh, setRefresh] = useMMKVString("authRefresh")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)
  const [profile, setProfile] = useMMKVObject<UserProfile>("userProfile")
  const [role, setRole] = useMMKVString("authRole")



  const logout = () => {
    setToken(undefined)
    setRefresh(undefined)
    setProfile(undefined)
    setRole(undefined)
  }

  const login = async (formData: LoginFormData) => {
    setLoading(true)
    setError("")
    const resp = await api.Login(formData)
    if (resp.kind === "ok") {
      setToken(resp.data.token)
      setRefresh(resp.data.refresh)
      setRole(resp.data.role)
    } else {
      setError(resp.kind)
    }
    setLoading(false)
  }

  const signup = async (formData: SignupFormData) => {
    setLoading(true)
    setError("")
    const resp = await api.Register(formData)
    console.debug(resp)
    if (resp.kind == "ok") {
      await login({ phone: formData.phone, password: formData.password })
    } else if (resp.kind === "validation") {
      setError((resp as FormValidationError).msg)
    } else {
      setError(resp.kind)
    }
    setLoading(false)
  }

  useEffect(() => {
    if (token) {
      api.setAuthHeader(token)
      api.getProfile().then(res => {
        if (res.kind === "ok") {
          console.log("set user", res.data)
          setProfile(res.data)
        } else if (res.kind === "expired-token" && refresh) {
          console.log("else", res)
          return api.Refresh({ refreshToken: refresh })
        } else {
          logout()
        }
      }).then(resp => {
        if (resp && resp.kind === "ok") {
          setToken(resp.data.token)
          setRefresh(resp.data.refresh)
          setRole(resp.data.role)
        } else {
          setError(resp?.kind || "bad data")
        }
      }).catch(e => {
        console.log("error:", e)
        setError(String(e))
      })
    }
  }, [token])

  // useEffect(() => {
  //   if (session) {
  //     supabase.auth.setSession(session)
  //   } else {
  //     supabase.auth.getSession().then((s) => {
  //       if (s.error) {
  //         setError(s.error.message)
  //         return
  //       }
  //       setSession(s.data.session)
  //     })
  //   }
  //   supabase.auth.onAuthStateChange((_event, session) => {
  //     setSession(session)
  //   })
  // }, [])

  const value = {
    isAuthenticated: !!token,
    logout,
    login,
    signup,
    role,
    profile,
    error,
    loading,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) throw new Error("useAuth must be used within an AuthProvider")
  return context
}
