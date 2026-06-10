package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// token ที่สร้างจะเก็บข้อมูลหลัก ๆ คือ user_id, email และเวลาหมดอายุ exp โดย token จะถูกเซ็นด้วย JWT_SECRET เพื่อป้องกันการปลอมแปลง
// หลังจากสร้างเสร็จ ระบบจะส่ง token นี้กลับไปให้ client ใช้แนบไปกับ request ที่ต้องการยืนยันตัวตน เช่น
// Authorization: Bearer <token>
// สรุป: ใช้สร้าง token สำหรับยืนยันตัวตนของผู้ใช้ในระบบ
func GenerateJWT(userID uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
