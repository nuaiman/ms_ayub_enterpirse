export type Role = 'admin' | 'manager' | 'accounts' | 'staff'

export type IDType =
  | 'nid'
  | 'passport'
  | 'driving_license'
  | 'birth_certificate'
  | 'trade_license'
  | 'other'

export interface User {
  id: number
  name: string
  username: string
  email?: string | null
  phone?: string | null
  address?: string | null
  id_type?: IDType | null
  id_number?: string | null
  image_url?: string | null
  role: Role
  is_active: boolean
  created_at: string
  updated_at: string
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data: T
}

export interface AuthSession {
  access_token: string
  user: User
}

export interface ChangePasswordPayload {
  current_password: string
  new_password: string
}
