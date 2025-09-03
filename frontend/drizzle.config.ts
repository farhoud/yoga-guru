import { DATABASE_URL } from "@/db"
import { defineConfig } from "drizzle-kit"

export default defineConfig({
  schema: "./app/db/schema.ts",
  out: "./supabase/migrations",
  dialect: "postgresql",
  dbCredentials: {
    url: DATABASE_URL!,
  },
})
