import {
  pgTable,
  uuid,
  text,
  timestamp,
  numeric,
  integer,
  date,
  time,
  boolean,
  pgSchema,
  pgEnum,
  pgPolicy,
  unique,
} from "drizzle-orm/pg-core"
import { relations, sql } from "drizzle-orm"
import { authenticatedRole, authUid } from "drizzle-orm/supabase" // Supabase RLS helpers

// --- Enums ---
export const dayOfWeekEnum = pgEnum("day_of_week", [
  "Sunday",
  "Monday",
  "Tuesday",
  "Wednesday",
  "Thursday",
  "Friday",
  "Saturday",
])

export const paymentStatusEnum = pgEnum("payment_status", ["pending", "paid", "failed", "refunded"])

export const enrollmentStatusEnum = pgEnum("enrollment_status", [
  "enrolled",
  "waitlisted",
  "cancelled",
  "no_show", // Added for clarity with attendance
])

// --- Supabase Auth Users Table (Reference Only) ---
export const authUsers = pgSchema("auth").table("users", {
  id: uuid("id").primaryKey(),
  email: text("email"),
})

// --- Profiles Table ---
export const profiles = pgTable(
  "profiles",
  {
    id: uuid("id")
      .primaryKey()
      .$defaultFn(() => sql`gen_random_uuid()`)
      .references(() => authUsers.id, { onDelete: "cascade" }),
    fullName: text("full_name"),
    avatarUrl: text("avatar_url"),
    isInstructor: boolean("is_instructor").notNull().default(false), // Added for instructor roles
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => {
    return [
      pgPolicy("select_own_profile", {
        for: "select",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.id}`,
      }),
      pgPolicy("update_own_profile", {
        for: "update",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.id}`,
        withCheck: sql`${authUid} = ${table.id}`,
      }),
    ]
  },
)

export const profilesRelations = relations(profiles, ({ one, many }) => ({
  user: one(authUsers, {
    fields: [profiles.id],
    references: [authUsers.id],
  }),
  instructedClasses: many(classes),
  memberships: many(memberships), // Relationship to new memberships table
  sessionEnrollments: many(sessionEnrollments), // Relationship to new session_enrollments
}))

// --- Classes Table ---
export const classes = pgTable(
  "classes",
  {
    id: uuid("id")
      .primaryKey()
      .$defaultFn(() => sql`gen_random_uuid()`),
    name: text("name").notNull(),
    description: text("description"),
    pricePerSession: numeric("price_per_session", {
      precision: 10,
      scale: 2,
    }).notNull(),
    instructorId: uuid("instructor_id").references(() => profiles.id, {
      onDelete: "set null",
    }),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => {
    return [
      pgPolicy("view_all_classes", {
        for: "select",
        to: authenticatedRole,
        using: sql`true`,
      }),
      // Instructors can create/update their own classes
      pgPolicy("manage_own_classes", {
        for: "all", // Apply to insert, update, delete
        to: authenticatedRole,
        using: sql`${authUid} = ${table.instructorId} AND (SELECT is_instructor FROM profiles WHERE id = ${authUid}) = TRUE`,
        withCheck: sql`${authUid} = ${table.instructorId} AND (SELECT is_instructor FROM profiles WHERE id = ${authUid}) = TRUE`,
      }),
    ]
  },
)

export const classesRelations = relations(classes, ({ one, many }) => ({
  instructor: one(profiles, {
    fields: [classes.instructorId],
    references: [profiles.id],
  }),
  recurringPatterns: many(recurringPatterns),
  classSessions: many(classSessions),
  memberships: many(memberships), // Relationship to memberships table
}))

// --- Recurring Patterns Table ---
export const recurringPatterns = pgTable("recurring_patterns", {
  id: uuid("id")
    .primaryKey()
    .$defaultFn(() => sql`gen_random_uuid()`),
  classId: uuid("class_id")
    .notNull()
    .references(() => classes.id, { onDelete: "cascade" }),
  dayOfWeek: dayOfWeekEnum("day_of_week").notNull(),
  startTime: time("start_time").notNull(),
  durationMinutes: integer("duration_minutes").notNull(),
  effectiveFromDate: date("effective_from_date").notNull(),
  effectiveToDate: date("effective_to_date"),
  createdAt: timestamp("created_at").notNull().defaultNow(),
  updatedAt: timestamp("updated_at").notNull().defaultNow(),
})

export const recurringPatternsRelations = relations(recurringPatterns, ({ one }) => ({
  class: one(classes, {
    fields: [recurringPatterns.classId],
    references: [classes.id],
  }),
}))

// --- Class Sessions Table ---
export const classSessions = pgTable(
  "class_sessions",
  {
    id: uuid("id")
      .primaryKey()
      .$defaultFn(() => sql`gen_random_uuid()`),
    classId: uuid("class_id")
      .notNull()
      .references(() => classes.id, { onDelete: "cascade" }),
    sessionDate: date("session_date").notNull(),
    startTime: time("start_time").notNull(),
    endTime: time("end_time").notNull(),
    maxCapacity: integer("max_capacity").notNull().default(20),
    instructorId: uuid("instructor_id").references(() => profiles.id, {
      onDelete: "set null",
    }),
    isCancelled: boolean("is_cancelled").notNull().default(false),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => {
    return [
      // Unique constraint to prevent duplicate sessions for the same class on the same date/time
      unique("unique_session_constraint").on(table.classId, table.sessionDate, table.startTime),
      pgPolicy("view_all_sessions", {
        for: "select",
        to: authenticatedRole,
        using: sql`true`,
      }),
    ]
  },
)

export const classSessionsRelations = relations(classSessions, ({ one, many }) => ({
  class: one(classes, {
    fields: [classSessions.classId],
    references: [classes.id],
  }),
  instructor: one(profiles, {
    fields: [classSessions.instructorId],
    references: [profiles.id],
  }),
  sessionEnrollments: many(sessionEnrollments), // New relationship to session_enrollments
}))

// --- NEW: Memberships Table ---
// Represents a user's subscription to a class for a period (e.g., monthly access)
export const memberships = pgTable(
  "memberships",
  {
    id: uuid("id")
      .primaryKey()
      .$defaultFn(() => sql`gen_random_uuid()`),
    userId: uuid("user_id")
      .notNull()
      .references(() => profiles.id, { onDelete: "cascade" }),
    classId: uuid("class_id")
      .notNull()
      .references(() => classes.id, { onDelete: "cascade" }),
    startDate: date("start_date").notNull(),
    endDate: date("end_date").notNull(),
    paymentStatus: paymentStatusEnum("payment_status").notNull().default("pending"),
    amountPaid: numeric("amount_paid", { precision: 10, scale: 2 }),
    transactionId: text("transaction_id"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => {
    return [
      // A user can only have one active membership for a given class at a time (simplified)
      unique("unique_membership_constraint").on(
        table.userId,
        table.classId,
        table.startDate,
        table.endDate, // Consider overlaps for more complex scenarios
      ),
      // RLS: Users can only manage their own memberships
      pgPolicy("select_own_membership", {
        for: "select",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.userId}`,
      }),
      pgPolicy("insert_own_membership", {
        for: "insert",
        to: authenticatedRole,
        withCheck: sql`${authUid} = ${table.userId}`,
      }),
      pgPolicy("update_own_membership", {
        for: "update",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.userId}`,
        withCheck: sql`${authUid} = ${table.userId}`,
      }),
    ]
  },
)

export const membershipsRelations = relations(memberships, ({ one, many }) => ({
  user: one(profiles, {
    fields: [memberships.userId],
    references: [profiles.id],
  }),
  class: one(classes, {
    fields: [memberships.classId],
    references: [classes.id],
  }),
  sessionEnrollments: many(sessionEnrollments), // Memberships lead to enrollments
}))

// --- NEW: Session Enrollments Table (formerly 'registrations') ---
// Represents a user explicitly signing up for a specific class session,
// assuming they have an active membership that covers it.
export const sessionEnrollments = pgTable(
  "session_enrollments",
  {
    id: uuid("id")
      .primaryKey()
      .$defaultFn(() => sql`gen_random_uuid()`),
    userId: uuid("user_id")
      .notNull()
      .references(() => profiles.id, { onDelete: "cascade" }),
    sessionId: uuid("session_id")
      .notNull()
      .references(() => classSessions.id, { onDelete: "cascade" }),
    membershipId: uuid("membership_id") // Link to the membership that covers this enrollment
      .notNull()
      .references(() => memberships.id, { onDelete: "cascade" }),
    enrollmentDate: timestamp("enrollment_date").notNull().defaultNow(),
    status: enrollmentStatusEnum("status").notNull().default("enrolled"),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => {
    return [
      unique("unique_enrollment_constraint").on(table.userId, table.sessionId),
      // RLS: Users can manage their own enrollments
      pgPolicy("select_own_enrollment", {
        for: "select",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.userId}`,
      }),
      pgPolicy("insert_own_enrollment", {
        for: "insert",
        to: authenticatedRole,
        withCheck: sql`${authUid} = ${table.userId} AND EXISTS (
          SELECT 1 FROM memberships m
          WHERE m.id = ${table.membershipId}
            AND m.user_id = ${authUid}
            AND m.class_id = (SELECT class_id FROM class_sessions cs WHERE cs.id = ${table.sessionId})
            AND m.start_date <= (SELECT session_date FROM class_sessions cs WHERE cs.id = ${table.sessionId})
            AND m.end_date >= (SELECT session_date FROM class_sessions cs WHERE cs.id = ${table.sessionId})
            AND m.payment_status = 'paid'
        )`, // Complex check for eligibility
      }),
      pgPolicy("update_own_enrollment", {
        for: "update",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.userId}`,
        withCheck: sql`${authUid} = ${table.userId}`,
      }),
      pgPolicy("delete_own_enrollment", {
        for: "delete",
        to: authenticatedRole,
        using: sql`${authUid} = ${table.userId} AND ${table.status} IN ('enrolled', 'waitlisted')`, // Can only delete if not a 'no_show'
      }),
    ]
  },
)

export const sessionEnrollmentsRelations = relations(sessionEnrollments, ({ one, many }) => ({
  user: one(profiles, {
    fields: [sessionEnrollments.userId],
    references: [profiles.id],
  }),
  session: one(classSessions, {
    fields: [sessionEnrollments.sessionId],
    references: [classSessions.id],
  }),
  membership: one(memberships, {
    fields: [sessionEnrollments.membershipId],
    references: [memberships.id],
  }),
  attendance: one(attendance), // One-to-one relationship with attendance
}))

// --- NEW: Attendance Table ---
// Tracks whether a user actually attended a specific session they were enrolled in.
export const attendance = pgTable(
  "attendance",
  {
    id: uuid("id")
      .primaryKey()
      .$defaultFn(() => sql`gen_random_uuid()`),
    sessionEnrollmentId: uuid("session_enrollment_id")
      .notNull()
      .references(() => sessionEnrollments.id, { onDelete: "cascade" }),
    attended: boolean("attended").notNull().default(false),
    checkInTime: timestamp("check_in_time"),
    checkedInBy: uuid("checked_in_by").references(() => profiles.id, {
      onDelete: "set null", // Who marked the attendance (e.g., instructor or staff)
    }),
    createdAt: timestamp("created_at").notNull().defaultNow(),
    updatedAt: timestamp("updated_at").notNull().defaultNow(),
  },
  (table) => {
    return [
      unique("unique_attendance_constraint").on(table.sessionEnrollmentId),
      // RLS: Users can view their own attendance (via their enrollment)
      pgPolicy("select_own_attendance", {
        for: "select",
        to: authenticatedRole,
        using: sql`${authUid} = (SELECT user_id FROM session_enrollments se WHERE se.id = ${table.sessionEnrollmentId})`,
      }),
      // RLS: Instructors can mark attendance for sessions they teach
      pgPolicy("mark_attendance", {
        for: "insert",
        to: authenticatedRole,
        withCheck: sql`
          (SELECT is_instructor FROM profiles WHERE id = ${authUid}) = TRUE
          AND ${authUid} = (
            SELECT cs.instructor_id
            FROM session_enrollments se
            JOIN class_sessions cs ON se.session_id = cs.id
            WHERE se.id = ${table.sessionEnrollmentId}
          )
        `,
      }),
      pgPolicy("update_attendance", {
        for: "update",
        to: authenticatedRole,
        using: sql`
          (SELECT is_instructor FROM profiles WHERE id = ${authUid}) = TRUE
          AND ${authUid} = (
            SELECT cs.instructor_id
            FROM session_enrollments se
            JOIN class_sessions cs ON se.session_id = cs.id
            WHERE se.id = ${table.sessionEnrollmentId}
          )
        `,
        withCheck: sql`
          (SELECT is_instructor FROM profiles WHERE id = ${authUid}) = TRUE
          AND ${authUid} = (
            SELECT cs.instructor_id
            FROM session_enrollments se
            JOIN class_sessions cs ON se.session_id = cs.id
            WHERE se.id = ${table.sessionEnrollmentId}
          )
        `,
      }),
    ]
  },
)

export const attendanceRelations = relations(attendance, ({ one }) => ({
  sessionEnrollment: one(sessionEnrollments, {
    fields: [attendance.sessionEnrollmentId],
    references: [sessionEnrollments.id],
  }),
  markedBy: one(profiles, {
    fields: [attendance.checkedInBy],
    references: [profiles.id],
  }),
}))
