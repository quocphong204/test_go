package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"order-system/config"
)

func SendOrderConfirmationEmail(to string, name string, orderID uint, total float64) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.SMTPUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Xác nhận đơn hàng #%d", orderID))
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Xin chào %s,</h2>
		<p>Đơn hàng #%d của bạn đã được tạo thành công.</p>
		<p>Tổng tiền: <strong>%.2f</strong></p>
		<p>Cảm ơn bạn đã mua hàng!</p>
	`, name, orderID, total))

	d := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	return d.DialAndSend(m)
}
