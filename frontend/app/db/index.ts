import { config } from "dotenv"
import { drizzle } from "drizzle-orm/postgres-js"
import postgres from "postgres"

export const DATABASE_URL =
  "postgresql://postgres.prejneqyouzchfnetywt:Bb8MGGCdheukX+xAaYc=@aws-0-us-east-2.pooler.supabase.com:6543/postgres"

config({ path: ".env" }) // or .env.local
const client = postgres(DATABASE_URL!)
export const db = drizzle({ client })
