import { Course } from "./types"

export const courses: Course[] = [
  {
    title: "یوگای هاتا برای مبتدیان",
    courseType: "هاتا",
    schedule: "دوشنبه و چهارشنبه‌ها، ساعت ۶ تا ۷ عصر (۸ هفته)",
    description: "یک معرفی ملایم از یوگا با تمرکز بر حرکات بدنی و تنفس.",
    level: "beginner",
    price: 50,
    capacity: true,
  },
  {
    title: "وین‌یاسا فلو – سطح متوسط",
    courseType: "وین‌یاسا",
    schedule: "سه‌شنبه و پنجشنبه‌ها، ساعت ۷ تا ۸ عصر (۶ هفته)",
    description: "حرکات پویا برای تقویت انعطاف‌پذیری و قدرت بدنی.",
    level: "intermediate",
    price: 75,
    capacity: true,
  },
  {
    title: "یوگای یین برای آرامش عمیق",
    courseType: "یین",
    schedule: "جمعه‌ها، ساعت ۵ تا ۶:۳۰ عصر (۴ هفته)",
    description: "کشش‌های آرام و طولانی برای بهبود گردش خون و رهایی ذهن.",
    level: "beginner",
    price: 40,
    capacity: false, // تکمیل ظرفیت
  },
  {
    title: "چالش پاور یوگا",
    courseType: "پاور",
    schedule: "شنبه‌ها، ساعت ۱۰ تا ۱۱:۳۰ صبح (۱۰ هفته)",
    description: "جلسات پرانرژی برای افزایش استقامت و قدرت.",
    level: "advanced",
    price: 120,
    capacity: true,
  },
  {
    title: "مدیتیشن صبحگاهی و فلو ملایم",
    courseType: "ترکیبی",
    schedule: "هر روز، ساعت ۷ تا ۷:۳۰ صبح (آنلاین)",
    description: "آغاز روز با آرامش مدیتیشن و حرکات ساده یوگا.",
    level: "beginner",
    price: 0, // رایگان
    capacity: true,
  },
]
