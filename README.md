# Order Cleaning System (Bewell Test)

## 📋 ภาพรวม
ระบบ Order Cleaning เป็นระบบสำหรับจัดการและประมวลผลคำสั่งซื้อของฟิล์มกันรอย โดยสามารถแปลงข้อมูลจาก `InputOrder` ให้เป็น `CleanedOrder` ที่มีข้อมูลสะอาดและสมบูรณ์

## 🏗️ โครงสร้างโปรเจค

```
bewell_test/
├── models/
│   └── order.go          # Data structures สำหรับ InputOrder และ CleanedOrder
├── constants/
│   └── constant.go       # ค่าคงที่ทั้งหมดของระบบ
├── services/
│   └── cleaner.go        # Logic หลักในการประมวลผลคำสั่งซื้อ
└── test/
    └── cleaner_test.go   # Test cases ครอบคลุมทุกกรณี
```

## 📊 Data Models

### InputOrder
```go
type InputOrder struct {
    No                int     // หมายเลขคำสั่งซื้อ
    PlatformProductId string  // รหัสสินค้าจาก platform (อาจมีข้อมูลเสีย)
    Qty               int     // จำนวนที่สั่ง
    UnitPrice         float64 // ราคาต่อหน่วย (ไม่ใช้ในการคำนวณ)
    TotalPrice        float64 // ราคารวม
}
```

### CleanedOrder
```go
type CleanedOrder struct {
    No         int     // หมายเลขลำดับ
    ProductId  string  // รหัสสินค้าที่สะอาด
    MaterialId string  // รหัสวัสดุ (เช่น FG0A-CLEAR, FG0A-MATTE)
    ModelId    string  // รหัสรุ่น (เช่น IPHONE16PROMAX, OPPOA3)
    Qty        int     // จำนวนจริง (รวม multiplier)
    UnitPrice  float64 // ราคาต่อหน่วยที่คำนวณใหม่
    TotalPrice float64 // ราคารวม
}
```

## ⚙️ การทำงานหลัก

### 1. Product ID Cleaning
ระบบจะทำความสะอาด `PlatformProductId` โดย:
- ลบข้อมูลขยะ (เช่น `x2-3&`, `%20x`)
- ดึงเฉพาะ Product ID จริง (รูปแบบ: `FG[0-9A-Z]+-[A-Z]+-[A-Z0-9-]+`)
- จัดการ Quantity Multiplier (เช่น `*3` หมายถึงคูณ 3)

### 2. Bundle Order Processing
สำหรับคำสั่งซื้อแบบ Bundle (มี `/` คั่น):
- แยกแต่ละสินค้าออกจากกัน
- คำนวณราคาต่อหน่วยโดยแบ่งจำนวนรวม
- สร้าง CleanedOrder แยกต่างหาก

### 3. Free Items Generation
ระบบจะสร้างของแถมอัตโนมัติ:
- **WIPING-CLOTH**: ผ้าเช็ด (จำนวนเท่ากับสินค้าหลัก)
- **[TEXTURE]-CLEANNER**: น้ำยาทำความสะอาด (ตามประเภทวัสดุ)

## 🧪 Test Cases

### Test Case 1: Basic Order
```
Input: "FG0A-CLEAR-IPHONE16PROMAX", Qty: 2
Output: 3 orders (หลัก + ผ้าเช็ด + น้ำยา)
```

### Test Case 2: Dirty Product ID
```
Input: "x2-3&FG0A-CLEAR-IPHONE16PROMAX", Qty: 2
Output: ทำความสะอาด ID แล้วประมวลผลตามปกติ
```

### Test Case 3: Quantity Multiplier
```
Input: "FG0A-MATTE-IPHONE16PROMAX*3", Qty: 1
Output: จำนวนจริง = 3, ราคาต่อหน่วย = 30
```

### Test Case 5: Bundle Order
```
Input: "FG0A-CLEAR-OPPOA3/FG0A-CLEAR-OPPOA3-B/FG0A-MATTE-OPPOA3"
Output: 6 orders (3 หลัก + ผ้าเช็ด + 2 น้ำยา)
```

## 🔧 การใช้งาน

### การรันโปรแกรม
```bash
go run cmd/main.go
```

### การรันเทส
```bash
go test ./...
```

### การรันเทสเฉพาะ
```bash
go test ./test -v
```

## ⚡ Constants

ระบบใช้ค่าคงที่จาก `constants/constant.go`:

```go
// Product related constants
const (
    WipingClothProductId = "WIPING-CLOTH"
    CleanerSuffix       = "-CLEANNER"
    BundleSeparator     = "/"
    ProductIdSeparator  = "-"
)

// Regex patterns
const (
    ProductIdPattern    = `(FG[0-9A-Z]+-[A-Z]+-[A-Z0-9\-]+)(\*\d+)?`
    QtyMultiplierPattern = `\*(\d+)`
)

// Default values
const (
    DefaultQtyMultiplier = 1
    DefaultUnitPrice     = 0.0
    DefaultTotalPrice    = 0.0
    MinProductParts      = 3
)
```

## 🚀 Features

- ✅ **Product ID Cleaning**: ลบข้อมูลขยะและดึงข้อมูลสำคัญ
- ✅ **Bundle Processing**: จัดการคำสั่งซื้อแบบรวม
- ✅ **Quantity Multiplier**: รองรับการคูณจำนวน
- ✅ **Free Items**: สร้างของแถมอัตโนมัติ
- ✅ **Price Calculation**: คำนวณราคาต่อหน่วยใหม่
- ✅ **Constants Management**: จัดการค่าคงที่อย่างเป็นระบบ
- ✅ **Comprehensive Testing**: Test cases ครอบคลุมทุกกรณี

## 📝 หมายเหตุ

- ระบบจะคำนวณ `UnitPrice` ใหม่จาก `TotalPrice / RealQty`
- `UnitPrice` จาก InputOrder จะไม่ถูกใช้ในการคำนวณ
- สำหรับ Bundle Order ราคาจะถูกแบ่งตามสัดส่วนจำนวนสินค้า
- ของแถมจะมี `UnitPrice` และ `TotalPrice` เป็น 0 