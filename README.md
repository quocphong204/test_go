📦 Order System – Hệ thống đặt hàng API với Golang
Dự án mô phỏng quy trình đặt hàng đơn giản với các chức năng hiện đại như: xác thực JWT, quản lý đơn hàng, thanh toán, xử lý hàng đợi với RabbitMQ, cache với Redis, và giao diện người dùng bằng HTML.

🛠 Công nghệ sử dụng
Backend: Golang + Gin-Gonic

Database: PostgreSQL (GORM ORM)

Authentication: JWT (golang-jwt)

Caching: Redis

Queue Processing: RabbitMQ (amqp091-go)

Email: gomail

Frontend: HTML + CSS + JS thuần

API Docs: Swaggo (có thể tích hợp)

🚀 Tính năng chính
✅ Đăng ký / Đăng nhập (JWT Token)

✅ Quản lý sản phẩm (chỉ Admin)

✅ Đặt hàng (người dùng)

✅ Gửi mail xác nhận đơn hàng

✅ Redis: cache đơn hàng gần nhất

✅ RabbitMQ: xử lý đơn hàng trễ (auto xử lý sau 5 phút)

✅ Quản lý trạng thái đơn: pending / processed

✅ Xử lý thanh toán đơn: unpaid / paid

✅ Giao diện đơn giản với HTML + JS

🧪 Cách chạy dự án
1. Clone & cấu hình
bash
Sao chép
Chỉnh sửa
git clone https://github.com/quocphong204/test_go.git
cd test_go
2. Cấu hình database
Tạo DB PostgreSQL tên: ordersystem

3. Cấu hình .env (nếu dùng) hoặc trong config/db.go
4. Cài đặt Redis & RabbitMQ (nên dùng Docker)
bash
Sao chép
Chỉnh sửa
docker run -d -p 6379:6379 redis
docker run -d -p 5672:5672 -p 15672:15672 rabbitmq:3-management
5. Cài Go modules
bash
Sao chép
Chỉnh sửa
go mod tidy
6. Chạy web server + consumer:
bash
Sao chép
Chỉnh sửa
go run main.go
# mở terminal khác:
go run cmd/consumer.go
🌐 Giao diện người dùng (frontend)
index.html: Đăng nhập

products.html: Danh sách sản phẩm

create_order.html: Đặt hàng

orders.html: Xem & xử lý đơn

Truy cập tại:
📍 http://localhost:8080/index.html

📂 Cấu trúc thư mục
css
Sao chép
Chỉnh sửa
.
├── cmd/                  # RabbitMQ consumer
├── config/               # Kết nối DB, Redis, MQ
├── controllers/          # Xử lý request API
├── models/               # Định nghĩa bảng
├── routes/               # Khai báo route
├── utils/                # Hàm phụ trợ (JWT, Mail, Hash)
├── frontend/             # HTML/CSS/JS giao diện người dùng
├── main.go               # Khởi động server
└── go.mod
📧 Liên hệ
Tác giả: Quốc Phong
Email: phongdev@example.com
GitHub: quocphong204
