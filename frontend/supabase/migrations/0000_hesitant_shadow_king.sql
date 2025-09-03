CREATE TYPE "public"."day_of_week" AS ENUM('Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday');--> statement-breakpoint
CREATE TYPE "public"."enrollment_status" AS ENUM('enrolled', 'waitlisted', 'cancelled', 'no_show');--> statement-breakpoint
CREATE TYPE "public"."payment_status" AS ENUM('pending', 'paid', 'failed', 'refunded');--> statement-breakpoint
CREATE TABLE "attendance" (
	"id" uuid PRIMARY KEY NOT NULL,
	"session_enrollment_id" uuid NOT NULL,
	"attended" boolean DEFAULT false NOT NULL,
	"check_in_time" timestamp,
	"checked_in_by" uuid,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "unique_attendance_constraint" UNIQUE("session_enrollment_id")
);
--> statement-breakpoint
ALTER TABLE "attendance" ENABLE ROW LEVEL SECURITY;--> statement-breakpoint
CREATE TABLE "auth"."users" (
	"id" uuid PRIMARY KEY NOT NULL,
	"email" text
);
--> statement-breakpoint
CREATE TABLE "class_sessions" (
	"id" uuid PRIMARY KEY NOT NULL,
	"class_id" uuid NOT NULL,
	"session_date" date NOT NULL,
	"start_time" time NOT NULL,
	"end_time" time NOT NULL,
	"max_capacity" integer DEFAULT 20 NOT NULL,
	"instructor_id" uuid,
	"is_cancelled" boolean DEFAULT false NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "unique_session_constraint" UNIQUE("class_id","session_date","start_time")
);
--> statement-breakpoint
ALTER TABLE "class_sessions" ENABLE ROW LEVEL SECURITY;--> statement-breakpoint
CREATE TABLE "classes" (
	"id" uuid PRIMARY KEY NOT NULL,
	"name" text NOT NULL,
	"description" text,
	"price_per_session" numeric(10, 2) NOT NULL,
	"instructor_id" uuid,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
ALTER TABLE "classes" ENABLE ROW LEVEL SECURITY;--> statement-breakpoint
CREATE TABLE "memberships" (
	"id" uuid PRIMARY KEY NOT NULL,
	"user_id" uuid NOT NULL,
	"class_id" uuid NOT NULL,
	"start_date" date NOT NULL,
	"end_date" date NOT NULL,
	"payment_status" "payment_status" DEFAULT 'pending' NOT NULL,
	"amount_paid" numeric(10, 2),
	"transaction_id" text,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "unique_membership_constraint" UNIQUE("user_id","class_id","start_date","end_date")
);
--> statement-breakpoint
ALTER TABLE "memberships" ENABLE ROW LEVEL SECURITY;--> statement-breakpoint
CREATE TABLE "profiles" (
	"id" uuid PRIMARY KEY NOT NULL,
	"full_name" text,
	"avatar_url" text,
	"is_instructor" boolean DEFAULT false NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
ALTER TABLE "profiles" ENABLE ROW LEVEL SECURITY;--> statement-breakpoint
CREATE TABLE "recurring_patterns" (
	"id" uuid PRIMARY KEY NOT NULL,
	"class_id" uuid NOT NULL,
	"day_of_week" "day_of_week" NOT NULL,
	"start_time" time NOT NULL,
	"duration_minutes" integer NOT NULL,
	"effective_from_date" date NOT NULL,
	"effective_to_date" date,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
CREATE TABLE "session_enrollments" (
	"id" uuid PRIMARY KEY NOT NULL,
	"user_id" uuid NOT NULL,
	"session_id" uuid NOT NULL,
	"membership_id" uuid NOT NULL,
	"enrollment_date" timestamp DEFAULT now() NOT NULL,
	"status" "enrollment_status" DEFAULT 'enrolled' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL,
	CONSTRAINT "unique_enrollment_constraint" UNIQUE("user_id","session_id")
);
--> statement-breakpoint
ALTER TABLE "session_enrollments" ENABLE ROW LEVEL SECURITY;--> statement-breakpoint
ALTER TABLE "attendance" ADD CONSTRAINT "attendance_session_enrollment_id_session_enrollments_id_fk" FOREIGN KEY ("session_enrollment_id") REFERENCES "public"."session_enrollments"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "attendance" ADD CONSTRAINT "attendance_checked_in_by_profiles_id_fk" FOREIGN KEY ("checked_in_by") REFERENCES "public"."profiles"("id") ON DELETE set null ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "class_sessions" ADD CONSTRAINT "class_sessions_class_id_classes_id_fk" FOREIGN KEY ("class_id") REFERENCES "public"."classes"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "class_sessions" ADD CONSTRAINT "class_sessions_instructor_id_profiles_id_fk" FOREIGN KEY ("instructor_id") REFERENCES "public"."profiles"("id") ON DELETE set null ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "classes" ADD CONSTRAINT "classes_instructor_id_profiles_id_fk" FOREIGN KEY ("instructor_id") REFERENCES "public"."profiles"("id") ON DELETE set null ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "memberships" ADD CONSTRAINT "memberships_user_id_profiles_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."profiles"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "memberships" ADD CONSTRAINT "memberships_class_id_classes_id_fk" FOREIGN KEY ("class_id") REFERENCES "public"."classes"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "profiles" ADD CONSTRAINT "profiles_id_users_id_fk" FOREIGN KEY ("id") REFERENCES "auth"."users"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "recurring_patterns" ADD CONSTRAINT "recurring_patterns_class_id_classes_id_fk" FOREIGN KEY ("class_id") REFERENCES "public"."classes"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session_enrollments" ADD CONSTRAINT "session_enrollments_user_id_profiles_id_fk" FOREIGN KEY ("user_id") REFERENCES "public"."profiles"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session_enrollments" ADD CONSTRAINT "session_enrollments_session_id_class_sessions_id_fk" FOREIGN KEY ("session_id") REFERENCES "public"."class_sessions"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
ALTER TABLE "session_enrollments" ADD CONSTRAINT "session_enrollments_membership_id_memberships_id_fk" FOREIGN KEY ("membership_id") REFERENCES "public"."memberships"("id") ON DELETE cascade ON UPDATE no action;--> statement-breakpoint
CREATE POLICY "select_own_attendance" ON "attendance" AS PERMISSIVE FOR SELECT TO "authenticated" USING ((select auth.uid()) = (SELECT user_id FROM session_enrollments se WHERE se.id = "attendance"."session_enrollment_id"));--> statement-breakpoint
CREATE POLICY "mark_attendance" ON "attendance" AS PERMISSIVE FOR INSERT TO "authenticated" WITH CHECK (
          (SELECT is_instructor FROM profiles WHERE id = (select auth.uid())) = TRUE
          AND (select auth.uid()) = (
            SELECT cs.instructor_id
            FROM session_enrollments se
            JOIN class_sessions cs ON se.session_id = cs.id
            WHERE se.id = "attendance"."session_enrollment_id"
          )
        );--> statement-breakpoint
CREATE POLICY "update_attendance" ON "attendance" AS PERMISSIVE FOR UPDATE TO "authenticated" USING (
          (SELECT is_instructor FROM profiles WHERE id = (select auth.uid())) = TRUE
          AND (select auth.uid()) = (
            SELECT cs.instructor_id
            FROM session_enrollments se
            JOIN class_sessions cs ON se.session_id = cs.id
            WHERE se.id = "attendance"."session_enrollment_id"
          )
        ) WITH CHECK (
          (SELECT is_instructor FROM profiles WHERE id = (select auth.uid())) = TRUE
          AND (select auth.uid()) = (
            SELECT cs.instructor_id
            FROM session_enrollments se
            JOIN class_sessions cs ON se.session_id = cs.id
            WHERE se.id = "attendance"."session_enrollment_id"
          )
        );--> statement-breakpoint
CREATE POLICY "view_all_sessions" ON "class_sessions" AS PERMISSIVE FOR SELECT TO "authenticated" USING (true);--> statement-breakpoint
CREATE POLICY "view_all_classes" ON "classes" AS PERMISSIVE FOR SELECT TO "authenticated" USING (true);--> statement-breakpoint
CREATE POLICY "manage_own_classes" ON "classes" AS PERMISSIVE FOR ALL TO "authenticated" USING ((select auth.uid()) = "classes"."instructor_id" AND (SELECT is_instructor FROM profiles WHERE id = (select auth.uid())) = TRUE) WITH CHECK ((select auth.uid()) = "classes"."instructor_id" AND (SELECT is_instructor FROM profiles WHERE id = (select auth.uid())) = TRUE);--> statement-breakpoint
CREATE POLICY "select_own_membership" ON "memberships" AS PERMISSIVE FOR SELECT TO "authenticated" USING ((select auth.uid()) = "memberships"."user_id");--> statement-breakpoint
CREATE POLICY "insert_own_membership" ON "memberships" AS PERMISSIVE FOR INSERT TO "authenticated" WITH CHECK ((select auth.uid()) = "memberships"."user_id");--> statement-breakpoint
CREATE POLICY "update_own_membership" ON "memberships" AS PERMISSIVE FOR UPDATE TO "authenticated" USING ((select auth.uid()) = "memberships"."user_id") WITH CHECK ((select auth.uid()) = "memberships"."user_id");--> statement-breakpoint
CREATE POLICY "select_own_profile" ON "profiles" AS PERMISSIVE FOR SELECT TO "authenticated" USING ((select auth.uid()) = "profiles"."id");--> statement-breakpoint
CREATE POLICY "update_own_profile" ON "profiles" AS PERMISSIVE FOR UPDATE TO "authenticated" USING ((select auth.uid()) = "profiles"."id") WITH CHECK ((select auth.uid()) = "profiles"."id");--> statement-breakpoint
CREATE POLICY "select_own_enrollment" ON "session_enrollments" AS PERMISSIVE FOR SELECT TO "authenticated" USING ((select auth.uid()) = "session_enrollments"."user_id");--> statement-breakpoint
CREATE POLICY "insert_own_enrollment" ON "session_enrollments" AS PERMISSIVE FOR INSERT TO "authenticated" WITH CHECK ((select auth.uid()) = "session_enrollments"."user_id" AND EXISTS (
          SELECT 1 FROM memberships m
          WHERE m.id = "session_enrollments"."membership_id"
            AND m.user_id = (select auth.uid())
            AND m.class_id = (SELECT class_id FROM class_sessions cs WHERE cs.id = "session_enrollments"."session_id")
            AND m.start_date <= (SELECT session_date FROM class_sessions cs WHERE cs.id = "session_enrollments"."session_id")
            AND m.end_date >= (SELECT session_date FROM class_sessions cs WHERE cs.id = "session_enrollments"."session_id")
            AND m.payment_status = 'paid'
        ));--> statement-breakpoint
CREATE POLICY "update_own_enrollment" ON "session_enrollments" AS PERMISSIVE FOR UPDATE TO "authenticated" USING ((select auth.uid()) = "session_enrollments"."user_id") WITH CHECK ((select auth.uid()) = "session_enrollments"."user_id");--> statement-breakpoint
CREATE POLICY "delete_own_enrollment" ON "session_enrollments" AS PERMISSIVE FOR DELETE TO "authenticated" USING ((select auth.uid()) = "session_enrollments"."user_id" AND "session_enrollments"."status" IN ('enrolled', 'waitlisted'));